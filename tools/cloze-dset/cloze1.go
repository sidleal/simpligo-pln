package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	senter "github.com/sidleal/simpligo-pln/tools/senter"
)

type Word struct {
	UID             string
	TextID          int
	TextPar         string
	TextSent        string
	WordNum         int
	SentenceNum     int
	WordInSent      int
	Word            string
	WordClean       string
	OrtoMatch       float64
	Certainty       float64
	PoS             string
	WordType        string
	PoSMatch        float64
	InflectionMatch float64
	FreqBra         int
	FreqBrWaC       int
	Genre           string
	Top10Resp       []string
}

type WordOrderUID []Word

func (a WordOrderUID) Len() int      { return len(a) }
func (a WordOrderUID) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a WordOrderUID) Less(i, j int) bool {
	return a[i].UID < a[j].UID
}

func readFile(path string) string {
	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}

	// r := charmap.ISO8859_1.NewDecoder().Reader(f)

	ret := ""

	buf := make([]byte, 32*1024)
	for {
		// n, err := r.Read(buf)
		n, err := f.Read(buf)
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

// ------------------------------------

func main() {

	path := "/home/sidleal/sid/usp/cloze_exps/"

	data := readFile(path + "cloze_predict19.tsv")

	sentList := map[string][]Word{}

	lines := strings.Split(data, "\n")
	lastSent := ""
	for i, line := range lines {

		// log.Println(line)

		if line == "" {
			continue
		}

		cols := strings.Split(line, "\t")
		if i == 0 {
			for j, col := range cols {
				log.Println(j, "->", col)
			}
			continue
		}

		word := Word{}
		word.UID = cols[0]
		word.TextID, _ = strconv.Atoi(cols[1])
		word.TextPar = cols[2]
		word.WordNum, _ = strconv.Atoi(cols[3])
		word.SentenceNum, _ = strconv.Atoi(cols[4])
		word.WordInSent, _ = strconv.Atoi(cols[5])
		word.Word = cols[6]
		word.WordClean = cols[7]
		word.OrtoMatch, _ = strconv.ParseFloat(cols[11], 64)
		word.Certainty, _ = strconv.ParseFloat(cols[16], 64)
		word.PoS = cols[17]
		word.WordType = cols[18]
		word.PoSMatch, _ = strconv.ParseFloat(cols[20], 64)
		word.InflectionMatch, _ = strconv.ParseFloat(cols[22], 64)
		word.FreqBra, _ = strconv.Atoi(cols[26])
		word.FreqBrWaC, _ = strconv.Atoi(cols[27])
		word.Genre = cols[39]

		resps := strings.ReplaceAll(cols[40], "{", "")
		resps = strings.ReplaceAll(resps, "}", "")
		word.Top10Resp = strings.Split(resps, ",")

		sent := senter.ParseText(word.TextPar)
		word.TextSent = sent.Paragraphs[0].Sentences[word.SentenceNum-1].Text

		thisSent := fmt.Sprintf("%v_%v", word.TextID, word.SentenceNum)
		if thisSent != lastSent {
			sentList[thisSent] = []Word{}
		}
		lastSent = thisSent

		sentList[thisSent] = append(sentList[thisSent], word)

	}

	// for k, v := range sentList {
	// 	log.Println("\n---------------", k, v[0].TextSent)
	// 	// for i, w := range v {
	// 	// 	log.Println(i, w.UID, w.Word, w.SentenceNum, w.TextSent)
	// 	// }
	// 	w := getWordLessFrequent(v)
	// 	log.Println("=====>", w.UID, w.WordInSent, w.WordInSent, w.Word, w.WordClean, w.PoS, w.FreqBrWaC)
	// }

	sentID := 1
	for _, v := range sentList {
		if len(v) < 2 {
			continue
		}
		w := getWordLessFrequent(v)
		maskedSent := strings.ReplaceAll(w.TextSent, w.Word, "_____")

		fmt.Print("\n")
		fmt.Println(sentID, " --- ", maskedSent)
		resps := getRespList(w.WordClean, w.Top10Resp)
		fmt.Println("--->", w.WordClean, "(", w.PoS, ")", resps)
		sentID++
	}

}

func getRespList(resp string, respList []string) []string {
	ret := []string{}
	for _, item := range respList {
		tokens := strings.Split(item, ":")
		tResp := strings.TrimSpace(tokens[0])
		tResp = strings.Trim(tResp, "'")
		if tResp != resp {
			ret = append(ret, tResp)
		}
	}
	if len(ret) > 4 {
		ret = ret[0:4]
	}
	return ret
}

func getWordLessFrequent(wordList []Word) Word {

	discardedWords := map[string]int{"UID_26_40": 1, "UID_46_18": 1}
	discardedPoS := map[string]int{"NPROP": 1, "ERR": 1}

	var regEx = regexp.MustCompile(`[0-9,\.]+`)

	ret := wordList[1]
	for _, item := range wordList {
		if item.WordNum == 1 {
			continue
		}
		if _, found := discardedWords[item.UID]; found {
			continue
		}
		if _, found := discardedPoS[item.PoS]; found {
			continue
		}
		matchContent := regEx.FindStringSubmatch(item.WordClean)
		if len(matchContent) > 0 {
			continue
		}

		if item.FreqBrWaC < ret.FreqBrWaC {
			ret = item
		}
	}
	return ret
}
