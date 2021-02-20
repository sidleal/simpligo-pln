package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// ------------------------------------

func main_freqs() {
	path := "/home/sidleal/sid/usp/cloze_exps/"
	data := readFile(path + "cloze_joao_20200926.tsv")

	f, err := os.Create(path + "cloze_joao_20200926_with_resp_freqs.tsv")
	if err != nil {
		log.Println("ERRO", err)
	}
	defer f.Close()

	mapFreq := map[string]int{}
	freqs := readFile(path + "freq_nathan.csv")
	lines := strings.Split(freqs, "\n")
	for i, line := range lines {
		if line == "" {
			continue
		}
		cols := strings.Split(line, ",")
		if i == 0 {
			for j, col := range cols {
				log.Println(j, col)
			}

			continue
			// log.Println(cols)
		}

		mapFreq[cols[0]], _ = strconv.Atoi(cols[20])

	}

	lines = strings.Split(data, "\n")
	for i, line := range lines {
		if line == "" {
			continue
		}
		// log.Println(i, line)

		cols := strings.Split(line, "\t")
		if i == 0 {
			for j, col := range cols {
				log.Println(j, col)
			}

			_, err = f.WriteString(fmt.Sprintf("%v\tresp_freq_br\n", line))
			if err != nil {
				log.Println("ERRO", err)
			}

			continue
			// log.Println(cols)
		}

		resp := cols[11]

		log.Println(resp, mapFreq[resp])

		newLine := fmt.Sprintf("%v\t%v\n", line, mapFreq[resp])
		_, err = f.WriteString(newLine)
		if err != nil {
			log.Println("ERRO", err)
		}

	}
}
