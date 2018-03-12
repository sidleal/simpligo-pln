package main

import (
	"io"
	"log"
	"os"
	"sort"
	"strings"
)

type Sentence struct {
	Producao  string
	Level     string
	Text      string
	Sentence  string
	Paragraph string
	RawText   string
}

type Align struct {
	Producao  string
	Level     string
	TextA     string
	SentenceA string
	TextB     string
	SentenceB string
}

//concat
func main() {

	rawSentences := readFile("/home/sidleal/usp/coling2018/v1/porsimples_sentences.txt")
	lines := strings.Split(rawSentences, "\n")

	sentences := []Sentence{}
	for _, line := range lines {
		if line == "" {
			break
		}
		tokens := strings.Split(line, "\t")
		// log.Println(tokens[0], tokens[1], tokens[2], tokens[3], tokens[4], tokens[5])
		sentence := Sentence{tokens[0], tokens[1], tokens[2], tokens[3], tokens[4], tokens[5]}
		sentences = append(sentences, sentence)

	}

	rawAligns := readFile("/home/sidleal/usp/coling2018/v1/porsimples_aligns.txt")
	lines = strings.Split(rawAligns, "\n")

	aligns := []Align{}
	for _, line := range lines {
		if line == "" {
			break
		}
		tokens := strings.Split(line, "\t")
		// log.Println(tokens[0], tokens[1], tokens[2], tokens[3], tokens[4], tokens[5])
		align := Align{tokens[0], tokens[1], tokens[2], tokens[3], tokens[4], tokens[5]}
		aligns = append(aligns, align)

	}

	// for i, item := range sentences {
	// 	log.Println(i, item.Level, item.Producao, item.Text, item.Sentence, item.RawText)
	// }

	// for i, item := range aligns {
	// 	log.Println(i, item.Level, item.Producao, item.TextA, item.SentenceA, item.TextB, item.SentenceB)
	// }

	f1, err := os.OpenFile("/home/sidleal/usp/coling2018/1_align_concat_nat.tsv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERRO", err)
	}
	defer f1.Close()

	oriTotalSentencas := 0
	oriUnchanged := 0
	oriChanged := 0
	oriDivided := 0
	oriUndivided := 0

	count := 1
	for _, item := range sentences {
		if item.Level == "NAT" && item.Producao == "67" {
			log.Println("-----------------------------------------------")

			// log.Println(count, item.Level, item.Producao, item.Text, item.Sentence, item.RawText)
			pairs := getPairs(aligns, sentences, item.Sentence)
			divided := ""
			pairStr := ""
			if len(pairs) > 1 {
				divided = "S"
				oriDivided++
			} else {
				divided = "N"
				oriUndivided++
			}

			for _, pair := range pairs {
				// log.Println(count, pair.Level, pair.Producao, pair.Text, pair.Sentence, pair.RawText)

				// if i > 0 && pair.RawText == sentences[i-1].RawText {
				// 	log.Println("-----------------------------------------------------------------DESALINHOU")
				// }

				// if i < len(sentences)-1 && pair.RawText == sentences[i+1].RawText {
				// 	log.Println("-----------------------------------------------------------------DESALINHOU")
				// }
				pairStr += pair.RawText + " "
			}
			count++

			item.RawText = strings.TrimSpace(item.RawText)
			pairStr = strings.TrimSpace(pairStr)

			if pairStr == "" {
				continue
			}

			oriTotalSentencas++

			log.Println(count, item.Level, item.Producao, item.Text, item.Sentence)
			log.Println(item.RawText)
			log.Println(pairStr)

			changed := ""
			if item.RawText == pairStr {
				changed = "N"
				oriUnchanged++
			} else {
				changed = "S"
				oriChanged++
			}
			level := "ORI->NAT"
			if item.Level == "NAT" {
				level = "NAT->STR"
			}

			line := item.Producao + "\t" + level + "\t" + changed + "\t" + divided + "\t" + item.RawText + "\t" + pairStr
			// _, err := f1.WriteString(line + "\n")
			// if err != nil {
			// 	log.Println("ERRO", err)
			// }

			log.Println(line)
		}
	}

	log.Println("------------------------------------")
	log.Println("Total sentenças : ", oriTotalSentencas)
	log.Println("Total sentenças Simplificadas: ", oriChanged)
	log.Println("Total sentenças NÃO Simplificadas: ", oriUnchanged)
	log.Println("Total sentenças Divididas: ", oriDivided)
	log.Println("Total sentenças NÃO Divididas: ", oriUndivided)

	log.Println("------------------------------------")

}

func getPairs(aligns []Align, sentences []Sentence, id string) []Sentence {
	ret := []Sentence{}
	for _, item := range aligns {
		if item.SentenceA == id {
			if ok, sentence := getSentence(sentences, item.SentenceB); ok {
				ret = append(ret, sentence)
			}
		}
	}
	return ret

}

func getSentence(sentences []Sentence, id string) (bool, Sentence) {
	for _, item := range sentences {
		if item.Sentence == id {
			return true, item
		}
	}
	return false, Sentence{}

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

//size
func main_size() {

	rawSentences := readFile("/home/sidleal/usp/coling2018/1_sentences.txt")
	lines := strings.Split(rawSentences, "\n")

	sentences := []Sentence{}
	for _, line := range lines {
		if line == "" {
			break
		}
		tokens := strings.Split(line, "\t")
		// log.Println(tokens[0], tokens[1], tokens[2], tokens[3], tokens[4], tokens[5])
		sentence := Sentence{tokens[0], tokens[1], tokens[2], tokens[3], tokens[4], tokens[5]}
		sentences = append(sentences, sentence)

	}

	rawAligns := readFile("/home/sidleal/usp/coling2018/1_aligns.txt")
	lines = strings.Split(rawAligns, "\n")

	aligns := []Align{}
	for _, line := range lines {
		if line == "" {
			break
		}
		tokens := strings.Split(line, "\t")
		// log.Println(tokens[0], tokens[1], tokens[2], tokens[3], tokens[4], tokens[5])
		align := Align{tokens[0], tokens[1], tokens[2], tokens[3], tokens[4], tokens[5]}
		aligns = append(aligns, align)

	}

	// for i, item := range sentences {
	// 	log.Println(i, item.Level, item.Producao, item.Text, item.Sentence, item.RawText)
	// }

	// for i, item := range aligns {
	// 	log.Println(i, item.Level, item.Producao, item.TextA, item.SentenceA, item.TextB, item.SentenceB)
	// }

	f1, err := os.OpenFile("/home/sidleal/usp/coling2018/1_align_size_nat.tsv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERRO", err)
	}
	defer f1.Close()

	oriTotalSentencas := 0
	oriUnchanged := 0
	oriChanged := 0
	oriDivided := 0
	oriUndivided := 0

	doubles := 0

	count := 1
	for _, item := range sentences {
		if item.Level == "NAT" { //&& item.Producao == "102" {
			log.Println("--------------------------------------------------------------------------------")

			// log.Println(count, item.Level, item.Producao, item.Text, item.Sentence, item.RawText)
			pairs := getPairs(aligns, sentences, item.Sentence)
			divided := ""
			if len(pairs) > 1 {
				divided = "S"
				oriDivided++
			} else {
				divided = "N"
				oriUndivided++
			}

			pairStrs := getBestPair(item.RawText, pairs)
			if len(pairStrs) > 1 {
				doubles++
			}
			for _, pairStr := range pairStrs {

				if pairStr == "" {
					continue
				}

				count++

				item.RawText = strings.TrimSpace(item.RawText)
				pairStr = strings.TrimSpace(pairStr)

				oriTotalSentencas++

				log.Println(count, item.Level, item.Producao, item.Text, item.Sentence)
				log.Println(item.RawText)
				log.Println(pairStr)

				changed := ""
				if item.RawText == pairStr {
					changed = "N"
					oriUnchanged++
				} else {
					changed = "S"
					oriChanged++
				}

				level := "ORI->NAT"
				if item.Level == "NAT" {
					level = "NAT->STR"
				}

				line := item.Producao + "\t" + level + "\t" + changed + "\t" + divided + "\t" + item.RawText + "\t" + pairStr
				_, err := f1.WriteString(line + "\n")
				if err != nil {
					log.Println("ERRO", err)
				}

				log.Println("-----")

				log.Println(line)
			}

		}
	}

	log.Println("------------------------------------")
	log.Println("Total PARES: ", oriTotalSentencas)
	log.Println("Sentenças repetidas (similaridades iguais): ", doubles)
	log.Println("------------------------------------")

}

type SentencePair struct {
	RawText    string
	Similarity int
}

type BySimilarity []SentencePair

func (a BySimilarity) Len() int           { return len(a) }
func (a BySimilarity) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySimilarity) Less(i, j int) bool { return a[i].Similarity > a[j].Similarity }

func getBestPair(sentA string, sentBList []Sentence) []string {
	pairs := []SentencePair{}
	for _, item := range sentBList {
		log.Println("---", item.Sentence, item.RawText)

		tokens := strings.Split(strings.TrimSpace(item.RawText), " ")
		tokensLen := len(tokens)
		commonTokens := getQtyCommonTokens(sentA, item.RawText)
		log.Println("------------------->", tokensLen, commonTokens, tokensLen+commonTokens)

		itemSimilarity := tokensLen + commonTokens
		pairs = append(pairs, SentencePair{item.RawText, itemSimilarity})
	}

	sort.Sort(BySimilarity(pairs))

	bestPairs := []string{}

	lastSim := 0
	for _, item := range pairs {
		if item.Similarity >= lastSim {
			bestPairs = append(bestPairs, item.RawText)
			lastSim = item.Similarity
		}
	}

	return bestPairs

}

func getQtyCommonTokens(a string, b string) int {
	tokensA := strings.Split(a, " ")
	tokensB := strings.Split(b, " ")
	ret := 0
	for _, tokenB := range tokensB {
		for _, tokenA := range tokensA {
			if tokenA == tokenB {
				ret++
				break
			}
		}
	}
	return ret
}
