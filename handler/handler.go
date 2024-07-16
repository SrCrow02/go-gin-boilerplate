package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"template/database"
	"template/model"
	"template/repository"
)

func Add(c *gin.Context) {
	client := database.GetClient()
	repo := repository.NewRepository(client, "templateDb", "templateCollection")

	var data model.Data
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := repo.Create(context.Background(), &data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data added successfully"})
}

func Delete(c *gin.Context) {
	client := database.GetClient()
	repo := repository.NewRepository(client, "templateDb", "templateCollection")

	id := c.Param("id")

	if err := repo.Delete(context.Background(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data deleted successfully"})
}

func Update(c *gin.Context) {
	client := database.GetClient()
	repo := repository.NewRepository(client, "templateDb", "templateCollection")

	var updatedData model.Data
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")

	if err := repo.Update(context.Background(), id, &updatedData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data updated successfully"})
}

func Get(c *gin.Context) {
	client := database.GetClient()
	repo := repository.NewRepository(client, "templateDb", "templateCollection")

	id := c.Param("id")

	data, err := repo.Get(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func GetAll(c *gin.Context) {
	client := database.GetClient()
	repo := repository.NewRepository(client, "templateDb", "templateCollection")

	data, err := repo.GetAll(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}
