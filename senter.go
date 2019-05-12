package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

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
