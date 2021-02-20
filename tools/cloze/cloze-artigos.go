package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

// ------------------------------------

func main_palavras_err() {

	path := "/home/sidleal/sid/usp/cloze_exps4/final/"

	clozeFrom := path + "cloze_predict_38_FULL.tsv"
	clozeTo := path + "cloze_predict_39_FULL.tsv"

	f2, err := os.Create(clozeTo)
	if err != nil {
		log.Println("ERRO", err)
	}
	defer f2.Close()

	data := readFile(clozeFrom)

	artigos := map[string]int{
		"o":    1,
		"a":    1,
		"os":   1,
		"as":   1,
		"um":   1,
		"uns":  1,
		"uma":  1,
		"umas": 1,
	}

	lines := strings.Split(data, "\n")
	for j, line := range lines {
		if line == "" {
			continue
		}

		cols := strings.Split(line, "\t")

		wordCleaned := cols[9]
		respCleaned := cols[12]

		pos := cols[14]
		wordPos := cols[16]
		responseTag := cols[17]
		posMatch := cols[18]
		wordInflection := cols[19]
		responseInflection := cols[20]
		inflectionMatch := cols[21]

		if strings.HasPrefix(wordInflection, "DET") {
			if _, found := artigos[wordCleaned]; found {
				log.Println("a----")
				log.Println(wordCleaned, pos, wordPos, responseTag, wordInflection, responseInflection)
				log.Println(pos, "-->", "DET")
				log.Println(wordPos, "-->", "Artigo")
				log.Println(wordInflection, "-->", wordInflection[4:len(wordInflection)])
				pos = "DET"
				wordPos = "Artigo"
				wordInflection = wordInflection[4:len(wordInflection)]
			}
		}

		if strings.HasPrefix(responseInflection, "DET") {
			if _, found := artigos[respCleaned]; found {
				log.Println("b----")
				log.Println(wordCleaned, pos, wordPos, responseTag, wordInflection, responseInflection)
				log.Println(responseTag, "-->", "DET")
				log.Println(responseInflection, "-->", responseInflection[4:len(responseInflection)])
				responseTag = "DET"
				responseInflection = responseInflection[4:len(responseInflection)]
			}
		}

		if wordInflection == responseInflection {
			inflectionMatch = "1"
		} else {
			inflectionMatch = "0"
		}
		if pos == responseTag {
			posMatch = "1"
		} else {
			posMatch = "0"
		}

		stts := cols[30]
		sttt := cols[31]
		stt := cols[32]
		if j > 0 {
			tts, _ := strconv.Atoi(cols[30])
			ttt, _ := strconv.Atoi(cols[31])
			if tts > 1000000 || ttt > 1000000 {
				stts, sttt, stt = "0", "0", "0"
			}
		}

		line := ""
		for i, val := range cols {
			newVal := val
			if i == 14 {
				newVal = pos
			}
			if i == 16 {
				newVal = wordPos
			}
			if i == 17 {
				newVal = responseTag
			}
			if i == 18 && j > 0 {
				newVal = posMatch
			}
			if i == 19 {
				newVal = wordInflection
			}
			if i == 20 {
				newVal = responseInflection
			}
			if i == 21 && j > 0 {
				newVal = inflectionMatch
			}
			if i == 30 {
				newVal = stts
			}
			if i == 31 {
				newVal = sttt
			}
			if i == 32 {
				newVal = stt
			}
			line += newVal + "\t"
		}

		log.Println(line)
		line = strings.TrimSuffix(line, "\t") + "\n"
		_, err = f2.WriteString(line)
		if err != nil {
			log.Println("ERRO", err)
		}

	}

}
