package main

import (
        "net/http"
        "github.com/labstack/echo"
		"github.com/labstack/echo/middleware"
		"time"
		"os"
		"git_limited_access/golang/git_mongo"
)
type RepoListJson struct {
	Username string   `json:"username"`
	Repolist []RepoInfoJson `json:"repolist"`
}
type RepoInfoJson struct {
	Name string `json:"name"`
	Expire string `json:"expire"`
	ExpireFlag bool `json:"expire_flag"`
	PasswordFlag bool `json:"password_flag"`
}


func main() {
	e := echo.New()
	e.Use(middleware.CORS())
	fp, err := os.OpenFile("/go/src/git_limited_access/golang/echo.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
	  panic(err)
	}
  
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
	  Output: fp,
	}))
	initRouting(e)

	e.Logger.Fatal(e.Start(":1313"))
}

func initRouting(e *echo.Echo) {
	e.GET("/api", alive)
	e.GET("/api/username", username)
	e.GET("/api/repolist", repolist)
	
}
func TimeToStr(t time.Time) string{
	const layout = "2006-01-02"
    return t.Format(layout)
}

func username(c echo.Context) error {
	username := os.Getenv("GIT_USERNAME")
	return c.JSON(http.StatusOK, map[string]string{"username": username})
}

func alive(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"api": "alive"})
}


func repolist(c echo.Context) error{
	repo ,err := git_mongo.RepoList()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "mongo error"})
	}

	repl := RepoListJson{}
	username := os.Getenv("GIT_USERNAME")
	repl.Username = username
	repl.Repolist = []RepoInfoJson{}
	for _ ,item := range repo {
		if item.ExpireFlag != false && time.Now().Unix() > item.Expire.Unix()  {
			continue
		}
		info := RepoInfoJson{}
		info.Name = item.Name
		info.Expire = TimeToStr(item.Expire)
		info.ExpireFlag = item.ExpireFlag
		if item.Password == "" {
			info.PasswordFlag = false
		}else{
			info.PasswordFlag = true
		}
		repl.Repolist = append(repl.Repolist,info)		
	}
	return c.JSON(http.StatusOK, repl)
	
}