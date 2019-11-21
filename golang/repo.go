package main
import (
	"fmt"
	"flag"
	"git_limited_access/git_mongo"
	"os"
)


func main(){
	reponame := flag.String("repo", "", "[Required] Repository name ")
	expire := flag.String("expire", "", "[Not required] repository access expire  '-expire 2019-11-01'")
	password := flag.String("password", "", "[Not required] password for repository  '-password test'")
	flag.Parse()
	//fmt.Println(*reponame, *expire, *password)
	username := os.Getenv("GIT_USERNAME")
	token := os.Getenv("GIT_TOKEN")
	if username == "" {
		fmt.Println("Set your username GIT_USERNAME")
		return
	}
	if token == "" {
		fmt.Println("Set your git token GIT_TOKEN")
		return
	}
	if *reponame == "" {
		fmt.Println("Put repository name. example '-repo project'")
		return 
	}
	git_mongo.RepoCrawl(username,token,*reponame,*expire,*password)
}