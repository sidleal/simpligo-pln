package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func PalavrasParseHandler(w http.ResponseWriter, r *http.Request) {
	err := validateSession(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	content := r.FormValue("content")
	retType := r.FormValue("type")
	options := r.FormValue("options")

	resp, err := http.PostForm("http://"+palavrasIP+":"+palavrasPort+"/"+retType,
		url.Values{"sentence": {content}, "options": {options}})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error: %v\n", err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Errorf("Error reading response: %v.", err))
	}

	bodyString := string(body)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text")
	fmt.Fprint(w, "SAÍDA: \n"+bodyString)

}
