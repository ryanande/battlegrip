package cobraui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/markbates/pkger"
	"github.com/pkg/browser"
	log "github.com/sirupsen/logrus"
)

var directoryPath string
var listeningAddr string
var rootCommand *cobra.Command

func init() {
    directoryPath = "/web/dist"
    listeningAddr = ":8080"
}

func Serve(rootCmd *cobra.Command) error {
	rootCommand = rootCmd
	// The router is now formed by calling the `newRouter` constructor function
	// that we defined above. The rest of the code stays the same
	router := newRouter()

    server := negroni.New(
		negroni.NewRecovery(),
		negroni.HandlerFunc(WebLogger),
	)
	server.UseHandler(router)

    log.Infof("Listening on address: %s", listeningAddr)
	log.Infof("Serving Path: %s", directoryPath)

	browser.OpenURL("http://localhost"+listeningAddr)

	err := http.ListenAndServe(listeningAddr, server)
	if err != nil {
		return err
	}
	return nil
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{ status:green }")
}

func cobraCommandHandler(resp http.ResponseWriter, req *http.Request) {

	cmds, err := GetCommandDetails(rootCommand)
	if err != nil {
		log.Error(err)
		_, _ = resp.Write([]byte(err.Error()))
	}

	resp.Header().Add("content-type", "application/json")
	jsonByteData, err := json.Marshal(cmds)
	if err != nil {
		_, _ = resp.Write([]byte(err.Error()))
	}
	_, _ = resp.Write(jsonByteData)
}

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/healthcheck", healthHandler).Methods("GET")
	r.HandleFunc("/commands", cobraCommandHandler).Methods("GET")

	staticFileDirectory := pkger.Dir(directoryPath) 
	staticFileHandler := http.StripPrefix("/", http.FileServer(staticFileDirectory))

	r.PathPrefix("/").Handler(staticFileHandler).Methods("GET")
	return r
}

func WebLogger(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	start := time.Now()
	next(rw, req)
	res := rw.(negroni.ResponseWriter)
	defer func() {
		elapsed := time.Since(start)
		log.WithFields(log.Fields{
			"elapsed": elapsed,
			"method":  req.Method,
			"host":    req.URL.Host,
			"path":    req.URL.Path,
			"query":   req.URL.RawQuery,
			"status":  res.Status(),
			"size":    res.Size(),
		}).Info(req.Method + " " + req.URL.Path)
	}()
}

