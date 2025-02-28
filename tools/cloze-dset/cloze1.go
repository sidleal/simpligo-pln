package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
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
	ParsedSent      senter.ParsedSentence
}

type WordOrderUID []Word

func (a WordOrderUID) Len() int      { return len(a) }
func (a WordOrderUID) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a WordOrderUID) Less(i, j int) bool {
	return a[i].UID < a[j].UID
}

type Phrase struct {
	PhraseNum  int
	PhraseText string
	PhraseType string
	Sentence   string
	ParInfo    Word
	SentData   []Word
}

type PhraseOrder []Phrase

func (a PhraseOrder) Len() int      { return len(a) }
func (a PhraseOrder) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a PhraseOrder) Less(i, j int) bool {
	if a[i].ParInfo.TextID == a[j].ParInfo.TextID {
		if a[i].ParInfo.SentenceNum == a[j].ParInfo.SentenceNum {
			return a[i].PhraseNum < a[j].PhraseNum
		}
		return a[i].ParInfo.SentenceNum < a[j].ParInfo.SentenceNum
	}
	return a[i].ParInfo.TextID < a[j].ParInfo.TextID
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

var mapRespPoS = map[string]string{}

var auditWords1 = map[string]int{}
var auditWords2 = map[string]int{}

func main1() {

	path := "/home/sidleal/sid/usp/cloze_exps/"

	data := readFile(path + "cloze_predict23b.tsv")

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

		auditWords1[cols[0]] = 1

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
		word.ParsedSent = sent.Paragraphs[0].Sentences[word.SentenceNum-1]

		thisSent := fmt.Sprintf("%v_%v", word.TextID, word.SentenceNum)
		if thisSent != lastSent {
			sentList[thisSent] = []Word{}
		}
		lastSent = thisSent

		sentList[thisSent] = append(sentList[thisSent], word)

	}

	dataPoS := readFile(path + "pos_delaf_resps23.tsv")
	lines = strings.Split(dataPoS, "\n")
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

		key := fmt.Sprintf("%v_%v", cols[1], cols[2])
		wPoS := cols[5]
		respPoS := cols[6]
		if wPoS == respPoS {
			mapRespPoS[key] = respPoS
		}
		// log.Println(cols)
	}

	wordGeneralIndex := 1

	phrases := []Phrase{}

	dataPhrases := readFile(path + "oracoes2.tsv")
	lines = strings.Split(dataPhrases, "\n")
	for j, line := range lines {
		// log.Println(i, line)
		if line == "" {
			continue
		}
		cols := strings.Split(line, "\t")
		// log.Println(i, "[", cols[0], "]")
		// log.Println(i, cols[1])
		// log.Println(i, cols[2])

		par := getParInfo(sentList, cols[0])
		// log.Println(par.TextID, par.Genre, par.TextPar)
		// par := sentList[cols[0]]
		// log.Println(par[0].TextID, par[0].TextPar)
		phrase := Phrase{}
		phrase.PhraseText = cols[2]
		phrase.PhraseType = cols[1]
		phrase.Sentence = cols[0]
		phrase.ParInfo = par
		phrase.PhraseNum = j

		thisSent := fmt.Sprintf("%v_%v", par.TextID, par.SentenceNum)
		phrase.SentData = sentList[thisSent]
		phrases = append(phrases, phrase)

		// log.Println("-------------------")

	}

	sort.Sort(PhraseOrder(phrases))

	lastParID := 0
	lastSentID := 0
	phraseCount := 0
	for _, it := range phrases {
		if lastParID != it.ParInfo.TextID {
			fmt.Println("\n================================================")
			fmt.Println("Texto: ", it.ParInfo.TextID, " - ", it.ParInfo.Genre, " - ", it.ParInfo.TextPar)
			lastSentID = 0
		}
		if lastSentID != it.ParInfo.SentenceNum {
			fmt.Println("-------------------------------------------------")
			fmt.Println("  Sentença: ", it.ParInfo.SentenceNum, " - ", it.Sentence)
			phraseCount = 1
		}
		fmt.Println("      Oração:", phraseCount, " - ", it.PhraseType, " -->", it.PhraseText)

		lastParID = it.ParInfo.TextID
		lastSentID = it.ParInfo.SentenceNum
		phraseCount++
	}

	log.Println("////////////////////////////////////////////////////////////////////")

	mapStats := map[string]int{}
	mapStats["word_tot"] = 0
	mapStats["word_tot_senter"] = 0
	mapStats["tokens_tot"] = 0
	mapStats["phrase_tot"] = 0
	mapStats["sent_tot"] = 0
	mapStats["text_tot"] = 0
	mapStats["min_sent_words"] = 1000
	mapStats["max_sent_words"] = 0
	mapStats["tot_sent_words"] = 0
	mapStats["min_text_sents"] = 1000
	mapStats["max_text_sents"] = 0
	mapStats["tot_text_sents"] = 0
	mapStats["min_text_words"] = 1000
	mapStats["max_text_words"] = 0
	mapStats["tot_text_words"] = 0
	mapStats["min_word_size"] = 1000
	mapStats["max_word_size"] = 0
	mapStats["tot_word_size"] = 0

	lastParID = 0
	lastSentID = 0
	phraseCount = 0
	maskedSent := ""
	totTextWords := 0
	totTextSents := 0
	respList := []string{}
	for _, it := range phrases {
		// log.Println(lastParID, i, item.ParInfo.TextID, item.PhraseText)
		if lastParID != it.ParInfo.TextID {

			fmt.Println(" -> ", maskedSent)
			for _, r := range respList {
				fmt.Println("    - ", r)
			}
			fmt.Print("\n")
			respList = []string{}

			fmt.Println("\n================================================")
			fmt.Println("Texto: ", it.ParInfo.TextID, " - ", it.ParInfo.Genre, " - ", it.ParInfo.TextPar, "\n")
			lastSentID = 0

			mapStats["text_tot"]++

			if totTextWords > mapStats["max_text_words"] {
				mapStats["max_text_words"] = totTextWords
			}
			if totTextWords != 0 && totTextWords < mapStats["min_text_words"] {
				mapStats["min_text_words"] = totTextWords
			}
			mapStats["tot_text_words"] += totTextWords
			totTextWords = 0

			if totTextSents > mapStats["max_text_sents"] {
				mapStats["max_text_sents"] = totTextSents
			}
			if totTextSents != 0 && totTextSents < mapStats["min_text_sents"] {
				mapStats["min_text_sents"] = totTextSents
			}
			mapStats["tot_text_sents"] += totTextSents
			totTextSents = 0

		}
		if lastSentID != it.ParInfo.SentenceNum {

			if lastSentID != 0 {
				fmt.Println(" -> ", maskedSent)
				for _, r := range respList {
					fmt.Println("    - ", r)
				}
				fmt.Print("\n")
				respList = []string{}
			}

			// fmt.Println("-------------------------------------------------")
			// fmt.Println("  Sentença: ", it.ParInfo.SentenceNum, " - ", it.Sentence)
			phraseCount = 1
			maskedSent = it.Sentence

			mapStats["sent_tot"]++

			wordLocalIndex := 1
			for _, w := range it.SentData {
				auditWords2[w.UID] = 1
				// log.Println(wordLocalIndex, wordGeneralIndex, idxw, w.Word, w.UID)
				mapStats["word_tot"]++
				wordGeneralIndex++
				wordLocalIndex++
				wsize := len(w.Word)
				if wsize > mapStats["max_word_size"] {
					mapStats["max_word_size"] = wsize
				}
				if wsize < mapStats["min_word_size"] {
					mapStats["min_word_size"] = wsize
				}
				mapStats["tot_word_size"] += wsize
			}
			totSentWords := int(it.SentData[0].ParsedSent.QtyWords)
			mapStats["word_tot_senter"] += totSentWords
			mapStats["tokens_tot"] += int(it.SentData[0].ParsedSent.QtyTokens)

			if totSentWords > mapStats["max_sent_words"] {
				mapStats["max_sent_words"] = totSentWords
			}
			if totSentWords < mapStats["min_sent_words"] {
				mapStats["min_sent_words"] = totSentWords
			}
			mapStats["tot_sent_words"] += totSentWords
			totTextWords += totSentWords
			totTextSents++

			// log.Println("tokens -> ", it.SentData[0].ParsedSent.QtyTokens)
			// log.Println("words ->", it.SentData[0].ParsedSent.QtyWords, "=", wordLocalIndex-1)
		}
		fmt.Println("      Oração:", phraseCount, " - ", it.PhraseType, " -->", it.PhraseText)

		// fmt.Println("----->", it.SentData)

		w := getPhraseWordLessFrequent(it.SentData, it.PhraseText)

		if w.Word != "" {
			maskedSent = strings.ReplaceAll(maskedSent, w.Word, "_____")
			// fmt.Println(" --- ", maskedSent)
			resps := getRespList(w.WordClean, w.PoS, w.Top10Resp, map[string]int{})
			// fmt.Println("--->", w.WordClean, "(", w.PoS, ")", resps)
			respList = append(respList, fmt.Sprintf("%v (%v) %v", w.WordClean, w.PoS, resps))
			// fmt.Print("\n")
		}

		lastParID = it.ParInfo.TextID
		lastSentID = it.ParInfo.SentenceNum
		phraseCount++
		mapStats["phrase_tot"]++

	}

	if lastSentID != 0 {
		fmt.Println(" -> ", maskedSent)
		for _, r := range respList {
			fmt.Println("    - ", r)
		}
		fmt.Print("\n")
	}

	if totTextWords > mapStats["max_text_words"] {
		mapStats["max_text_words"] = totTextWords
	}
	if totTextWords != 0 && totTextWords < mapStats["min_text_words"] {
		mapStats["min_text_words"] = totTextWords
	}
	mapStats["tot_text_words"] += totTextWords

	if totTextSents > mapStats["max_text_sents"] {
		mapStats["max_text_sents"] = totTextSents
	}
	if totTextSents != 0 && totTextSents < mapStats["min_text_sents"] {
		mapStats["min_text_sents"] = totTextSents
	}
	mapStats["tot_text_sents"] += totTextSents

	log.Println("----- text_tot ->", mapStats["text_tot"])
	log.Println("----- sent_tot ->", mapStats["sent_tot"])
	log.Println("----- phrase_tot ->", mapStats["phrase_tot"])
	log.Println("----- word_tot ->", mapStats["word_tot"]+2)
	log.Println("----- tokens_tot ->", mapStats["tokens_tot"])
	log.Println("----- word_tot_senter ->", mapStats["word_tot_senter"])
	log.Println("----- min_word_len ->", mapStats["min_word_size"])
	log.Println("----- max_word_len ->", mapStats["max_word_size"])
	log.Println("----- avg_word_len ->", mapStats["tot_word_size"]/mapStats["word_tot"])
	log.Println("----- min_sent_words ->", mapStats["min_sent_words"])
	log.Println("----- max_sent_words ->", mapStats["max_sent_words"])
	log.Println("----- avg_sent_words ->", mapStats["tot_sent_words"]/mapStats["sent_tot"])
	log.Println("----- min_text_words ->", mapStats["min_text_words"])
	log.Println("----- max_text_words ->", mapStats["max_text_words"])
	log.Println("----- avg_text_words ->", mapStats["tot_text_words"]/mapStats["text_tot"])
	log.Println("----- min_text_sents ->", mapStats["min_text_sents"])
	log.Println("----- max_text_sents ->", mapStats["max_text_sents"])
	log.Println("----- avg_text_sents ->", mapStats["tot_text_sents"]/mapStats["text_tot"])

	for k := range auditWords1 {
		if m, found := auditWords2[k]; !found {
			log.Println("--", k, m)
		}
	}
	// // for k, v := range sentList {
	// // 	log.Println("\n---------------", k, v[0].TextSent)
	// // 	// for i, w := range v {
	// // 	// 	log.Println(i, w.UID, w.Word, w.SentenceNum, w.TextSent)
	// // 	// }
	// // 	w := getWordLessFrequent(v)
	// // 	log.Println("=====>", w.UID, w.WordInSent, w.WordInSent, w.Word, w.WordClean, w.PoS, w.FreqBrWaC)
	// // }

	// sentID := 1
	// for _, v := range sentList {
	// 	if len(v) < 2 {
	// 		continue
	// 	}
	// 	w := getWordLessFrequent(v)
	// 	maskedSent := strings.ReplaceAll(w.TextSent, w.Word, "_____")

	// 	fmt.Print("\n")
	// 	fmt.Println(sentID, " --- ", maskedSent)
	// 	resps := getRespList(w.WordClean, w.PoS, w.Top10Resp)
	// 	fmt.Println("--->", w.WordClean, "(", w.PoS, ")", resps)
	// 	sentID++
	// }

}

func getParInfo(sentList map[string][]Word, sent string) Word {
	for _, v := range sentList {
		for _, w := range v {
			if w.TextSent == sent {
				return w
			}
		}
	}
	return Word{}
}

func getRespList(resp string, respPoS string, respList []string, synonyms map[string]int) []string {
	ret := []string{}
	for _, item := range respList {
		tokens := strings.Split(item, ":")
		tResp := strings.TrimSpace(tokens[0])
		tResp = strings.Trim(tResp, "'")
		if tResp != resp && !isSynonym(tResp, synonyms) {
			if _, found := mapRespPoS[fmt.Sprintf("%v_%v", resp, tResp)]; found {
				ret = append(ret, tResp)
			}
		}
	}
	// if len(ret) > 4 {
	// 	ret = ret[0:4]
	// }
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

func getPhraseWordLessFrequent(wordList []Word, phrase string) Word {
	// log.Println("-------xx-----", phrase)

	discardedWords := map[string]int{"UID_26_40": 1, "UID_46_18": 1}
	discardedPoS := map[string]int{"NPROP": 1, "ERR": 1}

	var regEx = regexp.MustCompile(`[0-9,\.]+`)

	phraseWordList := strings.Split(phrase, " ")

	ret := Word{}
	ret.FreqBrWaC = 1000000
	defaultWord := Word{}
	lastPhraseWord := phraseWordList[len(phraseWordList)-1]
	for _, item := range wordList {
		if strings.ToLower(item.Word) == strings.ToLower(lastPhraseWord) {
			defaultWord = item
		}
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

		wordInPhrase := false
		for _, pw := range phraseWordList {
			if strings.ToLower(item.Word) == strings.ToLower(pw) {
				wordInPhrase = true
			}
		}
		if !wordInPhrase {
			continue
		}
		if item.FreqBrWaC < ret.FreqBrWaC {
			ret = item
		}
	}
	if ret.Word == "" {
		ret = defaultWord
	}
	return ret
}
