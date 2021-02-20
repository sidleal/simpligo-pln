package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	ngroups = 57
)

type WordDP struct {
	Word        string
	Tag         string
	Annotation  string
	Freqs       [ngroups]int
	FreqTotal   int
	DP          float64
	DPNorm      float64
	FreqAvg     float64
	FreqCorr    float64
	Ranking     int
	FreqCorrPer float64
	FreqCorrAcu float64
}

type WordDPOrder []WordDP

func (a WordDPOrder) Len() int      { return len(a) }
func (a WordDPOrder) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a WordDPOrder) Less(i, j int) bool {
	return a[i].FreqCorr > a[j].FreqCorr
}

func main() {

	log.Println("Inicio")

	wordList := []WordDP{}

	data := readFile("/home/sidleal/sid/usp/brWaC/lista_nlpnet_groups_full_cleaned03.tsv")

	dataLines := strings.Split(data, "\n")
	for _, line := range dataLines {
		if line == "" {
			continue
		}
		tokens := strings.Split(line, "\t")

		w := WordDP{}
		w.Word = tokens[0]
		w.Tag = tokens[1]
		w.Annotation = tokens[59]

		qts := [ngroups]int{}
		for i := 0; i < ngroups; i++ {
			qty, _ := strconv.Atoi(tokens[i+2])
			qts[i] = qty
		}
		w.Freqs = qts

		wordList = append(wordList, w)

	}

	totals := [ngroups]int{}
	for i, w := range wordList {
		log.Println(i, w)
		for i := 0; i < ngroups; i++ {
			totals[i] += w.Freqs[i]
		}
	}

	totAll := 0
	for i := 0; i < ngroups; i++ {
		totAll += totals[i]
	}

	percents := []float64{}
	for i := 0; i < ngroups; i++ {
		percents = append(percents, float64(totals[i])/float64(totAll))
	}

	log.Println("Totals:", totals)
	log.Println("Total All:", totAll)
	log.Println("Percents:", percents)

	wordList2 := []WordDP{}

	var totFreqCorr float64 = 0

	for _, w := range wordList {
		// log.Println(k, w)

		v := []float64{}
		f := 0.0
		for i := 0; i < ngroups; i++ {
			v = append(v, float64(w.Freqs[i]))
			f += float64(w.Freqs[i])
		}

		// Deviation of proportions DP (Dispersion)
		div := []float64{}
		for i := 0; i < ngroups; i++ {
			div = append(div, v[i]/f)
		}

		sub := []float64{}
		for i := 0; i < ngroups; i++ {
			sub = append(sub, div[i]-percents[i])
		}

		sum := 0.0
		for i := 0; i < ngroups; i++ {
			sum += math.Abs(sub[i])
		}

		dp := sum / 2
		dpNorm := dp / (1 - min(percents))

		w.FreqTotal = int(f)
		w.DP = dp
		w.DPNorm = dpNorm
		w.FreqAvg = float64(w.FreqTotal) / float64(len(v))
		w.FreqCorr = w.FreqAvg*(1-w.DPNorm) + 0.0001

		totFreqCorr += w.FreqCorr

		wordList2 = append(wordList2, w)

		log.Println("------------>", f, dp, dpNorm)

	}

	sort.Sort(WordDPOrder(wordList2))

	wordList3 := []WordDP{}
	var freqCorrAcum float64 = 0
	for i, w := range wordList2 {
		w.Ranking = i + 1
		w.FreqCorrPer = w.FreqCorr / totFreqCorr * 100
		freqCorrAcum += w.FreqCorrPer
		w.FreqCorrAcu = freqCorrAcum
		wordList3 = append(wordList3, w)
	}

	f1, err := os.Create("/home/sidleal/sid/usp/brWaC/brwac_dispersion1.tsv")
	if err != nil {
		log.Println("ERRO", err)
	}

	defer f1.Close()

	header := "word\ttag\tannotation\t"
	for i := 0; i < ngroups; i++ {
		header += fmt.Sprintf("%v\t", i+1)
	}
	header += "total\taverage\tdp\tdp_norm\tfreq_corr\tfreq_perc\tfreq_acum\tranking\n"

	_, err = f1.WriteString(header)
	if err != nil {
		log.Println("ERRO", err)
	}

	for _, w := range wordList3 {
		annot := ""
		if w.Annotation != "-1" {
			annot = w.Annotation
		}
		line := fmt.Sprintf("%v\t%v\t%v\t", w.Word, strings.ToLower(w.Tag), annot)
		for i := 0; i < ngroups; i++ {
			line += fmt.Sprintf("%v\t", w.Freqs[i])
		}
		line += fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", w.FreqTotal, w.FreqAvg, w.DP, w.DPNorm, w.FreqCorr, w.FreqCorrPer, w.FreqCorrAcu, w.Ranking)
		_, err := f1.WriteString(line)
		if err != nil {
			log.Println("ERRO", err)
		}
	}

}
