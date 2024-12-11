package controllers

import (
	"database/sql"
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
	var animal models.CreateAnimalRequest
	err := c.BindJSON(&animal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = ac.AnimalService.CreateAnimal(animal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Animal created successfully"})
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
	if err == sql.ErrNoRows {
		var createAnimal models.CreateAnimalRequest
		createAnimal.Name = *animal.Name
		createAnimal.Class = *animal.Class
		createAnimal.Legs = *animal.Legs

		err = ac.AnimalService.CreateAnimal(createAnimal)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		} else {
			ctx.JSON(http.StatusCreated, gin.H{"message": "Animal created successfully"})
		}
		return
	}

	if err := ac.AnimalService.UpdateAnimal(animal); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "animal updated successfully"})
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
