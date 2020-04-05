package main

import (
	"log"
	"strconv"
	"strings"
)

func main() {
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
