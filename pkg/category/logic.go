package category

import (
	"errors"
	"fmt"
	"net/url"

	simutils "github.com/alifakhimi/simple-utils-go"
	"github.com/alifakhimi/simple-utils-go/simscheme"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/models"

	"github.com/sika365/admin-tools/context"
	"github.com/sika365/admin-tools/pkg/client"
	"github.com/sika365/admin-tools/pkg/excel"
)

type Logic interface {
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
	if uncategorizedNode, err := l.client.GetNodeByAlias(ctx, "Uncategorized"); err != nil {
		return err
	} else {
		if err := l.repo.Clear(ctx, l.conn.DB.WithContext(ctx.Request().Context())); err != nil {
			return err
		}

		for _, node := range doc.Nodes {
			var (
				rec = node.Data.(*CategoryRecord)
			)

			category, err := l.client.GetCategoryByAlias(ctx, rec.Title)
			if errors.Is(err, models.ErrNotFound) {
				if category, err = l.client.StoreCategory(ctx, &models.Category{
					Title: rec.Title,
					Alias: rec.Title,
				}, uncategorizedNode); err != nil {
					return err
				}
			} else if err != nil {
				return err
			}

			rec.LocalCategory = &LocalCategory{
				Title:    rec.Title,
				Alias:    rec.Title,
				Slug:     rec.Title,
				Category: category,
			}

			if err := l.repo.Create(ctx, l.conn.DB.WithContext(ctx.Request().Context()), rec); err != nil {
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

	if req.ScanRequest.CategoryHeaderMap.Title == "" {
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
			}
			catRecordDoc.AddNode(catRec)
		},
	); err != nil {
		return catRecordDoc, err
	} else if err := l.Store(ctx, req, catRecordDoc); err != nil {
		return nil, err
	} else {
		return catRecordDoc, nil
	}
}
