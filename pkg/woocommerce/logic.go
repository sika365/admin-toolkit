package woocommerce

import (
	"errors"
	"net/url"

	simutils "github.com/alifakhimi/simple-utils-go"
	"github.com/sirupsen/logrus"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/models"

	"github.com/sika365/admin-tools/context"
	"github.com/sika365/admin-tools/pkg/category"
	"github.com/sika365/admin-tools/pkg/client"
	"github.com/sika365/admin-tools/pkg/node"
	"github.com/sika365/admin-tools/pkg/product"
)

type Logic interface {
	Sync(ctx *context.Context) (prdRecs product.ProductRecords, err error)
}

type logic struct {
	conn      *simutils.DBConnection
	client    *client.Client
	repo      Repo
	catLogic  category.Logic
	catRepo   category.Repo
	prodLogic product.Logic
	prodRepo  product.Repo
}

func newLogic(conn *simutils.DBConnection, client *client.Client, repo Repo, catLogic category.Logic, catRepo category.Repo, prodLogic product.Logic, prodRepo product.Repo) (Logic, error) {
	l := &logic{
		conn:      conn,
		client:    client,
		repo:      repo,
		catLogic:  catLogic,
		catRepo:   catRepo,
		prodLogic: prodLogic,
		prodRepo:  prodRepo,
	}
	return l, nil
}

func (l *logic) Sync(ctx *context.Context) (prdRecs product.ProductRecords, err error) {
	if req, err := context.GetRequestModel[*SyncRequest](ctx); err != nil {
		return nil, err
	} else if dbconn := (&simutils.DBConnection{DBConfig: req.DatabaseConfig}); false {
		return nil, nil
	} else if err := simutils.Connect(dbconn); err != nil {
		return nil, err
	} else if err := l.SyncCategory(ctx, req, dbconn); err != nil {
		return nil, err
	} else if prdRecs, err := l.SyncProduct(ctx, req, dbconn); err != nil {
		return nil, err
	} else {
		return prdRecs, nil
	}
}

func (l *logic) SyncProduct(ctx *context.Context, req *SyncRequest, dbconn *simutils.DBConnection) (prdRecs product.ProductRecords, err error) {
	// Read all products from woocommerce db
	if posts, err := l.repo.ReadPosts(
		ctx,
		dbconn.DB.WithContext(ctx.Request().Context()),
		url.Values{
			"post_type":   []string{"product"},
			"post_status": []string{"publish"},
			"includes": []string{
				"Meta", "Parent", "TermTaxonomies.Term", "Posts",
				"Posts.Meta", "Posts.TermTaxonomies.Term", "Posts.Posts",
				"Posts.Posts.Meta", "Posts.Posts.TermTaxonomies.Term", "Posts.Posts.Posts",
			},
		},
	); err != nil {
		return nil, err
	} else {
		postDoc.AddNodes(posts)
		prdRecs = make(product.ProductRecords, 0, len(posts))

		for _, post := range posts {
			var (
				productGroup    *models.ProductGroup
				prd             = &models.Product{}
				topNodes        = models.Nodes{}
				categoryRecords = post.GetCategoryRecords()
				isProductGroup  = post.IsProductGroup()
			)

			if topNodes, err = l.GetTopNodes(ctx, categoryRecords); err != nil {
				return nil, err
			}

			if isProductGroup {
				// Product with variations
				reqLProductGroup := post.ToLocalProductGroup(topNodes)

				if lproductGroup, err := l.prodLogic.FindOrCreateProductGroup(
					ctx,
					reqLProductGroup,
				); err != nil {
					return nil, err
				} else if lproductGroup == nil {
					return nil, nil
				} else {
					for _, prd := range reqLProductGroup.ProductGroup.Products {
						prd.ProductGroupID, _ = lproductGroup.ID.ToNullPID()

						if storedPrdRecs, err := l.SaveProduct(ctx, post, prd, categoryRecords); err != nil && !errors.Is(err, models.ErrNotFound) { // check cover and gallery
							return nil, err
						} else if errors.Is(err, models.ErrNotFound) {
							continue
						} else {
							prdRecs = append(prdRecs, storedPrdRecs...)

							for _, storedPrdRec := range storedPrdRecs {
								lnodeRecDoc.AddNodes(storedPrdRec.LocalCategory.Nodes)
							}
						}
					}
				}
			} else {
				var (
					cover, gallery = post.GetAttachments()
				)

				prd = post.ToProduct(
					productGroup,
					cover,
					gallery,
					topNodes,
				)

				if storedPrdRecs, err := l.SaveProduct(ctx, post, prd, categoryRecords); err != nil && !errors.Is(err, models.ErrNotFound) { // check cover and gallery
					return nil, err
				} else if errors.Is(err, models.ErrNotFound) {
					continue
				} else {
					prdRecs = append(prdRecs, storedPrdRecs...)

					for _, storedPrdRec := range storedPrdRecs {
						lnodeRecDoc.AddNodes(storedPrdRec.LocalCategory.Nodes)
					}
				}
			}
		}

		return prdRecs, nil
	}
}

func (l *logic) GetTopNodes(ctx *context.Context, catRecs category.CategoryRecords) (models.Nodes, error) {
	var (
		topNodes models.Nodes
		slugs    []string
	)

	for _, catRec := range catRecs {
		if catRec.Slug.IsValid() {
			slugs = append(slugs, catRec.Slug.ToString())
		}
	}

	if lcats, err := l.catRepo.Read(ctx,
		l.conn.DB.WithContext(ctx.Request().Context()),
		url.Values{
			"slug": slugs,
		},
	); err != nil {
		return nil, err
	} else {
		for _, lcat := range lcats {
			topNodes = append(topNodes, lcat.Category.Nodes...)
		}

		return topNodes, nil
	}
}

func (l *logic) SaveProduct(ctx *context.Context, post *WpPost, prd *models.Product, categoryRecords category.CategoryRecords) (prdRecs product.ProductRecords, err error) {
	for _, catRec := range categoryRecords {
		postPrdRec := post.ToProductRecord(string(catRec.Slug), prd)
		if postPrdRec == nil {
			// return nil, models.ErrNotFound
			continue
		}

		if prdRec, err := l.prodLogic.Save(ctx, postPrdRec); err != nil {
			logrus.Errorf("register product failed %v", err)
			return nil, err
		} else {
			prdRecs = append(prdRecs, prdRec)
		}
	}

	return prdRecs, nil
}

func (l *logic) SyncCategory(ctx *context.Context, req *SyncRequest, dbconn *simutils.DBConnection) (err error) {
	// Read all categories from woocommerce db
	if tts, err := l.repo.ReadTermTaxonomies(
		ctx,
		dbconn.DB.WithContext(ctx.Request().Context()),
		url.Values{
			"taxonomy": []string{"product_cat"},
			"includes": []string{"Term", "ParentTermTaxonomy"},
		},
	); err != nil {
		return err
	} else {
		termTaxonomyDoc.AddNodes(tts)

		catRecs := make(category.CategoryRecords, 0, len(tts))
		for _, tt := range tts {
			catRec := tt.ToCategoryRecord()
			catRecs = append(catRecs, catRec)
			catRecDoc.AddNode(catRec)
			lnodeRecDoc.AddNodes(catRec.LocalCategory.Nodes)
		}

		nestedCatRecs := catRecs.ToNested()

		if err := l.StoreCategoryRecords(ctx, req, nestedCatRecs); err != nil {
			return err
		} else {
			catRecDoc.AddNodes(catRecs)
		}
	}

	return nil
}

func (l *logic) StoreCategoryRecords(ctx *context.Context, req *SyncRequest, srcNestedCatRec category.CategoryRecords) (err error) {
	// If ReplaceNodes is true, clear the existing categories in the database.
	if req.ReplaceNodes {
		if err := l.catRepo.Clear(ctx, l.conn.DB.WithContext(ctx.Request().Context())); err != nil {
			return err
		}
	}
	var (
		uncategorizedNode *models.Node
		// mainCategoriesNode *models.Node
	)

	// Retrieve the "uncategorized" node.
	if uncategorizedNode, err = l.client.GetNodeByAlias(ctx, "uncategorized"); err != nil {
		logrus.WithFields(logrus.Fields{
			"category_alias": "uncategorized",
		}).Errorln(err)
		return err
	}

	// Uncomment and use this if there is a "main_categories" node.
	// if mainCategoriesNode, err = l.client.GetNodeByAlias(ctx, "main_categories"); err != nil {
	// 	return err
	// }

	// Iterate over each category record in the source nested category records.
	for _, srcCatRec := range srcNestedCatRec {
		// Recursively call StoreCategoryRecord for each category record.
		if err := l.storeCategoryRecursive(ctx, srcCatRec, uncategorizedNode); err != nil {
			return err
		}
	}

	return nil
}

func (l *logic) storeCategoryRecursive(ctx *context.Context, srcCatRec *category.CategoryRecord, parentNode *models.Node) (err error) {
	// Store the current category
	if err := l.StoreCategoryRecord(ctx, srcCatRec, parentNode); err != nil {
		return err
	}

	srcCatRecNode := srcCatRec.LocalCategory.Nodes[0]
	// Recursively store each child category if it exists
	for _, childNode := range srcCatRecNode.SubNodes {
		// Find category record by slug
		catRecNode := catRecDoc.GetNode(&category.CategoryRecord{Slug: childNode.Slug})
		childCatRec := catRecNode.Data.(*category.CategoryRecord)
		// catNodes := childCatRec.LocalCategory.Category.Nodes
		// if len(catNodes) == 0 {
		// 	continue
		// }

		if err := l.storeCategoryRecursive(
			ctx,
			childCatRec,
			srcCatRec.LocalCategory.Category.Nodes[0],
		); err != nil {
			return err
		}
	}

	return nil
}

func (l *logic) StoreCategoryRecord(ctx *context.Context, srcCatRec *category.CategoryRecord, categoriesNode *models.Node) (err error) {
	if catRecs, err := l.catRepo.ReadCategoryRecords(
		ctx,
		l.conn.DB.WithContext(ctx.Request().Context()),
		url.Values{
			"slug": []string{srcCatRec.Slug.ToString()},
		},
	); err != nil {
		return err
	} else if len(catRecs) == 0 {
		// Check from api
		rcategory, err := l.client.GetCategoryByAlias(ctx, srcCatRec.Slug)
		// Store category if not found
		if errors.Is(err, models.ErrNotFound) {
			// Get parents
			parentLnodes := srcCatRec.LocalCategory.GetParentAliasByNodes()
			parentRnodes := models.Nodes{}

			if len(parentLnodes) == 0 {
				// Add to categoriesNode
				parentRnodes = append(parentRnodes, categoriesNode)
			} else {
				for _, lnodeParentAlias := range parentLnodes {
					lnParentNode := lnodeRecDoc.
						GetNode(&node.LocalNode{Alias: lnodeParentAlias})
					if lnParentNode == nil {
						// TODO: Fetch from db
						// Temporary dismissed
						continue
					}
					lnp := lnParentNode.Data.(*node.LocalNode)
					parentCatRecNode := catRecDoc.GetNode(&category.CategoryRecord{Slug: lnp.Slug})
					if parentCatRecNode == nil {
						// TODO: Fetch from db
						// Temporary dismissed
						continue
					}

					parentCatRec := parentCatRecNode.Data.(*category.CategoryRecord)
					parentRnodes = append(parentRnodes, parentCatRec.LocalCategory.Category.Nodes...)
				}
			}

			if rcategory, err = l.client.StoreCategory(
				ctx,
				srcCatRec.LocalCategory.Category,
				parentRnodes...,
			); err != nil {
				logrus.WithFields(logrus.Fields{
					"category_record": srcCatRec,
				}).Errorln(err)
				return err
			}
		} else if err != nil {
			return err
		}

		srcCatRec.LocalCategory.Category = rcategory

		if len(rcategory.Nodes) > 0 {
			for _, n := range rcategory.Nodes {
				if len(n.Nodes) > 0 {
					n.Nodes, _ = n.Nodes.ToNested()
				}
			}
		}

		if err := l.catRepo.Create(
			ctx,
			l.conn.DB.WithContext(ctx.Request().Context()),
			category.CategoryRecords{srcCatRec},
		); err != nil {
			return err
		}
	} else if len(catRecs) == 1 {
		*srcCatRec = *catRecs[0]
	}

	return nil
}
