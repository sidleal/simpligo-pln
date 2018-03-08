package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {

	raw := readFile("/home/sidleal/usp/coling2018/PorSimplesSent_v1.tsv")

	lines := strings.Split(raw, "\n")

	orinatFrom := [][229]string{}
	orinatTo := [][229]string{}
	natstrFrom := [][229]string{}
	natstrTo := [][229]string{}

	for _, line := range lines {
		if line == "" {
			break
		}
		tokens := strings.Split(line, "\t")
		// log.Println(tokens[0], tokens[1], tokens[2], tokens[3], tokens[229])
		from := [229]string{}
		to := [229]string{}

		if tokens[1] == "ORI->NAT" {
			for j := 0; j < 3; j++ {
				from[j] = tokens[j]
				to[j] = tokens[j]
			}
			for j := 3; j < 229; j++ {
				from[j] = tokens[j]
			}
			for j := 229; j < 455; j++ {
				to[j-226] = tokens[j]
			}
			orinatFrom = append(orinatFrom, from)
			orinatTo = append(orinatTo, to)
		} else if tokens[1] == "NAT->STR" {
			for j := 0; j < 3; j++ {
				from[j] = tokens[j]
				to[j] = tokens[j]
			}
			for j := 3; j < 229; j++ {
				from[j] = tokens[j]
			}
			for j := 229; j < 455; j++ {
				to[j-226] = tokens[j]
			}
			natstrFrom = append(natstrFrom, from)
			natstrTo = append(natstrTo, to)
		}

	}

	f, err := os.OpenFile("/home/sidleal/usp/coling2018/align-13.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERRO", err)
	}

	defer f.Close()

	f2, err := os.OpenFile("/home/sidleal/usp/coling2018/align3s-13.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERRO", err)
	}

	defer f2.Close()

	f3, err := os.OpenFile("/home/sidleal/usp/coling2018/align3-13.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERRO", err)
	}

	defer f3.Close()

	// ORI->NAT
	for i, item := range orinatFrom {
		log.Println(item[0], item[1], item[2], item[3])

		line := ""
		for _, field := range orinatFrom[i] {
			line += field + "\t"
		}
		for j, field := range orinatTo[i] {
			if j > 2 {
				line += field + "\t"
			}
		}
		line = strings.TrimSuffix(line, "\t") + "\n"

		n, err := f.WriteString(line)
		if err != nil {
			log.Println("ERRO", err)
		}

		fmt.Printf("wrote %d bytes\n", n)

	}

	// NAT->STR
	for i, item := range natstrFrom {
		log.Println(item[0], item[1], item[2], item[3])

		line := ""
		for _, field := range natstrFrom[i] {
			line += field + "\t"
		}
		for j, field := range natstrTo[i] {
			if j > 2 {
				line += field + "\t"
			}
		}
		line = strings.TrimSuffix(line, "\t") + "\n"

		n, err := f.WriteString(line)
		if err != nil {
			log.Println("ERRO", err)
		}

		fmt.Printf("wrote %d bytes\n", n)

	}

	log.Println("-------------------------------------------------------------")

	// ORI->STR
	for i, item := range natstrFrom {
		log.Println(item[0], item[1], item[2], item[3])

		ok, sentenceOri := findSentence(orinatTo, orinatFrom, natstrFrom[i])
		if ok {
			log.Println("ACHOU!!!!")

			sentenceOri[1] = "ORI->STR"

			line := ""
			for _, field := range sentenceOri {
				line += field + "\t"
			}
			for j, field := range natstrTo[i] {
				if j > 2 {
					line += field + "\t"
				}
			}
			line = strings.TrimSuffix(line, "\t") + "\n"

			n, err := f.WriteString(line)
			if err != nil {
				log.Println("ERRO", err)
			}

			fmt.Printf("wrote %d bytes\n", n)

			// 3 niveis
			sentenceOri[1] = "ORI->NAT->STR"

			line = ""
			for _, field := range sentenceOri {
				line += field + "\t"
			}
			for j, field := range natstrFrom[i] {
				if j > 2 {
					line += field + "\t"
				}
			}
			for j, field := range natstrTo[i] {
				if j > 2 {
					line += field + "\t"
				}
			}
			line = strings.TrimSuffix(line, "\t") + "\n"

			n3, err := f3.WriteString(line)
			if err != nil {
				log.Println("ERRO", err)
			}

			fmt.Printf("wrote %d bytes\n", n3)

			line2 := sentenceOri[3] + "\t" + natstrFrom[i][3] + "\t" + natstrTo[i][3] + "\n"
			n2, err := f2.WriteString(line2)
			if err != nil {
				log.Println("ERRO", err)
			}

			fmt.Printf("wrote %d bytes\n", n2)

		}

	}

}

func findSentence(orinatTo [][229]string, orinatFrom [][229]string, nat [229]string) (bool, [229]string) {

	for i, item := range orinatTo {
		if item[0] == nat[0] && item[3] == nat[3] {
			return true, orinatFrom[i]
		}
	}

	return false, [229]string{}

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
