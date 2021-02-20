package main

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/sidleal/simpligo-pln/tools/senter"
)

func main() {

	log.Println("starting")

	// path := "/home/sidleal/sid/usp/murilo/MuriloD3"
	path := "/home/sidleal/sid/usp/murilo/MuriloD1"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println("Erro", err)
	}

	// first := true
	totalTexts := 0
	totalSents := 0
	totalTokens := 0
	totalTypes := 0
	avgWordsPerSent := 0.0
	for _, f := range files {
		fileName := f.Name()
		log.Println("---------------------------------------------------------------")
		log.Println(fileName)

		files2, err := ioutil.ReadDir(path + "/" + fileName)
		if err != nil {
			log.Println("Erro", err)
		}
		for _, f2 := range files2 {
			fileName2 := f2.Name()

			if !strings.HasSuffix(fileName2, ".txt") {
				continue
			}

			raw := readFile(path + "/" + fileName + "/" + fileName2)

			// log.Println(raw)

			parsed := senter.ParseText(raw)

			totalTexts++

			mapTypes := map[string]int{}

			for _, p := range parsed.Paragraphs {
				for _, s := range p.Sentences {
					totalSents++
					for _, t := range s.Tokens {
						if t.IsWord == 1 {
							w := strings.ToLower(t.Token)
							if _, found := mapTypes[w]; found {
								mapTypes[w]++
							} else {
								mapTypes[w] = 1
							}
							totalTokens++
						}
					}
				}
			}
			avgWordsPerSent += float64(parsed.TotalWords) / float64(parsed.TotalSentences)

			totalTypes += len(mapTypes)
		}
	}

	log.Println("Total textos:", totalTexts)
	log.Println("Total sentenças:", totalSents)
	log.Println("Total types:", totalTypes)
	log.Println("Total tokens:", totalTokens)
	log.Println("Type Token Ration:", float32(totalTypes)/float32(totalTokens))
	log.Println("Avg words per sent:", avgWordsPerSent/float64(totalTexts))

}
