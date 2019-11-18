package main

import (
    "flag"
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strings"
	"encoding/base64"
)
//by https://mholt.github.io/json-to-go/
type GitJson struct {
	Name        string      `json:"name"`
	Path        string      `json:"path"`
	Sha         string      `json:"sha"`
	Size        int         `json:"size"`
	URL         string      `json:"url"`
	HTMLURL     string      `json:"html_url"`
	GitURL      string      `json:"git_url"`
	DownloadURL interface{} `json:"download_url"`
	Type        string      `json:"type"`
	Links       struct {
		Self string `json:"self"`
		Git  string `json:"git"`
		HTML string `json:"html"`
	} `json:"_links"`
}
type GitFIle struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Sha         string `json:"sha"`
	Size        int    `json:"size"`
	URL         string `json:"url"`
	HTMLURL     string `json:"html_url"`
	GitURL      string `json:"git_url"`
	DownloadURL string `json:"download_url"`
	Type        string `json:"type"`
	Content     string `json:"content"`
	Encoding    string `json:"encoding"`
	Links       struct {
		Self string `json:"self"`
		Git  string `json:"git"`
		HTML string `json:"html"`
	} `json:"_links"`
}


type FileOrDir struct {
	Name string
	Path string
	TypeFile bool
	Files []FileOrDir
}

var token string
var repourl string

func main() {
    flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Println("Put repository name")
		return 
	}
	reponame := flag.Args()[0]
	token = "ab242b969e1b170b4966267ff24fe9cf5b538596"

	var username string = "simoyama1333"
	repourl = "https://api.github.com/repos/" + username + "/" + reponame + "/"
	/*
	var gitFirtUrl := repourl + "contents?access_token=" + token
	fmt.Println("Now crawling")
	contents := GetContentsJson(gitFirstUrl)
	data := ContentsToDataRecursively(contents)
	fmt.Println(data)
	*/
	a := GetFileAndDecode("README.md")
	fmt.Println(a)


}
//gitのファイルを再帰的に取得
func ContentsToDataRecursively(contents []GitJson) []FileOrDir{
	var data []FileOrDir
	for _,item := range contents{
		info := FileOrDir{}
		info.Name = item.Name
		info.Path = item.Path
		info.Files = []FileOrDir{}
		if item.Type == "dir"{
			info.TypeFile = false
			url := item.URL
			url = strings.Replace(url, "ref=master", "", 1) + "access_token=" + token
			contents := GetContentsJson(url)
			info.Files = ContentsToDataRecursively(contents) 
		}else{
			info.TypeFile = true
		}
		data = append(data,info)
	}
	return data
}

func GetContentsJson(giturl string) []GitJson {
	resp, err := http.Get(giturl)
	if err != nil {
		panic(err)
	}
  
	defer resp.Body.Close()
	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	jsonBytes := ([]byte)(byteArray)
	
	var data []GitJson

    if err := json.Unmarshal(jsonBytes, &data); err != nil {
		panic(err)
	}
	return data
}


//gitのファイルをAPIから取得する場合、Base64デコードが必要となる
func GetFileAndDecode(path string) string{
	giturl := repourl + "contents/" + path + "?access_token=" + token
	fmt.Println(giturl)
	resp, err := http.Get(giturl)
	if err != nil {
		panic(err)
	}
  
	defer resp.Body.Close()
	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	jsonBytes := ([]byte)(byteArray)
	
	var data GitFIle

    if err := json.Unmarshal(jsonBytes, &data); err != nil {
		panic(err)
	}
	fmt.Println(data)

	encoded64 := strings.Replace(data.Content, "\n", "", -1) 
	fmt.Println(encoded64)
	retstr, err := base64.StdEncoding.DecodeString(encoded64)
	if err != nil {
			fmt.Println("nocontent :", err)
			return ""
	}
	return string(retstr)
}