package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	port     = ":8080"
	boredAPI = "https://www.boredapi.com/api/activity"
)

type activityResponse struct {
	Activity string `json:"activity"`
}

func getActivity() (string, error) {
	resp, err := http.Get(boredAPI)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var activityResp activityResponse
	err = json.Unmarshal(body, &activityResp)
	if err != nil {
		return "", err
	}

	return activityResp.Activity, nil
}

func handler(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Please provide a 'name' parameter in the URL."})
		return
	}

	activity, err := getActivity()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get activity from boredapi.com"})
		return
	}
	response := fmt.Sprintf("%s, you should %s.", name, activity)

	c.JSON(http.StatusOK, response)

}

func InitializeHandlers(router *gin.Engine) {
	router.GET("/:name", func(c *gin.Context) {
		handler(c)
	})
}

func main() {
	router := gin.Default()
	InitializeHandlers(router)
	fmt.Printf("Server listening on port %s\n", port)
	router.Run(":8080")
}
