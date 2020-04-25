package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

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

func main() {

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

	mapWordsGroupsClean := map[string][57]int32{}
	mapWordsGroupsLixo := map[string][57]int32{}

	for k, v := range mapWordsGroups1 {
		w := strings.Split(k, "_")[0]
		tag := strings.Split(k, "_")[1]
		if _, found := delaf[w]; found || tag == "NPROP" {
			mapWordsGroupsClean[k] = v
		} else {
			mapWordsGroupsLixo[k] = v
		}
	}

	log.Println("----------------> clean ok")

	f, err := os.Create("/home/sidleal/sid/usp/brWaC/lista_nlpnet_groups_full_cleaned.tsv")
	if err != nil {
		log.Println("ERRO", err)
	}

	for uniqueW, qts := range mapWordsGroupsClean {
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
	f.Close()

	f2, err := os.Create("/home/sidleal/sid/usp/brWaC/lista_nlpnet_groups_full_excluded.tsv")
	if err != nil {
		log.Println("ERRO", err)
	}

	for uniqueW, qts := range mapWordsGroupsLixo {
		fmt.Println(uniqueW, qts)
		tokens := strings.Split(uniqueW, "_")
		w := tokens[0]
		tag := tokens[1]

		line := fmt.Sprintf("%s\t%s", w, tag)
		for i := 0; i < 57; i++ {
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
