package sqlrepo

import (
	"github.com/b1994mi/test-deptech/internal/pkg/domain/helper"
	"github.com/b1994mi/test-deptech/internal/pkg/domain/sqlmodel"

	"gorm.io/gorm"
)

type AdminRepo interface {
	StartTx() *gorm.DB
	Create(m *sqlmodel.Admin, tx *gorm.DB) (*sqlmodel.Admin, error)
	Update(m *sqlmodel.Admin, tx *gorm.DB) error
	Delete(m *sqlmodel.Admin, tx *gorm.DB) error
	FindOneBy(criteria map[string]interface{}) (*sqlmodel.Admin, error)
	FindBy(criteria map[string]interface{}, page, size int) ([]*sqlmodel.Admin, error)
	Count(criteria map[string]interface{}) int64
}

type adminRepo struct {
	db *gorm.DB
}

func NewAdminRepo(db *gorm.DB) *adminRepo {
	return &adminRepo{
		db,
	}
}

func (rpo *adminRepo) StartTx() *gorm.DB {
	return rpo.db.Begin()
}

func (rpo *adminRepo) Create(m *sqlmodel.Admin, tx *gorm.DB) (*sqlmodel.Admin, error) {
	err := tx.Create(&m).Error
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (rpo *adminRepo) Update(m *sqlmodel.Admin, tx *gorm.DB) error {
	return tx.Save(&m).Error
}

func (rpo *adminRepo) Delete(m *sqlmodel.Admin, tx *gorm.DB) error {
	return tx.Delete(&m).Error
}

func (rpo *adminRepo) FindOneBy(criteria map[string]interface{}) (*sqlmodel.Admin, error) {
	var m sqlmodel.Admin

	err := rpo.db.Where(criteria).Take(&m).Error
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (rpo *adminRepo) FindBy(criteria map[string]interface{}, page, size int) ([]*sqlmodel.Admin, error) {
	var data []*sqlmodel.Admin
	if page == 0 || size == 0 {
		page, size = -1, -1
	}

	limit, offset := helper.GetLimitOffset(page, size)
	err := rpo.db.
		Where(criteria).
		Offset(offset).Limit(limit).
		Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (rpo *adminRepo) Count(criteria map[string]interface{}) int64 {
	var result int64

	if res := rpo.db.Model(sqlmodel.Admin{}).Where(criteria).Count(&result); res.Error != nil {
		return 0
	}

	return result
}
