package category

import (
	"errors"
	"fmt"
	"net/url"

	simutils "github.com/alifakhimi/simple-utils-go"
	"github.com/alifakhimi/simple-utils-go/simscheme"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/models"

	"github.com/sika365/admin-tools/context"
	"github.com/sika365/admin-tools/pkg/excel"
	"github.com/sika365/admin-tools/service/client"
)

type Logic interface {
	Find(ctx *context.Context) error
	Store(ctx *context.Context, req *SyncRequest, doc *simscheme.Document) error
	Sync(ctx *context.Context, req *SyncRequest, filters url.Values) (*simscheme.Document, error)
}

type logic struct {
	conn   *simutils.DBConnection
	client *client.Client
	repo   Repo
}

func newLogic(conn *simutils.DBConnection, client *client.Client, repo Repo) (Logic, error) {
	l := &logic{
		conn:   conn,
		client: client,
		repo:   repo,
	}
	return l, nil
}

func (l *logic) Find(ctx *context.Context) error {
	return nil
}

func (l *logic) Store(ctx *context.Context, req *SyncRequest, doc *simscheme.Document) error {
	if req.ReplaceNodes {
		if err := l.repo.Clear(ctx, l.conn.DB.WithContext(ctx.Request().Context())); err != nil {
			return err
		}
	}

	if uncategorizedNode, err := l.client.GetNodeByAlias(ctx, "Uncategorized"); err != nil {
		return err
	} else {
		for _, node := range doc.Nodes {
			var (
				rec = node.Data.(*CategoryRecord)
			)

			rcategory, err := l.client.GetCategoryByAlias(ctx, rec.Slug.ToString())
			if errors.Is(err, models.ErrNotFound) {
				if rcategory, err = l.client.StoreCategory(ctx, &models.Category{
					Title: rec.Title,
					Slug:  rec.Slug.ToString(),
					Alias: rec.Slug.ToString(),
				}, uncategorizedNode); err != nil {
					return err
				}
			} else if err != nil {
				return err
			}

			rec.LocalCategory = &LocalCategory{
				Title:    rec.Title,
				Alias:    rec.Slug,
				Slug:     rec.Slug,
				Category: rcategory,
			}

			if len(rcategory.Nodes) > 0 {
				for _, n := range rcategory.Nodes {
					if len(n.Nodes) > 0 {
						n.Nodes, _ = n.Nodes.ToNested()
					}
				}
			}

			if err := l.repo.Create(
				ctx,
				l.conn.DB.WithContext(ctx.Request().Context()),
				CategoryRecords{rec},
			); err != nil {
				return err
			}
		}
	}

	return nil
}

func (l *logic) Sync(ctx *context.Context, req *SyncRequest, filters url.Values) (*simscheme.Document, error) {
	var (
		catRecordDoc = simscheme.
			GetSchema().
			AddNewDocumentWithType(&CategoryRecord{})
	)

	if !simutils.IsSlug(req.ScanRequest.CategoryHeaderMap.Slug) {
		return nil, fmt.Errorf("key of title is not specified")
	} else if csvFiles, err := excel.LoadExcels(ctx, req.Root, req.MaxDepth); err != nil {
		return nil, err
		// Make CategoryNodes from the files
	} else if err := excel.FromFiles(
		csvFiles,
		req.Offset,
		func(header map[string]int, rec []string) {
			catRec := &CategoryRecord{
				Title: rec[header[req.CategoryHeaderMap.Title]],
				Slug:  simutils.MakeSlug(rec[header[req.CategoryHeaderMap.Slug]]),
			}

			if !catRecordDoc.LabelExists(*catRecordDoc.BuildNodeLabel(catRec)) {
				catRecordDoc.AddNode(catRec)
			}
		},
	); err != nil {
		return catRecordDoc, err
	} else if err := l.Store(ctx, req, catRecordDoc); err != nil {
		return nil, err
	} else {
		return catRecordDoc, nil
	}
}
