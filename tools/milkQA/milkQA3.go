package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
)

func mainx3() {

	log.Println("starting")

	mapTexts := map[int]string{}
	selected := map[string]string{}

	anno1 := 0
	anno2 := 0
	anno3 := 0

	path := "/home/sidleal/sid/usp/TopicosPLN/PerguntasMilkQA"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println("Erro", err)
	}

	fout, err := os.OpenFile(path+"/all_questions.txt", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERRO", err)
	}
	defer fout.Close()

	i := 0
	for _, f := range files {
		fileName := f.Name()
		if strings.HasSuffix(fileName, "txt") && fileName != "all_questions.txt" {
			// log.Println("---------------------------------------------------------------")
			// log.Println(fileName)
			tokens := strings.Split(fileName, "-")
			fileID := tokens[0]
			// raw := readFile(path + "/" + fileName)
			// log.Println(raw)

			mapTexts[i] = fileID
			i++
		}
	}

	// for code, _ := range mapTexts {
	// 	log.Println(code)
	// }

	for i := 0; i < 160; i++ {

		idx := rand.Intn(2657)

		sel := mapTexts[idx]
		if _, found := selected[sel]; found {
			i--
			continue
		}

		if anno1 < 50 {
			selected[sel] = "anno1"
			anno1++
		} else if anno2 < 50 {
			selected[sel] = "anno2"
			anno2++
		} else if anno3 < 50 {
			selected[sel] = "anno3"
			anno3++
		} else {
			selected[sel] = "all"
		}

		log.Println(sel, selected[sel])

	}

}
