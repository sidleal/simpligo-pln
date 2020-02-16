package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {

	log.Println("starting")

	carregarDELAF()

	path := "/home/sidleal/sid/usp/TopicosPLN/final/MilkQAE_curated_documents_2019-11-26/TSV3"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println("Erro", err)
	}

	fout, err := os.OpenFile(path+"/MilkQAE2.tsv", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERRO", err)
	}
	defer fout.Close()

	totFiles := 0
	totTokens := 0
	mapUniqueWords := map[string]int{}
	mapTotalErrors := map[string]int{}
	mapTextsTokens := map[string]int{}
	mapTextsErrors := map[string]int{}

	for _, f := range files {
		if f.Name() == "MilkQAE1x.tsv" {
			continue
		}
		fileName := f.Name() + "/CURATION_USER.tsv"
		nameTks := strings.Split(f.Name(), "-")
		textID := nameTks[0]
		// log.Println("---------------------------------------------------------------")
		// log.Println(textID, fileName)

		content := readFile(path + "/" + fileName)

		lines := strings.Split(content, "\n")
		cleanedLines := []string{}
		for _, line := range lines {
			if strings.HasPrefix(line, "#") || line == "" {
				continue
			}
			cleanedLines = append(cleanedLines, line)
		}

		for lineNum, line := range cleanedLines {
			if strings.HasPrefix(line, "#") || line == "" {
				continue
			}

			word, tag := getWordAndTag(line)
			if word == "" {
				continue
			}

			if _, found := mapUniqueWords[word]; !found {
				mapUniqueWords[word] = 0
			}
			mapUniqueWords[word]++

			errorsTks := strings.Split(tag, "|")
			for _, subtag := range errorsTks {
				subtagTks := strings.Split(subtag, "[")
				tagName := subtagTks[0]
				if _, found := mapTotalErrors[tagName]; !found {
					mapTotalErrors[tagName] = 0
				}
				mapTotalErrors[tagName]++

				inDELAF := 0
				delafDesc := ""
				if desc, found := delaf[word]; found {
					inDELAF = 1
					delafDesc = desc
				}

				log.Println(textID, word, tagName, inDELAF, delafDesc)
				// log.Println(textID, word, tagName)
				bag := getBagOfChars(word)
				bagPrev := [64]int{}
				bagNext := [64]int{}

				if lineNum > 0 {
					pWord, _ := getWordAndTag(cleanedLines[lineNum-1])
					// log.Println("prev:", pWord, pTag)
					if pWord != "" {
						bagPrev = getBagOfChars(pWord)
					}
				}
				if lineNum < len(cleanedLines)-1 {
					nWord, _ := getWordAndTag(cleanedLines[lineNum+1])
					// log.Println("next:", nWord, nTag)
					if nWord != "" {
						bagNext = getBagOfChars(nWord)
					}
				}

				log.Println(bag, bagPrev, bagNext)

				out := fmt.Sprintf("%v\t%v\t%v\t", textID, word, inDELAF)
				for _, it := range bag {
					out += fmt.Sprintf("%v\t", it)
				}
				for _, it := range bagPrev {
					out += fmt.Sprintf("%v\t", it)
				}
				for _, it := range bagNext {
					out += fmt.Sprintf("%v\t", it)
				}

				posArray := getCharPosArray(word)
				for _, it := range posArray {
					out += fmt.Sprintf("%v\t", it)
				}

				out = fmt.Sprintf("%v\t%v\n", out, tagName)
				log.Println("Gravando:", out)
				_, err := fout.WriteString(out)
				if err != nil {
					log.Println("ERRO", err)
				}

				if _, found := mapTextsTokens[textID]; !found {
					mapTextsTokens[textID] = 0
				}
				mapTextsTokens[textID]++
				if tagName != "OK" {
					if _, found := mapTextsErrors[textID]; !found {
						mapTextsErrors[textID] = 0
					}
					mapTextsErrors[textID]++
				}
			}
			totTokens++

		}

		totFiles++
	}

	log.Println("---------------------------------------------------------------")
	log.Println("Total textos: ", totFiles)
	log.Println("Total tokens: ", totTokens)
	log.Println("Total types: ", len(mapUniqueWords))
	for k, v := range mapTotalErrors {
		log.Println(k, v)
	}
	for k, v := range mapTextsTokens {
		errors := 0
		if nerrs, found := mapTextsErrors[k]; found {
			errors = nerrs
		}
		log.Println(fmt.Sprintf("%v\t%v\t%v", k, v, errors))
	}
	log.Println("---------------------------------------------------------------")

}

func getWordAndTag(line string) (string, string) {
	tsvTks := strings.Split(line, "\t")
	word := strings.ToLower(tsvTks[2])

	word = strings.ReplaceAll(word, "²", "2")
	word = strings.ReplaceAll(word, "ª", "a")

	tag := tsvTks[3]
	if tag == "_" || tag == "" {
		tag = "OK"
	}

	if word == "assunto" || word == "corpo" || isEscapeChar(word) || isNumber(word) {
		word, tag = "", ""
	}

	return word, tag
}

var mapEscapeChars = map[string]int{">": 1, "<": 1, ",": 1, "?": 1, "(": 1, ")": 1, "!": 1, "/": 1, ".": 1, "*": 1, "+": 1, "-": 1, ":": 1, "'": 1, "\\": 1, ";": 1, "\"": 1, "_": 1, "%": 1, "[": 1, "]": 1, "´": 1}

func isEscapeChar(char string) bool {
	if _, found := mapEscapeChars[char]; found {
		return true
	}
	return false
}

var regExNum = regexp.MustCompile(`^[$0-9]+.*`)

func isNumber(token string) bool {
	match := regExNum.FindAllStringSubmatch(token, -1)
	return len(match) > 0
}

var delaf = map[string]string{}

func carregarDELAF() {
	log.Println("carregando delaf")
	regEx1 := regexp.MustCompile(`(.+)\.(.+)`)

	path := "/home/sidleal/sid/usp/TopicosPLN/DELAF_PB_2018.dic"
	lines := readFileLines(path)
	for _, line := range lines {
		if line == "" {
			continue
		}
		cols := strings.Split(line, ",")

		info := cols[1]
		lema := ""
		tag := ""
		match := regEx1.FindStringSubmatch(info)
		if len(match) > 1 {
			lema = match[1]
			tag = match[2]
			info = lema + "|" + tag
		}

		if delafDesc, found := delaf[cols[0]]; found {
			tokens := strings.Split(delafDesc, "|")
			lemaB := tokens[0]
			tagB := tokens[1]

			pesoA := getTagWeight(tag)
			pesoB := getTagWeight(tagB)

			if pesoA < pesoB {
				info = lemaB + "|" + tagB + "," + tag
			} else {
				info = lema + "|" + tag + "," + tagB
			}
			// log.Println(info)
		}
		delaf[cols[0]] = info
	}

	log.Println("delaf em memoria")
}

func readFileLines(path string) []string {
	ret := []string{}

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ret = append(ret, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return ret

}

func getTagWeight(tag string) int {
	//Artigo > Preposição > Pronome > Numeral > Substantivo > Adjetivo > Verbo > Advérbio

	ret := 0
	if strings.HasPrefix(tag, "ADV") {
		ret = 1
	} else if strings.HasPrefix(tag, "V") {
		ret = 2
	} else if strings.HasPrefix(tag, "A") {
		ret = 3
	} else if strings.HasPrefix(tag, "N") {
		ret = 4
	} else if strings.HasPrefix(tag, "DET+Num") {
		ret = 5
	} else if strings.HasPrefix(tag, "PRO") {
		ret = 6
	} else if strings.HasPrefix(tag, "PREPXPRO") {
		ret = 8
	} else if strings.HasPrefix(tag, "PREPXDET") {
		ret = 9
	} else if strings.HasPrefix(tag, "PREP") {
		ret = 7
	} else if strings.HasPrefix(tag, "DET+Art") {
		ret = 10
	}
	return ret
}

var chars = "abcdefghijklmnopqrstuvwxyz0123456789áéíóúâêôãõàçñ-."

func getBagOfChars(word string) [64]int {
	ret := [64]int{}
	for _, char := range word {
		idx := strings.IndexRune(chars, char)
		if idx < 0 {
			log.Println("ERRO: ---------------------------> ", word, string(char))
			panic("oops")
		}
		// log.Println(string(char), idx)
		ret[idx]++
	}

	return ret
}

func getCharPosArray(word string) [64]int {
	ret := [64]int{}
	for i, char := range word {
		idx := strings.IndexRune(chars, char)
		if idx < 0 {
			log.Println("ERRO: ---------------------------> ", word, string(char))
			panic("oops")
		}
		// log.Println(string(char), idx)
		ret[i] = idx
	}

	return ret
}
