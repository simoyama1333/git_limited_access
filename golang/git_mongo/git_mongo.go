package git_mongo

import (
	"fmt"
	"time"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strings"
	"encoding/base64"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
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
type RepoInfo struct{
	Name string
	Expire time.Time
	ExpireFlag bool
	Password string
	Json string
}
const DBName string = "git_limited"
const ColRepoInfo string = "repoinfo"
const GitAPIURL string = "https://api.github.com/repos/"

var token string
var repourl string

func RepoCrawl(username string,_token string,reponame string,expire string,password string) {
	token = _token
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:mongo@mongodb:27017"))
    if err != nil {
        panic(err)
    }
    if err = client.Connect(context.Background()); err != nil {
        panic(err)
    }
	defer client.Disconnect(context.Background())


	repourl =  GitAPIURL + username + "/" + reponame + "/"
	gitFirstUrl := repourl + "contents?access_token=" + token
	fmt.Println("Now crawling repository ")
	contents := GetContentsJson(gitFirstUrl)
	data := ContentsToDataRecursively(contents)
	fmt.Println("Crawling repository is end")
	repojson, _ := json.Marshal(data)

	var expiretime time.Time
	var expireflag bool
	if expire == ""{
		expiretime = time.Now()
		expireflag = false
	}else{
		layout := "2006-01-02"
		expiretime, err = time.Parse(layout, expire)
		if err != nil {
			panic(err)
		}
		expireflag = true
	}

    doc := RepoInfo {
        reponame,
        expiretime,
        expireflag,
		password,
		string(repojson),
	}
	err = RepoCheckAndDelete(reponame)
	//nodocumentは正常
	if (err != mongo.ErrNoDocuments) && (err != nil) {
        panic(err)
	}
	repoinfo := client.Database(DBName).Collection(ColRepoInfo)
	_, err = repoinfo.InsertOne(context.Background(), doc)
    if err != nil {
        panic(err)
	}
	fmt.Println("End")

}
func RepoList() ([]RepoInfo,error){
	var doc []RepoInfo
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:mongo@mongodb:27017"))
    if err != nil {
        return doc,err
    }
    if err = client.Connect(context.Background()); err != nil {
        return doc,err
    }
	defer client.Disconnect(context.Background())

	repoinfo := client.Database(DBName).Collection(ColRepoInfo)
    findOptions := options.Find().SetSort(bson.D{{"name",1}})
    cur, err := repoinfo.Find(context.Background(), bson.D{{}},findOptions)
    if err != nil {
        return doc,err
	}

    if err = cur.All(context.Background(),&doc); err != nil {
        return doc,err
	}
	return doc,nil
}
//Repo単一を返す
func GetRepoInfo(reponame string) (RepoInfo,error){
	var doc RepoInfo
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:mongo@mongodb:27017"))
    if err != nil {
        return doc,err
    }
    if err = client.Connect(context.Background()); err != nil {
        return doc,err
    }
	defer client.Disconnect(context.Background())

	repoinfo := client.Database(DBName).Collection(ColRepoInfo)
	findOptions := options.FindOne()
    err = repoinfo.FindOne(context.Background(), bson.D{{"name",reponame}},findOptions).Decode(&doc)
    if err != nil {
        return doc,err
	}
	return doc,nil
}

func RepoCheckAndDelete(reponame string) error{
	var doc RepoInfo
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:mongo@mongodb:27017"))
    if err != nil {
        return err
    }
    if err = client.Connect(context.Background()); err != nil {
        return err
    }
	defer client.Disconnect(context.Background())

	repoinfo := client.Database(DBName).Collection(ColRepoInfo)
	findOptions := options.FindOne()
    err = repoinfo.FindOne(context.Background(), bson.D{{"name",reponame}},findOptions).Decode(&doc)
    if err != nil {
        return err
	}
	//存在している場合は削除 存在しない場合は上でエラーが帰っている
	_, err = repoinfo.DeleteOne(context.Background(), bson.D{{"name",reponame}})
	if err != nil {
		return err
	}
	return nil
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
func GetFileAndDecode(path string,token string) string{
	giturl := repourl + "contents/" + path + "?access_token=" + token
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