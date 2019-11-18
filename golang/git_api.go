package main

import (
    "flag"
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strings"
)
	
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
type FileOrDir struct {
	Name string
	Path string
	TypeFile bool
	Files []FileOrDir
}

var token string


func main() {
    flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Println("Put repository name")
		return 
	}
	reponame := flag.Args()[0]
	token = "ab242b969e1b170b4966267ff24fe9cf5b538596"

	var username string = "simoyama1333"
	var giturl string = "https://api.github.com/repos/" + username + "/" + reponame + "/contents?access_token=" + token

	fmt.Println("Now crawling")
	contents := GetContentsJson(giturl)
	data := ContentsToDataRecursively(contents)
	fmt.Println(data)
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