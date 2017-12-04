package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/olivere/elastic"
)

type PageInfo struct {
	Version string `json:"version"`
	Email   string `json:"email"`
}

var pageInfo PageInfo

var elAddress = "http://localhost:9200" // TODO: yml
var jwtKey = "a2lskdjf4jaks2dhfks"
var admEmail = "teste@teste.com"
var admKey = "simplifica"
var indexPrefix = "simpligo-pln-"

var elClient *elastic.Client
var err error

func Init() {

	pageInfo = PageInfo{
		Version: "0.5.1",
	}

	elClient, err = elastic.NewClient(
		elastic.SetURL(elAddress),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
	)
	if err != nil {
		panic(err)
	}

	createIndexIfNotExists("user")
	createAdminIfNotExists()

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

	return r
}

func main() {

	Init()

	srv := &http.Server{
		Handler:      Router(),
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	defer Finalize()

	log.Println("Listening for requests on 8080")

	log.Fatal(srv.ListenAndServe())
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	TemplateHandler(w, r, "menu")
}

func SenterHandler(w http.ResponseWriter, r *http.Request) {
	TemplateHandler(w, r, "senter")
}

func PalavrasHandler(w http.ResponseWriter, r *http.Request) {
	TemplateHandler(w, r, "palavras")
}

func AnotadorHandler(w http.ResponseWriter, r *http.Request) {
	TemplateHandler(w, r, "anotador")
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

func TemplateHandler(w http.ResponseWriter, r *http.Request, pageName string) {
	err := validateSession(w, r)
	if err != nil {
		log.Println(err)
		return
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

	return fmt.Errorf("Acesso negado.")
}

type User struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Pwd   string `json:"pwd"`
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

		pwdHash := GetHash(admEmail + admKey)

		admUser := User{admEmail, "Admin", pwdHash}
		put, err := elClient.Index().
			Index(indexPrefix + "user").
			Type("user").
			Id("1").
			BodyJson(admUser).
			Do(context.Background())
		if err != nil {
			panic(err)
		}
		log.Printf("Admin criado %s\n", put.Id)

	} else {
		log.Printf("Admin existe: %s", adm.Name)
	}

}

func getUser(email string) (User, error) {

	query := elastic.NewBoolQuery()
	query = query.Must(elastic.NewMatchQuery("email", email))

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
	email := r.FormValue("email")
	pwd := r.FormValue("pwd")

	if email == "" { // se nao veio como form, tenta pegar como json
		decoder := json.NewDecoder(r.Body)
		var u User
		err := decoder.Decode(&u)
		if err != nil {
			log.Printf("Erro ao tratar payload: %v", err)
		}
		email = u.Email
		pwd = u.Pwd
	}

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

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
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
	if searchResult.Hits.TotalHits > 0 {
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

	palavrasIP := "143.107.183.175"
	palavrasPort := "23380"

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
	Id     string `json:"id"`
	Name   string `json:"name"`
	Source string `json:"source"`
	Genre  string `json:"genre"`
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

	searchResult, err := elClient.Search().
		Index(indexPrefix + "corpus").
		Type("corpus").
		From(0).Size(100).
		Do(context.Background())
	if err != nil {
		log.Printf("Erro ao listar: %v", err)
	}

	ret := "{\"list\":["
	if searchResult.Hits.TotalHits > 0 {
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
