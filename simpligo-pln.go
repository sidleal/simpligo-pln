package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/olivere/elastic"
	"golang.org/x/crypto/acme/autocert"
)

type PageInfo struct {
	Version        string `json:"version"`
	Email          string `json:"email"`
	SessionExpired bool   `json:"sessionExp"`
	StaticHash     string `json:"shash"`
	LastPath       string `json:"path"`
}

var pageInfo PageInfo

var elAddress = "http://elasticsearch:9200" // TODO: yml
var jwtKey = "a2lskdjf4jaks2dhfks"
var admEmail = "admin@sidle.al"
var admKey = "simples"
var indexPrefix = "simpligo-pln-"
var abbrevList = []string{"Prof.", "A.C.", "a.C.", "prof."}

var elClient *elastic.Client
var err error

var (
	env          = "dev"
	palavrasIP   = "127.0.0.1"
	palavrasPort = "23080"
	faceSecret   = ""
	mainServerIP = "127.0.0.1"
)

const (
	httpPort = ":8080"
)

func Init() {

	parseFlags()

	pageInfo = PageInfo{
		Version:        "0.5.2",
		SessionExpired: false,
		StaticHash:     "026",
		LastPath:       "/",
	}

	elClient, err = elastic.NewClient(
		elastic.SetURL(elAddress),
		elastic.SetSniff(false),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
	)
	if err != nil {
		panic(err)
	}

	createIndexIfNotExists("user")
	createAdminIfNotExists()
	createAbbrevIfNotExists()

}

func Finalize() {
	elClient.Stop()
}

func Router() *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/login", LoginHandler)
	r.HandleFunc("/senter", SenterHandler)
	r.HandleFunc("/senter/abbrev/new", SenterAbbrevNewHandler).Methods("POST")
	r.HandleFunc("/senter/abbrev/list", SenterAbbrevListHandler)
	r.HandleFunc("/senter/abbrev/{id}", SenterAbbrevRemoveHandler).Methods("DELETE")
	r.HandleFunc("/palavras", PalavrasHandler).Methods("GET")
	r.HandleFunc("/palavras/parse", PalavrasParseHandler).Methods("POST")
	r.HandleFunc("/anotador", AnotadorHandler).Methods("GET")
	r.HandleFunc("/anotador/corpus/new", AnotadorCorpusNewHandler).Methods("POST")
	r.HandleFunc("/anotador/corpus/list", AnotadorCorpusListHandler)
	r.HandleFunc("/anotador/corpus/{id}", AnotadorCorpusRemoveHandler).Methods("DELETE")

	r.HandleFunc("/anotador/corpus/{corpusId}/text/new", AnotadorTextNewHandler).Methods("POST")
	r.HandleFunc("/anotador/corpus/{corpusId}/text/list", AnotadorTextListHandler)
	r.HandleFunc("/anotador/corpus/{corpusId}/text/{id}", AnotadorTextRemoveHandler).Methods("DELETE")
	r.HandleFunc("/anotador/corpus/{corpusId}/text/{id}", AnotadorTextGetHandler).Methods("GET")

	r.HandleFunc("/anotador/corpus/{corpusId}/simpl/new", AnotadorSimplNewHandler).Methods("POST")
	r.HandleFunc("/anotador/corpus/{corpusId}/simpl/list", AnotadorSimplListHandler)
	r.HandleFunc("/anotador/corpus/{corpusId}/simpl/{id}", AnotadorSimplRemoveHandler).Methods("DELETE")
	r.HandleFunc("/anotador/corpus/{corpusId}/simpl/{id}", AnotadorSimplGetHandler).Methods("GET")

	r.HandleFunc("/cloze", ClozeHandler)
	r.HandleFunc("/ranker", RankerHandler)
	r.HandleFunc("/privacidade", PrivacidadeHandler)

	r.HandleFunc("/ranker/ws", RankerWebSocketHandler)

	r.HandleFunc("/cloze/new", ClozeNewHandler).Methods("POST")
	r.HandleFunc("/cloze/list", ClozeListHandler)
	r.HandleFunc("/cloze/{id}", ClozeGetHandler).Methods("GET")
	r.HandleFunc("/cloze/export/{id}", ClozeExportHandler).Methods("GET")
	r.HandleFunc("/cloze/a/{code}", ClozeApplyHandler).Methods("GET")
	r.HandleFunc("/cloze/apply/new", ClozeApplyNewHandler).Methods("POST")
	r.HandleFunc("/cloze/apply/save", ClozeApplySaveHandler).Methods("POST")
	r.HandleFunc("/cloze/{id}", ClozeRemoveHandler).Methods("DELETE")

	r.HandleFunc("/api/v1/metrix/{subset}/{key}", MetrixAPIPostHandler).Methods("POST")

	r.HandleFunc("/nilcmetrix", MetrixHandler).Methods("GET")
	r.HandleFunc("/nilcmetrixdoc", MetrixDocHandler).Methods("GET")
	r.HandleFunc("/metrix/parse", MetrixParseHandler).Methods("POST")

	r.HandleFunc("/api/v1/sentence-ranker/{key}", SentenceRankerAPIPostHandler).Methods("POST")

	r.HandleFunc("/analysis", AnalysisHandler)
	r.HandleFunc("/analysis/new", AnalysisNewHandler).Methods("POST")
	r.HandleFunc("/analysis/list", AnalysisListHandler)
	r.HandleFunc("/analysis/{id}", AnalysisGetHandler).Methods("GET")
	r.HandleFunc("/analysis/{id}", AnalysisRemoveHandler).Methods("DELETE")

	return r
}

func main() {

	Init()

	var m *autocert.Manager

	var httpsSrv *http.Server
	if env == "prod" {
		hostPolicy := func(ctx context.Context, host string) error {
			allowedHost := "simpligo.sidle.al"
			if host != allowedHost {
				return fmt.Errorf("acme/autocert: only %s host is allowed", allowedHost)
			}
			return nil
		}

		dataDir := "/shared/certs"
		m = &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: hostPolicy,
			Cache:      autocert.DirCache(dataDir),
		}

		httpsSrv = makeHTTPServer()
		httpsSrv.Addr = ":443"
		httpsSrv.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}

		go func() {
			fmt.Printf("Starting HTTPS server on %s\n", httpsSrv.Addr)
			err := httpsSrv.ListenAndServeTLS("", "")
			if err != nil {
				log.Fatalf("httpsSrv.ListendAndServeTLS() failed with %s", err)
			}
		}()
	}

	var httpSrv *http.Server
	if env == "prod" {
		httpSrv = makeHTTPToHTTPSRedirectServer()
		// allow autocert handle Let's Encrypt callbacks over http
		httpSrv.Handler = m.HTTPHandler(httpSrv.Handler)
	} else {
		httpSrv = makeHTTPServer()
	}

	httpSrv.Addr = httpPort
	fmt.Printf("Starting HTTP server on %s\n", httpPort)
	err := httpSrv.ListenAndServe()
	if err != nil {
		log.Fatalf("httpSrv.ListenAndServe() failed with %s", err)
	}

	defer Finalize()

}

func parseFlags() {
	flag.StringVar(&env, "env", "dev", "Environment: dev or prod")
	flag.StringVar(&palavrasIP, "palavras-ip", "127.0.0.1", "IP Palavras")
	flag.StringVar(&palavrasPort, "palavras-port", "23080", "IP Palavras")
	flag.StringVar(&faceSecret, "face-secret", "", "Face App Secret")
	flag.StringVar(&mainServerIP, "main-server-ip", "127.0.0.1", "IP Main Server")
	flag.Parse()
}

func makeHTTPServer() *http.Server {
	mux := Router()
	return &http.Server{
		ReadTimeout:  300 * time.Second,
		WriteTimeout: 300 * time.Second,
		IdleTimeout:  360 * time.Second,
		Handler:      mux,
	}
}

func makeHTTPToHTTPSRedirectServer() *http.Server {
	handleRedirect := func(w http.ResponseWriter, r *http.Request) {
		newURI := "https://" + r.Host + r.URL.String()
		http.Redirect(w, r, newURI, http.StatusFound)
	}
	mux := &http.ServeMux{}
	mux.HandleFunc("/", handleRedirect)
	return &http.Server{
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	TemplateHandler(w, r, "menu", true)
}

func SenterHandler(w http.ResponseWriter, r *http.Request) {
	TemplateHandler(w, r, "senter", true)
}

func ClozeHandler(w http.ResponseWriter, r *http.Request) {
	TemplateHandler(w, r, "cloze", true)
}

func PrivacidadeHandler(w http.ResponseWriter, r *http.Request) {
	TemplateHandler(w, r, "privacidade", false)
}

func PalavrasHandler(w http.ResponseWriter, r *http.Request) {
	TemplateHandler(w, r, "palavras", true)
}

func AnotadorHandler(w http.ResponseWriter, r *http.Request) {
	TemplateHandler(w, r, "anotador", true)
}

func RankerHandler(w http.ResponseWriter, r *http.Request) {
	TemplateHandler(w, r, "ranker", true)
}

func MetrixHandler(w http.ResponseWriter, r *http.Request) {
	TemplateHandler(w, r, "nilcmetrix", true)
}

func MetrixDocHandler(w http.ResponseWriter, r *http.Request) {
	TemplateHandler(w, r, "nilcmetrixdoc", false)
}

func AnalysisHandler(w http.ResponseWriter, r *http.Request) {
	TemplateHandler(w, r, "analysis", true)
}

func validateSession(w http.ResponseWriter, r *http.Request) error {
	err := validateJWT(r)
	if err != nil {
		log.Printf("jwt validate: %v", err)

		// redireciona para login
		t, err := template.New("login.html").Delims("[[", "]]").ParseFiles("./templates/login.html")
		if err != nil {
			fmt.Fprintf(w, "Error openning template: %v", err)
		}
		err = t.Execute(w, pageInfo)
		if err != nil {
			fmt.Fprintf(w, "Error parsing template: %v.", err)
		}
		return fmt.Errorf("Sessao inválida")
	}
	return nil
}

func TemplateHandler(w http.ResponseWriter, r *http.Request, pageName string, checkSession bool) {
	if pageName == "menu" {
		pageInfo.LastPath = "/"
	} else {
		pageInfo.LastPath = "/" + pageName
	}

	if checkSession {
		err := validateSession(w, r)
		if err != nil {
			log.Println(err)
			return
		}
	}

	t, err := template.New(pageName+".html").Delims("[[", "]]").ParseFiles("./templates/" + pageName + ".html")
	if err != nil {
		fmt.Fprintf(w, "Error openning template: %v", err)
	}

	err = t.Execute(w, pageInfo)
	if err != nil {
		fmt.Fprintf(w, "Error parsing template: %v.", err)
	}

}
