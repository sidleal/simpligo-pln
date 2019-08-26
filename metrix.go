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

func callMetrix(text string) string {

	text = strings.Replace(text, "\"", "{{quotes}}", -1)
	text = strings.Replace(text, "\n", "{{enter}}", -1)
	text = strings.Replace(text, "!", "{{exclamation}}", -1)

	timeout := time.Duration(120 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Post("http://"+mainServerIP+":8008/metrics_all", "text", bytes.NewReader([]byte(text)))
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

	ret := callMetrix(string(text))

	if subset == "sel79" {
		feats := strings.Split(ret, ",")
		sel79 := []string{"adjective_ratio", "adverbs", "syllables_per_content_word", "words_per_sentence", "noun_ratio", "pronoun_ratio", "verbs", "negation_ratio", "cw_freq", "min_cw_freq", "first_person_pronouns", "ttr", "conn_ratio", "add_neg_conn_ratio", "add_pos_conn_ratio", "cau_neg_conn_ratio", "cau_pos_conn_ratio", "log_neg_conn_ratio", "log_pos_conn_ratio", "tmp_neg_conn_ratio", "tmp_pos_conn_ratio", "adjectives_ambiguity", "adverbs_ambiguity", "nouns_ambiguity", "verbs_ambiguity", "yngve", "frazier", "dep_distance", "words_before_main_verb", "mean_noun_phrase", "min_noun_phrase", "max_noun_phrase", "std_noun_phrase", "passive_ratio", "adj_arg_ovl", "arg_ovl", "adj_stem_ovl", "stem_ovl", "adj_cw_ovl", "third_person_pronouns", "concretude_mean", "concretude_std", "concretude_1_25_ratio", "concretude_25_4_ratio", "concretude_4_55_ratio", "concretude_55_7_ratio", "content_word_diversity", "familiaridade_mean", "familiaridade_std", "familiaridade_1_25_ratio", "familiaridade_25_4_ratio", "familiaridade_4_55_ratio", "familiaridade_55_7_ratio", "idade_aquisicao_mean", "idade_aquisicao_std", "idade_aquisicao_1_25_ratio", "idade_aquisicao_4_55_ratio", "idade_aquisicao_55_7_ratio", "idade_aquisicao_25_4_ratio", "imageabilidade_mean", "imageabilidade_std", "imageabilidade_1_25_ratio", "imageabilidade_25_4_ratio", "imageabilidade_4_55_ratio", "imageabilidade_55_7_ratio", "sentence_length_max", "sentence_length_min", "sentence_length_standard_deviation", "verb_diversity", "adj_mean", "adj_std", "all_mean", "all_std", "givenness_mean", "givenness_std", "span_mean", "span_std", "content_density", "ratio_function_to_content_words"}

		ret = ""
		for _, feat := range feats {
			kv := strings.Split(feat, ":")
			for _, selfeat := range sel79 {
				if len(kv) > 1 && kv[0] == selfeat {
					ret += kv[0] + ":" + kv[1] + ","
				}
			}

		}
	}

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

	fRet := callMetrix(content)
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
