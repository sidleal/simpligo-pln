package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
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
	Answers           map[string]int    `json:"answers"`
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
	query = elastic.NewTermQuery("code.keyword", c.Code)

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

	ret += "Código, Nome Teste, Quantidade Gêneros, Parágrafos por Participante, Nome Participante, Organização, Registro, Semestre, Parágrafos Lidos, Data Início, Hora Início, Parágrafo, Sentença, Índice Palavra, Palavra, Resposta, Tempo Início(ms), Tempo Digitação(ms), Tempo(ms), Tempo Acumulado Parágrafo(ms), Tempo Acumulado Teste(ms)\n"

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
			item.TargetWord = strings.ReplaceAll(item.TargetWord, ",", ".")
			item.GuessWord = strings.ReplaceAll(item.GuessWord, ",", ".")

			ret += fmt.Sprintf("%v,%v,%v,%v,", c.Code, c.Name, c.TotalClasses, c.QtyPerParticipant)
			ret += fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,", part.Name, part.Organization, part.RG, part.Semester, paragraphs, createdDate, createdTime)
			ret += fmt.Sprintf("%v,%v,%v,%v,%v,", item.ParagraphID, item.SentenceSeq, item.WordSeq, item.TargetWord, item.GuessWord)
			ret += fmt.Sprintf("%v,%v,%v,%v,%v\n", item.TimeToStart, item.TypingTime, item.ElapsedTime, item.TimeTotalPar, item.TimeTotal)
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
	RG           string `json:"rg"`
	Age          string `json:"age"`
	Gender       string `json:"gender"`
	Course       string `json:"course"`
	Languages    string `json:"lang"`
	Semester     string `json:"sem"`
	Organization string `json:"org"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	CPF          string `json:"cpf"`
	Created      string `json:"created"`
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
	TypingTime    int64  `json:"time_typing"`
	TimeToStart   int64  `json:"time_to_start"`
	TimeTotal     int64  `json:"time_total"`
	TimeTotalPar  int64  `json:"time_total_par"`
	Saved         string `json:"saved"`
	ParagraphID   int64  `json:"par_id"`
	SentenceID    int64  `json:"sen_id"`
	TokenID       int64  `json:"tok_id"`
	TotWords      int64  `json:"tot_words"`
	ClozeCode     string `json:"code"`
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
		fmt.Fprintf(w, "Código não encontrado: %v.", code)
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
	log.Printf("Dados salvos %s\n", put.Id)

	//atualiza qtd respostas do paragrafo
	if participantData.WordSeq == participantData.TotWords {
		clozeTest := getClozeTest(participantData.ClozeCode)
		if clozeTest.Answers == nil {
			clozeTest.Answers = map[string]int{}
		}
		parID := fmt.Sprintf("%v", participantData.ParagraphID)
		if _, found := clozeTest.Answers[parID]; !found {
			clozeTest.Answers[parID] = 0
		}

		clozeTest.Answers[parID]++
		// log.Println("-------", participantData.ClozeCode, clozeTest.Id, clozeTest.Answers)

		_, err := elClient.Update().
			Index(indexPrefix + "cloze").
			Refresh("true").
			Type("cloze").
			Id(clozeTest.Id).
			Doc(map[string]interface{}{"answers": clozeTest.Answers}).
			Do(context.Background())
		if err != nil {
			log.Println(err)
		}
	}

}

type ParagraphData struct {
	Index         int
	Class         string
	TotItensClass int
	QtyAnswer     int
	ParsedText    senter.ParsedParagraph
}

type ParagraphDataOrder []ParagraphData

func (a ParagraphDataOrder) Len() int      { return len(a) }
func (a ParagraphDataOrder) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ParagraphDataOrder) Less(i, j int) bool {
	if a[i].QtyAnswer == a[j].QtyAnswer {
		return a[i].TotItensClass > a[j].TotItensClass
	}
	return a[i].QtyAnswer < a[j].QtyAnswer
}

func getClozeTest(clozeCode string) ClozeTest {
	query := elastic.NewTermQuery("code.keyword", clozeCode)

	searchResult, err := elClient.Search().
		Index(indexPrefix + "cloze").
		Type("cloze").
		Query(query).
		From(0).Size(1).
		Do(context.Background())
	if err != nil {
		log.Printf("Não encontrado: %v", err)
	}

	clozeTest := ClozeTest{}
	if err == nil && searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			err := json.Unmarshal(*hit.Source, &clozeTest)
			if err != nil {
				log.Printf("Erro: %v", err)
			}
			clozeTest.Id = hit.Id
		}
	}

	return clozeTest
}
func ClozeApplyNewHandler(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var participant ClozeParticipant
	err = decoder.Decode(&participant)
	if err != nil {
		log.Printf("Erro ao tratar payload: %v", err)
	}

	clozeTest := getClozeTest(participant.ClozeCode)

	clozeData := ClozeData{}
	clozeData.Code = participant.ClozeCode
	clozeData.ID = clozeTest.Id

	selectedParagraphs := []int{}
	totClasses, _ := strconv.Atoi(clozeTest.TotalClasses)
	qtyPerPart, _ := strconv.Atoi(clozeTest.QtyPerParticipant)

	mapClassesCount := map[string]int{}
	allPars := []ParagraphData{}
	for _, par := range clozeTest.Parsed.Paragraphs {

		pdata := ParagraphData{}
		pdata.Index = int(par.Idx)
		pdata.QtyAnswer = clozeTest.Answers[fmt.Sprintf("%v", par.Idx)]
		if totClasses < 2 {
			pdata.ParsedText = par
			pdata.Class = "U"
		} else {
			tokens := strings.Split(par.Text, " ")
			group := tokens[0]
			newText := par.Text[len(group)+1:]
			pdata.ParsedText = senter.ParseText(newText).Paragraphs[0]
			pdata.Class = group
		}
		mapClassesCount[pdata.Class]++
		allPars = append(allPars, pdata)

	}

	for i, par := range allPars {
		allPars[i].TotItensClass = mapClassesCount[par.Class]
	}

	sort.Sort(ParagraphDataOrder(allPars))

	// um de cada classe
	mapClasses := map[string]int{}
	for _, par := range allPars {
		log.Println(par.Index, par.QtyAnswer, par.Class, par.ParsedText.QtyWords, par.TotItensClass)
		if _, found := mapClasses[par.Class]; !found {
			mapClasses[par.Class] = par.Index
		}
	}

	for _, v := range mapClasses {
		selectedParagraphs = append(selectedParagraphs, v)
	}

	//resto, escolhe apenas por ordem de menos respondida
	tail := qtyPerPart - totClasses
	if tail > 0 {
		i := 1
		for _, par := range allPars {
			if i > tail {
				break
			}
			if mapClasses[par.Class] != par.Index {
				selectedParagraphs = append(selectedParagraphs, par.Index)
				i++
			}
		}

	}

	participant.Paragraphs = selectedParagraphs
	participant = createClozeParticipantIfNotExists(participant)

	clozeData.Participant = participant

	log.Println("------", participant.Paragraphs)

	//train
	trainPar := "O nosso país, Brasil, é cheio de riquezas naturais e culturais. Não importa para onde formos, encontraremos belas paisagens e uma história rica a ser contada."
	trainSenter := senter.ParseText(trainPar)

	clozeData.Paragraphs = []senter.ParsedParagraph{}
	clozeData.Paragraphs = append(clozeData.Paragraphs, trainSenter.Paragraphs[0])
	for _, p := range allPars {
		for _, pn := range participant.Paragraphs {
			if p.Index == pn {
				p.ParsedText.Idx = int64(p.Index)
				clozeData.Paragraphs = append(clozeData.Paragraphs, p.ParsedText)
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
	// query = query.Must(elastic.NewTermQuery("org.keyword", participant.Organization))
	query = query.Must(elastic.NewTermQuery("rg.keyword", participant.RG))

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
