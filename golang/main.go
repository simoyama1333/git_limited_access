package main

import (
        "net/http"
        "github.com/labstack/echo"
		"time"
)
type RepoListJson struct {
	username string   `json:"username"`
	repolist []RepoInfoJson `json:"repolist"`
}
type RepoInfoJson struct {
	reponame string `json:"reponame"`
	expire time.Time `json:"expire"`
	expire_flag bool `json:"expire_flag"`
	password_flag bool `json:"password_flag"`
}

func main() {
	e := echo.New()

	initRouting(e)

	e.Logger.Fatal(e.Start(":80"))
}

func initRouting(e *echo.Echo) {
	e.GET("/", hello)
	e.GET("/username", username)
}

func username(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"username": "simoyama1333"})
}

func hello(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"api": "alive"})
}