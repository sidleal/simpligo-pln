package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	senter "github.com/sidleal/simpligo-pln/tools/senter"
)

// ------------------------------------

func main_quatro() {

	path := "/home/sidleal/sid/usp/cloze_exps/"

	f, err := os.Create(path + "dataset_v0_13c.txt")
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
	for _, it := range phrases {
		if lastParID != it.ParInfo.TextID {

			f.WriteString("\n=========================================================================\n")
			f.WriteString(fmt.Sprintf("Texto: %v - %v - %v\n", it.ParInfo.TextID, it.ParInfo.Genre, it.ParInfo.TextPar))
			lastSentID = 0
		}
		if lastSentID != it.ParInfo.SentenceNum {
			f.WriteString("---------------------------------------------------------------------------\n")
			f.WriteString(fmt.Sprintf("  Sentença: %v - %v\n", it.ParInfo.SentenceNum, it.Sentence))

		}

		lastParID = it.ParInfo.TextID
		lastSentID = it.ParInfo.SentenceNum
	}

	log.Println("////////////////////////////////////////////////////////////////////")

	lastParID = 0
	lastSentID = 0

	for _, it := range phrases {

		if lastParID != it.ParInfo.TextID {

			if it.ParInfo.TextID < 40 {
				continue
			}

			f.WriteString("\n===================================================================================\n")
			f.WriteString(fmt.Sprintf("Texto: %v - %v - %v\n\n", it.ParInfo.TextID, it.ParInfo.Genre, it.ParInfo.TextPar))

			log.Println(fmt.Sprintf("Texto: %v - %v - %v\n\n", it.ParInfo.TextID, it.ParInfo.Genre, it.ParInfo.TextPar))

			lastSentID = 0

			time.Sleep(2 * time.Second)

		}
		if lastSentID != it.ParInfo.SentenceNum {

			f.WriteString("---------------------------------------------------------------------------\n")
			f.WriteString(fmt.Sprintf("\n Sentença -> %v \n\n", it.Sentence))
			log.Println(fmt.Sprintf("\n Sentença -> %v \n\n", it.Sentence))

			candidateList := []WordCandidate{}

			// sentTokens := strings.Split(it.Sentence, " ")
			// sentLen := len(sentTokens)

			parsedSent := it.SentData[0].ParsedSent
			sentLen := len(parsedSent.Tokens)
			loopStart := sentLen/2 - 3
			if sentLen <= 8 {
				loopStart = 0
			}
			recWords := map[string]int{}
			for i := loopStart; i < sentLen; i++ {

				if parsedSent.Tokens[i].IsWord == 0 {
					continue
				}

				token := parsedSent.Tokens[i].Token
				if _, found := recWords[token]; !found {
					recWords[token] = 0
				} else {
					recWords[token]++
				}
				wo := getPhraseWordDisamb(it.SentData, token, recWords[token])
				if wo.Word != "" && wo.WordType == "Content" {

					// synonyms := map[string]int{}
					synonyms := callSinonimos(wo.WordClean)

					resps := getRespList(wo.WordClean, wo.PoS, wo.Top10Resp, synonyms)
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
				wordStr := wc.Word + strings.Repeat(" ", 25)[0:25-len(wc.Word)]
				f.WriteString(fmt.Sprintf("    %v\t%v\t- M %.2f\t- C %.2f\t- L %v\n", wordStr, wc.PoS, wc.OrtoMatch, wc.Certainty, len(wc.Resps)))
				// log.Print(fmt.Sprintf("    %v\t%v\t- M %.2f\t- C %.2f\t- L %v\n", wordStr, wc.PoS, wc.OrtoMatch, wc.Certainty, len(wc.Resps)))
			}

			for _, wc := range candidateList {
				f.WriteString(fmt.Sprintf("\n    ----> %v (%v) alts (Frq: %v - Match: %v - Cert %v):\n\n", wc.Word, wc.PoS, wc.FreqBrWaC, wc.OrtoMatch, wc.Certainty))
				if wc.Word != "" {
					alts := getSynAlternatives(wc.Resps, wc.PoS)
					for k, v := range alts {
						f.WriteString(fmt.Sprintf("         - %v ( %v )\n\n", k, v))
					}
				}
			}

		}

		fmt.Print("\n")

		lastParID = it.ParInfo.TextID
		lastSentID = it.ParInfo.SentenceNum
	}

	// syns := getSynAlternatives([]string{"divisei", "sorri", "cavalgou", "balançou"}, "V")
	// syns := getSynAlternatives([]string{"cavalgou"}, "V")
	// log.Println(syns)

}

func getPhraseWordDisamb(wordList []Word, word string, recWordsPos int) Word {
	ret := Word{}
	wordPos := 0
	for _, item := range wordList {
		if strings.ToLower(item.Word) == strings.ToLower(word) {
			log.Println("---------", wordPos, recWordsPos, word, item.Certainty)
			if wordPos == recWordsPos {
				ret = item
				break
			} else {
				wordPos++
			}
		}
	}
	return ret
}
