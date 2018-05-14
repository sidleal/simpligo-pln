package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// type WordD struct {
// 	Text         string
// 	Informative  int
// 	Encyclopedic int
// 	Spoken       int
// 	Prose        int
// 	Science      int
// 	ScienceCom   int

// 	Inf_G1    int
// 	Inf_GNews int
// 	Inf_Plnbr int
// 	Enc_Wiki  int
// 	Spo_Imdb  int
// 	Spo_Lex   int
// 	Pro_Por   int
// 	Pro_Dom   int
// 	Sci_ME    int
// 	Sci_CHC   int
// 	Sci2_Fap  int

// 	FreqTotal int
// 	FreqAvg   float64
// 	DP        float64
// 	DPNorm    float64
// 	FreqCorr  float64
// }

// func main() {

// 	list := map[string]WordD{}

// 	path := "/home/sidleal/usp/semestre5/propor_dp/lista_freq_nathan.csv"
// 	data := readFile(path)

// 	lines := strings.Split(data, "\n")
// 	for i, line := range lines {
// 		cols := strings.Split(line, ",")
// 		if i == 0 {
// 			log.Println(line)
// 			log.Println("-------------")
// 			log.Println(cols[6], cols[7], cols[14])
// 			log.Println("-------------")
// 			log.Println(cols[17])
// 			log.Println("-------------")
// 			log.Println(cols[16])
// 			log.Println("-------------")
// 			log.Println(cols[9], cols[11])
// 			log.Println("-------------")
// 			log.Println(cols[12], cols[2])
// 			log.Println("-------------")
// 			log.Println(cols[4])
// 			continue
// 		}

// 		// log.Println(i, cols[0])

// 		word := WordD{}
// 		word.Text = cols[0]

// 		a, _ := strconv.Atoi(cols[6])  //g1
// 		b, _ := strconv.Atoi(cols[7])  // gnews
// 		c, _ := strconv.Atoi(cols[14]) // plnbr
// 		word.Inf_G1 = a
// 		word.Inf_GNews = b
// 		word.Inf_Plnbr = c
// 		word.Informative = a + b + c

// 		a, _ = strconv.Atoi(cols[17]) // wiki
// 		word.Enc_Wiki = a
// 		word.Encyclopedic = a

// 		a, _ = strconv.Atoi(cols[16]) // imdb
// 		word.Spo_Imdb = a
// 		word.Spoken = a

// 		a, _ = strconv.Atoi(cols[9])  // livr dom pub
// 		b, _ = strconv.Atoi(cols[11]) // livr port
// 		word.Pro_Dom = a
// 		word.Pro_Por = b
// 		word.Prose = a + b

// 		a, _ = strconv.Atoi(cols[12]) //mundo estr
// 		b, _ = strconv.Atoi(cols[2])  // chc
// 		word.Sci_ME = a
// 		word.Sci_CHC = b
// 		word.Science = a + b

// 		a, _ = strconv.Atoi(cols[4]) // fapesp
// 		word.Sci2_Fap = a
// 		word.ScienceCom = a

// 		if (word.Informative + word.Encyclopedic + word.Spoken + word.Prose + word.Science + word.ScienceCom) > 0 {
// 			list[word.Text] = word
// 		}
// 	}

// 	path = "/home/sidleal/usp/semestre5/propor_dp/SUBTLEX-BRPOR.no-web.valid-char.csv"
// 	data = readFile(path)

// 	lines = strings.Split(data, "\n")
// 	for i, line := range lines {
// 		cols := strings.Split(line, "\t")
// 		if i == 0 {
// 			log.Println(line)
// 			continue
// 		}

// 		// log.Println(i, cols[0])

// 		a, _ := strconv.Atoi(cols[1])

// 		if word, found := list[cols[0]]; found {
// 			word.Spo_Lex = a
// 			word.Spoken = word.Spoken + a
// 			list[word.Text] = word
// 		} else {
// 			word := WordD{}
// 			word.Text = cols[0]
// 			word.Spo_Lex = a
// 			word.Spoken = a
// 			list[word.Text] = word
// 		}

// 	}

// 	// log.Println("------------------")
// 	// i := 0
// 	// for k, v := range list {
// 	// 	if i < 10 {
// 	// 		log.Println(k, v)
// 	// 	}
// 	// 	i++
// 	// }

// 	log.Println("carregando delaf")
// 	delaf := map[string]string{}

// 	path = "/home/sidleal/usp/semestre5/propor_dp/DELAF_PB_2018.dic"
// 	lines = readFileLines(path)
// 	for _, line := range lines {
// 		if line == "" {
// 			continue
// 		}
// 		cols := strings.Split(line, ",")
// 		delaf[cols[0]] = cols[1]
// 	}

// 	log.Println("delaf em memoria")

// 	exclude := map[string]bool{}
// 	for _, l := range []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"} {
// 		exclude[l] = true
// 	}

// 	listAll := map[string]WordD{}

// 	// i := 0
// 	regEx := regexp.MustCompile(`N\+Pr`)
// 	for k, v := range list {
// 		// if i < 100 {

// 		if _, found := exclude[k]; found {
// 			continue
// 		}

// 		if delafDesc, found := delaf[k]; found {
// 			matches := regEx.FindAllStringSubmatch(delafDesc, -1)
// 			if len(matches) == 0 {
// 				listAll[k] = v
// 			}
// 		}

// 		// }
// 		// i++
// 	}

// 	totInf_G1 := 0
// 	totInf_GNews := 0
// 	totInf_Plnbr := 0
// 	totEnc_Wiki := 0
// 	totSpo_Imdb := 0
// 	totSpo_Lex := 0
// 	totPro_Por := 0
// 	totPro_Dom := 0
// 	totSci_ME := 0
// 	totSci_CHC := 0
// 	totSci2_Fap := 0

// 	for _, v := range listAll {
// 		// log.Println(k, v)

// 		totInf_G1 += v.Inf_G1
// 		totInf_GNews += v.Inf_GNews
// 		totInf_Plnbr += v.Inf_Plnbr
// 		totEnc_Wiki += v.Enc_Wiki
// 		totSpo_Imdb += v.Spo_Imdb
// 		totSpo_Lex += v.Spo_Lex
// 		totPro_Por += v.Pro_Por
// 		totPro_Dom += v.Pro_Dom
// 		totSci_ME += v.Sci_ME
// 		totSci_CHC += v.Sci_CHC
// 		totSci2_Fap += v.Sci2_Fap

// 	}

// 	log.Println("Tots", totInf_G1, totInf_GNews, totInf_Plnbr, totEnc_Wiki, totSpo_Imdb, totSpo_Lex, totPro_Por, totPro_Dom, totSci_ME, totSci_CHC, totSci2_Fap)
// 	totAll := totInf_G1 + totInf_GNews + totInf_Plnbr + totEnc_Wiki + totSpo_Imdb + totSpo_Lex + totPro_Por + totPro_Dom + totSci_ME + totSci_CHC + totSci2_Fap

// 	perInf_G1 := float64(totInf_G1) / float64(totAll)
// 	perInf_GNews := float64(totInf_GNews) / float64(totAll)
// 	perInf_Plnbr := float64(totInf_Plnbr) / float64(totAll)
// 	perEnc_Wiki := float64(totEnc_Wiki) / float64(totAll)
// 	perSpo_Imdb := float64(totSpo_Imdb) / float64(totAll)
// 	perSpo_Lex := float64(totSpo_Lex) / float64(totAll)
// 	perPro_Por := float64(totPro_Por) / float64(totAll)
// 	perPro_Dom := float64(totPro_Dom) / float64(totAll)
// 	perSci_ME := float64(totSci_ME) / float64(totAll)
// 	perSci_CHC := float64(totSci_CHC) / float64(totAll)
// 	perSci2_Fap := float64(totSci2_Fap) / float64(totAll)

// 	log.Println("Perc", perInf_G1, perInf_GNews, perInf_Plnbr, perEnc_Wiki, perSpo_Imdb, perSpo_Lex, perPro_Por, perPro_Dom, perSci_ME, perSci_CHC, perSci2_Fap)

// 	s := []float64{perInf_G1, perInf_GNews, perInf_Plnbr, perEnc_Wiki, perSpo_Imdb, perSpo_Lex, perPro_Por, perPro_Dom, perSci_ME, perSci_CHC, perSci2_Fap}

// 	listFinal := map[string]WordD{}

// 	for k, w := range listAll {
// 		// log.Println(k, w)

// 		f := float64(w.Inf_G1 + w.Inf_GNews + w.Inf_Plnbr + w.Enc_Wiki + w.Spo_Imdb + w.Spo_Lex + w.Pro_Por + w.Pro_Dom + w.Sci_ME + w.Sci_CHC + w.Sci2_Fap)
// 		v := []float64{float64(w.Inf_G1), float64(w.Inf_GNews), float64(w.Inf_Plnbr), float64(w.Enc_Wiki), float64(w.Spo_Imdb), float64(w.Spo_Lex), float64(w.Pro_Por), float64(w.Pro_Dom), float64(w.Sci_ME), float64(w.Sci_CHC), float64(w.Sci2_Fap)}

// 		// Deviation of proportions DP

// 		div := []float64{}
// 		for i, _ := range v {
// 			div = append(div, v[i]/f)
// 		}

// 		sub := []float64{}
// 		for i, _ := range div {
// 			sub = append(sub, div[i]-s[i])
// 		}

// 		sum := 0.0
// 		for i, _ := range sub {
// 			sum += math.Abs(sub[i])
// 		}

// 		dp := sum / 2
// 		dp_norm := dp / (1 - min(s))

// 		w.FreqTotal = int(f)
// 		w.DP = dp
// 		w.DPNorm = dp_norm
// 		w.FreqAvg = float64(w.FreqTotal) / float64(len(v))
// 		w.FreqCorr = w.FreqAvg * (1 - w.DPNorm)

// 		listFinal[k] = w

// 		// log.Println("------------>", f, dp, dp_norm)

// 	}

// 	f1, err := os.Create("/home/sidleal/usp/semestre5/propor_dp/dispersion2_det.tsv")
// 	if err != nil {
// 		log.Println("ERRO", err)
// 	}

// 	defer f1.Close()

// 	_, err = f1.WriteString("word\tinf_G1\tinf_Gnews\tinf_plnbr\tenc_wiki\tspo_imdb\tspo_lexbr\tpro_port\tpro_dom_pub\tsci_me\tsci_chc\tsci_com_fap\ttotal\taverage\tdp\tdp_norm\tfreq_corr\n")
// 	if err != nil {
// 		log.Println("ERRO", err)
// 	}

// 	for k, w := range listFinal {
// 		_, err := f1.WriteString(fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", k, w.Inf_G1, w.Inf_GNews, w.Inf_Plnbr, w.Enc_Wiki, w.Spo_Imdb, w.Spo_Lex, w.Pro_Por, w.Pro_Dom, w.Sci_ME, w.Sci_CHC, w.Sci2_Fap, w.FreqTotal, w.FreqAvg, w.DP, w.DPNorm, w.FreqCorr))
// 		if err != nil {
// 			log.Println("ERRO", err)
// 		}
// 	}

// }

//---------------------------------------------------------------------------------------------

type Word struct {
	Text         string
	Lema         string
	Tag          string
	Informative  int
	Encyclopedic int
	Spoken       int
	Prose        int
	Science      int
	ScienceCom   int
	FreqTotal    int
	FreqAvg      float64
	DP           float64
	DPNorm       float64
	FreqCorr     float64
	FreqCorrPer  float64
	FreqCorrAcu  float64
	Ranking      int
	Extra        string
}

type WordOrder []Word

func (a WordOrder) Len() int      { return len(a) }
func (a WordOrder) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a WordOrder) Less(i, j int) bool {
	return a[i].FreqCorr > a[j].FreqCorr
}

func main_all() {

	list := map[string]Word{}

	path := "/home/sidleal/usp/semestre5/propor_dp/lista_freq_nathan.csv"
	data := readFile(path)

	lines := strings.Split(data, "\n")
	for i, line := range lines {
		cols := strings.Split(line, ",")
		if i == 0 {
			log.Println(line)
			log.Println("-------------")
			log.Println(cols[6], cols[7], cols[14])
			log.Println("-------------")
			log.Println(cols[17])
			log.Println("-------------")
			log.Println(cols[16])
			log.Println("-------------")
			log.Println(cols[9], cols[11])
			log.Println("-------------")
			log.Println(cols[12], cols[2])
			log.Println("-------------")
			log.Println(cols[4])
			continue
		}

		// log.Println(i, cols[0])

		word := Word{}
		word.Text = cols[0]

		a, _ := strconv.Atoi(cols[6])
		b, _ := strconv.Atoi(cols[7])
		c, _ := strconv.Atoi(cols[14])
		word.Informative = a + b + c

		a, _ = strconv.Atoi(cols[17])
		word.Encyclopedic = a

		a, _ = strconv.Atoi(cols[16])
		word.Spoken = a

		a, _ = strconv.Atoi(cols[9])
		b, _ = strconv.Atoi(cols[11])
		word.Prose = a + b

		a, _ = strconv.Atoi(cols[12])
		b, _ = strconv.Atoi(cols[2])
		word.Science = a + b

		a, _ = strconv.Atoi(cols[4])
		word.ScienceCom = a

		if (word.Informative + word.Encyclopedic + word.Spoken + word.Prose + word.Science + word.ScienceCom) > 0 {
			list[word.Text] = word
		}
	}

	path = "/home/sidleal/usp/semestre5/propor_dp/SUBTLEX-BRPOR.no-web.valid-char.csv"
	data = readFile(path)

	lines = strings.Split(data, "\n")
	for i, line := range lines {
		cols := strings.Split(line, "\t")
		if i == 0 {
			log.Println(line)
			continue
		}

		// log.Println(i, cols[0])

		a, _ := strconv.Atoi(cols[1])

		if word, found := list[cols[0]]; found {
			word.Spoken = word.Spoken + a
			list[word.Text] = word
		} else {
			word := Word{}
			word.Text = cols[0]
			word.Spoken = a
			list[word.Text] = word
		}

	}

	// log.Println("------------------")
	// i := 0
	// for k, v := range list {
	// 	if i < 10 {
	// 		log.Println(k, v)
	// 	}
	// 	i++
	// }

	log.Println("carregando delaf")
	delaf := map[string]string{}

	regEx1 := regexp.MustCompile(`(.+)\.(.+)`)

	path = "/home/sidleal/usp/semestre5/propor_dp/DELAF_PB_2018.dic"
	lines = readFileLines(path)
	for _, line := range lines {
		if line == "" {
			continue
		}
		cols := strings.Split(line, ",")

		info := cols[1]
		lema := ""
		tag := ""
		match := regEx1.FindStringSubmatch(info)
		if len(match) > 1 {
			lema = match[1]
			tag = match[2]
			info = lema + "|" + tag
		}

		if delafDesc, found := delaf[cols[0]]; found {
			tokens := strings.Split(delafDesc, "|")
			lemaB := tokens[0]
			tagB := tokens[1]

			pesoA := getTagWeight(tag)
			pesoB := getTagWeight(tagB)

			if pesoA < pesoB {
				info = lemaB + "|" + tagB + "," + tag
			} else {
				info = lema + "|" + tag + "," + tagB
			}
			// log.Println(info)
		}
		delaf[cols[0]] = info
	}

	log.Println("delaf em memoria")

	exclude := map[string]bool{}
	for _, l := range []string{"b", "c", "d", "f", "g", "h", "i", "j", "k", "l", "m", "n", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"} {
		exclude[l] = true
	}

	listAll := map[string]Word{}

	// i := 0
	// regEx := regexp.MustCompile(`N\+Pr`)
	regEx := regexp.MustCompile(`(PREP|PRO)X(DET|PRO|PREP|ADV).*`)
	for k, v := range list {
		// if i < 100 {

		if _, found := exclude[k]; found {
			continue
		}

		if delafDesc, found := delaf[k]; found {
			// log.Println(delafDesc)
			tokens := strings.Split(delafDesc, "|")
			lema := tokens[0]
			tag := tokens[1]
			if !strings.HasPrefix(tag, "N+Pr") {
				v.Lema = lema
				v.Tag = tag

				match := regEx.FindStringSubmatch(tag)
				if len(match) > 1 {
					v.Extra = "contração"
				}
				listAll[k] = v
			}
			// matches := regEx.FindAllStringSubmatch(delafDesc, -1)
			// if len(matches) == 0 {
			// listAll[k] = v
			// }
		}

		// }
		// i++
	}

	//aplica correções manuais
	path = "/home/sidleal/usp/semestre5/propor_dp/anotacao1.csv"
	data = readFile(path)
	lines = strings.Split(data, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		cols := strings.Split(line, ",")

		word := cols[0]
		action := cols[1]
		extra := cols[2]

		if action == "EP" || action == "EE" {
			delete(listAll, word)

		} else if action == "CJ" {
			w1 := listAll[word]
			w2 := listAll[extra]

			w2.Informative = w2.Informative + w1.Informative
			w2.Encyclopedic = w2.Encyclopedic + w1.Encyclopedic
			w2.Spoken = w2.Spoken + w1.Spoken
			w2.Prose = w2.Prose + w1.Prose
			w2.Science = w2.Science + w1.Science
			w2.ScienceCom = w2.ScienceCom + w1.ScienceCom

			if w2.Text == "" {
				if delafDesc, found := delaf[extra]; found {
					tokens := strings.Split(delafDesc, "|")
					w2.Lema = tokens[0]
					w2.Tag = tokens[1]
					w2.Text = extra
				} else {
					w2.Lema = extra
					w2.Text = extra
					w2.Tag = w1.Tag
				}
			}

			listAll[extra] = w2
			delete(listAll, word)

			// log.Println(extra, w2.Text)

		} else {
			w := listAll[word]
			w.Extra = action
			listAll[word] = w
		}

		// log.Println("----", cols[0], cols[1], cols[2])

	}

	totInformative := 0
	totEncyclopedic := 0
	totSpoken := 0
	totProse := 0
	totScience := 0
	totScienceCom := 0

	for k, v := range listAll {
		// log.Println(k, v)

		if v.Text == "" {
			log.Println(k, "nadaaaaaaaaaaaaaaaaaaaaaaaaaa")
		}

		totInformative += v.Informative
		totEncyclopedic += v.Encyclopedic
		totSpoken += v.Spoken
		totProse += v.Prose
		totScience += v.Science
		totScienceCom += v.ScienceCom

	}

	log.Println("Tots", totInformative, totEncyclopedic, totSpoken, totProse, totScience, totScienceCom)

	totAll := totInformative + totEncyclopedic + totSpoken + totProse + totScience + totScienceCom

	perInformative := float64(totInformative) / float64(totAll)
	perEncyclopedic := float64(totEncyclopedic) / float64(totAll)
	perSpoken := float64(totSpoken) / float64(totAll)
	perProse := float64(totProse) / float64(totAll)
	perScience := float64(totScience) / float64(totAll)
	perScienceCom := float64(totScienceCom) / float64(totAll)

	log.Println("Perc", perInformative, perEncyclopedic, perSpoken, perProse, perScience, perScienceCom)

	s := []float64{perInformative, perEncyclopedic, perSpoken, perProse, perScience, perScienceCom}

	// listFinal := map[string]Word{}
	listFinal := []Word{}

	var totFreqCorr float64 = 0

	for _, w := range listAll {
		// log.Println(k, w)

		f := float64(w.Informative + w.Encyclopedic + w.Spoken + w.Prose + w.Science + w.ScienceCom)
		v := []float64{float64(w.Informative), float64(w.Encyclopedic), float64(w.Spoken), float64(w.Prose), float64(w.Science), float64(w.ScienceCom)}
		// Deviation of proportions DP

		div := []float64{}
		for i, _ := range v {
			div = append(div, v[i]/f)
		}

		sub := []float64{}
		for i, _ := range div {
			sub = append(sub, div[i]-s[i])
		}

		sum := 0.0
		for i, _ := range sub {
			sum += math.Abs(sub[i])
		}

		dp := sum / 2
		dp_norm := dp / (1 - min(s))

		w.FreqTotal = int(f)
		w.DP = dp
		w.DPNorm = dp_norm
		w.FreqAvg = float64(w.FreqTotal) / float64(len(v))
		w.FreqCorr = w.FreqAvg*(1-w.DPNorm) + 0.0001

		totFreqCorr += w.FreqCorr

		listFinal = append(listFinal, w)

		// log.Println("------------>", f, dp, dp_norm)

	}

	sort.Sort(WordOrder(listFinal))

	listFinalMesmo := []Word{}
	var freqCorrAcum float64 = 0
	for i, w := range listFinal {
		w.Ranking = i + 1
		w.FreqCorrPer = w.FreqCorr / totFreqCorr * 100
		freqCorrAcum += w.FreqCorrPer
		w.FreqCorrAcu = freqCorrAcum
		listFinalMesmo = append(listFinalMesmo, w)
	}

	f1, err := os.Create("/home/sidleal/usp/semestre5/propor_dp/dispersion7.tsv")
	if err != nil {
		log.Println("ERRO", err)
	}

	defer f1.Close()

	_, err = f1.WriteString("word\tlema\ttag\textra\tinformative\tencyclopedic\tspoken\tprose\tscience\tscience_com\ttotal\taverage\tdp\tdp_norm\tfreq_corr\tfreq_perc\tfreq_acum\tranking\n")
	if err != nil {
		log.Println("ERRO", err)
	}

	for _, w := range listFinalMesmo {
		_, err := f1.WriteString(fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", w.Text, w.Lema, w.Tag, w.Extra, w.Informative, w.Encyclopedic, w.Spoken, w.Prose, w.Science, w.ScienceCom, w.FreqTotal, w.FreqAvg, w.DP, w.DPNorm, w.FreqCorr, w.FreqCorrPer, w.FreqCorrAcu, w.Ranking))
		if err != nil {
			log.Println("ERRO", err)
		}
	}

}

func getTagWeight(tag string) int {
	//Artigo > Preposição > Pronome > Numeral > Substantivo > Adjetivo > Verbo > Advérbio

	ret := 0
	if strings.HasPrefix(tag, "ADV") {
		ret = 1
	} else if strings.HasPrefix(tag, "V") {
		ret = 2
	} else if strings.HasPrefix(tag, "A") {
		ret = 3
	} else if strings.HasPrefix(tag, "N") {
		ret = 4
	} else if strings.HasPrefix(tag, "DET+Num") {
		ret = 5
	} else if strings.HasPrefix(tag, "PRO") {
		ret = 6
	} else if strings.HasPrefix(tag, "PREPXPRO") {
		ret = 8
	} else if strings.HasPrefix(tag, "PREPXDET") {
		ret = 9
	} else if strings.HasPrefix(tag, "PREP") {
		ret = 7
	} else if strings.HasPrefix(tag, "DET+Art") {
		ret = 10
	}
	return ret
}

func min(list []float64) float64 {
	var ret float64 = list[0]
	for _, v := range list {
		if v < ret {
			ret = v
		}
	}
	return ret
}

//---------------------------------------------------------------------

// type Word struct {
// 	Text string
// 	Freq int
// }

// func main() {

// 	listInformative := []Word{}
// 	listEncyclopedic := []Word{}
// 	listSpokenTemp := []Word{}
// 	listSpoken := []Word{}
// 	listProse := []Word{}
// 	listScience := []Word{}
// 	listScienceCom := []Word{}

// 	path := "/home/sidleal/usp/semestre5/propor_dp/lista_freq_nathan.csv"
// 	data := readFile(path)

// 	lines := strings.Split(data, "\n")
// 	for i, line := range lines {
// 		cols := strings.Split(line, ",")
// 		if i == 0 {
// 			log.Println(line)
// 			log.Println("-------------")
// 			log.Println(cols[6], cols[7], cols[14])
// 			log.Println("-------------")
// 			log.Println(cols[17])
// 			log.Println("-------------")
// 			log.Println(cols[16])
// 			log.Println("-------------")
// 			log.Println(cols[9], cols[11])
// 			log.Println("-------------")
// 			log.Println(cols[12], cols[2])
// 			log.Println("-------------")
// 			log.Println(cols[4])
// 			continue
// 		}

// 		log.Println(i, cols[0])

// 		a, _ := strconv.Atoi(cols[6])
// 		b, _ := strconv.Atoi(cols[7])
// 		c, _ := strconv.Atoi(cols[14])
// 		listInformative = append(listInformative, Word{cols[0], a + b + c})

// 		a, _ = strconv.Atoi(cols[17])
// 		listEncyclopedic = append(listEncyclopedic, Word{cols[0], a})

// 		a, _ = strconv.Atoi(cols[16])
// 		listSpokenTemp = append(listSpokenTemp, Word{cols[0], a})

// 		a, _ = strconv.Atoi(cols[9])
// 		b, _ = strconv.Atoi(cols[11])
// 		listProse = append(listProse, Word{cols[0], a + b})

// 		a, _ = strconv.Atoi(cols[12])
// 		b, _ = strconv.Atoi(cols[2])
// 		listScience = append(listScience, Word{cols[0], a + b})

// 		a, _ = strconv.Atoi(cols[4])
// 		listScienceCom = append(listScienceCom, Word{cols[0], a})

// 	}

// 	path = "/home/sidleal/usp/semestre5/propor_dp/SUBTLEX-BRPOR.no-web.valid-char.csv"
// 	data = readFile(path)

// 	lines = strings.Split(data, "\n")
// 	for i, line := range lines {
// 		cols := strings.Split(line, "\t")
// 		if i == 0 {
// 			log.Println(line)
// 			continue
// 		}

// 		if i == 10000 {
// 			break
// 		}

// 		log.Println(i, cols[0])

// 		a, _ := strconv.Atoi(cols[1])

// 		found := false
// 		for _, word := range listSpokenTemp {
// 			if word.Text == cols[0] {
// 				found = true
// 				listSpoken = append(listSpoken, Word{word.Text, word.Freq + a})
// 			}
// 		}
// 		if !found {
// 			listSpoken = append(listSpoken, Word{cols[0], a})
// 		}

// 	}

// 	log.Println("inf")
// 	for i, word := range listInformative {
// 		if i < 10 {
// 			log.Println(word)
// 		}
// 	}

// 	log.Println("enc")
// 	for i, word := range listEncyclopedic {
// 		if i < 10 {
// 			log.Println(word)
// 		}
// 	}

// 	log.Println("spk")
// 	for i, word := range listSpoken {
// 		if i < 10 {
// 			log.Println(word)
// 		}
// 	}

// 	log.Println("prs")
// 	for i, word := range listProse {
// 		if i < 10 {
// 			log.Println(word)
// 		}
// 	}

// 	log.Println("sci")
// 	for i, word := range listScience {
// 		if i < 10 {
// 			log.Println(word)
// 		}
// 	}

// 	log.Println("scicom")
// 	for i, word := range listScienceCom {
// 		if i < 10 {
// 			log.Println(word)
// 		}
// 	}

// }

func readFile(path string) string {
	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}

	// r := charmap.ISO8859_1.NewDecoder().Reader(f)

	ret := ""

	buf := make([]byte, 32*1024)
	for {
		// n, err := r.Read(buf)
		n, err := f.Read(buf)
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

func readFileLines(path string) []string {
	ret := []string{}

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ret = append(ret, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return ret

}

// ------------------------------------

func main() {

	list := map[string]Word{}

	path := "/home/sidleal/usp/semestre5/propor_dp/lista_freq_nathan.csv"
	data := readFile(path)

	lines := strings.Split(data, "\n")
	for i, line := range lines {
		cols := strings.Split(line, ",")
		if i == 0 {
			log.Println(line)
			log.Println("-------------")
			log.Println(cols[6], cols[7], cols[14])
			log.Println("-------------")
			log.Println(cols[17])
			log.Println("-------------")
			log.Println(cols[16])
			log.Println("-------------")
			log.Println(cols[9], cols[11])
			log.Println("-------------")
			log.Println(cols[12], cols[2])
			log.Println("-------------")
			log.Println(cols[4])
			continue
		}

		// log.Println(i, cols[0])

		word := Word{}
		word.Text = cols[0]

		a, _ := strconv.Atoi(cols[6])
		b, _ := strconv.Atoi(cols[7])
		c, _ := strconv.Atoi(cols[14])
		word.Informative = a + b + c

		a, _ = strconv.Atoi(cols[17])
		word.Encyclopedic = a

		a, _ = strconv.Atoi(cols[16])
		word.Spoken = a

		a, _ = strconv.Atoi(cols[9])
		b, _ = strconv.Atoi(cols[11])
		word.Prose = a + b

		a, _ = strconv.Atoi(cols[12])
		b, _ = strconv.Atoi(cols[2])
		word.Science = a + b

		a, _ = strconv.Atoi(cols[4])
		word.ScienceCom = a

		if (word.Informative + word.Encyclopedic + word.Spoken + word.Prose + word.Science + word.ScienceCom) > 0 {
			list[word.Text] = word
		}
	}

	path = "/home/sidleal/usp/semestre5/propor_dp/SUBTLEX-BRPOR.no-web.valid-char.csv"
	data = readFile(path)

	lines = strings.Split(data, "\n")
	for i, line := range lines {
		cols := strings.Split(line, "\t")
		if i == 0 {
			log.Println(line)
			continue
		}

		// log.Println(i, cols[0])

		a, _ := strconv.Atoi(cols[1])

		if word, found := list[cols[0]]; found {
			word.Spoken = word.Spoken + a
			list[word.Text] = word
		} else {
			word := Word{}
			word.Text = cols[0]
			word.Spoken = a
			list[word.Text] = word
		}

	}

	// log.Println("------------------")
	// i := 0
	// for k, v := range list {
	// 	if i < 10 {
	// 		log.Println(k, v)
	// 	}
	// 	i++
	// }

	log.Println("carregando delaf")
	delaf := map[string]string{}

	regEx1 := regexp.MustCompile(`(.+)\.(.+)`)

	path = "/home/sidleal/usp/semestre5/propor_dp/DELAF_PB_2018.dic"
	lines = readFileLines(path)
	for _, line := range lines {
		if line == "" {
			continue
		}
		cols := strings.Split(line, ",")

		info := cols[1]
		lema := ""
		tag := ""
		match := regEx1.FindStringSubmatch(info)
		if len(match) > 1 {
			lema = match[1]
			tag = match[2]
			info = lema + "|" + tag
		}

		if delafDesc, found := delaf[cols[0]]; found {
			tokens := strings.Split(delafDesc, "|")
			lemaB := tokens[0]
			tagB := tokens[1]

			pesoA := getTagWeight(tag)
			pesoB := getTagWeight(tagB)

			if pesoA < pesoB {
				info = lemaB + "|" + tagB + "," + tag
			} else {
				info = lema + "|" + tag + "," + tagB
			}
			// log.Println(info)
		}
		delaf[cols[0]] = info
	}

	log.Println("delaf em memoria")

	exclude := map[string]bool{}
	for _, l := range []string{"b", "c", "d", "f", "g", "h", "i", "j", "k", "l", "m", "n", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"} {
		exclude[l] = true
	}

	listAll := map[string]Word{}

	// i := 0
	// regEx := regexp.MustCompile(`N\+Pr`)
	regEx := regexp.MustCompile(`(PREP|PRO)X(DET|PRO|PREP|ADV).*`)
	for k, v := range list {
		// if i < 100 {

		if _, found := exclude[k]; found {
			continue
		}

		if delafDesc, found := delaf[k]; found {
			// log.Println(delafDesc)
			tokens := strings.Split(delafDesc, "|")
			lema := tokens[0]
			tag := tokens[1]
			if !strings.HasPrefix(tag, "N+Pr") {
				v.Lema = lema
				v.Tag = tag

				match := regEx.FindStringSubmatch(tag)
				if len(match) > 1 {
					v.Extra = "contração"
				}
				listAll[k] = v
			}
			// matches := regEx.FindAllStringSubmatch(delafDesc, -1)
			// if len(matches) == 0 {
			// listAll[k] = v
			// }
		}

		// }
		// i++
	}

	//aplica correções manuais
	path = "/home/sidleal/usp/semestre5/propor_dp/anotacao1.csv"
	data = readFile(path)
	lines = strings.Split(data, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		cols := strings.Split(line, ",")

		word := cols[0]
		action := cols[1]
		extra := cols[2]

		if action == "EP" || action == "EE" {
			delete(listAll, word)

		} else if action == "CJ" {
			w1 := listAll[word]
			w2 := listAll[extra]

			w2.Informative = w2.Informative + w1.Informative
			w2.Encyclopedic = w2.Encyclopedic + w1.Encyclopedic
			w2.Spoken = w2.Spoken + w1.Spoken
			w2.Prose = w2.Prose + w1.Prose
			w2.Science = w2.Science + w1.Science
			w2.ScienceCom = w2.ScienceCom + w1.ScienceCom

			if w2.Text == "" {
				if delafDesc, found := delaf[extra]; found {
					tokens := strings.Split(delafDesc, "|")
					w2.Lema = tokens[0]
					w2.Tag = tokens[1]
					w2.Text = extra
				} else {
					w2.Lema = extra
					w2.Text = extra
					w2.Tag = w1.Tag
				}
			}

			listAll[extra] = w2
			delete(listAll, word)

			// log.Println(extra, w2.Text)

		} else {
			w := listAll[word]
			w.Extra = action
			listAll[word] = w
		}

		// log.Println("----", cols[0], cols[1], cols[2])

	}

	listAllLema1 := map[string]Word{}

	for k, v := range listAll {
		// log.Println(k, v)

		tokens := strings.Split(k, "-")
		if len(tokens) > 1 {
			if w, found := listAllLema1[tokens[0]]; found {
				w.Informative += v.Informative
				w.Encyclopedic += v.Encyclopedic
				w.Spoken += v.Spoken
				w.Prose += v.Prose
				w.Science += v.Science
				w.ScienceCom += v.ScienceCom
				listAllLema1[tokens[0]] = w
			} else {
				v.Text = tokens[0]
				listAllLema1[k] = v
			}

			if w, found := listAllLema1[tokens[1]]; found {
				w.Informative += v.Informative
				w.Encyclopedic += v.Encyclopedic
				w.Spoken += v.Spoken
				w.Prose += v.Prose
				w.Science += v.Science
				w.ScienceCom += v.ScienceCom
				listAllLema1[tokens[1]] = w
			} else {
				v.Text = tokens[1]
				listAllLema1[k] = v
			}

		} else {
			if w, found := listAllLema1[k]; found {
				w.Informative += v.Informative
				w.Encyclopedic += v.Encyclopedic
				w.Spoken += v.Spoken
				w.Prose += v.Prose
				w.Science += v.Science
				w.ScienceCom += v.ScienceCom
				listAllLema1[k] = w
			} else {
				listAllLema1[k] = v
			}

		}

	}

	listAllLema := map[string]Word{}

	for _, v := range listAllLema1 {
		// log.Println(k, v)

		if lema, found := listAllLema[v.Lema]; found {
			lema.Informative += v.Informative
			lema.Encyclopedic += v.Encyclopedic
			lema.Spoken += v.Spoken
			lema.Prose += v.Prose
			lema.Science += v.Science
			lema.ScienceCom += v.ScienceCom

			listAllLema[v.Lema] = lema

		} else {
			v.Text = v.Lema
			listAllLema[v.Lema] = v
		}

	}

	totInformative := 0
	totEncyclopedic := 0
	totSpoken := 0
	totProse := 0
	totScience := 0
	totScienceCom := 0

	for _, v := range listAllLema {
		// log.Println(k, v)

		totInformative += v.Informative
		totEncyclopedic += v.Encyclopedic
		totSpoken += v.Spoken
		totProse += v.Prose
		totScience += v.Science
		totScienceCom += v.ScienceCom

	}

	log.Println("Tots", totInformative, totEncyclopedic, totSpoken, totProse, totScience, totScienceCom)

	totAll := totInformative + totEncyclopedic + totSpoken + totProse + totScience + totScienceCom

	perInformative := float64(totInformative) / float64(totAll)
	perEncyclopedic := float64(totEncyclopedic) / float64(totAll)
	perSpoken := float64(totSpoken) / float64(totAll)
	perProse := float64(totProse) / float64(totAll)
	perScience := float64(totScience) / float64(totAll)
	perScienceCom := float64(totScienceCom) / float64(totAll)

	log.Println("Perc", perInformative, perEncyclopedic, perSpoken, perProse, perScience, perScienceCom)

	s := []float64{perInformative, perEncyclopedic, perSpoken, perProse, perScience, perScienceCom}

	// listFinal := map[string]Word{}
	listFinal := []Word{}

	var totFreqCorr float64 = 0

	for _, w := range listAllLema {
		// log.Println(k, w)

		f := float64(w.Informative + w.Encyclopedic + w.Spoken + w.Prose + w.Science + w.ScienceCom)
		v := []float64{float64(w.Informative), float64(w.Encyclopedic), float64(w.Spoken), float64(w.Prose), float64(w.Science), float64(w.ScienceCom)}
		// Deviation of proportions DP

		div := []float64{}
		for i, _ := range v {
			div = append(div, v[i]/f)
		}

		sub := []float64{}
		for i, _ := range div {
			sub = append(sub, div[i]-s[i])
		}

		sum := 0.0
		for i, _ := range sub {
			sum += math.Abs(sub[i])
		}

		dp := sum / 2
		dp_norm := dp / (1 - min(s))

		w.FreqTotal = int(f)
		w.DP = dp
		w.DPNorm = dp_norm
		w.FreqAvg = float64(w.FreqTotal) / float64(len(v))
		w.FreqCorr = w.FreqAvg*(1-w.DPNorm) + 0.0001

		totFreqCorr += w.FreqCorr

		listFinal = append(listFinal, w)

		// log.Println("------------>", f, dp, dp_norm)

	}

	sort.Sort(WordOrder(listFinal))

	listFinalMesmo := []Word{}
	var freqCorrAcum float64 = 0
	for i, w := range listFinal {
		w.Ranking = i + 1
		w.FreqCorrPer = w.FreqCorr / totFreqCorr * 100
		freqCorrAcum += w.FreqCorrPer
		w.FreqCorrAcu = freqCorrAcum
		listFinalMesmo = append(listFinalMesmo, w)
	}

	f1, err := os.Create("/home/sidleal/usp/semestre5/propor_dp/dispersion7_lema.tsv")
	if err != nil {
		log.Println("ERRO", err)
	}

	defer f1.Close()

	_, err = f1.WriteString("word\tlema\ttag\textra\tinformative\tencyclopedic\tspoken\tprose\tscience\tscience_com\ttotal\taverage\tdp\tdp_norm\tfreq_corr\tfreq_perc\tfreq_acum\tranking\n")
	if err != nil {
		log.Println("ERRO", err)
	}

	for _, w := range listFinalMesmo {
		_, err := f1.WriteString(fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", w.Text, w.Lema, w.Tag, w.Extra, w.Informative, w.Encyclopedic, w.Spoken, w.Prose, w.Science, w.ScienceCom, w.FreqTotal, w.FreqAvg, w.DP, w.DPNorm, w.FreqCorr, w.FreqCorrPer, w.FreqCorrAcu, w.Ranking))
		if err != nil {
			log.Println("ERRO", err)
		}
	}

}
