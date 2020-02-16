package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type RegraOut struct {
	textID  string
	snippet string
	regra   string
}

func mainy() {

	log.Println("starting")

	f, err := os.OpenFile("/home/sidleal/sid/usp/TopicosPLN/regra/resultado_regra_parsed2.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERRO", err)
	}
	defer f.Close()

	totTokenMap := map[string]int{}
	totErrorsMap := map[string]int{}

	textMap := map[string]string{}
	buffer := readFile("/home/sidleal/sid/usp/TopicosPLN/regra/all_questions_regra.txt")
	lines := strings.Split(buffer, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		tokens := strings.Split(line, "\t")
		textMap[tokens[0]] = tokens[1]

		ttok := strings.Split(tokens[1], " ")
		totTokenMap[tokens[0]] = len(ttok)
	}

	mapOutput := map[string][]RegraOut{}

	regraOutput := readFile("/home/sidleal/sid/usp/TopicosPLN/regra/resultado_regra_utf.txt")

	var regEx = regexp.MustCompile(`(?s)\[S:([0-9]+)\].*?<b>(.*?)</b>.*?\[REGRA\](.*?)\[/REGRA\]`)
	matchContent := regEx.FindAllStringSubmatch(regraOutput, -1)

	for _, match := range matchContent {
		if len(match) < 4 {
			log.Println("Formato ruim", match)
			continue
		}
		regraOut := RegraOut{}
		regraOut.textID = match[1]
		regraOut.snippet = match[2]
		regraOut.regra = match[3]

		if len(regraOut.snippet) > 20 {
			regraOut.snippet = regraOut.snippet[0:20]
		}

		mapOutput[regraOut.textID] = append(mapOutput[regraOut.textID], regraOut)
	}

	countMap := map[string]int{}
	totalErros := 0

	for k, v := range mapOutput {
		log.Println("--------------")
		log.Println(k, textMap[k])
		_, err := f.WriteString(fmt.Sprintf("--------------\n%s - %s\n", k, textMap[k]))
		if err != nil {
			log.Println("ERRO", err)
		}
		errorsCount := 0
		for _, it := range v {
			log.Println(it.regra, "----> ", it.snippet)
			_, err := f.WriteString(fmt.Sprintf("%s --> %s\n", it.regra, it.snippet))
			totalErros++
			errorsCount++
			countMap[it.regra]++
			if err != nil {
				log.Println("ERRO", err)
			}
		}
		totErrorsMap[k] = errorsCount
	}

	log.Println("Total erros:", totalErros)
	for k, v := range countMap {
		log.Println(k, ":", v)
	}

	for k, _ := range mapOutput {
		// log.Println(k, textMap[k])
		// log.Println("Tokens:", totTokenMap[k], "Erros:", totErrorsMap[k], "Ratio:", float32(totErrorsMap[k])/float32(totTokenMap[k]))
		log.Println(k, totTokenMap[k], totErrorsMap[k], float32(totErrorsMap[k])/float32(totTokenMap[k]))
	}

}
