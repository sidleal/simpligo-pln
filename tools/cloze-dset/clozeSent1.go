package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
)

// ------------------------------------

func main_sents() {

	path := "/home/sidleal/sid/usp/cloze_exps/sent_compl/"

	f, err := os.Create(path + "dataset_v1.tsv")
	if err != nil {
		log.Println("ERRO", err)
	}
	defer f.Close()

	f.WriteString("id\tquestion\ta\tb\tc\td\te\n")

	f2, err := os.Create(path + "dataset_v1_answers.tsv")
	if err != nil {
		log.Println("ERRO", err)
	}
	defer f2.Close()

	f2.WriteString("id\tanswer\n")

	r := rand.New(rand.NewSource(2497))

	mapOpts := map[int]string{0: "a", 1: "b", 2: "c", 3: "d", 4: "e"}
	data := readFile(path + "sents_v0.txt")
	lines := strings.Split(data, "\n")
	for _, line := range lines {

		if line == "" {
			continue
		}

		cols := strings.Split(line, "|")
		textcols := strings.Split(cols[0], " -> ")

		sentID := textcols[0]
		text := textcols[1]
		log.Println(sentID)
		log.Println(text)
		resps := strings.Split(cols[1], " ")

		log.Println(resps)

		nuns := r.Perm(5)
		log.Println(nuns)

		f.WriteString(fmt.Sprintf("%v\t%v", sentID, text))
		for i, n := range nuns {
			log.Println(mapOpts[i], resps[n], n == 0)
			f.WriteString(fmt.Sprintf("\t%v", resps[n]))
			if n == 0 {
				f2.WriteString(fmt.Sprintf("%v\t%v\n", sentID, mapOpts[i]))
			}
		}
		f.WriteString("\n")

	}

	// f.WriteString("\n================================================\n")
}
