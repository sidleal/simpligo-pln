package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

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

var mapWords = map[string]int32{}

func main_clean() {

	carregarDELAF()

	var regEx = regexp.MustCompile(`<(.*)>`)

	path := "/home/sidleal/sid/usp/brWaC/"

	file, err := os.Open(path + "brwac.vert")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	i := 0
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			break
		}

		// log.Print(line)

		matchContent := regEx.FindStringSubmatch(line)

		if len(matchContent) < 1 {
			w := strings.ToLower(line)
			w = strings.ReplaceAll(w, "\n", "")
			if _, found := mapWords[w]; !found {
				mapWords[w] = 0
			}
			mapWords[w]++
			if i%500000 == 0 {
				fmt.Println(i)
			}
		}
		i++
		// if i > 1000 {
		// 	panic("quit")
		// }
	}

	f, err := os.Create(fmt.Sprintf("%v/lista.tsv", path))
	if err != nil {
		log.Println("ERRO", err)
	}

	for w, qty := range mapWords {
		fmt.Println(w, qty)
		if info, found := delaf[w]; found {
			tokens := strings.Split(info, "|")
			lema := tokens[0]
			tag := tokens[1]
			_, err = f.WriteString(fmt.Sprintf("%s\t%s\t%v\t%s\n", w, lema, qty, tag))
			if err != nil {
				log.Println("ERRO", err)
			}
		}

	}
	fmt.Println("total", len(mapWords))

	f.Close()
}

func mainAnnotation() {

	var regEx = regexp.MustCompile(`<(.*)>`)

	path := "/home/sidleal/sid/usp/brWaC/brwac-readability-annotation.csv"

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	i := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		matchContent := regEx.FindStringSubmatch(line)

		fmt.Print(line + " ")
		if len(matchContent) < 1 {
			// fmt.Print(line + " ")
		}

		if line == "</doc>" {
			fmt.Println("\n------------------------------------")
		}
		i++
		if i > 20 {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
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

func main_conll() {

	// carregarDELAF()

	path := "/home/sidleal/sid/usp/brWaC/"

	file, err := os.Open(path + "brwac.conll")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	i := 0
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			break
		}

		log.Print(line)
		tokens := strings.Split(line, "\t")
		if len(tokens) < 4 {
			continue
		}

		w := tokens[1]
		lemma := tokens[2]
		tag := tokens[3]

		// log.Println(w, lemma, tag)

		if tag != "PROP" {
			w = strings.ToLower(w)
		}
		w = strings.ReplaceAll(w, "\n", "")

		uniqueW := fmt.Sprintf("%v_%v_%v", w, lemma, tag)
		if _, found := mapWords[uniqueW]; !found {
			mapWords[uniqueW] = 0
		}
		mapWords[uniqueW]++

		if i%500000 == 0 {
			fmt.Println(i)
		}

		i++
		// if i > 1000 {
		// 	panic("quit")
		// }
	}

	f, err := os.Create(fmt.Sprintf("%v/lista_tag.tsv", path))
	if err != nil {
		log.Println("ERRO", err)
	}

	allowedTags := map[string]int{"ADJ": 1, "ALT": 1, "COND": 1, "DET": 1, "EC": 1, "GER": 1, "IMP": 1, "IN": 1, "IND": 1, "INF": 1, "KC": 1, "KS": 1, "N": 1, "PCP": 1, "PERS": 1, "PP": 1, "PR": 1, "PRON": 1, "PROP": 1, "PRP": 1, "PS": 1, "SPEC": 1, "V": 1, "VFIN": 1}
	cleanedWordQty := 0
	for uniqueW, qty := range mapWords {
		fmt.Println(uniqueW, qty)
		tokens := strings.Split(uniqueW, "_")
		w := tokens[0]
		lema := tokens[1]
		tag := tokens[2]

		if _, found := allowedTags[tag]; !found {
			continue
		}

		minToKeep := 9
		if tag == "PROP" {
			minToKeep = 99
		}
		if qty > int32(minToKeep) && w != "" {
			_, err = f.WriteString(fmt.Sprintf("%s\t%s\t%v\t%s\n", w, lema, qty, tag))
			if err != nil {
				log.Println("ERRO", err)
			}
			cleanedWordQty++
		}

	}
	fmt.Println("total:", len(mapWords), "cleaned", cleanedWordQty)

	f.Close()

}

func main_ver() {

	var regEx = regexp.MustCompile(`<([a-z\/]+)( docid="|)([^"]+|)`)

	path := "/home/sidleal/sid/usp/brWaC/"

	file, err := os.Open(path + "brwac.vert")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	f, err := os.Create(fmt.Sprintf("%v/brWaC_all_texts.tsv", path))
	if err != nil {
		log.Println("ERRO", err)
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("%s\t%s\t%s\n", "idx", "text_id", "text"))
	if err != nil {
		log.Println("ERRO", err)
	}

	i := 0
	j := 0
	text := ""
	textID := ""
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			break
		}

		matchContent := regEx.FindStringSubmatch(line)

		if len(matchContent) < 1 {
			w := strings.ReplaceAll(line, "\n", "")
			w = strings.ReplaceAll(w, "\t", "")
			//fmt.Print(w, " ")
			text += w + " "
		} else if matchContent[1] == "doc" {
			if textID != "" {
				_, err = f.WriteString(fmt.Sprintf("%v\t%s\t%s\n", j, textID, text))
				if err != nil {
					log.Println("ERRO", err)
				}
			}
			textID = matchContent[3]
			log.Print(j, " ", textID, " ", i)
			text = ""
			j++
		}
		i++
	}
}

func main_parsed() {

	path := "/home/sidleal/sid/usp/brWaC/brwac.vert.parsed"

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	i := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		fmt.Println(line + " ")

		i++
		if i > 200 {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

func main_nlpnet_full() {

	// carregarDELAF()

	path := "/home/sidleal/sid/usp/brWaC/"

	file, err := os.Open(path + "all_texts_nlpnet.tsv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var regEx2 = regexp.MustCompile(`([A-zÀ-ú\-]+?)\_([A-z_\+\-]+)`)

	i := 0
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			break
		}

		if i == 0 || line == "" {
			i++
			continue
		}

		matchContent := regEx2.FindAllStringSubmatch(line, -1)
		for _, item := range matchContent {
			w := item[0]
			tag := item[2]
			// log.Println(w)

			if len(w) > 50 {
				continue
			}

			if tag != "NPROP" {
				w = strings.ToLower(w)
			}

			if _, found := mapWords[w]; !found {
				mapWords[w] = 0
			}
			mapWords[w]++
		}

		i++

		if i%10000 == 0 {
			fmt.Println(i)
			runtime.GC()
		}
	}

	f, err := os.Create(fmt.Sprintf("%v/lista_nlpnet.tsv", path))
	if err != nil {
		log.Println("ERRO", err)
	}

	//"adj":1, "adv":1, "adv-ks":1, "art":1, "cur":1, "in":1, "kc":1, "ks":1, "n":1, "NPROP":1, "num":1, "pcp":1, "pden":1, "prep":1, "prep+adv":1, "prep+art":1, "prep+pro-ks":1, "prep+proadj":1, "prep+propess":1, "prep+prosub":1, "pro-ks":1, "proadj":1, "propess":1, "prosub":1, "v":1
	cleanedWordQty := 0
	for uniqueW, qty := range mapWords {
		fmt.Println(uniqueW, qty)
		tokens := strings.Split(uniqueW, "_")
		w := tokens[0]
		tag := tokens[1]

		minToKeep := 9
		if tag == "NPROP" {
			minToKeep = 99
		}
		if qty > int32(minToKeep) && w != "" {
			_, err = f.WriteString(fmt.Sprintf("%s\t%v\t%s\n", w, qty, tag))
			if err != nil {
				log.Println("ERRO", err)
			}
			cleanedWordQty++
		}

	}
	fmt.Println("total:", len(mapWords), "cleaned", cleanedWordQty)

	f.Close()
}

func main_groups_full() {

	mapWordsGroups := map[string][57]int32{}

	//lista grupos
	mapTextGroups := map[string]int8{}
	groupData := readFile("/home/sidleal/sid/usp/freq_exps/result-kmeans-full.csv")
	groupLines := strings.Split(groupData, "\n")
	for _, line := range groupLines {
		if line == "" {
			continue
		}
		tokens := strings.Split(line, ",")
		grp, _ := strconv.Atoi(tokens[1])
		mapTextGroups[tokens[0]] = int8(grp)
	}

	//all data
	path := "/home/sidleal/sid/usp/brWaC/"

	file, err := os.Open(path + "all_texts_nlpnet.tsv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var regEx2 = regexp.MustCompile(`([A-zÀ-ú\-]+?)\_([A-z_\+\-]+)`)

	i := 0
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			break
		}

		if i == 0 || line == "" {
			i++
			continue
		}

		if i <= 2000000 {
			i++
			continue
		}

		tokens := strings.Split(line, "\t")

		textID := tokens[1]
		textGroupID := -1
		if grp, found := mapTextGroups[textID]; found {
			textGroupID = int(grp)
		}
		if textGroupID < 0 {
			continue
		}

		// if textGroupID == 38 {
		// 	log.Println(textID, textGroupID)
		// 	log.Println(tokens[2])
		// } else {
		// 	continue
		// }

		matchContent := regEx2.FindAllStringSubmatch(tokens[2], -1)
		for _, item := range matchContent {
			w := item[0]
			tag := item[2]
			// log.Println(w)

			if len(w) > 50 {
				continue
			}

			if tag != "NPROP" {
				w = strings.ToLower(w)
			}

			if _, found := mapWordsGroups[w]; !found {
				mapWordsGroups[w] = [57]int32{}
			}
			qts := mapWordsGroups[w]
			qts[textGroupID] = qts[textGroupID] + 1
			mapWordsGroups[w] = qts
		}

		i++

		if i%10000 == 0 {
			fmt.Println(i)
			runtime.GC()
		}

	}

	f, err := os.Create(fmt.Sprintf("%v/lista_nlpnet_groups_2b.tsv", path))
	if err != nil {
		log.Println("ERRO", err)
	}

	//"adj":1, "adv":1, "adv-ks":1, "art":1, "cur":1, "in":1, "kc":1, "ks":1, "n":1, "NPROP":1, "num":1, "pcp":1, "pden":1, "prep":1, "prep+adv":1, "prep+art":1, "prep+pro-ks":1, "prep+proadj":1, "prep+propess":1, "prep+prosub":1, "pro-ks":1, "proadj":1, "propess":1, "prosub":1, "v":1
	cleanedWordQty := 0
	for uniqueW, qts := range mapWordsGroups {
		fmt.Println(uniqueW, qts)
		tokens := strings.Split(uniqueW, "_")
		w := tokens[0]
		tag := tokens[1]

		minToKeep := 9
		if tag == "NPROP" {
			minToKeep = 99
		}
		qty := sumArray(qts)
		if qty > int32(minToKeep) && w != "" {
			line := fmt.Sprintf("%s\t%s", w, tag)
			for i := 0; i < 57; i++ {
				line += fmt.Sprintf("\t%v", qts[i])
			}
			_, err = f.WriteString(line + "\n")
			if err != nil {
				log.Println("ERRO", err)
			}
			cleanedWordQty++
		}

	}
	fmt.Println("total:", len(mapWordsGroups), "cleaned", cleanedWordQty)

	f.Close()
}

func sumArray(arr [57]int32) int32 {
	ret := int32(0)
	for _, i := range arr {
		ret += i
	}
	return ret
}

func mainMergeResults() {

	mapWordsGroups1 := map[string][57]int32{}

	data1 := readFile("/home/sidleal/sid/usp/brWaC/lista_nlpnet_groups_2a.tsv")
	// data2 := readFile("/home/sidleal/sid/usp/brWaC/lista_nlpnet_groups_2a.tsv")

	data1Lines := strings.Split(data1, "\n")
	for _, line := range data1Lines {
		if line == "" {
			continue
		}
		tokens := strings.Split(line, "\t")

		w := tokens[0]
		tag := tokens[1]
		key := fmt.Sprintf("%s_%s", w, tag)

		if _, found := mapWordsGroups1[key]; !found {
			qts := [57]int32{}
			for i := 0; i < 57; i++ {
				qty, _ := strconv.Atoi(tokens[i+2])
				qts[i] = int32(qty)
			}
			mapWordsGroups1[key] = qts
		}

	}

	log.Println("----------------> Data1 Loaded")

	mapWordsGroups2 := map[string][57]int32{}

	data2 := readFile("/home/sidleal/sid/usp/brWaC/lista_nlpnet_groups_2b.tsv")

	data2Lines := strings.Split(data2, "\n")
	for _, line := range data2Lines {
		if line == "" {
			continue
		}
		tokens := strings.Split(line, "\t")

		w := tokens[0]
		tag := tokens[1]
		key := fmt.Sprintf("%s_%s", w, tag)

		if _, found := mapWordsGroups2[key]; !found {
			qts := [57]int32{}
			for i := 0; i < 57; i++ {
				qty, _ := strconv.Atoi(tokens[i+2])
				qts[i] = int32(qty)
			}
			mapWordsGroups2[key] = qts
		}

	}

	log.Println("----------------> Data2 Loaded")

	mapWordsGroups := map[string][57]int32{}

	for k, v := range mapWordsGroups1 {
		// log.Println(k, v)
		if v2, found := mapWordsGroups2[k]; found {
			qts := [57]int32{}
			for i := 0; i < 57; i++ {
				qts[i] = v[i] + v2[i]
			}
			mapWordsGroups[k] = qts
		}
	}

	log.Println("----------------> Sum OK")

	for k, v := range mapWordsGroups1 {
		if _, found := mapWordsGroups[k]; !found {
			mapWordsGroups[k] = v
		}
	}

	log.Println("----------------> Words from G1 OK")

	for k, v := range mapWordsGroups2 {
		if _, found := mapWordsGroups[k]; !found {
			mapWordsGroups[k] = v
		}
	}

	log.Println("----------------> Words from G2 OK")

	// for k, v := range mapWordsGroups {
	// 	log.Println(k, v)
	// }

	log.Println(len(mapWordsGroups1), len(mapWordsGroups2), len(mapWordsGroups))

	f, err := os.Create("/home/sidleal/sid/usp/brWaC/lista_nlpnet_groups_full.tsv")
	if err != nil {
		log.Println("ERRO", err)
	}

	for uniqueW, qts := range mapWordsGroups {
		fmt.Println(uniqueW, qts)
		tokens := strings.Split(uniqueW, "_")
		w := tokens[0]
		tag := tokens[1]

		line := fmt.Sprintf("%s\t%s", w, tag)
		for i := 0; i < 57; i++ {
			line += fmt.Sprintf("\t%v", qts[i])
		}
		_, err = f.WriteString(line + "\n")
		if err != nil {
			log.Println("ERRO", err)
		}

	}
	fmt.Println("total:", len(mapWordsGroups))

	f.Close()

}
