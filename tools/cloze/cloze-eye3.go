package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// ------------------------------------

type ClozeData struct {
	textID string
	wordID string
	sentID string
	word   string
}

type FeatureData struct {
	textID   string
	sentID   string
	lineText string
}

type EyeData struct {
	sessID     string
	textID     string
	wordID     string
	word       string
	firstPass  string
	regression string
	totalPass  string
}

type FirstMerge struct {
	textID     string
	wordID     string
	sentID     string
	word       string
	firstPass  int
	regression int
	totalPass  int
}

var mapClozeData = map[string]ClozeData{}
var mapFeatureData = map[string]FeatureData{}
var mapEyeData = map[string]EyeData{}
var mapFirstMergeData = map[string]FirstMerge{}

func mainx() {

	path := "/home/sidleal/sid/usp/cloze_exps/"

	sents := readFile(path + "120sent_features.tsv")

	lines := strings.Split(sents, "\n")
	for i, line := range lines {
		if line == "" {
			continue
		}
		log.Println(i, line)
		cols := strings.Split(line, "\t")

		featData := FeatureData{}
		featData.textID = cols[0]
		featData.sentID = cols[1]
		featData.lineText = line

		mapFeatureData[fmt.Sprintf("%v_%v", featData.textID, featData.sentID)] = featData
	}

	cloze := readFile(path + "out.tsv")
	lines = strings.Split(cloze, "\n")
	for i, line := range lines {
		if line == "" || i == 0 {
			log.Print(line)
			headers := strings.Split(line, "\t")
			for j, it := range headers {
				log.Println(j, it)
			}

			continue
		}
		log.Println(i, line)

		cols := strings.Split(line, "\t")

		clozeData := ClozeData{}
		clozeData.textID = cols[1]
		clozeData.wordID = cols[3]
		clozeData.sentID = cols[4]
		clozeData.word = cols[7]

		mapClozeData[fmt.Sprintf("%v_%v", clozeData.textID, clozeData.wordID)] = clozeData

		log.Println("---->", clozeData)
	}

	data := readFile(path + "RASTROS_with_quotations.xls")

	lines = strings.Split(data, "\n")
	for i, line := range lines {
		if line == "" || i == 0 {
			log.Print(line)
			headers := strings.Split(line, "\t")
			for j, it := range headers {
				log.Println(j, it)
			}

			continue
		}
		log.Println(i, line)

		cols := strings.Split(line, "\t")

		eyeData := EyeData{}
		eyeData.sessID = cols[0]
		eyeData.textID = cols[1]
		eyeData.wordID = cols[3]
		eyeData.word = cols[4]
		eyeData.firstPass = cols[19]
		eyeData.regression = cols[34]
		eyeData.totalPass = cols[24]

		mapEyeData[fmt.Sprintf("%v_%v_%v", eyeData.sessID, eyeData.textID, eyeData.wordID)] = eyeData

		log.Println("---->", eyeData)

	}

	for kc, vc := range mapClozeData {
		log.Println(kc, vc)

		firstMerge := FirstMerge{}
		firstMerge.textID = vc.textID
		firstMerge.wordID = vc.wordID
		firstMerge.sentID = vc.sentID
		firstMerge.word = vc.word

		totFirstPass := 0
		totRegression := 0
		totTotalPass := 0
		totSamples := 0
		for ke, ve := range mapEyeData {
			if ve.textID == vc.textID && ve.wordID == vc.wordID {
				log.Println("--->", ke, ve.wordID, ve)
				firstPass, err := strconv.Atoi(ve.firstPass)
				if err != nil {
					firstPass = 0
				}
				totFirstPass += firstPass

				regression, err := strconv.Atoi(ve.regression)
				if err != nil {
					regression = 0
				}
				totRegression += regression

				totalPass, err := strconv.Atoi(ve.totalPass)
				if err != nil {
					totalPass = 0
				}
				totTotalPass += totalPass

				totSamples++
			}
		}
		firstMerge.firstPass = totFirstPass / totSamples
		firstMerge.regression = totRegression / totSamples
		firstMerge.totalPass = totTotalPass / totSamples

		mapFirstMergeData[fmt.Sprintf("%v_%v", vc.textID, vc.wordID)] = firstMerge

		log.Println("---------------------")

	}

	f, err := os.Create(fmt.Sprintf("%v/120sent_eye_features.tsv", path))
	if err != nil {
		log.Println("ERRO", err)
	}
	defer f.Close()

	headers := ""
	for kf, vf := range mapFeatureData {
		log.Println(kf, vf)
		if kf == "id_sent" {
			headers = vf.lineText + "\tavg_first_pass\tavg_regression\tavg_total_pass\tsum_first_pass\tsum_regression\tsum_total_pass"
			break
		}
	}

	_, err = f.WriteString(headers + "\n")
	if err != nil {
		log.Println("ERRO", err)
	}

	for kf, vf := range mapFeatureData {
		log.Println(kf, vf)
		if kf == "id_sent" {
			continue
		}

		totFirstPass := 0
		totRegression := 0
		totTotalPass := 0
		totSamples := 0
		for km, vm := range mapFirstMergeData {
			if vm.textID == vf.textID && vm.sentID == vf.sentID {
				log.Println(km, vm)
				totFirstPass += vm.firstPass
				totRegression += vm.regression
				totTotalPass += vm.totalPass
				totSamples++
			}
		}
		avgFirstPass := totFirstPass / totSamples
		avgRegression := totRegression / totSamples
		avgTotalPass := totTotalPass / totSamples

		log.Println("---------------------", avgFirstPass, avgRegression, avgTotalPass, totFirstPass, totRegression, totTotalPass)
		line := fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\t%v\n", vf.lineText, avgFirstPass, avgRegression, avgTotalPass, totFirstPass, totRegression, totTotalPass)
		_, err := f.WriteString(line)
		if err != nil {
			log.Println("ERRO", err)
		}
	}

}
