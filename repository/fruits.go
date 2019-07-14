package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"golang-api-boilerplate-crud/database"
	"golang-api-boilerplate-crud/models"
	"time"
)

// FruitsRepositoryInterface ...
type FruitsRepositoryInterface interface {
	Create(m *models.Fruits) error
	GetOneByID(id int, m *models.Fruits) error
	GetAll(limit, offset, order, search string) ([]models.Fruits, error)
	UpdateOneByID(id int, m *models.Fruits) (rowAffected int, err error)
	UpdatePhotoOneByID(id int, url string) (rowAffected int, err error)
	DeleteOneByID(id int) (rowAffected int, err error)
	IsExistsByID(id int) (err error)
	Count() (totalRow int64, err error)
}

// Fruits Repository
type Fruits struct{}

// NewFruitsRepository ...
func NewFruitsRepository() FruitsRepositoryInterface {
	return Fruits{}
}

// Create ...
func (v Fruits) Create(m *models.Fruits) error {

	query := `
		INSERT INTO fruits (fruit_name, fruit_color)
		VALUES (?,?)`
	statement, err := database.GetMariaDBconn().Prepare(query)

	if err != nil {
		return fmt.Errorf("repository: could not prepare query to database: %s", err.Error())
	}

	defer statement.Close()

	r, err := statement.Exec(m.Name, m.Color)
	if err != nil {
		return fmt.Errorf("repository: could not exec query to database: %s", err.Error())
	}

	idInt64, err := r.LastInsertId()
	if err != nil {
		return fmt.Errorf("repository: could not get last inserted row id: %s", err.Error())
	}
	m.ID = int(idInt64) // assign last inserted row id to p.ID

	return nil
}

// GetOneByID ...
func (v Fruits) GetOneByID(id int, m *models.Fruits) error {

	query := `
		SELECT * FROM fruits
		WHERE fruit_id = ?`
	statement, err := database.GetMariaDBconn().Prepare(query)

	if err != nil {
		return fmt.Errorf("repository: could not prepare query to database: %s", err.Error())
	}

	defer statement.Close()

	var updatedAt sql.NullString
	err = statement.QueryRow(id).Scan(
		&m.ID,
		&m.Name,
		&m.Color,
		&m.ImageURL,
		&m.CreatedAt,
		&updatedAt)
	if err != nil {
		return fmt.Errorf("repository: could not query row to database: %s", err.Error())
	}

	if updatedAt.Valid {
		tu, err := time.Parse(time.RFC3339, updatedAt.String)
		if err != nil {
			return fmt.Errorf("repository: could not parse update_at time: %s", err.Error())
		}
		m.UpdatedAt = tu
	}

	return nil
}

// GetAll ...
func (v Fruits) GetAll(limit, offset, order, search string) ([]models.Fruits, error) {

	var query string
	var queryArgs []interface{}

	if search != "" {
		query = `
		SELECT * FROM fruits
		WHERE (fruit_name like ?) 
		OR (fruit_color like ?)`
		queryArgs = append(queryArgs, search, search)
	} else {
		query = `SELECT * FROM fruits ORDER BY fruit_created_at ` + order + ` LIMIT ? OFFSET ?`
		queryArgs = append(queryArgs, limit, offset)
	}

	statement, err := database.GetMariaDBconn().Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("repository: could not prepare query to database: %s", err.Error())
	}
	defer statement.Close()

	rows, err := statement.Query(queryArgs...)
	if err != nil {
		return nil, fmt.Errorf("repository: could not execute query to database: %s", err.Error())
	}
	defer rows.Close()

	var ml []models.Fruits
	for rows.Next() {
		var m models.Fruits
		var updatedAt sql.NullString
		err = rows.Scan(
			&m.ID,
			&m.Name,
			&m.Color,
			&m.ImageURL,
			&m.CreatedAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("repository: could not scan rows: %s", err.Error())
		}
		if updatedAt.Valid {
			tu, err := time.Parse(time.RFC3339, updatedAt.String)
			if err != nil {
				return nil, fmt.Errorf("repository: could not parse update_at time: %s", err.Error())
			}
			m.UpdatedAt = tu
		}
		ml = append(ml, m)
	}

	return ml, nil
}

// UpdateOneByID ...
func (v Fruits) UpdateOneByID(id int, m *models.Fruits) (int, error) {
	query := `
	UPDATE fruits
	SET fruit_name = ?, fruit_color = ?
	WHERE fruit_id = ?`

	statement, err := database.GetMariaDBconn().Prepare(query)
	if err != nil {
		return -1, fmt.Errorf("repository: could not prepare query to database: %s", err.Error())
	}
	defer statement.Close()

	r, err := statement.Exec(m.Name, m.Color, id)
	if err != nil {
		return -1, fmt.Errorf("repository: could not exec query to database: %s", err.Error())
	}

	idInt64, err := r.RowsAffected()
	if err != nil {
		return -1, fmt.Errorf("repository: could not get row affected: %s", err.Error())
	}

	return int(idInt64), nil
}

// UpdatePhotoOneByID ...
func (v Fruits) UpdatePhotoOneByID(id int, url string) (int, error) {

	query := `
	UPDATE fruits
	SET fruit_image = ?
	WHERE fruit_id = ?`
	statement, err := database.GetMariaDBconn().Prepare(query)

	if err != nil {
		return -1, fmt.Errorf("repository: could not prepare query to database: %s", err.Error())
	}

	defer statement.Close()

	r, err := statement.Exec(url, id)
	if err != nil {
		return -1, fmt.Errorf("repository: could not exec query to database: %s", err.Error())
	}

	idInt64, err := r.RowsAffected()
	if err != nil {
		return -1, fmt.Errorf("repository: could not get row affected: %s", err.Error())
	}

	return int(idInt64), nil
}

// DeleteOneByID ...
func (v Fruits) DeleteOneByID(id int) (int, error) {
	query := `
	DELETE FROM fruits 
	WHERE fruit_id = ?`
	statement, err := database.GetMariaDBconn().Prepare(query)
	if err != nil {
		return -1, fmt.Errorf("repository: could not prepare query to database: %s", err.Error())
	}

	defer statement.Close()

	r, err := statement.Exec(id)
	if err != nil {
		return -1, fmt.Errorf("repository: could not exec query to database: %s", err.Error())
	}

	idInt64, err := r.RowsAffected()
	if err != nil {
		return -1, fmt.Errorf("repository: could not get row affected: %s", err.Error())
	}

	return int(idInt64), nil
}

// IsExistsByID ...
func (v Fruits) IsExistsByID(id int) error {

	query := `SELECT EXISTS(SELECT true FROM fruits WHERE fruit_id = ?)`
	statement, err := database.GetMariaDBconn().Prepare(query)

	if err != nil {
		return fmt.Errorf("repository: could not prepare query to database: %s", err.Error())
	}

	defer statement.Close()

	var exists bool
	err = statement.QueryRow(id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("repository: could not query row to database: %s", err.Error())
	}

	if !exists {
		return errors.New("repository: row not found")
	}

	return nil
}

// Count ...
func (v Fruits) Count() (int64, error) {

	query := `SELECT COUNT(*) FROM fruits`
	statement, err := database.GetMariaDBconn().Prepare(query)

	if err != nil {
		return -1, fmt.Errorf("repository: could not prepare query to database: %s", err.Error())
	}

	defer statement.Close()

	var totalCount int64
	err = statement.QueryRow().Scan(&totalCount)
	if err != nil {
		return -1, fmt.Errorf("repository: could not query row to database: %s", err.Error())
	}

	return totalCount, nil
}
