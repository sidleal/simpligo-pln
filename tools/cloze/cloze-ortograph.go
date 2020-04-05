package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// ------------------------------------

func mainOrtogra() {

	carregarDELAF()

	path := "/home/sidleal/sid/usp/cloze_exps/"

	parsAll := readFile(path + "50pars_sid.txt")
	pars := strings.Split(parsAll, "\n")

	exportDate := "2020_03_19"
	inputFiles := []string{
		"cloze_puc_" + exportDate + ".csv",
		"cloze_usp_" + exportDate + ".csv",
		"cloze_ufc_" + exportDate + ".csv",
		"cloze_utfpr_" + exportDate + ".csv",
		"cloze_ufabc_" + exportDate + ".csv",
		"cloze_uerj_" + exportDate + ".csv",
	}

	f, err := os.Create(path + "cloze_erros_01.tsv")
	if err != nil {
		log.Println("ERRO", err)
	}
	defer f.Close()

	for i, rf := range inputFiles {
		data := readFile(path + rf)

		lines := strings.Split(data, "\n")
		for k, line := range lines {
			if k == 0 || line == "" {
				continue
			}
			// log.Println(k, line)
			tokens := strings.Split(line, ",")
			parID, _ := strconv.Atoi(tokens[11])
			wID, _ := strconv.Atoi(tokens[13])

			parText := pars[parID]
			target := tokens[14]
			resp := tokens[15]

			resp = strings.ToLower(resp)
			resp = strings.TrimSpace(resp)

			ortoErr := false
			if _, found := delaf[resp]; !found {
				ortoErr = true
			}

			excecoes := []string{"", "d'água", ".", "?"}
			for _, it := range excecoes {
				if resp == it {
					ortoErr = false
				}
			}

			var regEx = regexp.MustCompile(`[0-9,\.]+`)
			matchContent := regEx.FindStringSubmatch(resp)
			if len(matchContent) > 0 {
				ortoErr = false
			}

			if ortoErr {
				log.Println(i, k, parID, parText, target, resp)
				wTokens := strings.Split(parText+" _ _ _", " ")

				start := wID - 3
				if start < 0 {
					start = 0
				}

				if start+3 > len(wTokens) {
					continue
				}

				window := ""
				for j := start; j <= start+3; j++ {
					window += wTokens[j] + " "
				}
				window = "..." + window + "..."

				log.Println(i, k, parID, window, target, resp)

				_, err = f.WriteString(fmt.Sprintf("%v\t%v\t%v\t%v\t%v\n", parID, wID, window, target, resp))
				if err != nil {
					log.Println("ERRO", err)
				}
			}
		}
	}

}

var delaf = map[string]string{}

func carregarDELAF() {
	log.Println("carregando delaf")
	regEx1 := regexp.MustCompile(`(.+)\.(.+)`)

	path := "/home/sidleal/sid/usp/delaf/DELAF_PB_2018.dic"
	lines := readFileLines(path)
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
