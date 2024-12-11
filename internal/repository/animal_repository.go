package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"zoo-inventory/internal/models"
)

func CreateAnimal(db *sql.DB, payload models.CreateAnimalRequest) (err error) {
	sql := `INSERT INTO animals (name, class, legs) VALUES ($1, $2, $3)`
	_, err = db.Exec(sql, payload.Name, payload.Class, payload.Legs)
	if err != nil {
		return err
	}

	return
}

func UpdateAnimal(db *sql.DB, animal models.UpdateAnimalRequest) error {
	// Dynamic updated query
	setClauses := []string{}
	values := []interface{}{}
	counter := 1

	if animal.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", counter))
		values = append(values, *animal.Name)
		counter++
	}
	if animal.Class != nil {
		setClauses = append(setClauses, fmt.Sprintf("class = $%d", counter))
		values = append(values, *animal.Class)
		counter++
	}
	if animal.Legs != nil {
		setClauses = append(setClauses, fmt.Sprintf("legs = $%d", counter))
		values = append(values, *animal.Legs)
		counter++
	}

	// Ensure at least one field is being updated
	if len(setClauses) == 0 {
		return fmt.Errorf("no fields to update")
	}

	// Add the WHERE clause
	values = append(values, animal.ID)
	sqlQuery := fmt.Sprintf("UPDATE animals SET %s WHERE id = $%d", strings.Join(setClauses, ", "), counter)
	_, err := db.Exec(sqlQuery, values...)
	return err
}

func DeleteAnimal(db *sql.DB, id int) (err error) {
	sql := "DELETE FROM animals WHERE id = $1"
	_, err = db.Exec(sql, id)
	return
}

func GetAllAnimals(db *sql.DB) (result []models.Animal, err error) {
	sql := "SELECT id, name, class, legs FROM animals"

	rows, err := db.Query(sql)
	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		var animal models.Animal
		err = rows.Scan(
			&animal.ID,
			&animal.Name,
			&animal.Class,
			&animal.Legs)

		if err != nil {
			return
		}

		result = append(result, animal)
	}

	return
}

func GetAnimalByID(db *sql.DB, id int) (result models.Animal, err error) {
	sqlQuery := "SELECT id, name, class, legs FROM animals WHERE id = $1"

	// Eksekusi query dengan parameter id
	err = db.QueryRow(sqlQuery, id).Scan(
		&result.ID,
		&result.Name,
		&result.Class,
		&result.Legs,
	)

	// Tangani jika tidak ada data ditemukan
	if err == sql.ErrNoRows {
		return result, nil // Kembalikan struct kosong tanpa error
	} else if err != nil {
		return result, err // Kembalikan error lain jika terjadi
	}

	return result, nil
}

func GetAnimalsByClass(db *sql.DB, class string) (result []models.Animal, err error) {
	sqlQuery := `SELECT id, name, class, legs FROM animals WHERE class = $1`

	rows, err := db.Query(sqlQuery, class)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var animal models.Animal
		err = rows.Scan(
			&animal.ID,
			&animal.Name,
			&animal.Class,
			&animal.Legs,
		)

		if err != nil {
			return nil, err
		}

		result = append(result, animal)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return result, nil
}
