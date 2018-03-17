package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Sentence struct {
	Producao   string
	Level      string
	Text       string
	Splited    string
	Changed    string
	TextTarget string
}

type SentencePair struct {
	Producao string
	Level    string
	TextA    string
	TextB    string
	Splited  string
	Changed  string
}

func main_files() {

	f1, err := os.OpenFile("/home/sidleal/usp/coling2018/v3/validacao/validacao_ori3.tsv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERRO", err)
	}
	defer f1.Close()

	f2, err := os.OpenFile("/home/sidleal/usp/coling2018/v3/validacao/validacao_nat3.tsv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERRO", err)
	}
	defer f2.Close()

	f3, err := os.OpenFile("/home/sidleal/usp/coling2018/v3/validacao/validacao_str3.tsv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERRO", err)
	}
	defer f3.Close()

	sizeOriSent := readFile("/home/sidleal/usp/coling2018/v3/align_all_ori_nat.tsv")
	lines := strings.Split(sizeOriSent, "\n")

	oriNatPairs := []SentencePair{}
	for i, line := range lines {
		if i == 0 || line == "" {
			continue
		}
		tokens := strings.Split(line, "\t")
		sentence := SentencePair{}
		sentence.Producao = tokens[0]
		sentence.Level = tokens[1]
		sentence.Changed = tokens[2]
		sentence.Splited = tokens[3]
		sentence.TextA = tokens[4]
		sentence.TextB = tokens[5]
		oriNatPairs = append(oriNatPairs, sentence)

	}

	sizeNatSent := readFile("/home/sidleal/usp/coling2018/v3/align_all_nat_str.tsv")
	lines = strings.Split(sizeNatSent, "\n")

	natStrPairs := []SentencePair{}
	for i, line := range lines {
		if i == 0 || line == "" {
			continue
		}
		tokens := strings.Split(line, "\t")
		sentence := SentencePair{}
		sentence.Producao = tokens[0]
		sentence.Level = tokens[1]
		sentence.Changed = tokens[2]
		sentence.Splited = tokens[3]
		sentence.TextA = tokens[4]
		sentence.TextB = tokens[5]
		natStrPairs = append(natStrPairs, sentence)

	}

	header := ""

	metrics := readFile("/home/sidleal/usp/coling2018/v3/pss_sentences_features.tsv")
	lines = strings.Split(metrics, "\n")

	metricList := map[string][]string{}
	for i, line := range lines {
		if i == 0 {
			header = line + "\n"
			continue
		}
		if line == "" {
			continue
		}
		line = strings.TrimSuffix(line, "\t")
		tokens := strings.Split(line, "\t")
		metricList[tokens[0]] = tokens
	}

	_, err = f1.WriteString(header)
	if err != nil {
		log.Println("ERRO", err)
	}

	_, err = f2.WriteString(header)
	if err != nil {
		log.Println("ERRO", err)
	}

	_, err = f3.WriteString(header)
	if err != nil {
		log.Println("ERRO", err)
	}

	lastSent := ""
	for _, item1 := range oriNatPairs {
		log.Println("----------------")
		log.Println(item1.TextA)
		log.Println(item1.TextB)

		if item1.TextA != item1.TextB {
			if lastSent != item1.TextA {
				line := ""
				for i, token := range metricList[item1.TextA] {
					if i > 0 {
						token = round(token)
					}
					line += token + "\t"
				}
				line = strings.TrimSuffix(line, "\t")
				line += "\n"

				_, err := f1.WriteString(line)
				if err != nil {
					log.Println("ERRO", err)
				}
				lastSent = item1.TextA
			}

			line := ""
			for i, token := range metricList[item1.TextB] {
				if i > 0 {
					token = round(token)
				}
				line += token + "\t"
			}
			line = strings.TrimSuffix(line, "\t")
			line += "\n"

			_, err = f2.WriteString(line)
			if err != nil {
				log.Println("ERRO", err)
			}

		}

	}

	for _, item1 := range natStrPairs {
		log.Println("----------------")
		log.Println(item1.TextA)
		log.Println(item1.TextB)

		if item1.TextA != item1.TextB {
			line := ""
			for i, token := range metricList[item1.TextB] {
				if i > 0 {
					token = round(token)
				}
				line += token + "\t"
			}
			line = strings.TrimSuffix(line, "\t")
			line += "\n"

			_, err = f3.WriteString(line)
			if err != nil {
				log.Println("ERRO", err)
			}

		}

	}

	log.Println("----> ")

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

//estatísticas
func main() {

	mapRepeat := map[string]int{}
	oriNatFile := readFile("/home/sidleal/usp/coling2018/v3/align_size_ori_nat.tsv")
	lines := strings.Split(oriNatFile, "\n")

	sentencesZH := 0
	sentencesFSP := 0

	oriSizeSentences := []Sentence{}
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
		sentence.TextTarget = tokens[5]

		if _, ok := mapRepeat[sentence.Producao+"-"+sentence.Text]; !ok {
			oriSizeSentences = append(oriSizeSentences, sentence)
			mapRepeat[sentence.Producao+"-"+sentence.Text] = 1
			prod, _ := strconv.Atoi(sentence.Producao)
			if prod < 116 {
				sentencesZH++
			} else {
				sentencesFSP++
			}
		}

	}

	oriNatAllFile := readFile("/home/sidleal/usp/coling2018/v3/align_all_ori_nat.tsv")
	lines = strings.Split(oriNatAllFile, "\n")

	oriNatAllSentences := []Sentence{}
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
		sentence.TextTarget = tokens[5]

		oriNatAllSentences = append(oriNatAllSentences, sentence)

	}

	oriStrAllFile := readFile("/home/sidleal/usp/coling2018/v3/align_all_ori_str.tsv")
	lines = strings.Split(oriStrAllFile, "\n")

	oriStrAllSentences := []Sentence{}
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
		sentence.TextTarget = tokens[5]

		oriStrAllSentences = append(oriStrAllSentences, sentence)

	}

	mapRepeat = map[string]int{}
	natStrFile := readFile("/home/sidleal/usp/coling2018/v3/align_size_nat_str.tsv")
	lines = strings.Split(natStrFile, "\n")

	natSizeSentences := []Sentence{}
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
		sentence.TextTarget = tokens[5]

		if _, ok := mapRepeat[sentence.Producao+"-"+sentence.Text]; !ok {
			natSizeSentences = append(natSizeSentences, sentence)
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

	strSizeSentences := []Sentence{}
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
		sentence.TextTarget = tokens[5]

		if _, ok := mapRepeat[sentence.Producao+"-"+sentence.Text]; !ok {
			strSizeSentences = append(strSizeSentences, sentence)
			mapRepeat[sentence.Producao+"-"+sentence.Text] = 1
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
	natStrAllSentences := strSentences

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
	for _, itemOri := range oriSizeSentences {
		for _, itemNat := range natSizeSentences {
			if itemOri.Producao == itemNat.Producao && itemOri.Text == itemNat.Text {
				if itemOri.Changed == "S" || itemOri.Splited == "S" {
					log.Println("ORI", itemOri.Producao, itemOri.Text)
				}
				sameSentenceOriNat++
			}
		}
	}

	sameSentenceNatStr := 0
	for _, itemNat := range natSizeSentences {
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
	for _, item := range oriSizeSentences {
		if item.Splited == "S" {
			splitSentenceOriNat++
		}
	}

	splitSentenceNatStr := 0
	for _, item := range natSizeSentences {
		if item.Splited == "S" {
			splitSentenceNatStr++
		}
	}

	simpNoSplitSentenceOriNat := 0
	for _, item := range oriSizeSentences {
		if item.Splited == "N" && item.Changed == "S" {
			simpNoSplitSentenceOriNat++
		}
	}

	simpNoSplitSentenceNatStr := 0
	for _, item := range natSizeSentences {
		if item.Splited == "N" && item.Changed == "S" {
			simpNoSplitSentenceNatStr++
		}
	}

	numPairsOriNat := 0
	for _, item := range oriSizeSentences {
		if item.Changed == "S" {
			numPairsOriNat++
		}
	}

	numPairsNatStr := 0
	for _, item := range natSizeSentences {
		if item.Changed == "S" {
			numPairsNatStr++
		}
	}

	mediaSimplificacoesSemDivisaoOri := 0
	minSimplificacoesSemDivisaoOri := 0
	maxSimplificacoesSemDivisaoOri := 0
	for _, item := range oriSizeSentences {
		if item.Splited == "N" && item.Changed == "S" {
			tokens := tokenizeText(item.TextTarget)
			mediaSimplificacoesSemDivisaoOri = (mediaSimplificacoesSemDivisaoOri + len(tokens)) / 2
			if len(tokens) > maxSimplificacoesSemDivisaoOri {
				maxSimplificacoesSemDivisaoOri = len(tokens)
			}
			if minSimplificacoesSemDivisaoOri == 0 || len(tokens) < minSimplificacoesSemDivisaoOri {
				minSimplificacoesSemDivisaoOri = len(tokens)
			}
		}
	}

	mediaSimplificacoesComDivisaoOri := 0
	minSimplificacoesComDivisaoOri := 0
	maxSimplificacoesComDivisaoOri := 0
	for _, item := range oriSizeSentences {
		if item.Splited == "S" {
			tokens := tokenizeText(item.TextTarget)
			mediaSimplificacoesComDivisaoOri = (mediaSimplificacoesComDivisaoOri + len(tokens)) / 2
			if len(tokens) > maxSimplificacoesComDivisaoOri {
				maxSimplificacoesComDivisaoOri = len(tokens)
			}
			if minSimplificacoesComDivisaoOri == 0 || len(tokens) < minSimplificacoesComDivisaoOri {
				minSimplificacoesComDivisaoOri = len(tokens)
			}
		}
	}

	mediaSimplificacoesSemDivisaoNat := 0
	minSimplificacoesSemDivisaoNat := 0
	maxSimplificacoesSemDivisaoNat := 0
	for _, item := range natSizeSentences {
		if item.Splited == "N" && item.Changed == "S" {
			tokens := tokenizeText(item.TextTarget)
			mediaSimplificacoesSemDivisaoNat = (mediaSimplificacoesSemDivisaoNat + len(tokens)) / 2
			if len(tokens) > maxSimplificacoesSemDivisaoNat {
				maxSimplificacoesSemDivisaoNat = len(tokens)
			}
			if minSimplificacoesSemDivisaoNat == 0 || len(tokens) < minSimplificacoesSemDivisaoNat {
				minSimplificacoesSemDivisaoNat = len(tokens)
			}
		}
	}

	mediaSimplificacoesComDivisaoNat := 0
	minSimplificacoesComDivisaoNat := 0
	maxSimplificacoesComDivisaoNat := 0
	for _, item := range natSizeSentences {
		if item.Splited == "S" {
			tokens := tokenizeText(item.TextTarget)
			mediaSimplificacoesComDivisaoNat = (mediaSimplificacoesComDivisaoNat + len(tokens)) / 2
			if len(tokens) > maxSimplificacoesComDivisaoNat {
				maxSimplificacoesComDivisaoNat = len(tokens)
			}
			if minSimplificacoesComDivisaoNat == 0 || len(tokens) < minSimplificacoesComDivisaoNat {
				minSimplificacoesComDivisaoNat = len(tokens)
			}
		}
	}

	mediaDiffSimplificacoesSemDivisaoOri := 0
	minDiffSimplificacoesSemDivisaoOri := 0
	maxDiffSimplificacoesSemDivisaoOri := 0
	for _, item := range oriSizeSentences {
		if item.Splited == "N" && item.Changed == "S" {
			tokensO := tokenizeText(item.Text)
			tokensT := tokenizeText(item.TextTarget)
			if len(tokensO)-len(tokensT) > 0 {
				mediaDiffSimplificacoesSemDivisaoOri = (mediaDiffSimplificacoesSemDivisaoOri + (len(tokensO) - len(tokensT))) / 2
				if (len(tokensO) - len(tokensT)) > maxDiffSimplificacoesSemDivisaoOri {
					maxDiffSimplificacoesSemDivisaoOri = (len(tokensO) - len(tokensT))
				}
				if minDiffSimplificacoesSemDivisaoOri == 0 || (len(tokensO)-len(tokensT)) < minDiffSimplificacoesSemDivisaoOri {
					minDiffSimplificacoesSemDivisaoOri = (len(tokensO) - len(tokensT))
				}
			}
		}
	}

	mediaDiffSimplificacoesComDivisaoOri := 0
	minDiffSimplificacoesComDivisaoOri := 0
	maxDiffSimplificacoesComDivisaoOri := 0
	for _, item := range oriSizeSentences {
		if item.Splited == "S" {
			tokensO := tokenizeText(item.Text)
			tokensT := tokenizeText(item.TextTarget)
			if len(tokensO)-len(tokensT) > 0 {
				mediaDiffSimplificacoesComDivisaoOri = (mediaDiffSimplificacoesComDivisaoOri + (len(tokensO) - len(tokensT))) / 2
				if (len(tokensO) - len(tokensT)) > maxDiffSimplificacoesComDivisaoOri {
					maxDiffSimplificacoesComDivisaoOri = (len(tokensO) - len(tokensT))
				}
				if minDiffSimplificacoesComDivisaoOri == 0 || (len(tokensO)-len(tokensT)) < minDiffSimplificacoesComDivisaoOri {
					minDiffSimplificacoesComDivisaoOri = (len(tokensO) - len(tokensT))
				}
			}
		}
	}

	pss1OriSize := 0
	for _, item := range oriNatAllSentences {
		if item.TextTarget != item.Text {
			pss1OriSize++
		}
	}
	pss1NatSize := 0
	for _, item := range natStrAllSentences {
		if item.TextTarget != item.Text {
			pss1NatSize++
		}
	}
	pss1StrSize := 0
	for _, item := range oriStrAllSentences {
		if item.TextTarget != item.Text {
			pss1StrSize++
		}
	}
	pss1Total := pss1OriSize + pss1NatSize + pss1StrSize

	pss2OriSize := 0
	for _, item := range oriSizeSentences {
		if item.TextTarget != item.Text {
			pss2OriSize++
		}
	}
	pss2NatSize := 0
	for _, item := range natSizeSentences {
		if item.TextTarget != item.Text {
			pss2NatSize++
		}
	}
	pss2StrSize := 0
	for _, item := range strSizeSentences {
		if item.TextTarget != item.Text {
			pss2StrSize++
		}
	}
	pss2Total := pss2OriSize + pss2NatSize + pss2StrSize

	pss3OriSize := 0
	for _, item := range oriSizeSentences {
		if item.TextTarget != item.Text && item.Splited == "N" {
			pss3OriSize++
		}
	}
	pss3NatSize := 0
	for _, item := range natSizeSentences {
		if item.TextTarget != item.Text && item.Splited == "N" {
			pss3NatSize++
		}
	}
	pss3StrSize := 0
	for _, item := range strSizeSentences {
		if item.TextTarget != item.Text && item.Splited == "N" {
			pss3StrSize++
		}
	}
	pss3Total := pss3OriSize + pss3NatSize + pss3StrSize

	log.Println("-------------------------")
	log.Println("Total sentenças Original:", len(oriSizeSentences))
	log.Println("      Zero Hora:", sentencesZH)
	log.Println("      Caderno Ciencia FSP:", sentencesFSP)
	log.Println("Total sentenças Natural:", len(natSizeSentences))
	log.Println("Total sentenças Strong:", len(strSentences))
	log.Println("Total sentenças GERAL:", len(oriSizeSentences)+len(natSizeSentences)+len(strSentences))
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
	log.Println("")
	log.Println("Tamanho médio em tokens das sentenças simplificadas (sem divisão) - Ori->Nat:", mediaSimplificacoesSemDivisaoOri)
	log.Println("Tamanho mínimo em tokens das sentenças simplificadas (sem divisão) - Ori->Nat:", minSimplificacoesSemDivisaoOri)
	log.Println("Tamanho máximo em tokens das sentenças simplificadas (sem divisão) - Ori->Nat:", maxSimplificacoesSemDivisaoOri)
	log.Println("")
	log.Println("Tamanho médio em tokens das sentenças simplificadas (com divisão) - Ori->Nat:", mediaSimplificacoesComDivisaoOri)
	log.Println("Tamanho mínimo em tokens das sentenças simplificadas (com divisão) - Ori->Nat:", minSimplificacoesComDivisaoOri)
	log.Println("Tamanho máximo em tokens das sentenças simplificadas (com divisão) - Ori->Nat:", maxSimplificacoesComDivisaoOri)
	log.Println("")
	log.Println("Tamanho médio em tokens das sentenças simplificadas (sem divisão) - Nat->Str:", mediaSimplificacoesSemDivisaoNat)
	log.Println("Tamanho mínimo em tokens das sentenças simplificadas (sem divisão) - Nat->Str:", minSimplificacoesSemDivisaoNat)
	log.Println("Tamanho máximo em tokens das sentenças simplificadas (sem divisão) - Nat->Str:", maxSimplificacoesSemDivisaoNat)
	log.Println("")
	log.Println("Tamanho médio em tokens das sentenças simplificadas (com divisão) - Nat->Str:", mediaSimplificacoesComDivisaoNat)
	log.Println("Tamanho mínimo em tokens das sentenças simplificadas (com divisão) - Nat->Str:", minSimplificacoesComDivisaoNat)
	log.Println("Tamanho máximo em tokens das sentenças simplificadas (com divisão) - Nat->Str:", maxSimplificacoesComDivisaoNat)
	log.Println("")
	log.Println("Tamanho médio da diferença em tokens das sentenças originais vs simplificadas (sem divisão) - Ori->Nat:", mediaDiffSimplificacoesSemDivisaoOri)
	log.Println("Tamanho mínimo da diferença em tokens das sentenças originais vs simplificadas (sen divisão) - Ori->Nat:", minDiffSimplificacoesSemDivisaoOri)
	log.Println("Tamanho máximo da diferença em tokens das sentenças originais vs simplificadas (sem divisão) - Ori->Nat:", maxDiffSimplificacoesSemDivisaoOri)
	log.Println("")
	log.Println("Tamanho médio da diferença em tokens das sentenças originais vs simplificadas (com divisão) - Ori->Nat:", mediaDiffSimplificacoesComDivisaoOri)
	log.Println("Tamanho mínimo da diferença em tokens das sentenças originais vs simplificadas (com divisão) - Ori->Nat:", minDiffSimplificacoesComDivisaoOri)
	log.Println("Tamanho máximo da diferença em tokens das sentenças originais vs simplificadas (com divisão) - Ori->Nat:", maxDiffSimplificacoesComDivisaoOri)
	log.Println("")
	log.Println("Total PSS1 Original->Natural:", pss1OriSize)
	log.Println("Total PSS1 Natural->Strong:", pss1NatSize)
	log.Println("Total PSS1 Original->Strong:", pss1StrSize)
	log.Println("Total geral PSS1:", pss1Total)
	log.Println("")
	log.Println("Total PSS2 Original->Natural:", pss2OriSize)
	log.Println("Total PSS2 Natural->Strong:", pss2NatSize)
	log.Println("Total PSS2 Original->Strong:", pss2StrSize)
	log.Println("Total geral PSS2:", pss2Total)
	log.Println("")
	log.Println("Total PSS3 Original->Natural:", pss3OriSize)
	log.Println("Total PSS3 Natural->Strong:", pss3NatSize)
	log.Println("Total PSS3 Original->Strong:", pss3StrSize)
	log.Println("Total geral PSS3:", pss3Total)
	log.Println("")
	log.Println("-------------------------")

}

func tokenizeText(rawText string) []string {
	regEx := regexp.MustCompile(`([A-z]+)-([A-z]+)`)
	rawText = regEx.ReplaceAllString(rawText, "$1|hyp|$2")

	regEx = regexp.MustCompile(`\|gdot\|`)
	rawText = regEx.ReplaceAllString(rawText, ".")

	regEx = regexp.MustCompile(`\|gint\|`)
	rawText = regEx.ReplaceAllString(rawText, "?")

	regEx = regexp.MustCompile(`\|gexc\|`)
	rawText = regEx.ReplaceAllString(rawText, "!")

	regEx = regexp.MustCompile(`([\.\,"\(\)\[\]\{\}\?\!;:-]{1})`)
	rawText = regEx.ReplaceAllString(rawText, "  $1 ")

	regEx = regexp.MustCompile(`\s+`)
	rawText = regEx.ReplaceAllString(rawText, " ")

	return strings.Split(rawText, " ")
}
