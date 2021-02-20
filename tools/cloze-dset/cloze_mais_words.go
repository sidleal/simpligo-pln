package main

import (
	"log"
	"strconv"
	"strings"

	senter "github.com/sidleal/simpligo-pln/tools/senter"
)

// ------------------------------------

func main() {

	path := "/home/sidleal/sid/usp/cloze_exps/"

	data := readFile(path + "cloze_predict23b.tsv")

	wordList := map[string]Word{}

	textList := map[int][]Word{}

	lastText := 0

	lines := strings.Split(data, "\n")
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

		wordList[word.UID] = word

		if word.TextID != lastText {
			textList[word.TextID] = []Word{}
		}
		lastText = word.TextID

		textList[word.TextID] = append(textList[word.TextID], word)

	}

	repWords := map[string]int{}
	for k, v := range wordList {
		log.Println(k, v.TextID, v.WordClean)
		repWords[v.WordClean] = 0
	}

	for rw := range repWords {
		for k, v := range textList {
			repeated := false
			log.Println("Text", k)
			for _, w := range v {
				if w.WordClean == rw {
					repeated = true
				}
			}
			if repeated {
				repWords[rw]++
			}
		}

	}

	for k, v := range repWords {
		log.Println(k, "\t", v)
	}
	log.Println(len(wordList))
}
