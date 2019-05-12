package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/olivere/elastic"
)

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
