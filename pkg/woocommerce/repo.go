package woocommerce

import (
	"net/url"

	"github.com/sika365/admin-tools/context"
	"github.com/sika365/admin-tools/pkg/client"
	"github.com/sika365/admin-tools/utils"
	"gorm.io/gorm"
)

type Repo interface {
	ReadTermTaxonomies(ctx *context.Context, db *gorm.DB, filters url.Values) (termTaxonomies []*WpTermTaxonomy, err error)
	ReadPosts(ctx *context.Context, db *gorm.DB, filters url.Values) (posts []*WpPost, err error)
}

type repo struct {
	client *client.Client
}

func newRepo(client *client.Client) (Repo, error) {
	r := &repo{
		client: client,
	}
	return r, nil
}

func (i *repo) ReadTermTaxonomies(ctx *context.Context, db *gorm.DB, filters url.Values) (termTaxonomies []*WpTermTaxonomy, err error) {
	var stored []*WpTermTaxonomy
	if err = utils.
		BuildGormQuery(ctx, db, filters).
		Where("taxonomy in (?)", []string{"product_cat"}).
		Find(&stored).Error; err != nil {
		return nil, err
	} else {
		return stored, nil
	}
}

func (i *repo) ReadPosts(ctx *context.Context, db *gorm.DB, filters url.Values) (posts []*WpPost, err error) {
	var stored []*WpPost
	if err = utils.
		BuildGormQuery(ctx, db, filters).
		Find(&stored).Error; err != nil {
		return nil, err
	} else {
		return stored, nil
	}
}
