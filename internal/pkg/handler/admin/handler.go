package admin

import (
	"net/http"

	"github.com/b1994mi/test-deptech/internal/pkg/usecase/admin"
	"github.com/go-playground/validator/v10"
	"github.com/uptrace/bunrouter"
)

type Handler interface {
	Create(w http.ResponseWriter, bunReq bunrouter.Request) error
	// Read(w http.ResponseWriter, bunReq bunrouter.Request) error
	// ReadAll(w http.ResponseWriter, bunReq bunrouter.Request) error
	// Update(w http.ResponseWriter, bunReq bunrouter.Request) error
	// Delete(w http.ResponseWriter, bunReq bunrouter.Request) error
}

type handler struct {
	validate *validator.Validate
	uc       admin.Usecase
}

func NewHandler(validate *validator.Validate, uc admin.Usecase) *handler {
	return &handler{validate, uc}
}
