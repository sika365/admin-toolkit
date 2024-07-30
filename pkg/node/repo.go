package node

import (
	"net/url"

	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/database"
	"gorm.io/gorm"

	"github.com/sika365/admin-tools/context"
	"github.com/sika365/admin-tools/utils"
)

type Repo interface {
	Create(ctx *context.Context, db *gorm.DB, nodes LocalNodes) error
	Read(ctx *context.Context, db *gorm.DB, filters url.Values) (MapNodes, error)
	Update(ctx *context.Context, db *gorm.DB, node *LocalNode, filters url.Values) error
	Delete(ctx *context.Context, db *gorm.DB, id database.PID, filters url.Values) error
}

type repo struct {
}

func newRepo() (Repo, error) {
	r := &repo{}
	return r, nil
}

// Create implements Repo.
func (i *repo) Create(ctx *context.Context, db *gorm.DB, nodes LocalNodes) error {
	if len(nodes) == 0 {
		return nil
	} else if err := db.CreateInBatches(nodes, 100).Error; err != nil {
		return err
	} else {
		return nil
	}
}

// Read fetch nodes with filters
func (i *repo) Read(ctx *context.Context, db *gorm.DB, filters url.Values) (nodes MapNodes, err error) {
	var stored LocalNodes
	if err = utils.
		BuildGormQuery(ctx, db, filters).
		Find(&stored).Error; err != nil {
		return nil, err
	} else {
		return NewMapNodes(stored...), nil
	}
}

// Update implements Repo.
func (i *repo) Update(ctx *context.Context, db *gorm.DB, node *LocalNode, filters url.Values) error {
	panic("unimplemented")
}

// Delete implements Repo.
func (i *repo) Delete(ctx *context.Context, db *gorm.DB, id database.PID, filters url.Values) error {
	panic("unimplemented")
}
