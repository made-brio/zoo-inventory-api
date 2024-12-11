package models

type UpdateAnimalRequest struct {
	ID    int     `json:"id" binding:"required"`
	Name  *string `json:"name"`
	Class *string `json:"class"`
	Legs  *uint8  `json:"legs"`
}

type CreateAnimalRequest struct {
	Name  string `json:"name" binding:"required"`
	Class string `json:"class"  binding:"required"`
	Legs  uint8  `json:"legs"  binding:"required"`
}
