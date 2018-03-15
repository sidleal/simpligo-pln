package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type SentencePair struct {
	Producao string
	PairId   int
	Level    string
	Split    string
	TextA    string
	TextB    string
}

//generate
func main() {

	includeOriNat := true
	includeNatStr := true
	includeOriStr := true

	includeSplit := true
	onlySizeAligned := true

	pairId := 1

	oriAllSent := readFile("/home/sidleal/usp/coling2018/v3/align_all_ori_nat.tsv")
	lines := strings.Split(oriAllSent, "\n")

	oriAllSentences := []SentencePair{}
	for _, line := range lines {
		if line == "" {
			continue
		}
		tokens := strings.Split(line, "\t")
		sentence := SentencePair{}
		sentence.Producao = tokens[0]
		sentence.Level = tokens[1]
		sentence.Split = tokens[3]
		sentence.TextA = tokens[4]
		sentence.TextB = tokens[5]
		sentence.PairId = pairId
		pairId++

		oriAllSentences = append(oriAllSentences, sentence)

	}

	oriSent := readFile("/home/sidleal/usp/coling2018/v3/align_size_ori_nat.tsv")
	lines = strings.Split(oriSent, "\n")

	oriSentences := []SentencePair{}
	for _, line := range lines {
		if line == "" {
			continue
		}
		tokens := strings.Split(line, "\t")
		sentence := SentencePair{}
		sentence.Producao = tokens[0]
		sentence.Level = tokens[1]
		sentence.Split = tokens[3]
		sentence.TextA = tokens[4]
		sentence.TextB = tokens[5]
		sentence.PairId = pairId
		pairId++

		oriSentences = append(oriSentences, sentence)

	}

	natAllSent := readFile("/home/sidleal/usp/coling2018/v3/align_all_nat_str.tsv")
	lines = strings.Split(natAllSent, "\n")

	natAllSentences := []SentencePair{}
	for _, line := range lines {
		if line == "" {
			continue
		}
		tokens := strings.Split(line, "\t")
		sentence := SentencePair{}
		sentence.Producao = tokens[0]
		sentence.Level = tokens[1]
		sentence.Split = tokens[3]
		sentence.TextA = tokens[4]
		sentence.TextB = tokens[5]
		sentence.PairId = pairId
		pairId++

		natAllSentences = append(natAllSentences, sentence)

	}

	natSent := readFile("/home/sidleal/usp/coling2018/v3/align_size_nat_str.tsv")
	lines = strings.Split(natSent, "\n")

	natSentences := []SentencePair{}
	for _, line := range lines {
		if line == "" {
			continue
		}
		tokens := strings.Split(line, "\t")
		sentence := SentencePair{}
		sentence.Producao = tokens[0]
		sentence.Level = tokens[1]
		sentence.Split = tokens[3]
		sentence.TextA = tokens[4]
		sentence.TextB = tokens[5]
		sentence.PairId = pairId
		pairId++

		natSentences = append(natSentences, sentence)

	}

	oriStrAllSent := readFile("/home/sidleal/usp/coling2018/v3/align_all_ori_str.tsv")
	lines = strings.Split(oriStrAllSent, "\n")

	oriStrAllSentences := []SentencePair{}
	for _, line := range lines {
		if line == "" {
			continue
		}
		tokens := strings.Split(line, "\t")
		sentence := SentencePair{}
		sentence.Producao = tokens[0]
		sentence.Level = tokens[1]
		sentence.Split = tokens[3]
		sentence.TextA = tokens[4]
		sentence.TextB = tokens[5]
		sentence.PairId = pairId
		pairId++

		oriStrAllSentences = append(oriStrAllSentences, sentence)

	}

	oriStrSent := readFile("/home/sidleal/usp/coling2018/v3/align_size_ori_str.tsv")
	lines = strings.Split(oriStrSent, "\n")

	oriStrSentences := []SentencePair{}
	for _, line := range lines {
		if line == "" {
			continue
		}
		tokens := strings.Split(line, "\t")
		sentence := SentencePair{}
		sentence.Producao = tokens[0]
		sentence.Level = tokens[1]
		sentence.Split = tokens[3]
		sentence.TextA = tokens[4]
		sentence.TextB = tokens[5]
		sentence.PairId = pairId
		pairId++

		oriStrSentences = append(oriStrSentences, sentence)

	}

	metrics := readFile("/home/sidleal/usp/coling2018/v3/pss_sentences_features.tsv")
	lines = strings.Split(metrics, "\n")

	metricList := map[string]map[string]string{}
	fieldList := []string{}
	for i, line := range lines {
		if line == "" {
			continue
		}
		line = strings.TrimSuffix(line, "\t")
		if i == 0 { // headers
			header := strings.Split(line, "\t")
			for _, item := range header {
				fieldList = append(fieldList, item)
			}
			continue
		}

		tokens := strings.Split(line, "\t")
		fieldMap := map[string]string{}
		for i, item := range tokens {
			if fieldList[i] != "text" {
				fieldMap[fieldList[i]] = round(item)
			}
		}

		metricList[tokens[0]] = fieldMap
	}

	if onlySizeAligned {
		if includeOriNat {
			buildPairs(oriSentences, metricList, includeSplit)
		}
		if includeNatStr {
			buildPairs(natSentences, metricList, includeSplit)
		}
		if includeOriStr {
			buildPairs(oriStrSentences, metricList, includeSplit)
		}
	} else {

		if includeOriNat {
			buildPairs(oriAllSentences, metricList, includeSplit)
		}
		if includeNatStr {
			buildPairs(natAllSentences, metricList, includeSplit)
		}
		if includeOriStr {
			buildPairs(oriStrAllSentences, metricList, includeSplit)
		}
	}

	for i := 1; i < 12; i++ {
		main_split(fmt.Sprintf("%d", i), int64(7+i), 65)
	}
}

func buildPairs(sentences []SentencePair, metricList map[string]map[string]string, includeSplit bool) {

	f1, err := os.OpenFile("/home/sidleal/usp/coling2018/v3/svmrank/all.dat", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERRO", err)
	}
	defer f1.Close()

	for _, item := range sentences {

		if item.TextA == item.TextB {
			continue
		}

		if !includeSplit && item.Split == "S" {
			continue
		}

		if _, ok := metricList[item.TextA]; !ok {
			continue
		}
		if _, ok := metricList[item.TextB]; !ok {
			continue
		}

		log.Println("----------------")
		// log.Println(item.TextA)
		// log.Println(item.TextB)

		line1 := fmt.Sprintf("%d qid:%d %v # %v - %v", 2, item.PairId, getLine(item.TextA, metricList), item.Level, item.TextA)
		line2 := fmt.Sprintf("%d qid:%d %v # %v - %v", 1, item.PairId, getLine(item.TextB, metricList), item.Level, item.TextB)

		log.Print(line1)
		log.Print(line2)

		_, err := f1.WriteString(line1 + "\n" + line2 + "\n")
		if err != nil {
			log.Println("ERRO", err)
		}

	}
}

func getLine(str string, metricList map[string]map[string]string) string {

	tokens := tokenizeText(str)
	metricList[str]["qty_tokens"] = fmt.Sprintf("%d", len(tokens))
	features := []string{"qty_tokens"}

	// features := []string{"syllables_per_content_word"}
	// features := []string{"clauses_per_sentence"}
	// features := []string{"words_before_main_verb"}

	//features := []string{"syllables_per_content_word", "words_per_sentence", "brunet", "honore", "mcu", "yngve", "frazier", "dep_distance", "words_before_main_verb", "apposition_per_clause", "clauses_per_sentence", "max_noun_phrase", "mean_noun_phrase", "postponed_subject_ratio", "infinite_subordinate_clauses", "non-inflected_verbs", "subordinate_clauses"}
	// features := []string{"flesch", "brunet", "honore"}
	ret := ""

	for i, item := range features {
		val := metricList[str][item]
		if strings.TrimSpace(val) == "" {
			val = "0.000"
		}
		ret += fmt.Sprintf("%d:%v ", i+1, val)
	}

	return ret
}

func round(num string) string {
	ret := floatToStr(tofloat(num))
	// ret = strings.Replace(ret, ".", ",", -1)
	return ret
}

func tofloat(str string) float64 {
	if str == "" {
		str = "0"
	}
	ret, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Println("ERROOOOOOOOOOOOO", err)
	}
	return ret
}

func floatToStr(n float64) string {
	return fmt.Sprintf("%.3f", n)
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

//split train test
func main_split(fold string, seed int64, trainPercentSize int) {

	r := rand.New(rand.NewSource(seed))

	allDataFile := readFile("/home/sidleal/usp/coling2018/v3/svmrank/all.dat")
	lines := strings.Split(allDataFile, "\n")

	dataSize := (len(lines) - 1) / 2

	log.Println(dataSize)

	randomPerm := r.Perm(dataSize)
	log.Println(randomPerm)

	trainSize := dataSize * trainPercentSize / 100
	testSize := dataSize - trainSize

	log.Println("train:", trainSize, "test:", testSize)

	trainData := []string{}
	testData := []string{}

	regEx := regexp.MustCompile(`(qid):([0-9]+)`)

	for i, item := range randomPerm {
		line1 := lines[item*2]
		line2 := lines[item*2+1]
		line1 = regEx.ReplaceAllString(line1, fmt.Sprintf("qid:%d", i+1))
		line2 = regEx.ReplaceAllString(line2, fmt.Sprintf("qid:%d", i+1))
		if i <= trainSize {
			trainData = append(trainData, line1)
			trainData = append(trainData, line2)
		} else {
			testData = append(testData, line1)
			testData = append(testData, line2)
		}
	}

	dir := "/home/sidleal/usp/coling2018/v3/svmrank/t" + fold
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}

	f1, err := os.OpenFile(dir+"/train.dat", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERRO", err)
	}
	defer f1.Close()

	f2, err := os.OpenFile(dir+"/test.dat", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERRO", err)
	}
	defer f2.Close()

	for _, item := range trainData {
		_, err := f1.WriteString(item + "\n")
		if err != nil {
			log.Println("ERRO", err)
		}
	}

	for _, item := range testData {
		_, err := f2.WriteString(item + "\n")
		if err != nil {
			log.Println("ERRO", err)
		}
	}

}

func tokenizeText(rawText string) []string {
	regEx := regexp.MustCompile(`([A-z]+)-([A-z]+)`)
	rawText = regEx.ReplaceAllString(rawText, "$1|hyp|$2")

	regEx = regexp.MustCompile(`\|gdot\|`)
	rawText = regEx.ReplaceAllString(rawText, ".")

	regEx = regexp.MustCompile(`\|gint\|`)
	rawText = regEx.ReplaceAllString(rawText, "?")

	regEx = regexp.MustCompile(`\|gexc\|`)
	rawText = regEx.ReplaceAllString(rawText, "!")

	regEx = regexp.MustCompile(`([\.\,"\(\)\[\]\{\}\?\!;:-]{1})`)
	rawText = regEx.ReplaceAllString(rawText, "  $1 ")

	regEx = regexp.MustCompile(`\s+`)
	rawText = regEx.ReplaceAllString(rawText, " ")

	return strings.Split(rawText, " ")
}
