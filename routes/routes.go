package routes

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/akrylysov/pogreb"
	"github.com/gorilla/mux"
)

var PgName string

type imgResponse struct {
	ImgType          string `json:"imgType"`
	ImgBase64Content string `json:"imgBase64Content"`
}
type imgList struct {
	ImgName []string `json:"imgName"`
}

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", indexHandler).Methods("GET")
	router.HandleFunc("/img-upload", imgUploadHandlerHttp).Methods("POST")
	router.HandleFunc("/get-img/{imgName}", getImgHandlerHttp).Methods("GET")
	staticFileDirectory := http.Dir("./static")
	staticFileHandler := http.StripPrefix("/static/", http.FileServer(staticFileDirectory))
	router.PathPrefix("/static/").Handler(staticFileHandler).Methods("GET")
	return router
}

var templates *template.Template

func LoadTemplates(pattern string) {
	templates = template.Must(template.ParseGlob(pattern))
}

func executeTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	templates.ExecuteTemplate(w, tmpl, data)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	dbPgx := dbPgConnect()
	defer dbPgx.Close()
	executeTemplate(w, "index.html", struct {
		Title   string
		ImgList []string
	}{
		Title:   "Mustafar, test",
		ImgList: getImgList(*dbPgx),
	})
}

func imgUploadHandlerHttp(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)

	dbPg, err := pogreb.Open("pogreb.test", nil)
	defer dbPg.Close()
	if err != nil {
		log.Fatal(err)
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer                                           // nadefinujeme buffer pro pole bytů
	defer buf.Reset()                                              // vyprázníme buffer deferem
	io.Copy(&buf, file)                                            // zkopírujeme obsah souboru do bufferu
	strEnc := b64.StdEncoding.EncodeToString([]byte(buf.String())) //prekodujeme binarni data do base64
	err = dbPg.Put([]byte(header.Filename), []byte(strEnc))
	if err != nil {
		log.Fatal(err)
	}

	jsonResponse, err := json.Marshal(imgList{ImgName: getImgList(*dbPg)})
	if err != nil {
		log.Fatal(err)
	}

	w.Write(jsonResponse)
}

func ImgUploadHandlerSocket(nazev string, buf string) []byte {
	dbPg := dbPgConnect()
	defer dbPg.Close()

	err := dbPg.Put([]byte(nazev), []byte(buf))
	if err != nil {
		log.Fatal(err)
	}

	jsonResponse, err := json.Marshal(imgList{ImgName: getImgList(*dbPg)})
	if err != nil {
		log.Fatal(err)
	}

	return jsonResponse
}

func GetImgHandlerSocket(imgName string) []byte {
	dbPg := dbPgConnect()
	defer dbPg.Close()

	fileName := strings.Split(imgName, ".")
	imgType := fileName[1]
	val, err := dbPg.Get([]byte(imgName))
	if err != nil {
		log.Fatal(err)
	}

	jsonResponse, err := json.Marshal(imgResponse{ImgType: imgType, ImgBase64Content: string(val)})
	if err != nil {
		log.Fatal(err)
	}

	return jsonResponse
}

func getImgHandlerHttp(w http.ResponseWriter, r *http.Request) {
	dbPg := dbPgConnect()
	defer dbPg.Close()

	vars := mux.Vars(r)
	imgName := vars["imgName"]
	fileName := strings.Split(imgName, ".")
	val, err := dbPg.Get([]byte(imgName))
	if err != nil {
		log.Fatal(err)
	}

	jsonResponse, err := json.Marshal(imgResponse{ImgType: fileName[1], ImgBase64Content: string(val)})
	if err != nil {
		log.Fatal(err)
	}
	w.Write(jsonResponse)
}

func getImgList(dbPg pogreb.DB) []string {
	var imgArr []string
	it := dbPg.Items()
	for {
		key, _, err := it.Next()
		if err == pogreb.ErrIterationDone {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		imgArr = append(imgArr, string(key))
	}

	return imgArr
}

func dbPgConnect() *pogreb.DB {
	dbPg, err := pogreb.Open(PgName, nil)
	if err != nil {
		log.Fatal(err)
	}

	return dbPg
}

func TruncateDbPg() {
	dbPg := dbPgConnect()
	defer dbPg.Close()
	it := dbPg.Items()
	for {
		key, _, err := it.Next()
		if err == pogreb.ErrIterationDone {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		dbPg.Delete(key)
	}
}
