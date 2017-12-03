package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
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
var index = "simpligo-pln"

var elClient *elastic.Client
var err error

func Init() {

	elClient, err = elastic.NewClient(
		elastic.SetURL(elAddress),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
	)
	if err != nil {
		panic(err)
	}

	createIndexIfNotExists()
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
	//r.HandleFunc("/anotador", AnotadorHandler)
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
	pageInfo = PageInfo{
		Version: "0.5.1",
	}

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
		return
	}

	// carrega main page
	t, err := template.New("menu.html").Delims("[[", "]]").ParseFiles("./templates/menu.html")
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

		log.Println(claims["usr"])
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
			Index(index).
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
		Index(index).
		Query(query).
		Type("user").
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

func createIndexIfNotExists() {

	exists, err := elClient.IndexExists(index).Do(context.Background())
	if err != nil {
		panic(err)
	}
	if !exists {
		settings := `{
			"settings":{
				"number_of_shards":5,
				"number_of_replicas":1
			},
			"mappings" : {
				"user" : {
					"properties" : {
						"email" : { "type" : "text" },
						"name" : { "type" : "text" },
						"pwd" : { "type" : "text" }
					}
				}
			}			
		}`

		createIndex, err := elClient.CreateIndex(index).
			BodyString(settings).
			Do(context.Background())
		if err != nil {
			panic(err)
		}
		if !createIndex.Acknowledged {
			panic("Erro ao criar indice.")
		} else {
			log.Printf("Índice criado: %v", index)
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
