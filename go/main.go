package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ipreferwater/pikmin/go/db"
	"github.com/ipreferwater/pikmin/go/model"
)

func main() {
	db.InitDatabase()
	r := gin.Default()

	//database
	r.POST("/mock/mongo/pikmin", func(c *gin.Context) { createPikmin(c) })
	r.POST("/mock/mongo/pikmin/bomb", func(c *gin.Context) { givePikminsByColorBombs(c) })
	r.GET("/mock/mongo/pikmin", func(c *gin.Context) { GetPikminsByColor(c) })
	r.DELETE("/mock/mongo/pikmin", func(c *gin.Context) { KillPikminsByID(c) })

	r.Run(":8000")
}

func mockModel(c *gin.Context) {

	var input model.Input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't parse the payload to requiered model", "log": err.Error()})
		return
	}

	if err := work(input, "mock-model"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't work", "log": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "done",
	})
}

func mockAny(c *gin.Context) {

	var inputJson interface{}
	if err := c.ShouldBindJSON(&inputJson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't parse the payload to json", "log": err.Error()})
		return
	}

	if err := work(inputJson, "mock-any"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't work", "log": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "done",
	})
}

//CRUD
func createPikmin(c *gin.Context) {

	var input model.Pikmin
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't parse the payload to requiered model", "log": err.Error()})
		return
	}

	newID, err := db.PikminRepo.CreatePikmin(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't work", "log": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pikminID": newID,
	})

}

func givePikminsByColorBombs(c *gin.Context) {

	var input model.InputByColor
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't parse the payload to requiered model", "log": err.Error()})
		return
	}

	pikmins, err := db.PikminRepo.GetPikminsByColor(input.Color)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't work", "log": err.Error()})
		return
	}
	// not an error
	if len(pikmins) < 1 {
		c.JSON(http.StatusOK, gin.H{"log": "no pikmins of this color were found"})
		return
	}

	count, err := db.PikminRepo.GiveBombs(pikmins)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't work", "log": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"color":          input.Color,
		"obtained_bombs": count,
	})

}

func GetPikminsByColor(c *gin.Context) {
	var input model.InputByColor
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't parse the payload to requiered model", "log": err.Error()})
		return
	}

	pikmins, err := db.PikminRepo.GetPikminsByColor(input.Color)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't work", "log": err.Error()})
		return
	}
	// not an error
	if len(pikmins) < 1 {
		c.JSON(http.StatusOK, gin.H{"log": "no pikmins of this color were found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"color":   input.Color,
		"pikmins": pikmins,
	})

}
func KillPikminsByID(c *gin.Context) {
	var input model.InputByPikminsID
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't parse the payload to requiered model", "log": err.Error()})
		return
	}

	count, err := db.PikminRepo.DeletePikmins(input.IDS)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't work", "log": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pikminsKilled": count,
	})
}

// transform an interface to []byte that can be saved in a file
func work(body interface{}, prefix string) error {
	file, err := json.MarshalIndent(body, "", " ")
	if err != nil {
		return err
	}
	fileName := getFileName(prefix)
	return writeFile(fileName, file)
}

// increment the last digit of output file until it's free
func getFileName(name string) string {
	index := 1
	for {
		fileName := fmt.Sprintf("/output/%s-%d.json", name, index)

		if fileExists(fileName) {
			index++
			continue
		}
		return fileName
	}
}

func writeFile(fileName string, file []byte) error {
	if err := ioutil.WriteFile(fileName, file, 0644); err != nil {
		return err
	}
	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
