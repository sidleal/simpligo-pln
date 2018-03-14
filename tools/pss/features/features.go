package main

import (
	"io"
	"log"
	"os"
	"strings"
)

func main() {

	bigMap := map[string]int{}

	f1, err := os.OpenFile("/home/sidleal/usp/coling2018/v2/unique_sentences.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERRO", err)
	}
	defer f1.Close()

	concatSent := readFile("/home/sidleal/usp/coling2018/v2/align_concat_ori.tsv")
	lines := strings.Split(concatSent, "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}
		tokens := strings.Split(line, "\t")

		bigMap[tokens[4]] = 1
		bigMap[tokens[5]] = 1

	}

	allSent := readFile("/home/sidleal/usp/coling2018/v2/align_all_ori.tsv")
	lines = strings.Split(allSent, "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}
		tokens := strings.Split(line, "\t")

		bigMap[tokens[4]] = 1
		bigMap[tokens[5]] = 1

	}

	sizeSent := readFile("/home/sidleal/usp/coling2018/v2/align_size_ori.tsv")
	lines = strings.Split(sizeSent, "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}
		tokens := strings.Split(line, "\t")
		bigMap[tokens[4]] = 1
		bigMap[tokens[5]] = 1

	}

	concatNatSent := readFile("/home/sidleal/usp/coling2018/v2/align_concat_nat.tsv")
	lines = strings.Split(concatNatSent, "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}
		tokens := strings.Split(line, "\t")

		bigMap[tokens[4]] = 1
		bigMap[tokens[5]] = 1

	}

	allNatSent := readFile("/home/sidleal/usp/coling2018/v2/align_all_nat.tsv")
	lines = strings.Split(allNatSent, "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}
		tokens := strings.Split(line, "\t")

		bigMap[tokens[4]] = 1
		bigMap[tokens[5]] = 1

	}

	sizeNatSent := readFile("/home/sidleal/usp/coling2018/v2/align_size_nat.tsv")
	lines = strings.Split(sizeNatSent, "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}
		tokens := strings.Split(line, "\t")
		bigMap[tokens[4]] = 1
		bigMap[tokens[5]] = 1

	}

	count := 0
	for k, _ := range bigMap {
		count++
		log.Println(count, k)

		_, err := f1.WriteString(k + "\n")
		if err != nil {
			log.Println("ERRO", err)
		}
	}
}

func readFile(path string) string {
	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}

	// r := charmap.ISO8859_1.NewDecoder().Reader(f)
	r := io.Reader(f)

	ret := ""

	buf := make([]byte, 32*1024)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			ret += string(buf[:n])
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			break
		}
	}

	return ret

}
