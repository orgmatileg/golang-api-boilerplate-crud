package usecases

import (
	"bytes"
	"errors"
	"fmt"
	"golang-api-boilerplate-crud/models"
	"golang-api-boilerplate-crud/repository"
	"io"
	"mime/multipart"
)

// FruitsUsecasesInterface ...
type FruitsUsecasesInterface interface {
	Create(m *models.Fruits) error
	GetOneByID(id int, m *models.Fruits) error
	GetAll(limit, offset, order, search string) ([]models.Fruits, int64, error)
	UpdateOneByID(id int, m *models.Fruits) (rowAffected int, err error)
	UpdatePhotoOneByID(id int, f *multipart.File) (rowAffected int, err error)
	DeleteOneByID(id int) (rowAffected int, err error)
}

// Fruits Usecases
type Fruits struct{}

// NewFruitsUsecase ...
func NewFruitsUsecase() FruitsUsecasesInterface {
	return Fruits{}
}

// Create ...
func (v Fruits) Create(m *models.Fruits) error {
	if m.Name == "" {
		return errors.New("usecases: name cannot be empty string")
	}
	if m.Color == "" {
		return errors.New("usecases: color cannot be empty string")
	}
	return repository.NewFruitsRepository().Create(m)
}

// GetOneByID ...
func (v Fruits) GetOneByID(id int, m *models.Fruits) error {
	if id <= 0 {
		return fmt.Errorf("usecases: id cannot be 0 or below")
	}
	return repository.NewFruitsRepository().GetOneByID(id, m)
}

// GetAll ...
func (v Fruits) GetAll(limit, offset, order, search string) ([]models.Fruits, int64, error) {
	validOrder := func() bool {
		for _, v := range []string{"ASC", "DESC"} {
			if order == v {
				return true
			}
		}
		return false
	}()
	if !validOrder {
		return nil, -1, fmt.Errorf("usecases: order param must be string ASC or DESC")
	}

	ml, err := repository.NewFruitsRepository().GetAll(limit, offset, order, search)
	if err != nil {
		return nil, -1, fmt.Errorf("usecases: GetAll : %s", err.Error())
	}
	tc, err := repository.NewFruitsRepository().Count()
	if err != nil {
		return nil, -1, fmt.Errorf("usecases: Count : %s", err.Error())
	}

	return ml, tc, nil
}

// UpdateOneByID ...
func (v Fruits) UpdateOneByID(id int, m *models.Fruits) (int, error) {

	if err := repository.NewFruitsRepository().IsExistsByID(id); err != nil {
		return -1, fmt.Errorf("usecases: %s", err.Error())
	}

	if m.Name == "" {
		return -1, fmt.Errorf("usecases: name cannot be empty string")
	}
	if m.Color == "" {
		return -1, fmt.Errorf("usecases: name cannot be empty string")
	}
	return repository.NewFruitsRepository().UpdateOneByID(id, m)
}

// UpdatePhotoOneByID ...
func (v Fruits) UpdatePhotoOneByID(id int, f *multipart.File) (int, error) {

	if err := repository.NewFruitsRepository().IsExistsByID(id); err != nil {
		return -1, fmt.Errorf("usecases: %s", err.Error())
	}

	var buf bytes.Buffer
	_, err := io.Copy(&buf, *f)
	if err != nil {
		return -1, fmt.Errorf("usecases: could not copy multipart.file to buffer: %s", err.Error())
	}

	bs := buf.String()
	imgURL, err := repository.NewImgurRepository().Create(bs)
	if err != nil {
		return -1, fmt.Errorf("usecases: could not upload image to imgurl: %s", err.Error())
	}

	return repository.NewFruitsRepository().UpdatePhotoOneByID(id, imgURL)
}

// DeleteOneByID ...
func (v Fruits) DeleteOneByID(id int) (int, error) {
	if id <= 0 {
		return -1, fmt.Errorf("usecases: id cannot be 0 or below")
	}

	if err := repository.NewFruitsRepository().IsExistsByID(id); err != nil {
		return -1, fmt.Errorf("usecases: %s", err.Error())
	}

	return repository.NewFruitsRepository().DeleteOneByID(id)
}
