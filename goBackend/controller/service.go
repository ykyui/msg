package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"msg/auth"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var apiRouter = make(chan *mux.Router, 1)
var countInit sync.WaitGroup
var LongpollTimeout int
var ReverseProxy int

var (
	userSignupFormatError = 520
	errorUserSignupFormat = errors.New("errorUserSignupFormat")
)

func Run() {
	LongpollTimeout, _ = strconv.Atoi(os.Getenv("LONGPOLLTIMEOUT"))
	ReverseProxy, _ = strconv.Atoi(os.Getenv("ReverseProxy"))
	r := mux.NewRouter()
	// fs := http.FileServer(http.Dir("../reactForntend/build"))
	// Router.Handle("/", http.StripPrefix("/", fs))
	apiRouter <- r.PathPrefix("/api").Subrouter()
	r.PathPrefix("/").Handler(spaHandler{staticPath: "../reactForntend/build", indexPath: "index.html"})
	go func() {
		countInit.Wait()
		headersOk := handlers.AllowedHeaders([]string{"Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization", "X-Requested-With"})
		originsOk := handlers.AllowedOrigins([]string{"*", "http://localhost:3000"})
		methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
		(<-apiRouter).Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				fmt.Println(r.RequestURI)
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(string(body))
				r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
				next.ServeHTTP(rw, r)
			})
		})

		fmt.Println("api ready")
		log.Fatal(http.ListenAndServe(":8000", handlers.CORS(originsOk, headersOk, methodsOk)(r)))
	}()
}

type publicHandlerFunc func(http.ResponseWriter, *http.Request) (interface{}, int, error)

func addPublicApi(path string, next publicHandlerFunc, method []string) {
	r := <-apiRouter
	r.HandleFunc(path, func(rw http.ResponseWriter, r *http.Request) {
		result, statusCode, err := next(rw, r)
		rw.WriteHeader(statusCode)
		if err != nil {
			jsonStr, _ := json.Marshal(err)
			rw.Write(jsonStr)
			return
		}
		jsonStr, _ := json.Marshal(result)
		rw.Write(jsonStr)
	}).Methods(method...)
	apiRouter <- r
}

type privateHandlerFunc func(http.ResponseWriter, *http.Request, string) (interface{}, int, error)

func addPrivateApi(path string, next privateHandlerFunc, method []string) {
	r := <-apiRouter
	r.HandleFunc(path, func(rw http.ResponseWriter, r *http.Request) {
		jwtStr := strings.Split(r.Header.Get("Authorization"), " ")
		if len(jwtStr) != 2 {
			rw.WriteHeader(http.StatusForbidden)
			return
		}
		jwt, err := auth.ParseToken(jwtStr[1])
		if err != nil {
			rw.WriteHeader(http.StatusForbidden)
			return
		}
		id := jwt["jti"].(string)
		result, statusCode, err := next(rw, r, id)
		if result == "cancelled" {
			return
		}
		rw.WriteHeader(statusCode)
		if err != nil {
			jsonStr, _ := json.Marshal(err)
			rw.Write(jsonStr)
			return
		}
		jsonStr, _ := json.Marshal(result)
		rw.Write(jsonStr)
	}).Methods(method...)
	apiRouter <- r
}

type spaHandler struct {
	staticPath string
	indexPath  string
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// // get the absolute path to prevent directory traversal
	// path, err := filepath.Abs(r.URL.Path)
	// if err != nil {
	// 	// if we failed to get the absolute path respond with a 400 bad request
	// 	// and stop
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// prepend the path with the path to the static directory
	path := filepath.Join(h.staticPath)

	// check whether a file exists at the given path
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}
