package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/olivere/elastic"
	"github.com/sidleal/simpligo-pln/tools/senter"
	"github.com/thecodingmachine/gotenberg-go-client"
)

type ClozeTest struct {
	Id                string            `json:"id"`
	Name              string            `json:"name"`
	Code              string            `json:"code"`
	Content           string            `json:"content"`
	Term              string            `json:"term"`
	Parsed            senter.ParsedText `json:"parsed"`
	Owners            []string          `json:"owners"`
	QtyPerParticipant string            `json:"qtyPerPart"`
	TotalClasses      string            `json:"totClass"`
	Answers           map[string]int    `json:"answers"`
	FinalAnswers      map[string]int    `json:"f_answers"`
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
	//here
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

	ret += "Código\tNome Teste\tQuantidade Gêneros\tParágrafos por Participante\tNome Participante\tEmail\tIdade\tGênero\tRegistro\tSemestre\tOrganização\tCurso\tLínguas\tFone\tCPF\tParágrafos Lidos\tData Início\tHora Início\tParágrafo\tSentença\tÍndice Palavra\tPalavra Crua\tPalavra\tResposta\tTempo Início(ms)\tTempo Digitação(ms)\tTempo(ms)\tTempo Acumulado Parágrafo(ms)\tTempo Acumulado Teste(ms)\n"

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

		// part.Name = strings.ReplaceAll(part.Name, ",", "")
		// part.Organization = strings.ReplaceAll(part.Organization, ",", "")
		// part.RG = strings.ReplaceAll(part.RG, ",", "")
		// part.Semester = strings.ReplaceAll(part.Semester, ",", "")
		// part.Languages = strings.ReplaceAll(part.Languages, ",", "")
		// part.Course = strings.ReplaceAll(part.Course, ",", "")
		// part.CPF = strings.ReplaceAll(part.CPF, ",", "")
		// part.Phone = strings.ReplaceAll(part.Phone, ",", "")

		for _, item := range participantDataList {
			created, _ := isoToDate(part.Created)
			createdDate := created.Format("2006-01-02")
			createdTime := created.Format("15:04:05")
			// item.TargetWord = strings.ReplaceAll(item.TargetWord, ",", ".")
			// item.GuessWord = strings.ReplaceAll(item.GuessWord, ",", ".")

			rawWord := getRawWord(c, item)

			ret += fmt.Sprintf("%v\t%v\t%v\t%v\t", c.Code, c.Name, c.TotalClasses, c.QtyPerParticipant)
			ret += fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t", part.Name, part.Email, part.Age, part.Gender, part.RG, part.Semester, part.Organization, part.Course, part.Languages, part.Phone, part.CPF, paragraphs, createdDate, createdTime)
			ret += fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\t", item.ParagraphID, item.SentenceSeq, item.WordSeq, rawWord, item.TargetWord, item.GuessWord)
			ret += fmt.Sprintf("%v\t%v\t%v\t%v\t%v\n", item.TimeToStart, item.TypingTime, item.ElapsedTime, item.TimeTotalPar, item.TimeTotal)
		}

	}

	fmt.Fprintf(w, htmlSafeString(ret))

}

func getRawWord(c ClozeTest, partData ClozeParticipantData) string {
	p := c.Parsed.Paragraphs[partData.ParagraphID-1]
	// s := p.Sentences[partData.SentenceID-1]
	// log.Println(p.Text)
	pText := p.Text
	// if partData.ParagraphID-1 == 3 {
	pText = strings.ReplaceAll(pText, "prevê-se", "prevê se")
	// } else if partData.ParagraphID-1 == 49 {
	pText = strings.ReplaceAll(pText, "procurá-los", "procurá los")
	// } else if partData.ParagraphID-1 == 27 {
	pText = strings.ReplaceAll(pText, "descartá-los", "descartá los")
	// } else if partData.ParagraphID == 2 {
	pText = strings.ReplaceAll(pText, "- de", "-de")
	pText = strings.ReplaceAll(pText, "inédita -", "inédita-")
	// } else if partData.ParagraphID == 32 {
	pText = strings.ReplaceAll(pText, "-- os", "--os")
	pText = strings.ReplaceAll(pText, "natureza --", "natureza--")
	// }

	rawWords := strings.Split(pText, " ")
	// log.Println(rawWords)
	// log.Println(partData.ParagraphID-1, partData.WordSeq, partData.TargetWord)
	rawWord := rawWords[partData.WordSeq]
	// log.Println("---", rawWord)
	log.Println(partData.ParagraphID-1, partData.WordSeq, partData.TargetWord, "-", rawWord)
	return rawWord
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
	Term        template.HTML            `json:"term"`
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

	clozeTest := getClozeTest(code)

	if clozeTest.Code == "" {
		log.Printf("Não encontrado: %v", err)
		fmt.Fprintf(w, "Código não encontrado: %v.", code)
		return
	}

	clozeData := ClozeData{}
	clozeData.Code = code
	clozeData.StaticHash = pageInfo.StaticHash

	termHTML := clozeTest.Term
	termHTML = strings.ReplaceAll(termHTML, "<data-atual-extenso>", formatDateBRFull(time.Now()))
	clozeData.Term = template.HTML(termHTML)

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
		if clozeTest.FinalAnswers == nil {
			clozeTest.FinalAnswers = map[string]int{}
		}
		parID := fmt.Sprintf("%v", participantData.ParagraphID)
		if _, found := clozeTest.FinalAnswers[parID]; !found {
			clozeTest.FinalAnswers[parID] = 0
		}

		clozeTest.FinalAnswers[parID]++
		// log.Println("-------", participantData.ClozeCode, clozeTest.Id, clozeTest.Answers)

		_, err := elClient.Update().
			Index(indexPrefix + "cloze").
			Refresh("true").
			Type("cloze").
			Id(clozeTest.Id).
			Doc(map[string]interface{}{"f_answers": clozeTest.FinalAnswers}).
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

	// allPars, selectedParagraphs := getSelectedParagraphs(clozeTest)
	allPars, selectedParagraphs := getMultiCenterSelectedParagraphs(clozeTest)

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
				updateParagraphAnswers(clozeData.Code, p.Index)
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

func getSelectedParagraphs(clozeTest ClozeTest) ([]ParagraphData, []int) {
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
	return allPars, selectedParagraphs
}

func getMultiCenterSelectedParagraphs(clozeTest ClozeTest) ([]ParagraphData, []int) {
	selectedParagraphs := []int{}
	totClasses, _ := strconv.Atoi(clozeTest.TotalClasses)
	qtyPerPart, _ := strconv.Atoi(clozeTest.QtyPerParticipant)

	mapClassesCount := map[string]int{}
	allParsMap := map[string]ParagraphData{}
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
		allParsMap[fmt.Sprintf("%v", par.Idx)] = pdata
		// allPars = append(allPars, pdata)

	}

	allParsMap2 := map[string]ParagraphData{}
	for k, v := range allParsMap {
		v.TotItensClass = mapClassesCount[v.Class]
		allParsMap2[k] = v
	}

	allTests := []string{"rastros-ufabc", "rastros-usp", "rastros-ufc", "rastros-utfpr", "rastros-uerj", "rastros-puc-rj"}
	for _, test := range allTests {
		if test == clozeTest.Code {
			continue
		}
		clozeTestInfo := getClozeTest(test)
		for k, v := range allParsMap2 {
			if test == "rastros-ufabc" || test == "rastros-uerj" {
				kTmp, _ := strconv.Atoi(k)
				kTmp--
				if value, found := clozeTestInfo.FinalAnswers[fmt.Sprintf("%v", kTmp)]; found {
					v.QtyAnswer = allParsMap[k].QtyAnswer + value
				}
			} else {
				v.QtyAnswer = allParsMap[k].QtyAnswer + clozeTestInfo.FinalAnswers[k]
			}
			allParsMap[k] = v
			log.Println("-->-->", k, v.QtyAnswer, test)
		}
	}

	allPars := []ParagraphData{}
	for _, v := range allParsMap {
		allPars = append(allPars, v)
	}

	sort.Sort(ParagraphDataOrder(allPars))

	//escolhe apenas por ordem de menos respondida
	i := 1
	for _, par := range allPars {
		if i > qtyPerPart {
			break
		}
		selectedParagraphs = append(selectedParagraphs, par.Index)
		i++
	}

	return allPars, selectedParagraphs
}

func updateParagraphAnswers(clozeCode string, paragraphID int) {
	clozeTest := getClozeTest(clozeCode)
	if clozeTest.Answers == nil {
		clozeTest.Answers = map[string]int{}
	}
	parID := fmt.Sprintf("%v", paragraphID)
	if _, found := clozeTest.Answers[parID]; !found {
		clozeTest.Answers[parID] = 0
	}

	clozeTest.Answers[parID]++

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

func ClozeGetTermPDFHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	code := vars["code"]
	doc := vars["doc"]

	pdfPath := fmt.Sprintf("/shared/cloze-data/term/tcle-%s-%s.pdf", code, doc)
	pdf := readFileBytes(pdfPath)

	w.Write(pdf)

}

func ClozeSaveTermPDFHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	code := vars["code"]

	name := r.FormValue("name")
	doc := r.FormValue("doc")

	clozeTest := getClozeTest(code)

	termHTML := clozeTest.Term

	regEx := regexp.MustCompile(`<input.*id="name"[^>]+>`)
	termHTML = regEx.ReplaceAllString(termHTML, "<b>"+name+"</b>")

	regEx = regexp.MustCompile(`<input.*id="doc"[^>]+>`)
	termHTML = regEx.ReplaceAllString(termHTML, "<b>"+doc+"</b>")

	termHTML = strings.ReplaceAll(termHTML, "<data-atual-extenso>", formatDateBRFull(time.Now()))

	termHTML = "<html><head><meta http-equiv=\"Content-Type\" content=\"text/html; charset=UTF-8\"></head><body>\n" + termHTML + "\n</body></html>"

	tempDir := "/shared/cloze-data/tmp/"
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		os.MkdirAll(tempDir, os.ModePerm)
	}
	htmlTempFile := fmt.Sprintf("%s/termo-%s-%s.html", tempDir, code, doc)

	saveFile(htmlTempFile, termHTML)

	termDir := "/shared/cloze-data/term/"
	if _, err := os.Stat(termDir); os.IsNotExist(err) {
		os.MkdirAll(termDir, os.ModePerm)
	}

	c := &gotenberg.Client{Hostname: "http://" + mainServerIP + ":3000"}
	req, _ := gotenberg.NewHTMLRequest(htmlTempFile)
	pdfDestFile := fmt.Sprintf("%s/tcle-%s-%s.pdf", termDir, code, doc)
	req.PaperSize(gotenberg.A4)
	req.Margins(gotenberg.NormalMargins)
	req.Landscape(false)
	c.Store(req, pdfDestFile)

}

func formatDateBRFull(t time.Time) string {
	return fmt.Sprintf("%02d de %s de %4d",
		t.Day(), months[t.Month()-1], t.Year(),
	)
}

var months = [...]string{
	"janeiro", "fevereiro", "março", "abril", "maio", "junho",
	"julho", "agosto", "setembro", "outubro", "novembro", "dezembro",
}

func saveFile(path string, content string) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERRO", err)
	}
	defer f.Close()
	f.Write([]byte(content))
}

func removeFile(path string) {
	err := os.Remove(path)
	if err != nil {
		log.Println("Erro delete:", err)
	}
}

func readFile(path string) string {
	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}

	// r := charmap.ISO8859_1.NewDecoder().Reader(f)
	r := io.Reader(f)

	ret := ""

	buf := make([]byte, 32*1024)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			ret += string(buf[:n])
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			break
		}
	}

	return ret

}

func readFileBytes(path string) []byte {

	dat, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err)
	}
	return dat

}
