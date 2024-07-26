package file

import (
	"net/url"
	"regexp"

	simutils "github.com/alifakhimi/simple-utils-go"
	"github.com/alitto/pond"
	"github.com/sika365/admin-tools/context"
	"github.com/sirupsen/logrus"
)

type Logic interface {
	Find(ctx *context.Context, filters url.Values) (files MapFiles, err error)
	Create(ctx *context.Context, files MapFiles, batchSize int) error
	Sync(ctx *context.Context, files MapFiles, loadBatchSize int, writeBatchSize int) (err error)
	Load(ctx *context.Context, files MapFiles, batchSize int) (unsaves Files, err error)
	ReadFiles(ctx *context.Context, root string, maxDepth int, reContentType *regexp.Regexp, filters url.Values) (files MapFiles, err error)
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

func (l *logic) Find(ctx *context.Context, filters url.Values) (files MapFiles, err error) {
	q := l.conn.DB.WithContext(ctx.Request().Context())
	if stores, err := l.repo.Read(ctx, q, filters); err != nil {
		return nil, err
	} else {
		return stores, nil
	}
}

func (l *logic) Create(ctx *context.Context, files MapFiles, batchSize int) error {
	pool := pond.New(batchSize, 0)

	for _, file := range files {
		pool.Submit(func() {
			logrus.Infof("Running task for %s", file)
			// Upload files
			// ...
			// Write uploaded files into the database
			tx := l.conn.DB.WithContext(ctx.Request().Context())
			if err := l.repo.Create(ctx, tx, file); err != nil {
				logrus.Infof("writing file %s in db failed", file)
			}
		})
	}

	pool.StopAndWait()

	return nil
}

func (l *logic) Sync(ctx *context.Context, files MapFiles, loadBatchSize int, writeBatchSize int) (err error) {
	if _, err := l.Load(ctx, files, loadBatchSize); err != nil {
		return err
		// Write into the database
	} else if err := l.Create(ctx, files, writeBatchSize); err != nil {
		return err
	}

	return err
}

func (l *logic) Load(ctx *context.Context, files MapFiles, batchSize int) (unsaves Files, err error) {
	hashes := files.GetKeys()

	tx := l.conn.DB

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

	unsaves = make(Files, 0, len(files)-found)
	for _, f := range files {
		if f.Stored == nil {
			unsaves = append(unsaves, f)
		}
	}

	return unsaves, nil
}

func (l *logic) ReadFiles(ctx *context.Context, root string, maxDepth int, reContentType *regexp.Regexp, filters url.Values) (files MapFiles, err error) {
	if logrus.Infof("walk directory %s", root); false {
		return nil, nil
		// Filter files
	} else if filteredFiles, _ := WalkDir(root, maxDepth, reContentType, nil); len(filteredFiles) == 0 {
		logrus.Info("!!! no files found !!!")
		return nil, nil
	} else {
		return filteredFiles, nil
	}
}
