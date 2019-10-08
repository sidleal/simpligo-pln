package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/olivere/elastic"
	"github.com/sidleal/simpligo-pln/tools/senter"
)

type AnalysisData struct {
	ID      string            `json:"id"`
	Title   string            `json:"title"`
	Content string            `json:"content"`
	Parsed  senter.ParsedText `json:"parsed"`
	Owners  []string          `json:"owners"`
}

func AnalysisNewHandler(w http.ResponseWriter, r *http.Request) {
	err := validateSession(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var analysis AnalysisData
	err = decoder.Decode(&analysis)
	if err != nil {
		log.Printf("Erro ao tratar payload: %v", err)
	}

	analysis.Content = strings.Replace(analysis.Content, "\n\n", "\n", -1)

	analysis.Parsed = senter.ParseText(analysis.Content)

	analysis.Owners = []string{normalizeEmail(pageInfo.Email), normalizeEmail(admEmail)}

	createIndexIfNotExists("analysis")

	put, err := elClient.Index().
		Refresh("true").
		Type("analysis").
		Index(indexPrefix + "analysis").
		BodyJson(analysis).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	log.Printf("Analysis criado %s\n", put.Id)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}

func AnalysisListHandler(w http.ResponseWriter, r *http.Request) {
	err := validateSession(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	query := elastic.NewTermQuery("owners", normalizeEmail(pageInfo.Email))
	searchResult, err := elClient.Search().
		Index(indexPrefix + "analysis").
		Type("analysis").
		Query(query).
		From(0).Size(100).
		Do(context.Background())
	if err != nil {
		log.Printf("Erro ao listar: %v", err)
	}

	ret := "{\"list\":[ "
	if err == nil && searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			var c AnalysisData
			err := json.Unmarshal(*hit.Source, &c)
			if err != nil {
				log.Printf("Erro: %v", err)
			}
			c.ID = hit.Id
			c.Parsed = senter.ParsedText{}
			cJSON, err := json.Marshal(c)
			if err != nil {
				log.Printf("Erro: %v", err)
			}
			ret += string(cJSON) + ","
		}
	}
	ret = ret[0 : len(ret)-1]
	ret += "]}"

	fmt.Fprintf(w, ret)

}

func AnalysisGetHandler(w http.ResponseWriter, r *http.Request) {
	err := validateSession(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	query := elastic.NewTermQuery("_id", id)

	searchResult, err := elClient.Search().
		Index(indexPrefix + "analysis").
		Type("analysis").
		Query(query).
		From(0).Size(1).
		Do(context.Background())
	if err != nil {
		log.Printf("Não encontrado: %v", err)
	}

	ret := ""
	if err == nil && searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			var c AnalysisData

			err := json.Unmarshal(*hit.Source, &c)
			if err != nil {
				log.Printf("Erro: %v", err)
			}
			c.ID = hit.Id
			cJSON, err := json.Marshal(c)
			if err != nil {
				log.Printf("Erro: %v", err)
			}
			ret = string(cJSON)

		}
	}

	fmt.Fprintf(w, htmlSafeString(ret))

}

func getAnalysisData(analysisId string) AnalysisData {
	query := elastic.NewTermQuery("id.keyword", analysisId)

	searchResult, err := elClient.Search().
		Index(indexPrefix + "analysis").
		Type("analysis").
		Query(query).
		From(0).Size(1).
		Do(context.Background())
	if err != nil {
		log.Printf("Não encontrado: %v", err)
	}

	analysisData := AnalysisData{}
	if err == nil && searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			err := json.Unmarshal(*hit.Source, &analysisData)
			if err != nil {
				log.Printf("Erro: %v", err)
			}
			analysisData.ID = hit.Id
		}
	}

	return analysisData
}

func AnalysisRemoveHandler(w http.ResponseWriter, r *http.Request) {
	err := validateSession(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	_, err = elClient.Delete().
		Refresh("true").
		Index(indexPrefix + "analysis").
		Type("analysis").
		Id(id).
		Do(context.Background())
	if err != nil {
		log.Printf("Erro ao remover: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")

}
