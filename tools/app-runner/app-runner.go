package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

const (
	httpPort = ":8008"
)

func Router() *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	r.HandleFunc("/ranker", RankerHandler).Methods("POST")
	r.HandleFunc("/ranking", RankingHandler).Methods("POST")
	r.HandleFunc("/metrics/{subset}", MetricsHandler).Methods("POST")

	return r
}

func main() {

	var httpSrv *http.Server
	httpSrv = makeHTTPServer()

	httpSrv.Addr = httpPort
	fmt.Printf("Starting HTTP server on %s\n", httpPort)
	err := httpSrv.ListenAndServe()
	if err != nil {
		log.Fatalf("httpSrv.ListenAndServe() failed with %s", err)
	}

}

func makeHTTPServer() *http.Server {
	mux := Router()
	return &http.Server{
		ReadTimeout:  300 * time.Second,
		WriteTimeout: 300 * time.Second,
		IdleTimeout:  360 * time.Second,
		Handler:      mux,
	}
}

func RankerHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text")

	ret := ""

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ret += "Error reading req: %v." + err.Error()
		log.Println(ret)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, ret)
		return
	}

	text := string(body)
	text = strings.Replace(text, "\"", "'", -1)
	text = strings.Replace(text, "!", "! ", -1)
	ret += text + "\n"
	ret += "-------------------------------" + "\n"

	// log.Println(text)
	log.Println("/bin/bash", "-c", "coh-metrix-nilc/run_sentence.sh \""+text+"\"")

	cmd := exec.Command("/bin/bash", "-c", "coh-metrix-nilc/run_sentence.sh \""+text+"\"")
	out, err := cmd.CombinedOutput()
	if err != nil {
		ret += "cmd.Run() failed with " + err.Error()
		log.Println(ret)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, ret)
		return
	}
	fmt.Printf("combined out:\n%s\n", string(out))

	features := ""

	regEx := regexp.MustCompile(`([^\[]*)\[(.*)\]`)
	match := regEx.FindStringSubmatch(string(out))
	if len(match) < 1 {
		ret += "Erro resposta cohmetrix: \n" + string(out)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, ret)
		return
	}

	log.Println(match[2])
	features = match[2]

	ret += match[1] + "\n"
	ret += "-------------------------------" + "\n"

	log.Println("/bin/bash", "-c", "simpligo-ranker/run.sh "+features)

	cmd2 := exec.Command("/bin/bash", "-c", "simpligo-ranker/run.sh "+features)
	out2, err := cmd2.CombinedOutput()
	if err != nil {
		ret += "cmd.Run() failed with " + err.Error()
		log.Println(ret)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, ret)
		return
	}
	fmt.Printf("combined out:\n%s\n", string(out2))

	complexity := ""

	regEx2 := regexp.MustCompile(`result: ([0-9\.]+)`)
	match2 := regEx2.FindStringSubmatch(string(out2))
	if len(match2) > 0 {
		log.Println(match2[1])
		ret += match2[1] + "\n"

		floatVal, _ := strconv.ParseFloat(match2[1], 32)
		complexity = fmt.Sprintf("%.2f", floatVal*100)

	}

	ret += "-------------------------------" + "\n"

	w.WriteHeader(http.StatusOK)

	result := ""
	result += "====================================\n"
	result += "Complexidade --> " + complexity + "\n"
	result += "====================================\n"
	result += "\n\n" + ret

	w.Write([]byte(result))
	// fmt.Fprint(w, ret)
}

func MetricsHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	subset := vars["subset"]

	if subset == "all" {
		subset = ""
	}

	w.Header().Set("Content-Type", "text")

	ret := ""

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ret += "Error reading req: %v." + err.Error()
		log.Println(ret)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, ret)
		return
	}

	text := string(body)
	ret += text + "\n"
	ret += "-------------------------------" + "\n"

	// log.Println(text)
	log.Println("/bin/bash", "-c", "coh-metrix-nilc/run"+subset+".sh \""+text+"\"")

	cmd := exec.Command("/bin/bash", "-c", "coh-metrix-nilc/run"+subset+".sh \""+text+"\"")
	out, err := cmd.CombinedOutput()
	if err != nil {
		ret += "cmd.Run() failed with " + err.Error()
		log.Println(ret)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, ret)
		return
	}
	fmt.Printf("combined out:\n%s\n", string(out))

	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func RankingHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text")

	ret := ""

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ret += "Error reading req: %v." + err.Error()
		log.Println(ret)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, ret)
		return
	}

	text := string(body)
	text = strings.Replace(text, "\"", "'", -1)
	text = strings.Replace(text, "!", "! ", -1)
	ret += text + "\n"
	ret += "-------------------------------" + "\n"

	// log.Println(text)
	log.Println("/bin/bash", "-c", "coh-metrix-nilc/run_sentence156.sh \""+text+"\"")

	cmd := exec.Command("/bin/bash", "-c", "coh-metrix-nilc/run_sentence156.sh \""+text+"\"")
	out, err := cmd.CombinedOutput()
	if err != nil {
		ret += "cmd.Run() failed with " + err.Error()
		log.Println(ret)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, ret)
		return
	}
	fmt.Printf("combined out:\n%s\n", string(out))

	features := ""

	regEx := regexp.MustCompile(`([^\[]*)\[(.*)\]`)
	match := regEx.FindStringSubmatch(string(out))
	if len(match) < 1 {
		ret += "Erro resposta cohmetrix: \n" + string(out)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, ret)
		return
	}

	log.Println(match[2])
	features = match[2]

	ret += match[1] + "\n"
	ret += "-------------------------------" + "\n"

	log.Println("/bin/bash", "-c", "simpligo-ranking/run.sh "+features)

	cmd2 := exec.Command("/bin/bash", "-c", "simpligo-ranking/run.sh "+features)
	out2, err := cmd2.CombinedOutput()
	if err != nil {
		ret += "cmd.Run() failed with " + err.Error()
		log.Println(ret)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, ret)
		return
	}
	fmt.Printf("combined out:\n%s\n", string(out2))

	complexity := ""

	regEx2 := regexp.MustCompile(`result: ([0-9\.]+)`)
	match2 := regEx2.FindStringSubmatch(string(out2))
	if len(match2) > 0 {
		log.Println(match2[1])
		ret += match2[1] + "\n"

		floatVal, _ := strconv.ParseFloat(match2[1], 32)
		complexity = fmt.Sprintf("%.2f", floatVal*100)

	}

	ret += "-------------------------------" + "\n"

	w.WriteHeader(http.StatusOK)

	result := ""
	result += "====================================\n"
	result += "Complexidade --> " + complexity + "\n"
	result += "====================================\n"
	result += "\n\n" + ret

	w.Write([]byte(result))
	// fmt.Fprint(w, ret)
}
