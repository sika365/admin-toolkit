package image

import (
	"net/http"
	"regexp"

	simutils "github.com/alifakhimi/simple-utils-go"
	"github.com/labstack/echo/v4"

	"github.com/sika365/admin-tools/context"
	"github.com/sika365/admin-tools/pkg/file"
	"github.com/sika365/admin-tools/utils"
)

type Rest interface {
	Scan(ctx echo.Context) error
	Sync(ectx echo.Context) error
}

type rest struct {
	logic Logic
}

func newRest(h *simutils.HttpServer, logic Logic) (Rest, error) {
	r := &rest{
		logic: logic,
	}

	sg := h.
		// use prefix group
		PrefixGroup().
		// [prefixgroup path]/images
		Group("/images")
	{
		sg.POST("", r.Scan)
		sg.POST("/sync", r.Sync)
	}

	return r, nil
}

func (r *rest) Scan(ectx echo.Context) error {
	var request ScanRequest
	if ctx, ok := ectx.(*context.Context); !ok {
		return nil
	} else if err := ctx.Bind(&request); err != nil {
		return err
	} else if filters := ctx.QueryParams(); false {
		return nil
	} else if request.MaxDepth = utils.Max(request.MaxDepth, -1); request.MaxDepth > 10 {
		return file.ErrInvalidMaxDepthValue
	} else if request.Root = utils.DefaultIfZero(request.Root, DefaultRoot); request.Root == "" {
		return file.ErrEmptyRootNotValid
	} else if request.ContentTypes = utils.DefaultIfZero(request.ContentTypes, ImageContentTypeRegex); request.ContentTypes == "" {
		return file.ErrInvalidContentTypesValue
	} else if reContentType, err := regexp.Compile(request.ContentTypes); err != nil {
		return err
	} else if images, err := r.logic.ReadFiles(ctx, request.Root, request.MaxDepth, reContentType, nil, filters); err != nil {
		return err
	} else {
		return simutils.ReplyTemplate(ctx, http.StatusOK, nil,
			&ScanResponse{Data: images},
			simutils.CreatePaginateTemplate(len(images), 0, len(images)),
		)
	}
}

func (r *rest) Sync(ectx echo.Context) error {
	var request SyncRequest
	if ctx, ok := ectx.(*context.Context); !ok {
		return nil
	} else if err := ctx.Bind(&request); err != nil {
		return err
	} else if filters := ctx.QueryParams(); false {
		return nil
	} else if request.MaxDepth = utils.Max(request.MaxDepth, -1); request.MaxDepth > 10 {
		return file.ErrInvalidMaxDepthValue
	} else if request.Root = utils.DefaultIfZero(request.Root, DefaultRoot); request.Root == "" {
		return file.ErrEmptyRootNotValid
	} else if request.ContentTypes = utils.DefaultIfZero(request.ContentTypes, ImageContentTypeRegex); request.ContentTypes == "" {
		return file.ErrInvalidContentTypesValue
	} else if images, err := r.logic.Sync(ctx, request.Root, request.MaxDepth, request.Replace, filters); err != nil {
		return err
	} else {
		return simutils.ReplyTemplate(ctx, http.StatusOK, nil,
			&SyncResponse{Data: images},
			simutils.CreatePaginateTemplate(len(images), 0, len(images)),
		)
	}
}
