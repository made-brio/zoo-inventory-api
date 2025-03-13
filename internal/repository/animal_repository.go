package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"zoo-inventory/internal/models"
)

func CreateAnimal(db *sql.DB, payload models.CreateAnimalRequest) error {
	sql := `INSERT INTO animals (id, name, class, legs) VALUES (?, ?, ?, ?)`
	_, err := db.Exec(sql, payload.ID, payload.Name, payload.Class, payload.Legs)
	return err
}

func UpdateAnimal(db *sql.DB, animal models.UpdateAnimalRequest) error {
	setClauses := []string{}
	values := []interface{}{}

	if animal.Name != nil {
		setClauses = append(setClauses, "name = ?")
		values = append(values, *animal.Name)
	}
	if animal.Class != nil {
		setClauses = append(setClauses, "class = ?")
		values = append(values, *animal.Class)
	}
	if animal.Legs != nil {
		setClauses = append(setClauses, "legs = ?")
		values = append(values, *animal.Legs)
	}

	// Pastikan ada field yang di-update
	if len(setClauses) == 0 {
		return fmt.Errorf("no fields to update")
	}

	// Tambahkan WHERE clause di akhir
	values = append(values, animal.ID)
	sqlQuery := fmt.Sprintf("UPDATE animals SET %s WHERE id = ?", strings.Join(setClauses, ", "))

	_, err := db.Exec(sqlQuery, values...)
	return err
}

func DeleteAnimal(db *sql.DB, id int) error {
	sql := "DELETE FROM animals WHERE id = ?"
	_, err := db.Exec(sql, id)
	return err
}

func GetAllAnimals(db *sql.DB) ([]models.Animal, error) {
	sql := "SELECT id, name, class, legs FROM animals"
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Animal
	for rows.Next() {
		var animal models.Animal
		if err := rows.Scan(&animal.ID, &animal.Name, &animal.Class, &animal.Legs); err != nil {
			return nil, err
		}
		result = append(result, animal)
	}

	return result, rows.Err()
}

func GetAnimalByID(db *sql.DB, id int) (models.Animal, error) {
	sqlQuery := "SELECT id, name, class, legs FROM animals WHERE id = ?"
	var result models.Animal
	err := db.QueryRow(sqlQuery, id).Scan(&result.ID, &result.Name, &result.Class, &result.Legs)
	if err == sql.ErrNoRows {
		return result, sql.ErrNoRows
	}
	return result, err
}

func GetAnimalsByClass(db *sql.DB, class string) ([]models.Animal, error) {
	sqlQuery := `SELECT id, name, class, legs FROM animals WHERE class = ?`
	rows, err := db.Query(sqlQuery, class)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Animal
	for rows.Next() {
		var animal models.Animal
		if err := rows.Scan(&animal.ID, &animal.Name, &animal.Class, &animal.Legs); err != nil {
			return nil, err
		}
		result = append(result, animal)
	}

	return result, rows.Err()
}
