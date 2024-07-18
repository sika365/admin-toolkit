package file

import (
	"errors"
	"net/http"
	"regexp"

	simutils "github.com/alifakhimi/simple-utils-go"
	"github.com/labstack/echo/v4"
	"github.com/sika365/admin-tools/context"
)

var (
	ErrInvalidMaxDepthValue     = errors.New("invalid max depth value")
	ErrInvalidContentTypesValue = errors.New("invalid content types value")
	ErrEmptyRootNotValid        = errors.New("empty root is not valid")
)

type Rest interface {
	ReadFiles(ctx echo.Context) error
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
		// [prefixgroup path]/files
		Group("/files")
	{
		sg.POST("", r.ReadFiles)
	}

	return r, nil
}

func (r *rest) ReadFiles(ectx echo.Context) error {
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
		// } else if maxDepth := cast.ToInt(utils.PopQueryParam[string](filters, "max_depth")); maxDepth < 0 {
		// 	return ErrInvalidMaxDepthValue
		// } else if reContentType, err := regexp.Compile(utils.PopQueryParam[string](filters, "content_types")); err != nil {
		// 	return err
		// } else if root := utils.PopQueryParam[string](filters, "root"); root == "" {
		// 	return ErrEmptyRootNotValid
	} else if request.MaxDepth == nil || *request.MaxDepth < 0 {
		return ErrInvalidMaxDepthValue
	} else if request.Root == nil || *request.Root == "" {
		return ErrEmptyRootNotValid
	} else if request.ContentTypes == nil || *request.ContentTypes == "" {
		return ErrInvalidContentTypesValue
	} else if reContentType, err := regexp.Compile(*request.ContentTypes); err != nil {
		return err
	} else if files, err := r.logic.ReadFiles(ctx, *request.Root, *request.MaxDepth, reContentType, filters); err != nil {
		return err
	} else {
		return simutils.Reply(ctx, http.StatusOK, nil,
			map[string]interface{}{
				"files": files,
			},
			simutils.CreatePaginateTemplate(len(files), 0, len(files)),
		)
	}
}
