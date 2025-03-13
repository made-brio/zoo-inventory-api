package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"zoo-inventory/internal/models"
	"zoo-inventory/internal/services"

	"github.com/gin-gonic/gin"
)

type AnimalController struct {
	AnimalService *services.AnimalService
}

func NewAnimalController(service *services.AnimalService) *AnimalController {
	return &AnimalController{AnimalService: service}
}

func (ac *AnimalController) CreateAnimal(c *gin.Context) {
	var animals []models.CreateAnimalRequest

	// Read raw JSON data
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Try unmarshaling as an array
	if err := json.Unmarshal(body, &animals); err != nil {
		// If it fails, try unmarshaling as a single object
		var singleAnimal models.CreateAnimalRequest
		if err := json.Unmarshal(body, &singleAnimal); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
			return
		}
		animals = append(animals, singleAnimal)
	}

	// Process each animal
	for _, animal := range animals {
		if err := ac.AnimalService.CreateAnimal(animal); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Animal(s) created successfully"})
}

func (ac *AnimalController) UpdateAnimal(ctx *gin.Context) {
	var animal models.UpdateAnimalRequest

	animalId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid animal ID"})
		return
	}

	animal.ID = animalId
	if err := ctx.ShouldBindJSON(&animal); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the animal exists
	_, err = ac.AnimalService.GetAnimalByID(animal.ID)
	if errors.Is(err, sql.ErrNoRows) { // Correct error checking
		if animal.Name == nil || animal.Class == nil || animal.Legs == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}

		createAnimal := models.CreateAnimalRequest{
			ID: animal.ID,
			Name:  *animal.Name,
			Class: *animal.Class,
			Legs:  *animal.Legs,
		}

		err = ac.AnimalService.CreateAnimal(createAnimal)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "Animal created successfully"})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update existing animal
	if err := ac.AnimalService.UpdateAnimal(animal); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Animal updated successfully"})
}

func (ac *AnimalController) GetAllAnimals(c *gin.Context) {
	animals, err := ac.AnimalService.GetAllAnimals()

	if err != nil {
		//If no animal is found then the API should return 404 Not Found
		c.JSON(http.StatusNotFound, gin.H{"result": "There is no animal."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": animals})
}

func (ac *AnimalController) GetAnimalByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	result, err := ac.AnimalService.GetAnimalByID(id)

	switch {
	case err == sql.ErrNoRows:
		c.JSON(http.StatusNotFound, gin.H{"error": "Animal not found"})
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch animal details"})
	default:
		c.JSON(http.StatusOK, result)
	}
}

func (ac *AnimalController) GetAnimalsByClass(c *gin.Context) {
	class := c.Param("classId")
	result, err := ac.AnimalService.GetAnimalsByClass(class)

	switch {
	case err == sql.ErrNoRows:
		c.JSON(http.StatusNotFound, gin.H{"error": "Animal not found"})
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch class details"})
	default:
		c.JSON(http.StatusOK, gin.H{"result": result})
	}
}

func (ac *AnimalController) DeleteAnimal(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := ac.AnimalService.DeleteAnimal(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Animal has been deleted"})
}
