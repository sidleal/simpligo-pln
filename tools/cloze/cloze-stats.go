package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	senter "github.com/sidleal/simpligo-pln/tools/senter"
)

// ------------------------------------

func mainz() {

	path := "/home/sidleal/sid/usp/cloze_exps/"

	data := readFile(path + "50pars.txt")

	f, err := os.Create(fmt.Sprintf(path + "stats.tsv"))
	if err != nil {
		log.Println("ERRO", err)
	}
	defer f.Close()

	header := "id\ttext\ttotal_words\ttotal_tokens\ttotal_chars\tmin_word\tmax_word\tavg_word\n"
	_, err = f.WriteString(header)
	if err != nil {
		log.Println("ERRO", err)
	}

	uniqueWords := map[string]int{}

	totPalavrasGeral := 0
	totCharsPalavrasGeral := 0
	lines := strings.Split(data, "\n")
	for i, line := range lines {
		if line == "" || i == 0 {
			continue
		}
		log.Println(i, line)

		sent := senter.ParseText(line)

		log.Println("Tot words:", sent.TotalWords)
		log.Println("Tot Tokens:", sent.TotalTokens)
		log.Println("caracteres por par:", len(line))
		totalCharsWord := 0
		maxWord := 0
		minWord := 1000
		for _, s := range sent.Paragraphs[0].Sentences {
			for _, w := range s.Tokens {
				if w.IsWord == 1 {
					totalCharsWord += len(w.Token)
					if maxWord < len(w.Token) {
						maxWord = len(w.Token)
					}
					if minWord > len(w.Token) {
						minWord = len(w.Token)
					}
					totPalavrasGeral++
					if _, found := uniqueWords[strings.ToLower(w.Token)]; !found {
						uniqueWords[strings.ToLower(w.Token)] = 1
					} else {
						uniqueWords[strings.ToLower(w.Token)]++
					}
				}
			}
		}
		log.Println("max palavra:", maxWord)
		log.Println("min palavra:", minWord)
		log.Println("avg palavra:", totalCharsWord/int(sent.TotalWords))
		totCharsPalavrasGeral += totalCharsWord

		newLine := fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", i, line, sent.TotalWords, sent.TotalTokens, len(line), maxWord, minWord, totalCharsWord/int(sent.TotalWords))
		_, err = f.WriteString(newLine)
		if err != nil {
			log.Println("ERRO", err)
		}

	}

	log.Println("Total geral palavras:", totPalavrasGeral)
	log.Println("Avg geral palavras:", totCharsPalavrasGeral/totPalavrasGeral)
	log.Println("Total palavras únicas:", len(uniqueWords))

	for k, v := range uniqueWords {
		log.Println(k, v)
	}

}
