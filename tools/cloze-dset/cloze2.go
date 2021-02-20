package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	senter "github.com/sidleal/simpligo-pln/tools/senter"
	"golang.org/x/text/encoding/charmap"
)

// ------------------------------------

type WordCandidate struct {
	Word      string
	PoS       string
	FreqBrWaC int
	OrtoMatch float64
	Certainty float64
	Resps     []string
}

type WordCandidateOrder []WordCandidate

func (a WordCandidateOrder) Len() int      { return len(a) }
func (a WordCandidateOrder) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a WordCandidateOrder) Less(i, j int) bool {
	return a[i].FreqBrWaC < a[j].FreqBrWaC
}

func main_2() {

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
		ortoStr := strings.ReplaceAll(cols[11], ",", ".")
		ortoStr = (ortoStr + "000")[0:4]
		word.OrtoMatch, _ = strconv.ParseFloat(ortoStr, 64)
		certStr := strings.ReplaceAll(cols[16], ",", ".")
		certStr = (certStr + "000")[0:4]
		word.Certainty, _ = strconv.ParseFloat(certStr, 64)
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
	// maskedSent := ""
	totTextWords := 0
	totTextSents := 0
	respList := []string{}
	for _, it := range phrases {
		// log.Println(lastParID, i, item.ParInfo.TextID, item.PhraseText)
		if lastParID != it.ParInfo.TextID {

			if it.ParInfo.TextID < 41 {
				continue
			}

			// fmt.Println(" -> ", maskedSent)
			// for _, r := range respList {
			// 	fmt.Println("    - ", r)
			// }
			// fmt.Print("\n")
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

			fmt.Println("\n -> ", it.Sentence, "\n")

			// if lastSentID != 0 {
			// 	fmt.Println(" -> ", maskedSent)
			// 	for _, r := range respList {
			// 		fmt.Println("    - ", r)
			// 	}
			// 	fmt.Print("\n")
			// 	respList = []string{}
			// }

			// fmt.Println("-------------------------------------------------")
			// fmt.Println("  Sentença: ", it.ParInfo.SentenceNum, " - ", it.Sentence)
			phraseCount = 1
			// maskedSent = it.Sentence

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

		candidateList := []WordCandidate{}
		phraseTokens := strings.Split(it.PhraseText, " ")
		phraseLen := len(phraseTokens)
		for i := phraseLen/2 - 1; i < phraseLen; i++ {
			token := phraseTokens[i]
			wo := getPhraseWord(it.SentData, token)
			if wo.Word != "" && wo.WordType == "Content" {
				// log.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxx", token, wo.FreqBrWaC, wo.OrtoMatch, wo.Certainty)

				synonyms := map[string]int{}
				// synonyms := callSinonimos(wo.WordClean)
				// log.Println(synonyms)

				resps := getRespList(wo.WordClean, wo.PoS, wo.Top10Resp, synonyms)
				respList = append(respList, fmt.Sprintf("%v (%v) Frq: %v Match: %v Cert %v  -- %v", wo.WordClean, wo.PoS, wo.FreqBrWaC, wo.OrtoMatch, wo.Certainty, resps))
				candidate := WordCandidate{}
				candidate.Word = wo.WordClean
				candidate.PoS = wo.PoS
				candidate.FreqBrWaC = wo.FreqBrWaC
				candidate.OrtoMatch = wo.OrtoMatch
				candidate.Certainty = wo.Certainty
				candidate.Resps = resps
				candidateList = append(candidateList, candidate)

				// time.Sleep(2 * time.Second)
			}
		}

		sort.Sort(WordCandidateOrder(candidateList))

		for _, wc := range candidateList {
			fmt.Println("    - ", fmt.Sprintf("%v (%v) Frq: %v Match: %v Cert %v -- %v", wc.Word, wc.PoS, wc.FreqBrWaC, wc.OrtoMatch, wc.Certainty, wc.Resps))
		}

		// for _, r := range respList {
		// 	fmt.Println("    - ", r)
		// }
		fmt.Print("\n")
		respList = []string{}

		// w := getPhraseWordLessFrequent(it.SentData, it.PhraseText)

		// if w.Word != "" {
		// 	maskedSent = strings.ReplaceAll(maskedSent, w.Word, "_____")
		// 	// fmt.Println(" --- ", maskedSent)
		// 	resps := getRespList(w.WordClean, w.PoS, w.Top10Resp)
		// 	// fmt.Println("--->", w.WordClean, "(", w.PoS, ")", resps)
		// 	respList = append(respList, fmt.Sprintf("%v (%v) %v", w.WordClean, w.PoS, resps))
		// 	// fmt.Print("\n")
		// }

		lastParID = it.ParInfo.TextID
		lastSentID = it.ParInfo.SentenceNum
		phraseCount++
		mapStats["phrase_tot"]++

	}

	// if lastSentID != 0 {
	// 	fmt.Println(" -> ", maskedSent)
	// 	for _, r := range respList {
	// 		fmt.Println("    - ", r)
	// 	}
	// 	fmt.Print("\n")
	// }

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

	// for k := range auditWords1 {
	// 	if m, found := auditWords2[k]; !found {
	// 		log.Println("--", k, m)
	// 	}
	// }
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

	// synonyms := callSinonimos("divisei")
	// synonyms := callSinonimos("espécie")
	// log.Println(synonyms)

}

func isSynonym(word string, list map[string]int) bool {
	if _, found := list[word]; found {
		return true
	}
	return false
}

func getPhraseWord(wordList []Word, word string) Word {
	ret := Word{}
	for _, item := range wordList {
		if strings.ToLower(item.Word) == strings.ToLower(word) {
			ret = item
		}
	}
	return ret
}

func callSinonimos(word string) map[string]int {

	ret := map[string]int{}

	// return ret
	time.Sleep(800 * time.Millisecond)

	timeout := time.Duration(300 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	word, err := charmap.ISO8859_1.NewEncoder().String(word)
	if err != nil {
		log.Println("Erro utf to iso", err)
	}
	// log.Println("----", word)
	resp, err := client.Get("https://www.sinonimos.com.br/" + word + "/")
	if err != nil {
		log.Println(fmt.Sprintf("Error: %v", err))
	}

	rdrBody := io.Reader(resp.Body)
	rdrBody = charmap.ISO8859_1.NewDecoder().Reader(rdrBody)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(rdrBody)
	// body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(fmt.Sprintf("Error: %v", err))
	}

	bodyStr := string(body)

	// log.Println(bodyStr)

	regEx := regexp.MustCompile(`class="sinonimos"\>(.*?)\<\/p\>`)
	regEx2 := regexp.MustCompile(`sinonimo"\>(.*?)\<\/a\>`)
	regEx3 := regexp.MustCompile(`\<span\>(.*?)\<\/span\>`)

	matches := regEx.FindAllStringSubmatch(bodyStr, -1)
	for _, match := range matches {
		// log.Println("-------", match[1])
		matches2 := regEx2.FindAllStringSubmatch(match[1], -1)
		for _, match2 := range matches2 {
			ret[strings.TrimSpace(match2[1])] = 1
		}
		matches3 := regEx3.FindAllStringSubmatch(match[1], -1)
		for _, match3 := range matches3 {
			ret[strings.TrimSpace(match3[1])] = 1
		}

		// tokens := strings.Split(match[1], ",")
		// for _, w := range tokens {
		// 	ret[strings.TrimSpace(w)] = 1
		// }
	}

	return ret
}

func callSinonimos1(word string) map[string]int {

	ret := map[string]int{}

	timeout := time.Duration(300 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get("https://www.xn--sinnimo-v0a.com/" + word + ".html")
	if err != nil {
		log.Println(fmt.Sprintf("Error: %v", err))
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(fmt.Sprintf("Error: %v", err))
	}

	bodyStr := string(body)

	// log.Println(bodyStr)

	regEx := regexp.MustCompile(`\<div class="item"\>\<p.*?\>(.*?)\<span`)
	matches := regEx.FindAllStringSubmatch(bodyStr, -1)
	for _, match := range matches {
		// log.Println("-------", match[1])
		tokens := strings.Split(match[1], ",")
		for _, w := range tokens {
			ret[strings.TrimSpace(w)] = 1
		}
	}

	return ret
}
