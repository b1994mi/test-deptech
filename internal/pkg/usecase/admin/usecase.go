package admin

import (
	"github.com/b1994mi/test-deptech/internal/pkg/domain/sqlrepo"
)

type Usecase interface {
	Create(req Request) (interface{}, error)
	// Read() (interface{}, error)
	// ReadAll() (interface{}, error)
	// Update() (interface{}, error)
	// Delete() (interface{}, error)
}

type usecase struct {
	adminRepo sqlrepo.AdminRepo
}

func NewUsecase(
	adminRepo sqlrepo.AdminRepo,
) *usecase {
	return &usecase{
		adminRepo,
	}
}
