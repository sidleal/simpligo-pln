package main

import (
	"log"
	"sort"
	"strconv"
	"strings"
)

// ------------------------------------

func mainresp() {

	// list := map[string]Word{}

	rawFiles := []string{"cloze_baJwn24BmRI7xu8F5xxd_data4.csv", "cloze_bqJxn24BmRI7xu8Fqhz6_data4.csv", "cloze_bKJwn24BmRI7xu8FSBzU_data4.csv"}

	path := "/home/sidleal/sid/usp/cloze_exps/"

	for _, rf := range rawFiles {
		data := readFile(path + rf)

		mapParticipants := map[string]map[int]int{}
		wordList := []Word{}

		lines := strings.Split(data, "\n")
		for i, line := range lines {
			if line == "" {
				continue
			}

			cols := strings.Split(line, ",")
			if i == 0 {
				log.Println(line)
			}

			word := Word{}
			word.Code = cols[0]
			word.Name = cols[1]
			word.QtdGenre = cols[2]
			word.Pars = cols[3]
			word.Part = cols[4]
			word.Email = cols[5]
			word.Age = cols[6]
			word.Gender = cols[7]
			word.Reg = cols[8]
			word.Sem = cols[9]
			word.ParsRead = cols[10]
			word.DTBegin = cols[11]
			word.HRBegin = cols[12]
			word.ParID, _ = strconv.Atoi(cols[13])
			word.SentID, _ = strconv.Atoi(cols[14])
			word.WordID, _ = strconv.Atoi(cols[15])
			word.Word = cols[16]
			word.Resp = cols[17]
			word.TBegin, _ = strconv.Atoi(cols[18])
			word.TDig, _ = strconv.Atoi(cols[19])
			word.TTot, _ = strconv.Atoi(cols[20])
			word.TPar, _ = strconv.Atoi(cols[21])
			word.TTest, _ = strconv.Atoi(cols[22])

			wordList = append(wordList, word)

			if _, found := mapParticipants[word.Part]; !found {
				mapParticipants[word.Part] = map[int]int{}
			}
			if _, found := mapParticipants[word.Part][word.ParID]; !found {
				mapParticipants[word.Part][word.ParID] = 0
			}
			mapParticipants[word.Part][word.ParID]++

		}

		sort.Sort(WordOrder(wordList))
		for _, w := range wordList {
			log.Println(w.ParID, w.SentID, w.Word, w.Resp)
		}

	}

}
