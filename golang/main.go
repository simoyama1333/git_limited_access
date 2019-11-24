package main

import (
        "net/http"
        "github.com/labstack/echo"
		"github.com/labstack/echo/middleware"
		"time"
		"os"
		"encoding/json"
		"git_limited_access/golang/git_mongo"
		"go.mongodb.org/mongo-driver/mongo"
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

type PostedJson struct {
	Reponame string `json:"reponame"`
	Password string `json:"password"`
	Token string  `json:"token"`
}

type PostedRequest struct{
	Reponame string `json:"reponame"`
	Token string  `json:"token"`
	Path string `json:"path"`
}

type RetAuth struct {
	Newtoken string `json:"newtoken"`
	AuthResult bool `json:"authresult"`
	Readme string `json:"readme"`
	Path string `json:"path"`
	Tree []git_mongo.FileOrDir `json:"tree"`
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
	e.POST("/api/auth", auth)
	e.POST("/api/request", request)
	
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

func auth(c echo.Context) error {
	post := new(PostedJson)
    if err := c.Bind(post); err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "json parse error"})
	}
	
	info,err := git_mongo.GetRepoInfo(post.Reponame)
	if err == mongo.ErrNoDocuments{
		return c.JSON(http.StatusNotFound, map[string]string{"error": "404"})
	}
	if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "mongo error"})
	}
	if info.ExpireFlag != false && time.Now().Unix() > info.Expire.AddDate(0,0,1).Unix()  {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "404"})
	}
	// tokenがない場合はパスワード
	var authok bool
	newtoken := ""
	if post.Token == "" {
		if git_mongo.PasswordAuth(post.Reponame,post.Password) == true{
			newtoken,err = git_mongo.InsertToken(post.Reponame)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "mongo error"})
			}
			authok = true
		}
	}else{
		authok = git_mongo.CheckExistLoginToken(post.Reponame,post.Token)
	}
	//パスワード設定なしならauth通過
	if info.Password == "" {
		authok = true
	}
	if authok == true {
		path := "README.md"
		readme := git_mongo.GetFileAndDecode(path,os.Getenv("GIT_USERNAME"),post.Reponame,os.Getenv("GIT_TOKEN"))
		
		var tree []git_mongo.FileOrDir
		if err := json.Unmarshal([]byte(info.Json), &tree); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "json error"})
		}
		jso := RetAuth{newtoken,true,readme,path,tree}
		return c.JSON(http.StatusOK,jso)
	}
	return c.JSON(http.StatusUnauthorized, map[string]string{"authresult": "false"})
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
		//その日の開始0時に期限切れになるので、その日中はという日本語的解釈において実質的+1日する
		if item.ExpireFlag == true && time.Now().Unix() > item.Expire.AddDate(0,0,1).Unix()  {
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
func request(c echo.Context) error {
	post := new(PostedRequest)
    if err := c.Bind(post); err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "json parse error"})
	}

	info,err := git_mongo.GetRepoInfo(post.Reponame)
	if err == mongo.ErrNoDocuments{
		return c.JSON(http.StatusNotFound, map[string]string{"error": "404"})
	}
	if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "mongo error"})
	}
	if info.ExpireFlag == true && time.Now().Unix() > info.Expire.AddDate(0,0,1).Unix()  {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "404"})
	}

	authok := git_mongo.CheckExistLoginToken(post.Reponame,post.Token)
	//パスワード設定なしならauth通過
	if info.Password == "" {
		authok = true
	}
	if authok == true {
		path := post.Path
		file := git_mongo.GetFileAndDecode(path,os.Getenv("GIT_USERNAME"),post.Reponame,os.Getenv("GIT_TOKEN"))
		if file == ""{
			return c.JSON(http.StatusOK, map[string]string{"error": "404"})
		}else{
			return c.JSON(http.StatusOK, map[string]string{"code": file})
		}
	}else{
		return c.JSON(http.StatusUnauthorized, map[string]string{"authresult": "false"})
	}
	
}