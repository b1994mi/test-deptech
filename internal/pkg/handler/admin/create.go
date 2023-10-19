package admin

import (
	"net/http"

	"github.com/b1994mi/test-deptech/internal/pkg/domain/helper"
	"github.com/b1994mi/test-deptech/internal/pkg/usecase/admin"
	"github.com/uptrace/bunrouter"
)

func (h *handler) Create(w http.ResponseWriter, bunReq bunrouter.Request) error {
	var req admin.Request
	err := helper.ShouldBindJSON(&req, bunReq)
	if err != nil {
		helper.NewErrRes(w, http.StatusBadRequest, err)
		return nil
	}

	err = h.validate.Struct(req)
	if err != nil {
		helper.NewErrRes(w, http.StatusUnprocessableEntity, err)
	}

	res, err := h.uc.Create(req)
	if err != nil {
		switch e := err.(type) {
		case helper.StatusError:
			helper.NewErrRes(w, e.HTTPCode, e.InternalCode, e.Err)
		default:
			helper.NewErrRes(w, http.StatusUnprocessableEntity, err)
		}

		return nil
	}

	bunrouter.JSON(w, bunrouter.H{
		"data": res,
	})
	return nil
}
