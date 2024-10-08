package category

import (
	"net/http"

	simutils "github.com/alifakhimi/simple-utils-go"
	"github.com/labstack/echo/v4"

	"github.com/sika365/admin-tools/context"
)

// Errors
var ()

type Rest interface {
	Find(ctx echo.Context) error
	Scan(ctx echo.Context) error
	Sync(ctx echo.Context) error
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
		// [prefixgroup path]/products
		Group("/categories")
	{
		sg.GET("", r.Find)
		sg.POST("/scan", r.Scan)
		sg.POST("/sync", r.Sync)
	}

	return r, nil
}

func (r *rest) Find(ectx echo.Context) error {
	return nil
}

func (r *rest) Scan(ctx echo.Context) error {
	return nil
}

func (r *rest) Sync(ectx echo.Context) error {
	var req SyncRequest
	if ctx, ok := ectx.(*context.Context); !ok {
		return nil
	} else if err := ctx.Bind(&req); err != nil {
		return err
	} else if filters := ctx.QueryParams(); false {
		return nil
	} else if catDoc, err := r.logic.Sync(ctx, &req, filters); err != nil {
		return err
	} else {
		return simutils.ReplyTemplate(ctx, http.StatusOK, nil,
			&SyncResponse{Data: catDoc},
			simutils.CreatePaginateTemplate(len(catDoc.Nodes), 0, len(catDoc.Nodes)),
		)
	}
}
