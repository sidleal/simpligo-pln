package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/olivere/elastic"
)

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
