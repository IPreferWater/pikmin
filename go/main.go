package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ipreferwater/api-output/db"
	"github.com/ipreferwater/api-output/model"
)

func main() {
	db.InitDatabase()
	r := gin.Default()

	//output
	r.POST("/mock/model", func(c *gin.Context) { mockModel(c) })
	r.POST("/mock/any", func(c *gin.Context) { mockAny(c) })

	//database
	r.POST("/mock/mongo/pikmin", func(c *gin.Context) { createPikmin(c) })
	r.POST("/mock/mongo/pikmin/bomb", func(c *gin.Context) { givePikminsByColorBombs(c) })

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

	newID, err := db.CreatePikmin(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't work", "log": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pikminID": newID,
	})

}

func givePikminsByColorBombs(c *gin.Context) {

	var input model.GiveBombs
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't parse the payload to requiered model", "log": err.Error()})
		return
	}

	pikmins, err := db.GetPikminsByColor(input.Color)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't work", "log": err.Error()})
		return
	}
	// not an error
	if len(pikmins) < 1 {
		c.JSON(http.StatusOK, gin.H{"log": "no pikmins of this color were found"})
		return
	}

	count, err := db.GiveBombs(pikmins)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't work", "log": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"color":          "newID",
		"obtained_bombs": count,
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
