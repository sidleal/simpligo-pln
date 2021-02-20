package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

//teste
func mainxxxx() {
	//lista grupos
	mapTextGroups := map[string]int{}
	groupData := readFile("/home/sidleal/sid/usp/freq_exps/result-kmeans-full.csv")
	groupLines := strings.Split(groupData, "\n")
	for _, line := range groupLines {
		if line == "" {
			continue
		}
		tokens := strings.Split(line, ",")
		grp, _ := strconv.Atoi(tokens[1])
		mapTextGroups[tokens[0]] = grp
	}

	textsPerGroup := [57]int{}

	for _, v := range mapTextGroups {
		textsPerGroup[v]++
	}

	log.Println(textsPerGroup)

}

type AnnotatedWord struct {
	Word       string
	Tag        string
	Annotation int
	Correction string
	Freq       int
}

//clean words
func main_clean_words_last() {

	carregarDELAF()

	mapWordsGroups1 := map[string][57]int32{}

	data := readFile("/home/sidleal/sid/usp/brWaC/lista_nlpnet_groups_full.tsv")

	dataLines := strings.Split(data, "\n")
	for _, line := range dataLines {
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

	mapAnnotatedWords := map[string]AnnotatedWord{}
	data2 := readFile("/home/sidleal/sid/usp/brWaC/palavras_excluidas_anotadas_v1_ate_freq_900.tsv")
	dataLines2 := strings.Split(data2, "\n")
	for i, line := range dataLines2 {
		if i == 0 || line == "" {
			continue
		}
		tokens := strings.Split(line, "\t")

		aw := AnnotatedWord{}

		aw.Word = tokens[0]
		aw.Annotation, _ = strconv.Atoi(tokens[1])
		aw.Correction = tokens[2]
		aw.Tag = tokens[3]
		aw.Freq, _ = strconv.Atoi(tokens[4])

		key := fmt.Sprintf("%s_%s", aw.Word, aw.Tag)
		mapAnnotatedWords[key] = aw

	}

	mapWordsGroupsClean := map[string][]int32{}
	mapWordsGroupsLixo := map[string][]int32{}

	for k, v := range mapWordsGroups1 {
		w := strings.Split(k, "_")[0]
		tag := strings.Split(k, "_")[1]

		hasHyphen := false
		if strings.Index(w, "-") > 1 {
			hasHyphen = true
			hypTokens := strings.Split(w, "-")
			for _, sub := range hypTokens {
				if _, found := delaf[sub]; !found {
					hasHyphen = false
				}
			}
		}

		_, foundInDELAF := delaf[w]

		newV := v[0:57]
		cleanAfterAnnotation := false
		if aw, found := mapAnnotatedWords[k]; found {
			newV = append(newV, int32(aw.Annotation))
			if aw.Annotation == 1 || aw.Annotation == 2 || aw.Annotation == 4 ||
				aw.Annotation == 7 || aw.Annotation == 10 {
				cleanAfterAnnotation = true
			}
		} else {
			newV = append(newV, int32(-1))
		}

		if foundInDELAF || tag == "NPROP" || hasHyphen || cleanAfterAnnotation {
			mapWordsGroupsClean[k] = newV
		} else {
			mapWordsGroupsLixo[k] = newV
		}
	}

	log.Println("----------------> clean ok")

	f, err := os.Create("/home/sidleal/sid/usp/brWaC/lista_nlpnet_groups_full_cleaned03.tsv")
	if err != nil {
		log.Println("ERRO", err)
	}

	for uniqueW, qts := range mapWordsGroupsClean {
		fmt.Println(uniqueW, qts)
		tokens := strings.Split(uniqueW, "_")
		w := tokens[0]
		tag := tokens[1]

		line := fmt.Sprintf("%s\t%s", w, tag)
		for i := 0; i < 58; i++ {
			line += fmt.Sprintf("\t%v", qts[i])
		}
		_, err = f.WriteString(line + "\n")
		if err != nil {
			log.Println("ERRO", err)
		}

	}
	f.Close()

	f2, err := os.Create("/home/sidleal/sid/usp/brWaC/lista_nlpnet_groups_full_excluded03.tsv")
	if err != nil {
		log.Println("ERRO", err)
	}

	for uniqueW, qts := range mapWordsGroupsLixo {
		fmt.Println(uniqueW, qts)
		tokens := strings.Split(uniqueW, "_")
		w := tokens[0]
		tag := tokens[1]

		line := fmt.Sprintf("%s\t%s", w, tag)
		for i := 0; i < 58; i++ {
			line += fmt.Sprintf("\t%v", qts[i])
		}
		_, err = f2.WriteString(line + "\n")
		if err != nil {
			log.Println("ERRO", err)
		}

	}
	f2.Close()

	log.Println(len(mapWordsGroups1), len(mapWordsGroupsClean), len(mapWordsGroupsLixo))

}

//tamanhos dos grupos
func mainGroupSizeStats() {

	mapGroupsStats := map[int]map[string]int{}

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

	f1, err := os.Create(fmt.Sprintf("%v/groups_text_sizes.tsv", path))
	if err != nil {
		log.Println("ERRO", err)
	}

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

		// if i <= 2000000 {
		// 	i++
		// 	continue
		// }

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

		if _, found := mapGroupsStats[textGroupID]; !found {
			mapGroupsStats[textGroupID] = map[string]int{}
			mapGroupsStats[textGroupID]["total"] = 0
			mapGroupsStats[textGroupID]["min"] = 1000000
			mapGroupsStats[textGroupID]["max"] = 0
			mapGroupsStats[textGroupID]["qtd"] = 0
		}

		totTokens := len(matchContent)
		mapGroupsStats[textGroupID]["total"] += totTokens

		if totTokens < mapGroupsStats[textGroupID]["min"] {
			mapGroupsStats[textGroupID]["min"] = totTokens
		}

		if totTokens > mapGroupsStats[textGroupID]["max"] {
			mapGroupsStats[textGroupID]["max"] = totTokens
		}

		mapGroupsStats[textGroupID]["qtd"]++

		i++

		_, err = f1.WriteString(fmt.Sprintf("%v\t%v\t%v\n", textGroupID, textID, totTokens))
		if err != nil {
			log.Println("ERRO", err)
		}

		if i%10000 == 0 {
			fmt.Println(i)
			runtime.GC()
		}

	}
	f1.Close()

	f, err := os.Create(fmt.Sprintf("%v/stats_groups.tsv", path))
	if err != nil {
		log.Println("ERRO", err)
	}

	for grp, val := range mapGroupsStats {
		fmt.Println(grp, val)
		line := fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v", grp, val["total"], val["min"], val["max"], val["qtd"], val["total"]/val["qtd"])
		_, err = f.WriteString(line + "\n")
		if err != nil {
			log.Println("ERRO", err)
		}

	}

	f.Close()
}

type GroupTextStat struct {
	GroupID string
	TextID  string
	Qty     int
}

type GroupTextStatOrder []GroupTextStat

func (a GroupTextStatOrder) Len() int      { return len(a) }
func (a GroupTextStatOrder) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a GroupTextStatOrder) Less(i, j int) bool {
	return a[i].Qty < a[j].Qty
}

// amostra de cada grupo
func mainSamples() {

	mapGroupStats := map[string][]GroupTextStat{}
	mapGroupStatsSelected := map[string]GroupTextStat{}

	//all data
	path := "/home/sidleal/sid/usp/brWaC/"

	file, err := os.Open(path + "groups_text_sizes.tsv")
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

		if i == 0 || line == "" {
			i++
			continue
		}

		// if i <= 200 {
		// 	i++
		// 	continue
		// }

		line = strings.ReplaceAll(line, "\n", "")

		tokens := strings.Split(line, "\t")
		// log.Println(tokens)
		textID := tokens[1]
		textGroupID := tokens[0]
		qty, _ := strconv.Atoi(tokens[2])

		if _, found := mapGroupStats[textGroupID]; !found {
			mapGroupStats[textGroupID] = []GroupTextStat{}
		}

		mapGroupStats[textGroupID] = append(mapGroupStats[textGroupID], GroupTextStat{textGroupID, textID, qty})
		i++

		if i%10000 == 0 {
			fmt.Println(i)
			runtime.GC()
		}

	}

	for _, val := range mapGroupStats {
		sort.Sort(GroupTextStatOrder(val))
		for i, v := range val {
			// log.Println("------", i, grp, v.TextID, v.Qty)
			mapGroupStatsSelected[v.TextID] = v

			if i > 30 {
				break
			}
		}

	}

	for tID, val := range mapGroupStatsSelected {
		log.Println("------", i, val.GroupID, tID, val.Qty)

	}

	file2, err := os.Open(path + "all_texts_nlpnet.tsv")
	if err != nil {
		log.Fatal(err)
	}
	defer file2.Close()

	var regEx2 = regexp.MustCompile(`([A-zÀ-ú\-]+?)\_([A-z_\+\-]+)`)

	f1, err := os.Create(fmt.Sprintf("%v/groups_top_10.tsv", path))
	if err != nil {
		log.Println("ERRO", err)
	}

	i = 0
	reader = bufio.NewReader(file2)
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

		tokens := strings.Split(line, "\t")
		textID := tokens[1]

		if val, found := mapGroupStatsSelected[textID]; found {
			matchContent := regEx2.FindAllStringSubmatch(tokens[2], -1)
			text := ""
			for _, w := range matchContent {
				text += w[1] + " "
			}

			log.Println("-----------> ", val, text)
			line := fmt.Sprintf("%v\t%v\t%v\t%v\n", val.GroupID, val.TextID, val.Qty, text)
			_, err = f1.WriteString(line)
			if err != nil {
				log.Println("ERRO", err)
			}
		}
	}

	f1.Close()

}
