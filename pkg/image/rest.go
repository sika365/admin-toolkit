package image

import (
	"net/http"
	"regexp"

	simutils "github.com/alifakhimi/simple-utils-go"
	"github.com/labstack/echo/v4"
	"github.com/sika365/admin-tools/context"
	"github.com/sika365/admin-tools/pkg/file"
)

// Errors
var ()

type Rest interface {
	Get(ctx echo.Context) error
}

type rest struct {
	logic Logic
}

func newRest(logic Logic, h *simutils.HttpServer) (Rest, error) {
	r := &rest{
		logic: logic,
	}

	sg := h.
		// use prefix group
		PrefixGroup().
		// [prefixgroup path]/images
		Group("/images")
	{
		sg.POST("/sync", r.Sync)
		sg.POST("", r.Get)
	}

	return r, nil
}

func (r *rest) Sync(ectx echo.Context) error {
	var request = struct {
		Root     *string `json:"root,omitempty"`
		MaxDepth *int    `json:"max_depth,omitempty"`
	}{}
	if ctx, ok := ectx.(*context.Context); !ok {
		return nil
	} else if err := ctx.Bind(&request); err != nil {
		return err
	} else if filters := ctx.QueryParams(); false {
		return nil
	} else if request.MaxDepth == nil || *request.MaxDepth < 0 {
		return file.ErrInvalidMaxDepthValue
	} else if request.Root == nil || *request.Root == "" {
		return file.ErrEmptyRootNotValid
	} else if images, err := r.logic.Sync(ctx, *request.Root, *request.MaxDepth, filters); err != nil {
		return err
	} else {
		return simutils.Reply(ctx, http.StatusOK, nil,
			map[string]interface{}{
				"images": images,
			},
			simutils.CreatePaginateTemplate(len(images), 0, len(images)),
		)
	}
}

func (r *rest) Get(ectx echo.Context) error {
	var request = struct {
		Root         *string `json:"root,omitempty"`
		MaxDepth     *int    `json:"max_depth,omitempty"`
		ContentTypes *string `json:"content_types,omitempty"`
	}{}
	if ctx, ok := ectx.(*context.Context); !ok {
		return nil
	} else if err := ctx.Bind(&request); err != nil {
		return err
	} else if filters := ctx.QueryParams(); false {
		return nil
	} else if request.MaxDepth == nil || *request.MaxDepth < 0 {
		return file.ErrInvalidMaxDepthValue
	} else if request.Root == nil || *request.Root == "" {
		return file.ErrEmptyRootNotValid
	} else if request.ContentTypes == nil || *request.ContentTypes == "" {
		return file.ErrInvalidContentTypesValue
	} else if reContentType, err := regexp.Compile(*request.ContentTypes); err != nil {
		return err
	} else if images, err := r.logic.ReadFiles(ctx, *request.Root, *request.MaxDepth, reContentType, nil, filters); err != nil {
		return err
	} else {
		return simutils.Reply(ctx, http.StatusOK, nil,
			map[string]interface{}{
				"images": images,
			},
			simutils.CreatePaginateTemplate(len(images), 0, len(images)),
		)
	}
}
