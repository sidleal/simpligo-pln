package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	senter "github.com/sidleal/simpligo-pln/tools/senter"
)

// ------------------------------------

var freqBRWaC = map[string]int{}

func main3() {

	path := "/home/sidleal/sid/usp/cloze_exps/"

	f, err := os.Create(path + "dataset_v0_09xxxxxx.txt")
	if err != nil {
		log.Println("ERRO", err)
	}
	defer f.Close()

	listBRWaC := readFile(path + "lista_brWaC_geral_v3_nlpnet.tsv")
	lines := strings.Split(listBRWaC, "\n")
	for i, line := range lines {

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
		freq, _ := strconv.Atoi(cols[1])
		freqBRWaC[fmt.Sprintf("%v_%v", cols[0], cols[2])] = freq

	}

	data := readFile(path + "cloze_predict23b.tsv")

	sentList := map[string][]Word{}

	lines = strings.Split(data, "\n")
	lastSent := ""
	for i, line := range lines {

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

	}

	wordGeneralIndex := 1

	phrases := []Phrase{}

	dataPhrases := readFile(path + "oracoes2.tsv")
	lines = strings.Split(dataPhrases, "\n")
	for j, line := range lines {

		if line == "" {
			continue
		}
		cols := strings.Split(line, "\t")

		par := getParInfo(sentList, cols[0])

		phrase := Phrase{}
		phrase.PhraseText = cols[2]
		phrase.PhraseType = cols[1]
		phrase.Sentence = cols[0]
		phrase.ParInfo = par
		phrase.PhraseNum = j

		thisSent := fmt.Sprintf("%v_%v", par.TextID, par.SentenceNum)
		phrase.SentData = sentList[thisSent]
		phrases = append(phrases, phrase)

	}

	sort.Sort(PhraseOrder(phrases))

	lastParID := 0
	lastSentID := 0
	phraseCount := 0
	for _, it := range phrases {
		if lastParID != it.ParInfo.TextID {

			f.WriteString("\n================================================\n")
			f.WriteString(fmt.Sprintf("Texto: %v - %v - %v\n", it.ParInfo.TextID, it.ParInfo.Genre, it.ParInfo.TextPar))
			lastSentID = 0
		}
		if lastSentID != it.ParInfo.SentenceNum {
			f.WriteString("-------------------------------------------------\n")
			f.WriteString(fmt.Sprintf("  Sentença: %v - %v\n", it.ParInfo.SentenceNum, it.Sentence))
			phraseCount = 1
		}
		f.WriteString(fmt.Sprintf("      Oração: %v - %v --> %v\n", phraseCount, it.PhraseType, it.PhraseText))

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

	totTextWords := 0
	totTextSents := 0
	respList := []string{}
	for _, it := range phrases {

		if lastParID != it.ParInfo.TextID {

			// if it.ParInfo.TextID < 41 {
			// 	continue
			// }

			respList = []string{}

			f.WriteString("\n================================================\n")
			f.WriteString(fmt.Sprintf("Texto: %v - %v - %v\n\n", it.ParInfo.TextID, it.ParInfo.Genre, it.ParInfo.TextPar))

			log.Println(fmt.Sprintf("Texto: %v - %v - %v\n\n", it.ParInfo.TextID, it.ParInfo.Genre, it.ParInfo.TextPar))

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

			f.WriteString("\n====================\n")
			f.WriteString(fmt.Sprintf("\n Sentença -> %v \n\n", it.Sentence))

			phraseCount = 1

			mapStats["sent_tot"]++

			wordLocalIndex := 1
			for _, w := range it.SentData {
				auditWords2[w.UID] = 1
				mapStats["word_tot"]++
				wordGeneralIndex++
				wordLocalIndex++
				runas := []rune(w.Word)
				wsize := len(runas)
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

		}
		f.WriteString("\n-----------------------------------------------\n\n")
		f.WriteString(fmt.Sprintf("      Oração: %v - %v --> %v\n\n", phraseCount, it.PhraseType, it.PhraseText))

		candidateList := []WordCandidate{}
		phraseTokens := strings.Split(it.PhraseText, " ")
		phraseLen := len(phraseTokens)
		for i := phraseLen/2 - 1; i < phraseLen; i++ {
			token := phraseTokens[i]
			wo := getPhraseWord(it.SentData, token)
			if wo.Word != "" && wo.WordType == "Content" {

				synonyms := map[string]int{}
				// synonyms := callSinonimos(wo.WordClean)

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
			f.WriteString(fmt.Sprintf("    - %v (%v) Frq: %v Match: %v Cert %v -- %v\n", wc.Word, wc.PoS, wc.FreqBrWaC, wc.OrtoMatch, wc.Certainty, wc.Resps))
			if wc.Word != "" && false {
				alts := getSynAlternatives(wc.Resps, wc.PoS)
				f.WriteString("          alts:\n")
				for k, v := range alts {
					f.WriteString(fmt.Sprintf("          %v ( %v )\n\n", k, v))
				}
			}
		}

		fmt.Print("\n")
		respList = []string{}

		lastParID = it.ParInfo.TextID
		lastSentID = it.ParInfo.SentenceNum
		phraseCount++
		mapStats["phrase_tot"]++

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

	// syns := getSynAlternatives([]string{"divisei", "sorri", "cavalgou", "balançou"}, "V")
	// syns := getSynAlternatives([]string{"cavalgou"}, "V")
	// log.Println(syns)
}

func getSynAlternatives(words []string, tag string) map[string][]string {

	tag = strings.ToLower(tag)
	ret := map[string][]string{}
	cMap := map[string][]WordCandidate{}

	for _, w := range words {
		cMap[w] = []WordCandidate{}
		syns := callSinonimos(w)
		// log.Println(syns)
		for k := range syns {
			if strings.Index(k, " ") > 0 {
				continue
			}
			freq := freqBRWaC[fmt.Sprintf("%v_%v", k, tag)]
			c := WordCandidate{}
			c.Word = k
			c.FreqBrWaC = freq
			cMap[w] = append(cMap[w], c)
		}

		// time.Sleep(500 * time.Millisecond)
	}
	for k, v := range cMap {
		ret[k] = []string{}
		sort.Sort(WordCandidateOrder(v))
		for _, c := range v {
			ret[k] = append(ret[k], fmt.Sprintf("%v (%v), ", c.Word, c.FreqBrWaC))
		}

	}

	// if len(ret) > 30 {
	// 	ret = ret[:30]
	// }

	return ret
}
