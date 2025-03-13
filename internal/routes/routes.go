package routes

import (
	"database/sql"
	"zoo-inventory/internal/controllers"
	"zoo-inventory/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, db *sql.DB) {

	// Initialize Services
	animalService := services.NewAnimalService(db)

	//animals
	animalController := controllers.NewAnimalController(animalService)
	animalRoutes := router.Group("/api/animals")
	{
		animalRoutes.POST("/", animalController.CreateAnimal)
		animalRoutes.PUT("/:id", animalController.UpdateAnimal)
		animalRoutes.GET("/", animalController.GetAllAnimals)
		animalRoutes.GET("/:id", animalController.GetAnimalByID)
		animalRoutes.GET("/class/:classId", animalController.GetAnimalsByClass)
		animalRoutes.DELETE("/:id", animalController.DeleteAnimal)
	}

}
