package battlegrip

import (
	// embed is required
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/pkg/browser"
	log "github.com/sirupsen/logrus"
)

//go:embed index.html
var indexPage []byte

//go:embed favicon.ico
var favIcon []byte

//go:embed favicon.svg
var favIconSvg []byte

var (
	directoryPath string
	listeningAddr string
	rootCommand   *cobra.Command
)

func init() {
	directoryPath = "/web"
	listeningAddr = ":8080"
}

// Serve starts up and runs the http server.
func Serve(cmd *cobra.Command) error {
	rootCommand = cmd
	// The router is now formed by calling `newRouter` defined above.
	router := newRouter()

	server := negroni.New(
		negroni.NewRecovery(),
		negroni.HandlerFunc(webLogger),
	)
	server.UseHandler(router)

	log.Infof("Listening on address: %s", listeningAddr)
	log.Infof("Serving Path: %s", directoryPath)

	err := browser.OpenURL("http://localhost" + listeningAddr)
	if err != nil {
		return err
	}

	err = http.ListenAndServe(listeningAddr, server)
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

	app := ApplicationDetails{
		AssemblyName: filepath.Base(os.Args[0]),
		Command:      *cmds,
	}

	resp.Header().Add("content-type", "application/json")
	jsonByteData, err := json.Marshal(app)
	if err != nil {
		_, _ = resp.Write([]byte(err.Error()))
	}
	_, _ = resp.Write(jsonByteData)
}

func cobraRootCommandHandler(resp http.ResponseWriter, req *http.Request) {

	resp.Header().Add("content-type", "application/json")
	jsonByteData, err := json.Marshal(rootCommand.Commands())
	if err != nil {
		_, _ = resp.Write([]byte(err.Error()))
	}
	_, _ = resp.Write(jsonByteData)
}

func indexCommandHandler(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(200)
	_, err := resp.Write(indexPage)
	if err != nil {
		log.Error(err)
	}
}

func favIconCommandHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("content-type", "image/x-icon")
	resp.WriteHeader(200)
	_, err := resp.Write(favIcon)
	if err != nil {
		log.Error(err)
	}
}

func favIconSvgCommandHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("content-type", "image/svg+xml")
	resp.WriteHeader(200)
	_, err := resp.Write(favIconSvg)
	if err != nil {
		log.Error(err)
	}
}

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/healthcheck", healthHandler).Methods("GET")
	r.HandleFunc("/commands", cobraCommandHandler).Methods("GET")
	r.HandleFunc("/root", cobraRootCommandHandler).Methods("GET")

	r.HandleFunc("/favicon.ico", favIconCommandHandler)
	r.HandleFunc("/favicon.svg", favIconSvgCommandHandler)

	r.HandleFunc("/", indexCommandHandler)

	return r
}

func webLogger(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
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
