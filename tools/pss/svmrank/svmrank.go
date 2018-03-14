package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type SentencePair struct {
	Producao string
	PairId   int
	Level    string
	TextA    string
	TextB    string
}

func main() {

	pairId := 1

	oriSent := readFile("/home/sidleal/usp/coling2018/v2/align_size_ori_nat.tsv")
	lines := strings.Split(oriSent, "\n")

	oriSentences := []SentencePair{}
	for _, line := range lines {
		if line == "" {
			continue
		}
		tokens := strings.Split(line, "\t")
		sentence := SentencePair{}
		sentence.Producao = tokens[0]
		sentence.Level = tokens[1]
		sentence.TextA = tokens[4]
		sentence.TextB = tokens[5]
		sentence.PairId = pairId
		pairId++

		oriSentences = append(oriSentences, sentence)

	}

	natSent := readFile("/home/sidleal/usp/coling2018/v2/align_size_nat_str.tsv")
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
		sentence.TextA = tokens[4]
		sentence.TextB = tokens[5]
		sentence.PairId = pairId
		pairId++

		natSentences = append(natSentences, sentence)

	}

	oriStrSent := readFile("/home/sidleal/usp/coling2018/v2/align_size_ori_str.tsv")
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
		sentence.TextA = tokens[4]
		sentence.TextB = tokens[5]
		sentence.PairId = pairId
		pairId++

		oriStrSentences = append(oriStrSentences, sentence)

	}

	metrics := readFile("/home/sidleal/usp/coling2018/v2/pss_sentences_features.tsv")
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

	buildPairs(oriSentences, metricList)
	buildPairs(natSentences, metricList)
	buildPairs(oriStrSentences, metricList)

}

func buildPairs(sentences []SentencePair, metricList map[string]map[string]string) {

	f1, err := os.OpenFile("/home/sidleal/usp/coling2018/v2/svmrank_all.dat", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERRO", err)
	}
	defer f1.Close()

	for _, item := range sentences {

		if item.TextA == item.TextB {
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

	features := []string{"syllables_per_content_word", "words_per_sentence", "brunet", "honore", "mcu", "yngve", "frazier", "dep_distance", "words_before_main_verb", "apposition_per_clause", "clauses_per_sentence", "max_noun_phrase", "mean_noun_phrase", "postponed_subject_ratio", "infinite_subordinate_clauses", "non-inflected_verbs", "subordinate_clauses"}
	// features := []string{"flesch", "brunet", "honore"}
	ret := ""

	for i, item := range features {
		ret += fmt.Sprintf("%d:%v ", i+1, metricList[str][item])
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
