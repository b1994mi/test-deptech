package admin

import (
	"fmt"
	"log"
	"time"

	"github.com/b1994mi/test-deptech/internal/pkg/domain/sqlmodel"
)

type Request struct {
	Firstname   string    `validate:"required" json:"firstname"`
	Lastname    string    `validate:"required" json:"lastname"`
	Email       string    `validate:"required,email" json:"email"`
	DateOfBirth time.Time `validate:"required" json:"date_of_birth"`
	Gender      string    `validate:"oneof=male female prefer_not_to" json:"gender"`
	Password    string    `validate:"required" json:"password"`
}

func (uc *usecase) Create(req Request) (interface{}, error) {
	tx := uc.adminRepo.StartTx()
	defer tx.Rollback()

	_, err := uc.adminRepo.Create(&sqlmodel.Admin{
		ID:          0,
		Firstname:   "",
		Lastname:    "",
		Email:       "",
		DateOfBirth: time.Time{},
		Gender:      "",
		Password:    "",
	}, tx)
	if err != nil {
		log.Printf("unable to create admin: %v", err)
		return nil, fmt.Errorf("failed to create admin")
	}

	tx.Commit()

	return map[string]interface{}{
		"acknowledge": true,
	}, nil
}
