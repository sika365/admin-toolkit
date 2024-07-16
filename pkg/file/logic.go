package file

import (
	"net/url"
	"regexp"

	simutils "github.com/alifakhimi/simple-utils-go"
	"github.com/alitto/pond"
	"github.com/sika365/admin-tools/context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Logic interface {
	ReadFiles(ctx *context.Context, root string, maxDepth int, reContentType *regexp.Regexp, filters url.Values) (files MapFiles, err error)
	Sync(ctx *context.Context, files MapFiles, batchSize int) error
}

type logic struct {
	conn *simutils.DBConnection
	repo Repo
}

func newLogic(repo Repo, conn *simutils.DBConnection) (Logic, error) {
	l := &logic{
		conn: conn,
		repo: repo,
	}
	return l, nil
}

func (l *logic) Sync(ctx *context.Context, files MapFiles, batchSize int) error {
	hashes := files.GetKeys()

	err := l.conn.DB.
		WithContext(ctx.Request().Context()).
		Transaction(func(tx *gorm.DB) error {
			// Create a buffered (blocking) pool that can scale up to 10 workers
			pool := pond.New(10, 0)
			found := 0

			for i := 0; i < len(hashes); i += batchSize {
				start := i
				end := i + batchSize
				if end > len(hashes) {
					end = len(hashes)
				}

				filters := url.Values{
					"hash": hashes[start:end],
				}

				pool.Submit(func() {
					logrus.Infof("Running task from %d to %d", start, end)
					if stores, err := l.repo.Read(ctx, tx, filters); err != nil {
						return
					} else {
						found += len(stores)
						// set stored files in input files
						for _, f := range stores {
							files.Get(f.Hash).Stored = f
						}
					}
				})
			}

			// Stop the pool and wait for all submitted tasks to complete
			pool.StopAndWait()

			// Write into the database
			unsaved := make(Files, 0, len(files)-found)
			for _, f := range files {
				if f.Stored == nil {
					unsaved = append(unsaved, f)
				}
			}
			if err := l.repo.Create(ctx, tx, unsaved); err != nil {
				return err
			}

			return nil
		})

	return err
}

func (l *logic) ReadFiles(ctx *context.Context, root string, maxDepth int, reContentType *regexp.Regexp, filters url.Values) (files MapFiles, err error) {
	if logrus.Infof("walk directory %s", root); false {
		return nil, nil
		// Filter images
	} else if filteredFiles, _ := WalkDir(root, maxDepth, reContentType); len(filteredFiles) == 0 {
		logrus.Info("!!! no image found !!!")
		return nil, nil
	} else if err := l.Sync(ctx, filteredFiles, 10); err != nil {
		return nil, err
	} else {
		return filteredFiles, nil
	}
}
