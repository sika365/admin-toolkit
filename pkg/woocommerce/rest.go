package woocommerce

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
	SyncByWoocommerce(ctx echo.Context) error
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
		Group("/woocommerce")
	{
		sg.GET("", r.Find)
		sg.POST("/sync", r.SyncByWoocommerce)
	}

	return r, nil
}

func (r *rest) Find(ectx echo.Context) error {
	return nil
}

func (r *rest) SyncByWoocommerce(ctx echo.Context) error {
	if ctx, err := context.Binder(ctx, &SyncRequest{}); err != nil {
		return err
	} else if err := r.logic.Sync(ctx); err != nil {
		return err
	} else {
		return simutils.ReplyTemplate(ctx, http.StatusOK, nil,
			&SyncResponse{},
			nil,
		)
	}
}
