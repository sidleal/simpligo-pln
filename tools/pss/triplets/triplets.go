package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Pair struct {
	Production string
	TextA      string
	TextB      string
}

func main() {

	// f1, err := os.OpenFile("/home/sidleal/usp/coling2018/v2/validacao_ori.tsv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	log.Println("ERRO", err)
	// }
	// defer f1.Close()

	// f2, err := os.OpenFile("/home/sidleal/usp/coling2018/v2/validacao_nat.tsv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	log.Println("ERRO", err)
	// }
	// defer f2.Close()

	// f3, err := os.OpenFile("/home/sidleal/usp/coling2018/v2/validacao_str.tsv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	log.Println("ERRO", err)
	// }
	// defer f3.Close()

	f4, err := os.OpenFile("/home/sidleal/usp/coling2018/v3/triplets.tsv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERRO", err)
	}
	defer f4.Close()

	sizeSentOri := readFile("/home/sidleal/usp/coling2018/v3/align_size_ori_nat.tsv")
	lines := strings.Split(sizeSentOri, "\n")

	oriNatPairs := []Pair{}
	for _, line := range lines {
		if line == "" {
			continue
		}
		tokens := strings.Split(line, "\t")
		pair := Pair{}
		pair.Production = tokens[0]
		pair.TextA = tokens[4]
		pair.TextB = tokens[5]
		oriNatPairs = append(oriNatPairs, pair)

	}

	sizeSentNat := readFile("/home/sidleal/usp/coling2018/v3/align_size_nat_str.tsv")
	lines = strings.Split(sizeSentNat, "\n")

	natStrPairs := []Pair{}
	for _, line := range lines {
		if line == "" {
			continue
		}
		tokens := strings.Split(line, "\t")
		pair := Pair{}
		pair.Production = tokens[0]
		pair.TextA = tokens[4]
		pair.TextB = tokens[5]
		natStrPairs = append(natStrPairs, pair)

	}

	// metrics := readFile("/home/sidleal/usp/coling2018/v3/pss_sentences_metrics.tsv")
	// lines = strings.Split(metrics, "\n")

	// metricList := map[string][]string{}
	// for _, line := range lines {
	// 	if line == "" {
	// 		continue
	// 	}
	// 	tokens := strings.Split(line, "\t")
	// 	metricList[tokens[0]] = tokens
	// }

	_, err = f4.WriteString("production\tlevel\tchanged_ori_nat\tchanged_nat_str\toriginal\tnatural\tstrong\n")
	if err != nil {
		log.Println("ERRO", err)
	}

	for _, itemNat := range natStrPairs {
		for _, itemOri := range oriNatPairs {
			if itemNat.TextA == itemOri.TextB {

				oriNatChanged := itemOri.TextA != itemOri.TextB
				natStrChanged := itemNat.TextA != itemNat.TextB

				log.Println("----------------")
				log.Println(itemOri.TextA)
				log.Println(itemOri.TextB)
				log.Println(itemNat.TextB)

				line := itemOri.Production + "\t" + "ORI->NAT->STR\t"

				if oriNatChanged {
					line += "S\t"
				} else {
					line += "N\t"
				}

				if natStrChanged {
					line += "S\t"
				} else {
					line += "N\t"
				}

				line += itemOri.TextA + "\t"
				line += itemOri.TextB + "\t"
				line += itemNat.TextB + "\n"

				_, err := f4.WriteString(line)
				if err != nil {
					log.Println("ERRO", err)
				}

			}
		}
	}

	// for _, item1 := range oriNatPairs {
	// 	log.Println("----------------")
	// 	log.Println(item1.TextA)
	// 	log.Println(item1.TextB)

	// 	if item1.TextA != item1.TextB {
	// 		line := ""
	// 		for i, token := range metricList[item1.TextA] {
	// 			if i > 0 {
	// 				token = round(token)
	// 			}
	// 			line += token + "\t"
	// 		}
	// 		line = strings.TrimSuffix(line, "\t")
	// 		line += "\n"

	// 		// _, err := f1.WriteString(line)
	// 		// if err != nil {
	// 		// 	log.Println("ERRO", err)
	// 		// }

	// 	}

	// }

	// for _, item1 := range natStrPairs {
	// 	log.Println("----------------")
	// 	log.Println(item1.TextA)
	// 	log.Println(item1.TextB)

	// 	if item1.TextA != item1.TextB {
	// 		line := ""
	// 		for i, token := range metricList[item1.TextA] {
	// 			if i > 0 {
	// 				token = round(token)
	// 			}
	// 			line += token + "\t"
	// 		}
	// 		line = strings.TrimSuffix(line, "\t")
	// 		line += "\n"

	// 		// _, err := f2.WriteString(line)
	// 		// if err != nil {
	// 		// 	log.Println("ERRO", err)
	// 		// }

	// 		line = ""
	// 		for i, token := range metricList[item1.TextB] {
	// 			if i > 0 {
	// 				token = round(token)
	// 			}
	// 			line += token + "\t"
	// 		}
	// 		line = strings.TrimSuffix(line, "\t")
	// 		line += "\n"

	// 		// _, err = f3.WriteString(line)
	// 		// if err != nil {
	// 		// 	log.Println("ERRO", err)
	// 		// }

	// 	}

	// }

	//---------

	// for _, item1 := range oriNatPairs {
	// 	for _, item2 := range oriNatPairs {
	// 		if itemConcatOri.TextA == itemOri.TextA && itemConcatOri.TextA != itemConcatOri.TextB && itemConcatOri.TextB != itemOri.TextB {
	// 			// log.Println("----------------")
	// 			// log.Println(itemConcatOri.TextA)
	// 			// log.Println(itemConcatOri.TextB)
	// 			// log.Println(itemOri.TextB)
	// 			// log.Println(metricList[itemOri.TextB])

	// 			line := ""
	// 			for _, token := range metricList[itemConcatOri.TextA] {
	// 				line += token + "\t"
	// 			}
	// 			line = strings.TrimSuffix(line, "\t")
	// 			line += "\n"

	// 			for _, token := range metricList[itemConcatOri.TextB] {
	// 				line += token + "\t"
	// 			}
	// 			line = strings.TrimSuffix(line, "\t")
	// 			line += "\n"

	// 			for _, token := range metricList[itemOri.TextB] {
	// 				line += token + "\t"
	// 			}
	// 			line = strings.TrimSuffix(line, "\t")
	// 			line += "\n"

	// 			if itemConcatOri.TextB == "Há um total de 81 senadores." {
	// 				log.Println(line)
	// 			}

	// 			// _, err := f1.WriteString(line + "\n")
	// 			// if err != nil {
	// 			// 	log.Println("ERRO", err)
	// 			// }

	// 		}

	// 	}
	// }

}

func round(num string) string {
	ret := floatToStr(tofloat(num))
	ret = strings.Replace(ret, ".", ",", -1)
	return ret
}

func tofloat(str string) float64 {
	if str == "" {
		str = "0"
	}
	ret, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Println("ERROOOOOOOOOOOOO", err)
	}
	return ret
}

func floatToStr(n float64) string {
	return fmt.Sprintf("%.3f", n)
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
