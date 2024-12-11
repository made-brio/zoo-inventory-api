package services

import (
	"database/sql"
	"zoo-inventory/internal/models"
	"zoo-inventory/internal/repository"
)

type AnimalService struct {
	DB *sql.DB
}

func NewAnimalService(db *sql.DB) *AnimalService {
	return &AnimalService{DB: db}
}

// create a new animal entry
func (as *AnimalService) CreateAnimal(animal models.CreateAnimalRequest) error {
	return repository.CreateAnimal(as.DB, animal)
}

// update an existing animal
func (as *AnimalService) UpdateAnimal(animal models.UpdateAnimalRequest) error {
	return repository.UpdateAnimal(as.DB, animal)
}

// delete an existing animal
func (as *AnimalService) DeleteAnimal(id int) error {
	return repository.DeleteAnimal(as.DB, id)
}

// request to get a list of all currently existing animals
func (as *AnimalService) GetAllAnimals() ([]models.Animal, error) {
	return repository.GetAllAnimals(as.DB)
}

// request specifically for an animal by using its ID
func (as *AnimalService) GetAnimalByID(id int) (models.Animal, error) {
	return repository.GetAnimalByID(as.DB, id)
}

// request specifically for an animals by using its class
func (as *AnimalService) GetAnimalsByClass(class string) ([]models.Animal, error) {
	return repository.GetAnimalsByClass(as.DB, class)
}
