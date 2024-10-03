package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {

	log.Println("starting ---")

	path := "/home/sidleal/sources/500redacoes"
	files, err := os.ReadDir(path)
	if err != nil {
		log.Println("Erro", err)
	}

	fout, err := os.OpenFile(path+"/jeane_500redacoes_metricas.csv", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERRO", err)
	}
	defer fout.Close()

	for i, f := range files {
		fileName := f.Name()
		log.Println("---------------------------------------------------------------")
		log.Println(fileName)
		idTexto := strings.Split(fileName, ".")[0]
		log.Println(i, idTexto)
		if !strings.HasSuffix(fileName, ".txt") {
			continue
		}
		raw := readFile(path + "/" + fileName)

		log.Println(raw)

		text := cleanText(raw)

		log.Println(text)
		fRet := callMetrixNilc(text)
		//fRet := "adjective_ratio:0.03181,adverbs:0.08549,content_words:0.60968,flesch:-4.75802,function_words:0.39032,sentences_per_paragraph:46.0,syllables_per_content_word:2.79022,words_per_sentence:65.6087,noun_ratio:0.2674,paragraphs:1,sentences:46,words:3018,pronoun_ratio:0.08449,verbs:0.22498,logic_operators:0.05268,and_ratio:0.01988,if_ratio:0.00364,or_ratio:0.00431,negation_ratio:0.02353,cw_freq:674556.87337,min_cw_freq:417.91304,hypernyms_verbs:0.5538,brunet:25.0203,honore:350.76688,personal_pronouns:0.01195,ttr:0.76095,conn_ratio:0.08449,add_neg_conn_ratio:0.00265,add_pos_conn_ratio:0.03512,cau_neg_conn_ratio:0.0,cau_pos_conn_ratio:0.04506,log_neg_conn_ratio:0.00199,log_pos_conn_ratio:0.03744,tmp_neg_conn_ratio:0.0,tmp_pos_conn_ratio:0.00994,adjectives_ambiguity:2.16304,adverbs_ambiguity:2.69118,nouns_ambiguity:2.79336,verbs_ambiguity:9.54367,yngve:3.70354,frazier:8.6087,dep_distance:228.0,content_density:1.56197,words_before_main_verb:13.34783,adjacent_refs:0.8,anaphoric_refs:4.71111,adj_arg_ovl:6.62222,arg_ovl:8.01449,adj_stem_ovl:18.62222,stem_ovl:23.32174,adj_cw_ovl:11.6,adj_mean:0.32564,adj_std:0.19653,all_mean:0.35555,all_std:0.26293,paragraph_mean:0,paragraph_std:0,givenness_mean:0.53904,givenness_std:0.18303,span_mean:0.86771,span_std:0.24434,apposition_per_clause:0.00396,clauses_per_sentence:4.46903,prepositions_per_clause:0.91485,adjunct_per_clause:0.88317,prepositions_per_sentence:4.0885,relative_clauses:0.21386,aux_plus_PCP_per_sentence:0.74336,coordinate_conjunctions_per_clauses:0.14059,ratio_coordinate_conjunctions:0.55469,first_person_possessive_pronouns:0,first_person_pronouns:0.11429,gerund_verbs:0.03134,infinitive_verbs:0.31791,inflected_verbs:0.4806,non-inflected_verbs:0.5194,participle_verbs:0.17015,passive_ratio:0.14059,second_person_possessive_pronouns:0,second_person_pronouns:0.62857,sentences_with_five_clauses:0.0,sentences_with_four_clauses:0.23894,sentences_with_one_clause:0.17699,sentences_with_seven_more_clauses:0.17699,sentences_with_six_clauses:0.0708,sentences_with_three_clauses:0.31858,sentences_with_two_clauses:0.0177,sentences_with_zero_clause:0.0,simple_word_ratio:0.75326,ratio_subordinate_conjunctions:0.44531,third_person_possessive_pronouns:0,third_person_pronouns:0.25714,adjective_diversity_ratio:0.14583,adjectives_max:0.08,adjectives_min:0.0,adjectives_standard_deviation:0.0208,adverbs_diversity_ratio:0.16629,adverbs_max:0.13636,adverbs_min:0.0,adverbs_standard_deviation:0.03671,concretude_mean:3.83787,concretude_std:0.77188,concretude_1_25_ratio:0.00926,concretude_25_4_ratio:0.5933,concretude_4_55_ratio:0.3661,concretude_55_7_ratio:0.03134,content_word_diversity:0.78778,content_word_max:0.625,content_word_min:0.4,content_word_standard_deviation:0.07538,content_words_ambiguity:5.23289,dalechall_adapted:12.22727,verbal_time_moods_diversity:5,easy_conjunctions_ratio:0.07621,familiaridade_mean:5.02373,familiaridade_std:0.729,familiaridade_1_25_ratio:0.0,familiaridade_25_4_ratio:0.09046,familiaridade_4_55_ratio:0.5812,familiaridade_55_7_ratio:0.32835,function_word_diversity:0.37196,gunning_fox:26.37416,hard_conjunctions_ratio:0.01524,idade_aquisicao_mean:4.61723,idade_aquisicao_std:1.7372,idade_aquisicao_1_25_ratio:0.10826,idade_aquisicao_4_55_ratio:0.22293,idade_aquisicao_55_7_ratio:0.32835,idade_aquisicao_25_4_ratio:0.30769,imageabilidade_mean:4.15621,imageabilidade_std:0.62518,imageabilidade_1_25_ratio:0.0,imageabilidade_25_4_ratio:0.40385,imageabilidade_4_55_ratio:0.5755,imageabilidade_55_7_ratio:0.02066,indefinite_pronouns_diversity:0.08621,medium_long_sentence_ratio:0.0,max_noun_phrase:104,mean_noun_phrase:6.62274,medium_short_sentence_ratio:0.04348,min_noun_phrase:1,named_entity_ratio_sentence:0.03095,named_entity_ratio_text:0.03019,noun_diversity:0.68655,nouns_max:0.35659,nouns_min:0.21127,nouns_standard_deviation:0.03292,subtitles:0.0,postponed_subject_ratio:0.22896,preposition_diversity:0.20447,pronoun_diversity:0.18147,pronouns_max:0.33333,pronouns_min:0.0,pronouns_standard_deviation:0.04586,dialog_pronoun_ratio:0.74286,punctuation_diversity:0.02275,punctuation_ratio:0.11348,abstract_nouns_ratio:0.04009,adverbs_before_main_verb_ratio:0.32079,subjunctive_future_ratio:0.0,indefinite_pronoun_ratio:0.22745,indicative_condition_ratio:0.05689,indicative_future_ratio:0.0,indicative_imperfect_ratio:0.01198,indicative_pluperfect_ratio:0.05988,indicative_present_ratio:0.78743,indicative_preterite_perfect_ratio:0.08383,infinite_subordinate_clauses:0.24752,oblique_pronouns_ratio:0.0902,relative_pronouns_ratio:0.42353,subjunctive_imperfect_ratio:0.0,subjunctive_present_ratio:0.0,subordinate_clauses:0.57426,temporal_adjunct_ratio:0.02316,demonstrative_pronoun_ratio:3.0,coreference_pronoum_ratio:0.46667,non_svo_ratio:0.37228,relative_pronouns_diversity_ratio:0.06,sentence_length_max:138,sentence_length_min:13,sentence_length_standard_deviation:47.12632,short_sentence_ratio:0.0,std_noun_phrase:17.44632,verb_diversity:0.482,verbs_max:0.24242,verbs_min:0.10078,verbs_standard_deviation:0.05057,long_sentence_ratio:0.95652,ratio_function_to_content_words:0.64022,"
		log.Println(fRet)

		feats := strings.Split(fRet, ",")

		header := "id_texto,"
		ret := idTexto + ","
		for _, feat := range feats {
			kv := strings.Split(feat, ":")
			if len(kv) > 1 {
				if i == 0 {
					header += kv[0] + ","
				}
				ret += kv[1] + ","
			}
		}

		if i == 0 {
			header = strings.TrimRight(header, ",")
			_, err := fout.WriteString(header + "\n")
			if err != nil {
				log.Println("ERRO", err)
			}
		}

		ret = strings.TrimRight(ret, ",")
		_, err := fout.WriteString(ret + "\n")
		if err != nil {
			log.Println("ERRO", err)
		}

	}

}

func callMetrixNilc(text string) string {

	timeout := time.Duration(300 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	url := "http://fw.nilc.icmc.usp.br:23380/api/v1/metrix/all/m3tr1x01?format=plain"
	resp, err := client.Post(url, "text", bytes.NewReader([]byte(text)))
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}

	print(url, resp.Body)

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("Error parsing response. %v", err)
	}

	ret := string(body)

	ret = strings.Split(ret, "++")[1]

	return ret
}

func readFile(path string) string {
	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}

	// r := charmap.ISO8859_1.NewDecoder().Reader(f)
	r := io.Reader(f)

	ret := ""

	buf := make([]byte, 32*1024)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			ret += string(buf[:n])
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			break
		}
	}

	return ret

}

func cleanText(text string) string {
	text = strings.ReplaceAll(text, "\\r\\n", "{{enter}}")
	text = strings.ReplaceAll(text, "\\n", "{{enter}}")
	text = strings.ReplaceAll(text, "\\\"", "\"")
	text = strings.ReplaceAll(text, "è", "e")
	text = strings.ReplaceAll(text, "ì", "i")
	text = strings.ReplaceAll(text, "ò", "o")
	text = strings.ReplaceAll(text, "ù", "u")

	text = strings.ReplaceAll(text, " à ", "{{crase}}")
	text = strings.ReplaceAll(text, "à", "a")
	text = strings.ReplaceAll(text, "{{crase}}", " à ")
	text = strings.ReplaceAll(text, "`", "\"")
	text = strings.ReplaceAll(text, "´", "\"")
	return text
}
