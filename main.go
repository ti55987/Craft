package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"gopkg.in/yaml.v2"
)

type metaData struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type temporaryLink struct {
	MetaData *metaData `json:"metadata"`
	Link     string    `json:"link"`
}

type conf struct {
	AccessToken string `yaml:"access_token"`
}

var tpl *template.Template
var auth string
var netClient *http.Client

func init() {
	tpl = template.Must(template.ParseGlob("view/template/*.html"))
	config, err := getConf()
	if err != nil {
		panic(err)
	}
	auth = fmt.Sprintf("Bearer %s", config.AccessToken)

	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	netClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
}

func getConf() (*conf, error) {
	yamlFile, err := ioutil.ReadFile("config/secret.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err `%v`", err)
		return nil, err
	}

	c := conf{}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return nil, err
	}

	return &c, nil
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/about", about)
	http.HandleFunc("/menu", menu)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// the home page
func home(w http.ResponseWriter, req *http.Request) {
	// rb, err := ioutil.ReadAll(response.Body)
	// if err != nil {
	// 	panic(err)
	// }
	// log.Print(string(rb))
	varMap := map[string]interface{}{
		"link1": getImageLink("/Craft/1.jpg"),
		"link2": getImageLink("/Craft/2.jpg"),
		"link3": getImageLink("/Craft/3.jpg"),
		"link4": getImageLink("/Craft/4.jpg"),
		"link5": getImageLink("/Craft/5.jpg"),
		"link6": getImageLink("/Craft/6.jpg"),
	}
	err := tpl.ExecuteTemplate(w, "index.html", varMap)
	if err != nil {
		log.Fatal(err)
	}
}

// the about page
func about(w http.ResponseWriter, req *http.Request) {
	//io.WriteString(w, "about")
	err := tpl.ExecuteTemplate(w, "about.html", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func menu(w http.ResponseWriter, req *http.Request) {
	err := tpl.ExecuteTemplate(w, "menu.html", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func getImageLink(pathName string) string {
	body := struct {
		Path string `json:"path"`
	}{
		Path: pathName,
	}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(body)
	imgReq, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/files/get_temporary_link", b)
	if err != nil {
		panic(err)
	}
	imgReq.Header.Add("Content-Type", "application/json")
	imgReq.Header.Add("Authorization", auth)

	response, err := netClient.Do(imgReq)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	var dropBoxResp *temporaryLink
	json.NewDecoder(response.Body).Decode(&dropBoxResp)

	if dropBoxResp == nil {
		return ""
	}

	return dropBoxResp.Link
}
