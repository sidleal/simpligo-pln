package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/olivere/elastic"
	"github.com/sidleal/simpligo-pln/tools/senter"
)

type ClozeTest struct {
	Id                string            `json:"id"`
	Name              string            `json:"name"`
	Code              string            `json:"code"`
	Content           string            `json:"content"`
	Parsed            senter.ParsedText `json:"parsed"`
	Owners            []string          `json:"owners"`
	QtyPerParticipant string            `json:"qtyPerPart"`
	TotalClasses      string            `json:"totClass"`
}

func ClozeNewHandler(w http.ResponseWriter, r *http.Request) {
	err := validateSession(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var cloze ClozeTest
	err = decoder.Decode(&cloze)
	if err != nil {
		log.Printf("Erro ao tratar payload: %v", err)
	}

	cloze.Content = strings.Replace(cloze.Content, "\n\n", "\n", -1)

	cloze.Parsed = senter.ParseText(cloze.Content)

	cloze.Owners = []string{normalizeEmail(pageInfo.Email), normalizeEmail(admEmail)}

	createIndexIfNotExists("cloze")

	put, err := elClient.Index().
		Refresh("true").
		Type("cloze").
		Index(indexPrefix + "cloze").
		BodyJson(cloze).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	log.Printf("Cloze criado %s\n", put.Id)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}

func ClozeListHandler(w http.ResponseWriter, r *http.Request) {
	err := validateSession(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	query := elastic.NewTermQuery("owners", normalizeEmail(pageInfo.Email))
	searchResult, err := elClient.Search().
		Index(indexPrefix + "cloze").
		Type("cloze").
		Query(query).
		From(0).Size(100).
		Do(context.Background())
	if err != nil {
		log.Printf("Erro ao listar: %v", err)
	}

	ret := "{\"list\":[ "
	if err == nil && searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			var c ClozeTest
			err := json.Unmarshal(*hit.Source, &c)
			if err != nil {
				log.Printf("Erro: %v", err)
			}
			c.Id = hit.Id
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

func ClozeGetHandler(w http.ResponseWriter, r *http.Request) {
	err := validateSession(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	query := elastic.NewTermQuery("_id", id)

	searchResult, err := elClient.Search().
		Index(indexPrefix + "cloze").
		Type("cloze").
		Query(query).
		From(0).Size(1).
		Do(context.Background())
	if err != nil {
		log.Printf("Não encontrado: %v", err)
	}

	ret := ""
	if err == nil && searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			var c ClozeTest

			err := json.Unmarshal(*hit.Source, &c)
			if err != nil {
				log.Printf("Erro: %v", err)
			}
			c.Id = hit.Id
			cJSON, err := json.Marshal(c)
			if err != nil {
				log.Printf("Erro: %v", err)
			}
			ret = string(cJSON)

		}
	}

	fmt.Fprintf(w, htmlSafeString(ret))

}

func ClozeExportHandler(w http.ResponseWriter, r *http.Request) {
	err := validateSession(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	query := elastic.NewTermQuery("_id", id)

	searchResult, err := elClient.Search().
		Index(indexPrefix + "cloze").
		Type("cloze").
		Query(query).
		From(0).Size(1).
		Do(context.Background())
	if err != nil {
		log.Printf("Não encontrado: %v", err)
	}

	ret := ""
	var c ClozeTest
	if err == nil && searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			err := json.Unmarshal(*hit.Source, &c)
			if err != nil {
				log.Printf("Erro: %v", err)
			}
			c.Id = hit.Id
		}
	}

	//participantes
	query = elastic.NewTermQuery("code", c.Code)

	searchResult, err = elClient.Search().
		Index(indexPrefix + "cloze-participants").
		Type("participant").
		Query(query).
		From(0).Size(1000).
		Do(context.Background())
	if err != nil {
		log.Printf("Não encontrado: %v", err)
	}

	participantList := []ClozeParticipant{}
	if err == nil && searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			var part ClozeParticipant
			err := json.Unmarshal(*hit.Source, &part)
			if err != nil {
				log.Printf("Erro: %v", err)
			}
			part.ID = hit.Id
			participantList = append(participantList, part)
		}
	}

	ret += "Código, Nome Teste, Quantidade Gêneros, Parágrafos por Participante, Nome Participante, Organização, Registro, Semestre, Parágrafos Lidos, Data Início, Hora Início, Parágrafo, Sentença, Índice Palavra, Palavra, Resposta, Tempo(ms)\n"

	for _, part := range participantList {
		paragraphs := ""
		for _, par := range part.Paragraphs {
			paragraphs += fmt.Sprintf("%v ", par)
		}

		//dados resposta participantes
		query = elastic.NewTermQuery("part.keyword", part.ID)

		searchResult, err = elClient.Search().
			Index(indexPrefix+"cloze-participant-data").
			Type("data").
			Query(query).
			Sort("wseq", true).
			From(0).Size(1000).
			Do(context.Background())
		if err != nil {
			log.Printf("Não encontrado: %v", err)
		}

		participantDataList := []ClozeParticipantData{}
		if err == nil && searchResult.Hits.TotalHits > 0 {
			for _, hit := range searchResult.Hits.Hits {
				var partData ClozeParticipantData
				err := json.Unmarshal(*hit.Source, &partData)
				if err != nil {
					log.Printf("Erro: %v", err)
				}
				participantDataList = append(participantDataList, partData)
			}
		}

		for _, item := range participantDataList {
			created, _ := isoToDate(part.Created)
			createdDate := created.Format("2006-01-02")
			createdTime := created.Format("15:04:05")
			ret += fmt.Sprintf("%v,%v,%v,%v,", c.Code, c.Name, c.TotalClasses, c.QtyPerParticipant)
			ret += fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,", part.Name, part.Organization, part.RegNumber, part.Semester, paragraphs, createdDate, createdTime)
			ret += fmt.Sprintf("%v,%v,%v,%v,%v,%v\n", item.ParagraphID, item.SentenceSeq, item.WordSeq, item.TargetWord, item.GuessWord, item.ElapsedTime)
		}

	}

	fmt.Fprintf(w, htmlSafeString(ret))

}

func dateToISO(date time.Time) string {
	return date.Format("2006-01-02T15:04:05.000Z07:00")
}

func isoToDate(date string) (time.Time, error) {
	ret, err := time.Parse("2006-01-02T15:04:05.000Z07:00", date)
	if err != nil {
		return ret, fmt.Errorf("erro convertendo data: %v", err)
	}
	return ret, nil
}

func htmlSafeString(str string) string {
	str = strings.ReplaceAll(str, "%", "%%")
	return str
}

type ClozeData struct {
	ID          string                   `json:"id"`
	Code        string                   `json:"code"`
	Participant ClozeParticipant         `json:"part"`
	Paragraphs  []senter.ParsedParagraph `json:"prgphs"`
	StaticHash  string                   `json:"shash"`
}

type ClozeParticipant struct {
	ID           string `json:"id"`
	ClozeCode    string `json:"code"`
	Name         string `json:"name"`
	Organization string `json:"org"`
	RegNumber    string `json:"ra"`
	Semester     string `json:"sem"`
	Created      string `json:"created"`
	Birthdate    string `json:"birth"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	RG           string `json:"rg"`
	CPF          string `json:"cpf"`
	Paragraphs   []int  `json:"prgphs"`
}

type ClozeParticipantData struct {
	ParticipantID string `json:"part"`
	ParagraphSeq  int64  `json:"para"`
	SentenceSeq   int64  `json:"sent"`
	WordSeq       int64  `json:"wseq"`
	TargetWord    string `json:"tword"`
	GuessWord     string `json:"word"`
	ElapsedTime   int64  `json:"time"`
	Saved         string `json:"saved"`
	ParagraphID   int64  `json:"par_id"`
	SentenceID    int64  `json:"sen_id"`
	TokenID       int64  `json:"tok_id"`
}

func ClozeApplyHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	code := vars["code"]

	query := elastic.NewTermQuery("code.keyword", code)

	searchResult, err := elClient.Search().
		Index(indexPrefix + "cloze").
		Type("cloze").
		Query(query).
		From(0).Size(1).
		Do(context.Background())

	if err != nil || searchResult.Hits.TotalHits < 1 {
		log.Printf("Não encontrado: %v", err)
		fmt.Fprintf(w, "Cóidigo não encontrado: %v.", code)
		return
	}

	clozeData := ClozeData{}
	clozeData.Code = code
	clozeData.StaticHash = pageInfo.StaticHash

	t, err := template.New("cloze_apply.html").Delims("[[", "]]").ParseFiles("./templates/cloze_apply.html")
	if err != nil {
		fmt.Fprintf(w, "Error openning template: %v", err)
	}

	err = t.Execute(w, clozeData)
	if err != nil {
		fmt.Fprintf(w, "Error parsing template: %v.", err)
	}

}

func ClozeApplySaveHandler(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var participantData ClozeParticipantData
	err = decoder.Decode(&participantData)
	if err != nil {
		log.Printf("Erro ao tratar payload: %v", err)
	}

	participantData.Saved = dateToISO(time.Now())

	createIndexIfNotExists("cloze-participant-data")

	put, err := elClient.Index().
		Refresh("true").
		Type("data").
		Index(indexPrefix + "cloze-participant-data").
		BodyJson(participantData).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	log.Printf("Dados savos %s\n", put.Id)

}

func ClozeApplyNewHandler(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var participant ClozeParticipant
	err = decoder.Decode(&participant)
	if err != nil {
		log.Printf("Erro ao tratar payload: %v", err)
	}

	query := elastic.NewTermQuery("code.keyword", participant.ClozeCode)

	searchResult, err := elClient.Search().
		Index(indexPrefix + "cloze").
		Type("cloze").
		Query(query).
		From(0).Size(1).
		Do(context.Background())
	if err != nil {
		log.Printf("Não encontrado: %v", err)
	}

	clozeData := ClozeData{}
	clozeData.Code = participant.ClozeCode

	clozeTest := ClozeTest{}
	if err == nil && searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			err := json.Unmarshal(*hit.Source, &clozeTest)
			if err != nil {
				log.Printf("Erro: %v", err)
			}
			clozeData.ID = hit.Id
		}
	}

	selectedParagraphs := []int{}
	totClasses, _ := strconv.Atoi(clozeTest.TotalClasses)
	qtyPerPart, _ := strconv.Atoi(clozeTest.QtyPerParticipant)

	classes := map[int][]int{}
	parPerClass := int(clozeTest.Parsed.TotalParagraphs) / totClasses

	classNum := 1
	for i := 1; i <= int(clozeTest.Parsed.TotalParagraphs); i++ {
		classes[classNum] = append(classes[classNum], i)
		if i%parPerClass == 0 && classNum < totClasses {
			classNum++
		}
	}

	randomTail := qtyPerPart % totClasses
	for i := 1; i <= qtyPerPart-randomTail; i++ {
		for _, v := range classes {
			// log.Println(k, v)
			randI := rand.Intn(len(v))
			// log.Println("len:", len(v), "rand:", randI)
			randPar := v[randI]
			for isDuplicated(selectedParagraphs, randPar) {
				randI = rand.Intn(len(v))
				randPar = v[randI]
			}
			selectedParagraphs = append(selectedParagraphs, randPar)
			i++
		}
	}

	i := 1
	for i <= randomTail {
		randC := rand.Intn(len(classes)) + 1
		randClass := classes[randC]
		randI := rand.Intn(len(randClass))
		randPar := randClass[randI]
		// log.Println("---------------", selectedParagraphs, randPar)
		if isDuplicated(selectedParagraphs, randPar) {
			// log.Println("--------------- duplicated")
			continue
		}
		// log.Println("--------------- nao duplicou")

		selectedParagraphs = append(selectedParagraphs, randPar)
		i++
	}

	participant.Paragraphs = selectedParagraphs
	participant = createClozeParticipantIfNotExists(participant)

	clozeData.Participant = participant

	log.Println("------", participant.Paragraphs)

	//train
	trainPar := "O rato roeu a roupa do rei de Roma. A rainha ruim resolveu remendar."
	trainSenter := senter.ParseText(trainPar)

	clozeData.Paragraphs = []senter.ParsedParagraph{}
	clozeData.Paragraphs = append(clozeData.Paragraphs, trainSenter.Paragraphs[0])
	for _, p := range clozeTest.Parsed.Paragraphs {
		for _, pn := range participant.Paragraphs {
			if int(p.Idx) == pn {
				clozeData.Paragraphs = append(clozeData.Paragraphs, p)
			}
		}
	}

	cJSON, err := json.Marshal(clozeData)
	if err != nil {
		log.Printf("Erro: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(cJSON))
}

func isDuplicated(selectedParagraphs []int, newVal int) bool {
	duplicated := false
	for _, p := range selectedParagraphs {
		// log.Println("---------------", newVal, p)
		if newVal == p {
			duplicated = true
			break
		}
	}
	return duplicated
}

func createClozeParticipantIfNotExists(participant ClozeParticipant) ClozeParticipant {

	query := elastic.NewBoolQuery()
	query = query.Must(elastic.NewTermQuery("code.keyword", participant.ClozeCode))
	query = query.Must(elastic.NewTermQuery("org.keyword", participant.Organization))
	query = query.Must(elastic.NewTermQuery("ra.keyword", participant.RegNumber))

	searchResult, err := elClient.Search().
		Index(indexPrefix + "cloze-participants").
		Type("participant").
		Query(query).
		From(0).Size(1).
		Do(context.Background())

	if err == nil && searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			err := json.Unmarshal(*hit.Source, &participant)
			if err != nil {
				log.Printf("Erro: %v", err)
			}
			participant.ID = hit.Id
		}

	} else {

		participant.Created = dateToISO(time.Now())

		createIndexIfNotExists("cloze-participants")

		put, err := elClient.Index().
			Refresh("true").
			Type("participant").
			Index(indexPrefix + "cloze-participants").
			BodyJson(participant).
			Do(context.Background())
		if err != nil {
			panic(err)
		}
		participant.ID = put.Id
		log.Printf("Participante criado %s\n", put.Id)
	}

	return participant

}

func ClozeRemoveHandler(w http.ResponseWriter, r *http.Request) {
	err := validateSession(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	_, err = elClient.Delete().
		Refresh("true").
		Index(indexPrefix + "cloze").
		Type("cloze").
		Id(id).
		Do(context.Background())
	if err != nil {
		log.Printf("Erro ao remover: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")

}
