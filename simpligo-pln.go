package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
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
		Version:        "0.5.1",
		SessionExpired: false,
		StaticHash:     "002",
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

	r.HandleFunc("/ranker/eval", RankerEvalHandler).Methods("POST")

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
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
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
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
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

func validateJWT(r *http.Request) error {

	tokenString := r.Header.Get("Authorization")

	if tokenString == "" {
		return fmt.Errorf("Token não encontrado no header")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtKey), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		pageInfo.Email = claims["usr"].(string)
		return nil

	} else {
		log.Println(err)
	}

	pageInfo.SessionExpired = true

	return fmt.Errorf("Acesso negado.")
}

type User struct {
	Email  string `json:"email"`
	Name   string `json:"name"`
	Pwd    string `json:"pwd"`
	Source string `json:"src"`
}

func GetHash(s string) string {
	hasher := sha256.New()
	hasher.Write([]byte(s))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}

func createAdminIfNotExists() {

	adm, err := getUser(admEmail)
	if err != nil {
		log.Println(err)

		createUser(admEmail, admKey, "Admin", "raw")

	} else {
		log.Printf("Admin existe: %s", adm.Name)
	}

}

func createUser(email string, key string, name string, userType string) {

	pwdHash := GetHash(email + key)

	user := User{email, name, pwdHash, userType}
	put, err := elClient.Index().
		Index(indexPrefix + "user").
		Type("user").
		BodyJson(user).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	log.Printf("Usuario criado %s - %s\n", email, put.Id)
}

func createAbbrevIfNotExists() {

	_, err := elClient.Search().
		Index(indexPrefix + "abbrev").
		Type("abbrev").
		From(0).Size(100).
		Do(context.Background())
	if err != nil {
		log.Printf("Erro ao listar abbrevs: %v", err)

		for _, a := range abbrevList {
			SaveAbbrev(Abbreviation{Name: a})
		}
	}
}

func getUser(email string) (User, error) {

	query := elastic.NewBoolQuery()
	query = query.Must(elastic.NewTermQuery("email.keyword", email))

	searchResult, err := elClient.Search().
		Index(indexPrefix + "user").
		Type("user").
		Query(query).
		From(0).Size(1).
		Do(context.Background())
	if err != nil {
		return User{}, fmt.Errorf("Erro ao listar usuário: %v", err)
	}

	var user User
	for _, item := range searchResult.Each(reflect.TypeOf(user)) {
		user = item.(User)
		return user, nil
	}

	return User{}, fmt.Errorf("Usuário não encontrado: %v", email)

}

func createIndexIfNotExists(indexSuffix string) {

	exists, err := elClient.IndexExists(indexPrefix + indexSuffix).Do(context.Background())
	if err != nil {
		panic(err)
	}
	if !exists {
		settings := `{
			"settings":{
				"number_of_shards":5,
				"number_of_replicas":1
			}
		}`

		createIndex, err := elClient.CreateIndex(indexPrefix + indexSuffix).
			BodyString(settings).
			Do(context.Background())
		if err != nil {
			panic(err)
		}
		if !createIndex.Acknowledged {
			panic("Erro ao criar indice.")
		} else {
			log.Printf("Índice criado: %v", indexPrefix+indexSuffix)
		}
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	source := r.FormValue("src")
	email := r.FormValue("email")
	pwd := r.FormValue("pwd")
	name := r.FormValue("name")

	if email == "" { // se nao veio como form, tenta pegar como json
		decoder := json.NewDecoder(r.Body)
		var u User
		err := decoder.Decode(&u)
		if err != nil {
			log.Printf("Erro ao tratar payload: %v", err)
		}
		email = u.Email
		pwd = u.Pwd
		source = u.Source
		name = u.Name
	}

	log.Println(source, email, pwd)

	if source == "raw" {
		user, err := getUser(email)
		if err != nil {
			log.Printf("Erro ao obter usuário: %v", err)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		pwdHash := GetHash(email + pwd)

		if pwdHash != user.Pwd {
			log.Printf("Senha invalida para: %v.", email)
			w.WriteHeader(http.StatusForbidden)
			return
		}

	} else if source == "face" {

		urlFace := fmt.Sprintf("https://graph.facebook.com/v3.1/debug_token?input_token=%v&access_token=%v|%v", pwd, "346173842588743", faceSecret)
		client := &http.Client{}
		req, err := http.NewRequest("GET", urlFace, nil)

		resp, err := client.Do(req)
		if err != nil {
			log.Println("Failed to do request:", err)
			w.WriteHeader(http.StatusForbidden)
			return
		}
		defer resp.Body.Close()

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Failed to read response: ", err)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// log.Println("----------", string(respBody))
		if !strings.Contains(string(respBody), `"is_valid":true`) {
			log.Println("Invalid Facebook token: ", string(respBody))
			w.WriteHeader(http.StatusForbidden)
			return
		}

		user, err := getUser(email)
		if err != nil {
			log.Println("Primeiro acesso. Criando.")
			createUser(email, pwd, name, "face")
			user, _ = getUser(email)
		}

		log.Println("Login with Facebook: ", user.Email)

	} else if source == "google" {

		urlFace := fmt.Sprintf("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=%v", pwd)
		client := &http.Client{}
		req, err := http.NewRequest("GET", urlFace, nil)

		resp, err := client.Do(req)
		if err != nil {
			log.Println("Failed to do request:", err)
			w.WriteHeader(http.StatusForbidden)
			return
		}
		defer resp.Body.Close()

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Failed to read response: ", err)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		log.Println("----------", string(respBody))
		if !strings.Contains(string(respBody), email) {
			log.Println("Invalid Google token: ", string(respBody))
			w.WriteHeader(http.StatusForbidden)
			return
		}

		user, err := getUser(email)
		if err != nil {
			log.Println("Primeiro acesso. Criando.")
			createUser(email, pwd, name, "google")
			user, _ = getUser(email)
		}

		log.Println("Login with Google: ", user.Email)

	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"usr": email,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		log.Printf("Error: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(tokenString))
}

type Abbreviation struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func SenterAbbrevNewHandler(w http.ResponseWriter, r *http.Request) {
	err := validateSession(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var abbrev Abbreviation
	err = decoder.Decode(&abbrev)
	if err != nil {
		log.Printf("Erro ao tratar payload: %v", err)
	}

	SaveAbbrev(abbrev)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}

func SaveAbbrev(abbrev Abbreviation) {

	createIndexIfNotExists("abbrev")

	put, err := elClient.Index().
		Refresh("true").
		Index(indexPrefix + "abbrev").
		Type("abbrev").
		BodyJson(abbrev).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	log.Printf("Abreviação criada %s\n", put.Id)

}

func SenterAbbrevRemoveHandler(w http.ResponseWriter, r *http.Request) {
	err := validateSession(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	vars := mux.Vars(r)
	abbrevId := vars["id"]

	_, err = elClient.Delete().
		Refresh("true").
		Index(indexPrefix + "abbrev").
		Type("abbrev").
		Id(abbrevId).
		Do(context.Background())
	if err != nil {
		log.Printf("Erro ao remover: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")

}

func SenterAbbrevListHandler(w http.ResponseWriter, r *http.Request) {
	err := validateSession(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	searchResult, err := elClient.Search().
		Index(indexPrefix + "abbrev").
		Type("abbrev").
		From(0).Size(100).
		Do(context.Background())
	if err != nil {
		log.Printf("Erro ao listar abbrevs: %v", err)
	}

	ret := "{\"list\":["
	if err == nil && searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			var a Abbreviation
			err := json.Unmarshal(*hit.Source, &a)
			if err != nil {
				log.Printf("Erro: %v", err)
			}
			a.Id = hit.Id
			aJson, err := json.Marshal(a)
			if err != nil {
				log.Printf("Erro: %v", err)
			}
			ret += string(aJson) + ","
		}
	}
	ret = ret[0 : len(ret)-1]
	ret += "]}"

	fmt.Fprintf(w, ret)

}

func PalavrasParseHandler(w http.ResponseWriter, r *http.Request) {
	err := validateSession(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	content := r.FormValue("content")
	retType := r.FormValue("type")
	options := r.FormValue("options")

	resp, err := http.PostForm("http://"+palavrasIP+":"+palavrasPort+"/"+retType,
		url.Values{"sentence": {content}, "options": {options}})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error: %v\n", err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Errorf("Error reading response: %v.", err))
	}

	bodyString := string(body)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text")
	fmt.Fprint(w, "SAÍDA: \n"+bodyString)

}

// -----

type Corpus struct {
	Id     string   `json:"id"`
	Name   string   `json:"name"`
	Source string   `json:"source"`
	Genre  string   `json:"genre"`
	Owners []string `json:"owners"`
}

type Text struct {
	Id        string `json:"id"`
	CorpusId  string `json:"corpusId"`
	Name      string `json:"name"`
	Title     string `json:"title"`
	Source    string `json:"source"`
	Level     int    `json:"level"`
	Published string `json:"published"`
}

type TextFull struct {
	Id        string `json:"id"`
	CorpusId  string `json:"corpusId"`
	Name      string `json:"name"`
	Title     string `json:"title"`
	SubTitle  string `json:"subTitle"`
	Source    string `json:"source"`
	Level     int    `json:"level"`
	Published string `json:"published"`
	Author    string `json:"author"`
	Content   string `json:"content"`
	Parsed    struct {
		Paragraphs []struct {
			Idx       int    `json:"idx"`
			Text      string `json:"text"`
			Sentences []struct {
				Idx    int    `json:"idx"`
				Text   string `json:"text"`
				Qtt    int    `json:"qtt"`
				Qtw    int    `json:"qtw"`
				Tokens []struct {
					Idx   int    `json:"idx"`
					Token string `json:"token"`
				} `json:"tokens"`
			} `json:"sentences"`
		} `json:"paragraphs"`
		TotP int `json:"totP"`
		TotS int `json:"totS"`
		TotT int `json:"totT"`
		TotW int `json:"totW"`
	} `json:"parsed"`
}

type Simplification struct {
	Id       string `json:"id"`
	CorpusId string `json:"corpusId"`
	Name     string `json:"name"`
	Title    string `json:"title"`
	Tags     string `json:"tags"`
	From     string `json:"from"`
	To       string `json:"to"`
	Updated  string `json:"updated"`
}

type SimplificationFull struct {
	Id        string `json:"id"`
	CorpusId  string `json:"corpusId"`
	Name      string `json:"name"`
	Title     string `json:"title"`
	Tags      string `json:"tags"`
	From      string `json:"from"`
	To        string `json:"to"`
	Updated   string `json:"updated"`
	Sentences []struct {
		From       string `json:"from"`
		To         string `json:"to"`
		Operations string `json:"operations"`
	} `json:"sentences"`
}

func normalizeEmail(email string) string {
	return strings.Replace(email, "@", "_at_", -1)
}

func AnotadorCorpusNewHandler(w http.ResponseWriter, r *http.Request) {
	err := validateSession(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var corpus Corpus
	err = decoder.Decode(&corpus)
	if err != nil {
		log.Printf("Erro ao tratar payload: %v", err)
	}

	corpus.Owners = []string{normalizeEmail(pageInfo.Email), normalizeEmail(admEmail)}

	createIndexIfNotExists("corpus")

	put, err := elClient.Index().
		Refresh("true").
		Index(indexPrefix + "corpus").
		Type("corpus").
		BodyJson(corpus).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	log.Printf("Corpus criado %s\n", put.Id)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}

func AnotadorCorpusListHandler(w http.ResponseWriter, r *http.Request) {
	err := validateSession(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	query := elastic.NewTermQuery("owners", normalizeEmail(pageInfo.Email))
	searchResult, err := elClient.Search().
		Index(indexPrefix + "corpus").
		Type("corpus").
		Query(query).
		From(0).Size(100).
		Do(context.Background())
	if err != nil {
		log.Printf("Erro ao listar: %v", err)
	}

	ret := "{\"list\":[ "
	if err == nil && searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			var c Corpus
			err := json.Unmarshal(*hit.Source, &c)
			if err != nil {
				log.Printf("Erro: %v", err)
			}
			c.Id = hit.Id
			cJson, err := json.Marshal(c)
			if err != nil {
				log.Printf("Erro: %v", err)
			}
			ret += string(cJson) + ","
		}
	}
	ret = ret[0 : len(ret)-1]
	ret += "]}"

	fmt.Fprintf(w, ret)

}

func AnotadorCorpusRemoveHandler(w http.ResponseWriter, r *http.Request) {
	err := validateSession(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	_, err = elClient.Delete().
		Refresh("true").
		Index(indexPrefix + "corpus").
		Type("corpus").
		Id(id).
		Do(context.Background())
	if err != nil {
		log.Printf("Erro ao remover: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")

}

// TEXTS

func AnotadorTextNewHandler(w http.ResponseWriter, r *http.Request) {
	err := validateSession(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	//debug
	// requestDump, err := httputil.DumpRequest(r, true)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(string(requestDump))

	// decoder := json.NewDecoder(r.Body)
	// var corpus Corpus
	// err = decoder.Decode(&corpus)
	// if err != nil {
	// 	log.Printf("Erro ao tratar payload: %v", err)
	// }

	// corpus.Owners = []string{normalizeEmail(pageInfo.Email), normalizeEmail(admEmail)}

	createIndexIfNotExists("corpus-text")

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(fmt.Errorf("Error reading req: %v.", err))
	}

	put, err := elClient.Index().
		Refresh("true").
		Index(indexPrefix + "corpus-text").
		Type("text").
		BodyString(string(body)).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	log.Printf("Text criado %s\n", put.Id)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, put.Id)
}

func AnotadorTextListHandler(w http.ResponseWriter, r *http.Request) {
	err := validateSession(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	vars := mux.Vars(r)
	corpusId := vars["corpusId"]

	query := elastic.NewTermQuery("corpusId.keyword", corpusId)
	searchResult, err := elClient.Search().
		Index(indexPrefix + "corpus-text").
		Type("text").
		Query(query).
		From(0).Size(100).
		Do(context.Background())
	if err != nil {
		log.Printf("Erro ao listar: %v", err)
	}

	ret := "{\"list\":[ "
	if err == nil && searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			var t Text
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				log.Printf("Erro: %v", err)
			}
			t.Id = hit.Id
			tJson, err := json.Marshal(t)
			if err != nil {
				log.Printf("Erro: %v", err)
			}
			ret += string(tJson) + ","
		}
	}
	ret = ret[0 : len(ret)-1]
	ret += "]}"

	fmt.Fprintf(w, ret)

}

func AnotadorTextRemoveHandler(w http.ResponseWriter, r *http.Request) {
	err := validateSession(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	_, err = elClient.Delete().
		Refresh("true").
		Index(indexPrefix + "corpus-text").
		Type("text").
		Id(id).
		Do(context.Background())
	if err != nil {
		log.Printf("Erro ao remover: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")

}

func AnotadorTextGetHandler(w http.ResponseWriter, r *http.Request) {
	err := validateSession(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	vars := mux.Vars(r)
	// corpusId := vars["corpusId"]
	id := vars["id"]

	query := elastic.NewTermQuery("_id", id)

	searchResult, err := elClient.Search().
		Index(indexPrefix + "corpus-text").
		Type("text").
		Query(query).
		From(0).Size(1).
		Do(context.Background())
	if err != nil {
		log.Printf("Não encontrado: %v", err)
	}

	ret := ""
	if err == nil && searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			var t TextFull
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				log.Printf("Erro: %v", err)
			}
			t.Id = hit.Id
			tJson, err := json.Marshal(t)
			if err != nil {
				log.Printf("Erro: %v", err)
			}
			ret = string(tJson)
		}
	}

	fmt.Fprintf(w, ret)

}

// SIMPLIFICATIONS

func AnotadorSimplNewHandler(w http.ResponseWriter, r *http.Request) {
	err := validateSession(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	createIndexIfNotExists("corpus-simpl")

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(fmt.Errorf("Error reading req: %v.", err))
	}

	put, err := elClient.Index().
		Refresh("true").
		Index(indexPrefix + "corpus-simpl").
		Type("simplification").
		BodyString(string(body)).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	log.Printf("Simplification created %s\n", put.Id)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, put.Id)
}

func AnotadorSimplListHandler(w http.ResponseWriter, r *http.Request) {
	err := validateSession(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	vars := mux.Vars(r)
	corpusId := vars["corpusId"]

	query := elastic.NewTermQuery("corpusId.keyword", corpusId)
	searchResult, err := elClient.Search().
		Index(indexPrefix + "corpus-simpl").
		Type("simplification").
		Query(query).
		From(0).Size(100).
		Do(context.Background())
	if err != nil {
		log.Printf("Erro ao listar: %v", err)
	}

	ret := "{\"list\":[ "
	if err == nil && searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			var s Simplification
			err := json.Unmarshal(*hit.Source, &s)
			if err != nil {
				log.Printf("Erro: %v", err)
			}
			s.Id = hit.Id
			sJson, err := json.Marshal(s)
			if err != nil {
				log.Printf("Erro: %v", err)
			}
			ret += string(sJson) + ","
		}
	}
	ret = ret[0 : len(ret)-1]
	ret += "]}"

	fmt.Fprintf(w, ret)

}

func AnotadorSimplRemoveHandler(w http.ResponseWriter, r *http.Request) {
	err := validateSession(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	_, err = elClient.Delete().
		Refresh("true").
		Index(indexPrefix + "corpus-simpl").
		Type("simplification").
		Id(id).
		Do(context.Background())
	if err != nil {
		log.Printf("Erro ao remover: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")

}

func AnotadorSimplGetHandler(w http.ResponseWriter, r *http.Request) {
	err := validateSession(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	vars := mux.Vars(r)
	// corpusId := vars["corpusId"]
	id := vars["id"]

	query := elastic.NewTermQuery("_id", id)

	searchResult, err := elClient.Search().
		Index(indexPrefix + "corpus-simpl").
		Type("simplification").
		Query(query).
		From(0).Size(1).
		Do(context.Background())
	if err != nil {
		log.Printf("Não encontrado: %v", err)
	}

	ret := ""
	if err == nil && searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			var s SimplificationFull

			// j, _ := hit.Source.MarshalJSON()
			// log.Println(string(j))

			err := json.Unmarshal(*hit.Source, &s)
			if err != nil {
				log.Printf("Erro: %v", err)
			}
			s.Id = hit.Id
			sJson, err := json.Marshal(s)
			if err != nil {
				log.Printf("Erro: %v", err)
			}
			ret = string(sJson)
		}
	}

	fmt.Fprintf(w, ret)

}

func RankerEvalHandler(w http.ResponseWriter, r *http.Request) {
	err := validateSession(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	content := r.FormValue("content")

	resp, err := http.Post("http://"+mainServerIP+":8008/ranker", "text", bytes.NewReader([]byte(content)))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error: %v\n", err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Errorf("Error reading response: %v.", err))
	}

	bodyString := string(body)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text")
	fmt.Fprint(w, bodyString)

}
