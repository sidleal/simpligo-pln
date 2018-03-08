package main

import (
	"io"
	"log"
	"os"
	"strings"
)

func main() {

	raw := readFile("/home/sidleal/usp/coling2018/align-12.txt")

	lines := strings.Split(raw, "\n")

	orinatFrom := [][6]string{}
	orinatTo := [][6]string{}
	natstrFrom := [][6]string{}
	natstrTo := [][6]string{}

	for _, line := range lines {
		if line == "" {
			break
		}
		tokens := strings.Split(line, "\t")
		log.Println(tokens[0], tokens[1], tokens[2], tokens[3], tokens[4])
		from := [6]string{}
		to := [6]string{}

		if tokens[2] == "ORI->NAT" {
			for j := 0; j < 5; j++ {
				from[j] = tokens[j]
				to[j] = tokens[j]
			}
			from[5] = tokens[5]
			to[5] = tokens[6]
			orinatFrom = append(orinatFrom, from)
			orinatTo = append(orinatTo, to)
		} else if tokens[2] == "NAT->STR" {

			for j := 0; j < 5; j++ {
				from[j] = tokens[j]
				to[j] = tokens[j]
			}
			from[5] = tokens[5]
			to[5] = tokens[6]
			natstrFrom = append(natstrFrom, from)
			natstrTo = append(natstrTo, to)
		}

	}

	countNSimpl := 0
	for i, _ := range natstrTo {
		// log.Println(item[0], item[1], item[2], item[3], item[4])

		if natstrFrom[i][5] == natstrTo[i][5] {
			ok, sentenceOri := findSentenceOriginal(orinatTo, orinatFrom, natstrTo[i])
			if ok && sentenceOri[5] == natstrTo[i][5] {
				countNSimpl++
				log.Println(sentenceOri[5])
				// log.Println(natstrFrom[i][5])
				// log.Println(natstrTo[i][5])
			}
		}
	}

	log.Println("Nao simplificadas em nenhum nivel: ", countNSimpl)

	log.Println("==================================================================")

	sentListSimpNat := map[string]int{}
	for i, item := range orinatFrom {
		// log.Println(item[0], item[1], item[2], item[3], item[4])

		if item[5] != orinatTo[i][5] {
			sentListSimpNat[item[5]] = 1
		}
	}

	f1, err := os.OpenFile("/home/sidleal/usp/coling2018/sentListSimpNat.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERRO", err)
	}
	defer f1.Close()

	countSimplNat := 0
	for k, _ := range sentListSimpNat {
		countSimplNat++
		log.Println(k)

		_, err := f1.WriteString(k + "\n")
		if err != nil {
			log.Println("ERRO", err)
		}

	}
	log.Println("Simplificadas nivel natural: ", countSimplNat)
	log.Println("==================================================================")

	sentListSimpStr := map[string]int{}
	for i, _ := range natstrTo {
		// log.Println(item[0], item[1], item[2], item[3], item[4])

		if natstrFrom[i][5] != natstrTo[i][5] {
			ok, sentenceOri := findSentenceOriginal(orinatTo, orinatFrom, natstrFrom[i])
			if ok && sentenceOri[5] == natstrFrom[i][5] {
				sentListSimpStr[sentenceOri[5]] = 1
			}
		}
	}

	countSimplStr := 0
	for k, _ := range sentListSimpStr {
		countSimplStr++
		log.Println(k)
		// log.Println(natstrFrom[i][5])
		// log.Println(natstrTo[i][5])

	}

	log.Println("Simplificadas apenas do nat p/ str: ", countSimplStr)

	// for i, item := range natstrFrom {
	// 	log.Println(item[0], item[1], item[2], item[3], item[4])

	// 	ok, sentenceOri := findSentence(orinatTo, orinatFrom, natstrFrom[i])
	// 	if ok {
	// 		log.Println("ACHOU!!!!", natstrFrom[i][5], sentenceOri[5])

	// 	}
	// }

	// f, err := os.OpenFile("/home/sidleal/usp/coling2018/align-13.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	log.Println("ERRO", err)
	// }

	// defer f.Close()

	// f2, err := os.OpenFile("/home/sidleal/usp/coling2018/align3s-13.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	log.Println("ERRO", err)
	// }

	// defer f2.Close()

	// f3, err := os.OpenFile("/home/sidleal/usp/coling2018/align3-13.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	log.Println("ERRO", err)
	// }

	// defer f3.Close()

	// // ORI->NAT
	// for i, item := range orinatFrom {
	// 	log.Println(item[0], item[1], item[2], item[3])

	// 	line := ""
	// 	for _, field := range orinatFrom[i] {
	// 		line += field + "\t"
	// 	}
	// 	for j, field := range orinatTo[i] {
	// 		if j > 2 {
	// 			line += field + "\t"
	// 		}
	// 	}
	// 	line = strings.TrimSuffix(line, "\t") + "\n"

	// 	n, err := f.WriteString(line)
	// 	if err != nil {
	// 		log.Println("ERRO", err)
	// 	}

	// 	fmt.Printf("wrote %d bytes\n", n)

	// }

	// // NAT->STR
	// for i, item := range natstrFrom {
	// 	log.Println(item[0], item[1], item[2], item[3])

	// 	line := ""
	// 	for _, field := range natstrFrom[i] {
	// 		line += field + "\t"
	// 	}
	// 	for j, field := range natstrTo[i] {
	// 		if j > 2 {
	// 			line += field + "\t"
	// 		}
	// 	}
	// 	line = strings.TrimSuffix(line, "\t") + "\n"

	// 	n, err := f.WriteString(line)
	// 	if err != nil {
	// 		log.Println("ERRO", err)
	// 	}

	// 	fmt.Printf("wrote %d bytes\n", n)

	// }

	// log.Println("-------------------------------------------------------------")

	// // ORI->STR
	// for i, item := range natstrFrom {
	// 	log.Println(item[0], item[1], item[2], item[3])

	// 	ok, sentenceOri := findSentence(orinatTo, orinatFrom, natstrFrom[i])
	// 	if ok {
	// 		log.Println("ACHOU!!!!")

	// 		sentenceOri[1] = "ORI->STR"

	// 		line := ""
	// 		for _, field := range sentenceOri {
	// 			line += field + "\t"
	// 		}
	// 		for j, field := range natstrTo[i] {
	// 			if j > 2 {
	// 				line += field + "\t"
	// 			}
	// 		}
	// 		line = strings.TrimSuffix(line, "\t") + "\n"

	// 		n, err := f.WriteString(line)
	// 		if err != nil {
	// 			log.Println("ERRO", err)
	// 		}

	// 		fmt.Printf("wrote %d bytes\n", n)

	// 		// 3 niveis
	// 		sentenceOri[1] = "ORI->NAT->STR"

	// 		line = ""
	// 		for _, field := range sentenceOri {
	// 			line += field + "\t"
	// 		}
	// 		for j, field := range natstrFrom[i] {
	// 			if j > 2 {
	// 				line += field + "\t"
	// 			}
	// 		}
	// 		for j, field := range natstrTo[i] {
	// 			if j > 2 {
	// 				line += field + "\t"
	// 			}
	// 		}
	// 		line = strings.TrimSuffix(line, "\t") + "\n"

	// 		n3, err := f3.WriteString(line)
	// 		if err != nil {
	// 			log.Println("ERRO", err)
	// 		}

	// 		fmt.Printf("wrote %d bytes\n", n3)

	// 		line2 := sentenceOri[3] + "\t" + natstrFrom[i][3] + "\t" + natstrTo[i][3] + "\n"
	// 		n2, err := f2.WriteString(line2)
	// 		if err != nil {
	// 			log.Println("ERRO", err)
	// 		}

	// 		fmt.Printf("wrote %d bytes\n", n2)

	// 	}

	// }

}

func findSentenceOriginal(orinatTo [][6]string, orinatFrom [][6]string, nat [6]string) (bool, [6]string) {

	for i, item := range orinatTo {
		if item[1] == nat[1] && item[5] == nat[5] {
			return true, orinatFrom[i]
		}
	}

	return false, [6]string{}

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
