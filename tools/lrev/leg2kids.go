package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// func main() {

// 	log.Println("starting")

// 	f3, err := os.Create("/home/sidleal/sid/usp/lrev/leg2kids_recorte.tsv")
// 	if err != nil {
// 		log.Println("ERRO", err)
// 	}
// 	defer f3.Close()

// 	_, err = f3.WriteString("text_id\tfile\ttext\n")
// 	if err != nil {
// 		log.Println("ERRO", err)
// 	}

// 	path := "/home/sidleal/sid/usp/lrev"

// 	indexData := readFile(path+"/lef2kids_txt_tokenized_tokens.csv", false)
// 	lines := strings.Split(indexData, "\n")
// 	for i, line := range lines {
// 		if i == 0 || line == "" {
// 			continue
// 		}
// 		line := strings.ReplaceAll(line, "\r", "")
// 		tokens := strings.Split(line, ",")
// 		fileName := tokens[0]
// 		size, _ := strconv.Atoi(tokens[1])

// 		if size > 100 && i < 7500 {
// 			text := readFile(path+"/leg2kids/"+fileName, false)
// 			text = strings.ReplaceAll(text, "\r", "")
// 			text = strings.ReplaceAll(text, "\n", "<br>")
// 			_, err = f3.WriteString(fmt.Sprintf("%v\t%v\t%v\n", i, fileName, text))
// 			if err != nil {
// 				log.Println("ERRO", err)
// 			}

// 			log.Println(fileName, size)
// 		}

// 	}

// }

type LegFile struct {
	Name      string
	QtyTokens int
}

type SizeOrder []LegFile

func (a SizeOrder) Len() int      { return len(a) }
func (a SizeOrder) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SizeOrder) Less(i, j int) bool {
	return a[i].QtyTokens < a[j].QtyTokens
}

var t = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

func main_one_recorte() {

	log.Println("starting")

	f3, err := os.Create("/home/sidleal/sid/usp/lrev/leg2kids_recorte_v5.tsv")
	if err != nil {
		log.Println("ERRO", err)
	}
	defer f3.Close()

	_, err = f3.WriteString("text_id\tfile\tqty_tokens\ttext\n")
	if err != nil {
		log.Println("ERRO", err)
	}

	regex := regexp.MustCompile(`[^ \n]+`)
	regexBR := regexp.MustCompile(`([^\.\?\!])<br>`)

	path := "/home/sidleal/sid/usp/lrev"

	mapMovies := map[string]string{}
	metadata := readFile(path+"/leg2kids_export.tsv", false)
	lines := strings.Split(metadata, "\n")
	for i, line := range lines {
		if i == 0 || line == "" {
			continue
		}
		tokens := strings.Split(line, "\t")
		fileName := tokens[1] + ".txt"
		movieName := tokens[7]
		mapMovies[fileName] = movieName
	}

	legFiles := []LegFile{}

	files := listAllSubTxts(path + "/leg2kids")
	for _, f := range files {
		log.Println(f)

		fileTokens := strings.Split(f, "/")
		fileName := fileTokens[len(fileTokens)-1]
		text := readFile(f, false)

		matches := regex.FindAllStringIndex(text, -1)
		qtyTokens := len(matches)

		legFiles = append(legFiles, LegFile{fileName, qtyTokens})

	}

	sort.Sort(SizeOrder(legFiles))

	selectedMovies := map[string]string{}
	r := 0
	lastTexts := []string{}
	for i, f := range legFiles {
		log.Println(i, f.Name, f.QtyTokens)
		if f.QtyTokens > 600 && r < 7150 {
			text := readFile(path+"/leg2kids/"+f.Name, false)
			text = strings.ReplaceAll(text, "\r", "")
			text = strings.ReplaceAll(text, "\n", "<br>")

			text = regexBR.ReplaceAllString(text, "$1 ")

			movieTitle := mapMovies[f.Name]
			if _, found := selectedMovies[movieTitle]; found {
				continue
			}

			repetead := false
			for _, lastText := range lastTexts {
				ts := cleanToCompare(text)
				lts := cleanToCompare(lastText)
				if ts == lts {
					repetead = true
					break
				}
			}
			if repetead {
				continue
			}

			lastTexts = append(lastTexts, text)
			if len(lastTexts) > 10 {
				lastTexts = lastTexts[1:11]
			}

			selectedMovies[movieTitle] = f.Name
			r++
			_, err = f3.WriteString(fmt.Sprintf("%v\t%v\t%v\t%v\n", r, f.Name, f.QtyTokens, text))
			if err != nil {
				log.Println("ERRO", err)
			}
		}

	}

}

func cleanToCompare(text string) string {
	ts, _, _ := transform.String(t, text)
	ret := strings.ToLower(ts)
	ret = strings.ReplaceAll(ret, " ", "")
	ret = strings.ReplaceAll(ret, "<br>", "")
	ret = strings.ReplaceAll(ret, "y", "i")
	ret = strings.ReplaceAll(ret, "w", "v")
	ret = strings.ReplaceAll(ret, "k", "c")
	return ret
}

func main() { //_leg_xtract() { //

	log.Println("starting")

	fout, err := os.OpenFile("/home/sidleal/sid/usp/lrev/leg2kids_metrics_last4.tsv", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERRO", err)
	}
	defer fout.Close()

	processed := readFile("/home/sidleal/sid/usp/lrev/leg2kids_metrics.tsv", false)
	mapAlreadyProcessed := map[string]int{}
	lines := strings.Split(processed, "\n")
	for i, line := range lines {
		if line == "" || i == 0 {
			continue
		}
		tokens := strings.Split(line, "\t")
		mapAlreadyProcessed[tokens[0]] = 1
	}

	data := readFile("/home/sidleal/sid/usp/lrev/leg2kids_recorte_v5_clean.tsv", false)

	lines = strings.Split(data, "\n")
	for i, line := range lines {
		if line == "" || i == 0 {
			continue
		}

		// if i < 0 || i > 3500 {
		// 	continue
		// }
		// if i < 3500 || i > 8000 {
		// 	continue
		// }
		tokens := strings.Split(line, "\t")
		textID := tokens[0]
		fileName := tokens[1]
		rawText := tokens[3]

		if _, found := mapAlreadyProcessed[textID]; found {
			continue
		}

		rawText = strings.ReplaceAll(rawText, "\\N", "<br>")
		rawText = strings.ReplaceAll(rawText, "\r", "")

		text := strings.ReplaceAll(rawText, "<br>", "\n")
		text = cleanText(text)

		log.Println("-------------------------------------------------------")
		log.Println("-->", i)
		fRet := callMetrix(text)
		log.Println(fRet)
		log.Println("-->", i)

		feats := strings.Split(fRet, ",")

		header := "text_id\tfile\ttext\t"
		ret := fmt.Sprintf("%v\t%v\t%v\t", textID, fileName, rawText)
		for _, feat := range feats {
			kv := strings.Split(feat, ":")
			if len(kv) > 1 {
				if i == 1 {
					header += kv[0] + "\t"
				}
				ret += kv[1] + "\t"
			}
		}

		if i == 1 {
			header = strings.TrimRight(header, "\t")
			_, err := fout.WriteString(header + "\n")
			if err != nil {
				log.Println("ERRO", err)
			}
		}

		ret = strings.TrimRight(ret, "\t")
		_, err := fout.WriteString(ret + "\n")
		if err != nil {
			log.Println("ERRO", err)
		}

	}

}
