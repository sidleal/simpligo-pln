package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var charsToReplace = map[string]string{"è": "e", "ì": "i", "ò": "o", "ù": "u", "`": "'", "´": "'"}

func callMetrix(subset string, text string) string {

	text = strings.Replace(text, "\"", "{{quotes}}", -1)
	text = strings.Replace(text, "\n", "{{enter}}", -1)
	text = strings.Replace(text, "!", "{{exclamation}}", -1)
	text = strings.Replace(text, "#", "{{sharp}}", -1)
	text = strings.Replace(text, "&", "{{ampersand}}", -1)
	text = strings.Replace(text, "%", "{{percent}}", -1)
	text = strings.Replace(text, "$", "{{dollar}}", -1)

	text = strings.Replace(text, " à ", "{{crase}}", -1)
	text = strings.Replace(text, "à", "a", -1)
	text = strings.Replace(text, "{{crase}}", " à ", -1)

	for k, v := range charsToReplace {
		text = strings.ReplaceAll(text, k, v)
	}

	timeout := time.Duration(300 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	log.Println("MetrixAPIPostHandler - calling", "http://"+mainServerIP+":8008/metrics/"+subset)
	resp, err := client.Post("http://"+mainServerIP+":8008/metrics"+subset, "text", bytes.NewReader([]byte(text)))
	if err != nil {
		return fmt.Sprintf("Error extracting metrics: %v", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("Error parsing response. %v", err)
	}

	ret := string(body)
	return ret
}

func MetrixAPIPostHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	subset := vars["subset"]
	key := vars["key"]

	if key != "m3tr1x01" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	defer r.Body.Close()
	text, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Error parsing text. %v", err)
		return
	}

	ret := callMetrix(subset, string(text))

	log.Println("MetrixAPIPostHandler - Ret", ret)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, ret)

}

type MetrixResult struct {
	RawList string             `json:"raw"`
	List    []MetrixResultItem `json:"list"`
}

type MetrixResultItem struct {
	Metric string `json:"name"`
	Val    string `json:"val"`
}

func MetrixParseHandler(w http.ResponseWriter, r *http.Request) {
	err := validateSession(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	content := r.FormValue("content")
	// options := r.FormValue("options")

	fRet := callMetrix("all", content)
	feats := strings.Split(fRet, ",")

	ret := MetrixResult{}

	sRet := ""
	for _, feat := range feats {
		kv := strings.Split(feat, ":")
		if len(kv) > 1 {
			sRet += kv[0] + " : " + kv[1] + "\n"
			ret.List = append(ret.List, MetrixResultItem{kv[0], kv[1]})
		}
	}
	ret.RawList = sRet

	cJSON, err := json.Marshal(ret)
	if err != nil {
		log.Printf("Erro: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(cJSON))

	// w.WriteHeader(http.StatusOK)
	// w.Header().Set("Content-Type", "text")
	// fmt.Fprint(w, ret)

}
