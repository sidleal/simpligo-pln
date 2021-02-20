package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/sidleal/simpligo-pln/tools/senter"
)

func main_adolesendo() {

	log.Println("starting ---")

	mapParts := map[string]string{}

	meta := readFile("/home/sidleal/sid/usp/adole-sendo/metadados.tsv")
	lines := strings.Split(meta, "\n")
	for i, line := range lines {
		tokens := strings.Split(line, "\t")
		if i == 0 {
			for j, t := range tokens {
				log.Println(j, t)
			}
			continue
		}
		partID := tokens[0]
		partYear := tokens[9]
		mapParts[fmt.Sprintf("VF_%v", partID)] = partYear
	}

	// for k, v := range mapParts {
	// 	log.Println(k, v)
	// }

	mapTotals := map[string]map[string]int{}

	raw := readFile("/home/sidleal/sid/usp/adole-sendo/versao_entrega_corrigida_sandra.txt")
	lines = strings.Split(raw, "\n")

	textID := ""
	for _, line := range lines {
		if line == "" {
			continue
		}
		if strings.Index(line, " ") < 1 {
			textID = strings.ReplaceAll(line, "\r", "")
			continue
		}

		text := line

		// log.Println(i, textID, text)

		year := strings.TrimSpace(mapParts[textID])

		if _, found := mapTotals[year]; !found {
			mapTotals[year] = map[string]int{}
			mapTotals[year]["sents"] = 0
		}

		parsed := senter.ParseText(text)

		mapTypes := map[string]int{}

		mapTotals[year]["texts"]++

		for _, p := range parsed.Paragraphs {
			for _, s := range p.Sentences {
				mapTotals[year]["sents"]++
				for _, t := range s.Tokens {
					if t.IsWord == 1 {
						w := strings.ToLower(t.Token)
						if _, found := mapTypes[w]; found {
							mapTypes[w]++
						} else {
							mapTypes[w] = 1
						}
						mapTotals[year]["tokens"]++
					}
				}
			}
		}
		mapTotals[year]["types"] += len(mapTypes)

		// log.Println("---------------------------")

	}

	log.Println("-------")

	for k, v := range mapTotals {
		for k2, v2 := range v {
			log.Println(k, k2, v2)
		}
		log.Println(k, "ttr", float32(v["types"])/float32(v["tokens"]))
		log.Println(k, "mean sents", float32(v["tokens"])/float32(v["sents"]))
		log.Println("-------------------------")
	}

}
