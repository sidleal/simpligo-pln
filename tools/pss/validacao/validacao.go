package main

import (
	"io"
	"log"
	"os"
	"strings"
)

type Sentence struct {
	Producao string
	Level    string
	Text     string
	Splited  string
	Changed  string
}

func main_val1() {

	concatSent := readFile("/home/sidleal/usp/coling2018/1_align_concat_ori.tsv")
	lines := strings.Split(concatSent, "\n")

	concatSentences := []Sentence{}
	for _, line := range lines {
		if line == "" {
			continue
		}
		tokens := strings.Split(line, "\t")
		sentence := Sentence{}
		sentence.Producao = tokens[0]
		sentence.Level = tokens[1]
		sentence.Changed = tokens[2]
		sentence.Splited = tokens[3]
		sentence.Text = tokens[4]
		concatSentences = append(concatSentences, sentence)

	}

	sizeSent := readFile("/home/sidleal/usp/coling2018/1_align_size_ori.tsv")
	lines = strings.Split(sizeSent, "\n")

	sizeSentences := []Sentence{}
	for _, line := range lines {
		if line == "" {
			continue
		}
		tokens := strings.Split(line, "\t")
		sentence := Sentence{}
		sentence.Producao = tokens[0]
		sentence.Level = tokens[1]
		sentence.Changed = tokens[2]
		sentence.Splited = tokens[3]
		sentence.Text = tokens[4]
		sizeSentences = append(sizeSentences, sentence)

	}

	log.Println(len(concatSentences), len(sizeSentences), len(sizeSentences)-len(concatSentences))

	for _, item := range sizeSentences {
		achou := false
		for _, itemC := range concatSentences {
			if itemC.Text == item.Text {
				achou = true
				break
			}
		}
		if !achou {
			log.Println("-----------------------", item.Text)
		}
	}

	lastSent := ""
	count := 0
	for _, item := range sizeSentences {
		if lastSent == item.Text {
			count++
			log.Println(item.Text)
		}
		lastSent = item.Text
	}

	log.Println("----> ", count)

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

//estatísticas
func main() {

	mapRepeat := map[string]int{}
	oriNatFile := readFile("/home/sidleal/usp/coling2018/v3/align_size_ori_nat.tsv")
	lines := strings.Split(oriNatFile, "\n")

	oriSentences := []Sentence{}
	for i, line := range lines {
		if i == 0 || line == "" {
			continue
		}
		tokens := strings.Split(line, "\t")
		sentence := Sentence{}
		sentence.Producao = tokens[0]
		sentence.Level = tokens[1]
		sentence.Changed = tokens[2]
		sentence.Splited = tokens[3]
		sentence.Text = tokens[4]

		if _, ok := mapRepeat[sentence.Producao+"-"+sentence.Text]; !ok {
			oriSentences = append(oriSentences, sentence)
			mapRepeat[sentence.Producao+"-"+sentence.Text] = 1
		}

	}

	mapRepeat = map[string]int{}
	natStrFile := readFile("/home/sidleal/usp/coling2018/v3/align_size_nat_str.tsv")
	lines = strings.Split(natStrFile, "\n")

	natSentences := []Sentence{}
	for i, line := range lines {
		if i == 0 || line == "" {
			continue
		}
		tokens := strings.Split(line, "\t")
		sentence := Sentence{}
		sentence.Producao = tokens[0]
		sentence.Level = tokens[1]
		sentence.Changed = tokens[2]
		sentence.Splited = tokens[3]
		sentence.Text = tokens[4]

		if _, ok := mapRepeat[sentence.Producao+"-"+sentence.Text]; !ok {
			natSentences = append(natSentences, sentence)
			mapRepeat[sentence.Producao+"-"+sentence.Text] = 1
		}

	}

	mapRepeat = map[string]int{}
	oriStrFile := readFile("/home/sidleal/usp/coling2018/v3/align_size_ori_str.tsv")
	lines = strings.Split(oriStrFile, "\n")

	numPairsOriStr := 0
	for i, line := range lines {
		if i == 0 || line == "" {
			continue
		}
		tokens := strings.Split(line, "\t")
		if tokens[2] == "S" {
			numPairsOriStr++
		}
	}

	strFile := readFile("/home/sidleal/usp/coling2018/v3/align_all_nat_str.tsv")
	lines = strings.Split(strFile, "\n")

	strSentences := []Sentence{}
	for i, line := range lines {
		if i == 0 || line == "" {
			continue
		}
		tokens := strings.Split(line, "\t")
		sentence := Sentence{}
		sentence.Producao = tokens[0]
		sentence.Level = tokens[1]
		//sentence.Changed = tokens[2]
		//sentence.Splited = tokens[3]
		sentence.Text = tokens[5]

		strSentences = append(strSentences, sentence)

	}

	triFile := readFile("/home/sidleal/usp/coling2018/v3/triplets.tsv")
	lines = strings.Split(triFile, "\n")

	triSameSentence := 0
	triSimplOriNat := 0
	triSimplNatStr := 0
	triSimplAll := 0
	triTotal := 0
	for i, line := range lines {
		if i == 0 || line == "" {
			continue
		}
		tokens := strings.Split(line, "\t")

		if tokens[2] == "N" && tokens[3] == "N" {
			triSameSentence++
		}
		if tokens[2] == "S" && tokens[3] == "N" {
			triSimplOriNat++
		}
		if tokens[2] == "N" && tokens[3] == "S" {
			triSimplNatStr++
		}
		if tokens[2] == "S" && tokens[3] == "S" {
			triSimplAll++
		}
		triTotal++

	}

	allOriFile := readFile("/home/sidleal/usp/coling2018/v3/align_all_ori_nat.tsv")
	lines = strings.Split(allOriFile, "\n")

	allNatFromSplit := 0
	for i, line := range lines {
		if i == 0 || line == "" {
			continue
		}
		tokens := strings.Split(line, "\t")

		if tokens[3] == "S" {
			allNatFromSplit++
		}

	}

	allNatFile := readFile("/home/sidleal/usp/coling2018/v3/align_all_nat_str.tsv")
	lines = strings.Split(allNatFile, "\n")

	allStrFromSplit := 0
	for i, line := range lines {
		if i == 0 || line == "" {
			continue
		}
		tokens := strings.Split(line, "\t")

		if tokens[3] == "S" {
			allStrFromSplit++
		}

	}

	sameSentenceOriNat := 0
	for _, itemOri := range oriSentences {
		for _, itemNat := range natSentences {
			if itemOri.Producao == itemNat.Producao && itemOri.Text == itemNat.Text {
				if itemOri.Changed == "S" || itemOri.Splited == "S" {
					log.Println("ORI", itemOri.Producao, itemOri.Text)
				}
				sameSentenceOriNat++
			}
		}
	}

	sameSentenceNatStr := 0
	for _, itemNat := range natSentences {
		for _, itemStr := range strSentences {
			if itemNat.Producao == itemStr.Producao && itemNat.Text == itemStr.Text {
				if itemNat.Changed == "S" || itemNat.Splited == "S" {
					log.Println("NAT", itemNat.Producao, itemNat.Text)
				}
				sameSentenceNatStr++
			}
		}
	}

	splitSentenceOriNat := 0
	for _, item := range oriSentences {
		if item.Splited == "S" {
			splitSentenceOriNat++
		}
	}

	splitSentenceNatStr := 0
	for _, item := range natSentences {
		if item.Splited == "S" {
			splitSentenceNatStr++
		}
	}

	simpNoSplitSentenceOriNat := 0
	for _, item := range oriSentences {
		if item.Splited == "N" && item.Changed == "S" {
			simpNoSplitSentenceOriNat++
		}
	}

	simpNoSplitSentenceNatStr := 0
	for _, item := range natSentences {
		if item.Splited == "N" && item.Changed == "S" {
			simpNoSplitSentenceNatStr++
		}
	}

	numPairsOriNat := 0
	for _, item := range oriSentences {
		if item.Changed == "S" {
			numPairsOriNat++
		}
	}

	numPairsNatStr := 0
	for _, item := range natSentences {
		if item.Changed == "S" {
			numPairsNatStr++
		}
	}

	log.Println("-------------------------")
	log.Println("Total sentenças Original:", len(oriSentences))
	log.Println("Total sentenças Natural:", len(natSentences))
	log.Println("Total sentenças Strong:", len(strSentences))
	log.Println("Total sentenças GERAL:", len(oriSentences)+len(natSentences)+len(strSentences))
	log.Println("")
	log.Println("Total sentenças IGUAIS Original->Natural:", sameSentenceOriNat)
	log.Println("Total sentenças IGUAIS Natural->Strong:", sameSentenceNatStr)
	log.Println("")
	log.Println("Total sentenças DIVIDIDAS Original->Natural:", splitSentenceOriNat)
	log.Println("Total sentenças DIVIDIDAS Natural->Strong:", splitSentenceNatStr)
	log.Println("")
	log.Println("Total sentenças Natural resultado de Divisão:", allNatFromSplit)
	log.Println("Total sentenças Strong resultado de Divisão:", allStrFromSplit)
	log.Println("")
	log.Println("Total sentenças SIMPLIFICADAS (sem divisão) Original->Natural:", simpNoSplitSentenceOriNat)
	log.Println("Total sentenças SIMPLIFICADAS (sem divisão) Natural->Strong:", simpNoSplitSentenceNatStr)
	log.Println("")
	log.Println("Total pares simplificados Original->Natural:", numPairsOriNat)
	log.Println("Total pares simplificados Natural->Strong:", numPairsNatStr)
	log.Println("Total pares simplificados Original->Strong:", numPairsOriStr)
	log.Println("Total geral pares simplificados:", numPairsOriStr+numPairsNatStr+numPairsOriNat)
	log.Println("")
	log.Println("Total trios IGUAIS 3 Níveis:", triSameSentence)
	log.Println("Total trios Simplificação Apenas Original->Natural:", triSimplOriNat)
	log.Println("Total trios Simplificação Apenas Natural->Strong:", triSimplNatStr)
	log.Println("Total trios Simplificação 3 Níveis:", triSimplAll)
	log.Println("Total trios:", triTotal)
	log.Println("-------------------------")

}
