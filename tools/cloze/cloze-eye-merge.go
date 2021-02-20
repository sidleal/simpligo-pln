package main

import (
	"log"
	"os"
	"strings"
)

// ------------------------------------

func main() { //_eyemerge() {

	path := "/home/sidleal/sid/usp/cloze_exps4/final/"

	clozeFile := path + "/cloze_predict39_consolidada.tsv"
	eyeFile := path + "/All_part_AGO20.tsv"

	clozeMap := map[string]map[string]string{}
	clozeFields := map[int]string{
		0:  "Word_Unique_ID",
		1:  "Text_ID",
		2:  "Text",
		3:  "Genre",
		4:  "Word_Number",
		5:  "Sentence_Number",
		6:  "Word_In_Sentence_Number",
		7:  "Word_Place_In_Sent",
		8:  "Word",
		9:  "Word_Cleaned",
		10: "Word_Length",
		11: "Total_Response_Count",
		12: "Unique_Count",
		13: "OrthographicMatch",
		14: "IsModalResponse",
		15: "ModalResponse",
		16: "ModalResponseCount",
		17: "Certainty",
		18: "PoS",
		19: "Word_Content_Or_Function",
		20: "Word_PoS",
		21: "POSMatch",
		22: "Word_Inflection",
		23: "InflectionMatch",
		24: "Semantic_Word_Context_Score",
		25: "Semantic_Response_Match_Score",
		26: "Semantic_Response_Context_Score",
		27: "Freq_brWaC_fpm",
		28: "Freq_Brasileiro_fpm",
		29: "Freq_brWaC_log",
		30: "Freq_Brasileiro_log",
		31: "Time_to_Start",
		32: "Typing_Time",
		33: "Total_time",
		34: "Top_10_resp",
	}

	eyeItemMap := map[string]string{}
	eyeFields := map[int]string{
		0: "RECORDING_SESSION_LABEL",
		// 1:  "text",
		2:  "Trial_Index_",
		3:  "trial", // text_id
		4:  "IA_ID", // word_number
		5:  "IA_LABEL",
		6:  "TRIAL_INDEX",
		7:  "IA_LEFT",
		8:  "IA_RIGHT",
		9:  "IA_TOP",
		10: "IA_BOTTOM",
		11: "IA_AREA",
		12: "IA_FIRST_FIXATION_DURATION",
		13: "IA_FIRST_FIXATION_INDEX",
		14: "IA_FIRST_FIXATION_VISITED_IA_COUNT",
		15: "IA_FIRST_FIXATION_X",
		16: "IA_FIRST_FIXATION_Y",
		17: "IA_FIRST_FIX_PROGRESSIVE",
		18: "IA_FIRST_FIXATION_RUN_INDEX",
		19: "IA_FIRST_FIXATION_TIME",
		20: "IA_FIRST_RUN_DWELL_TIME",
		21: "IA_FIRST_RUN_FIXATION_COUNT",
		22: "IA_FIRST_RUN_START_TIME",
		23: "IA_FIRST_RUN_END_TIME",
		24: "IA_FIRST_RUN_FIXATION_%",
		25: "IA_DWELL_TIME",
		26: "IA_FIXATION_COUNT",
		27: "IA_RUN_COUNT",
		28: "IA_SKIP",
		29: "IA_REGRESSION_IN",
		30: "IA_REGRESSION_IN_COUNT",
		31: "IA_REGRESSION_OUT",
		32: "IA_REGRESSION_OUT_COUNT",
		33: "IA_REGRESSION_OUT_FULL",
		34: "IA_REGRESSION_OUT_FULL_COUNT",
		35: "IA_REGRESSION_PATH_DURATION",
		36: "IA_FIRST_SACCADE_AMPLITUDE",
		37: "IA_FIRST_SACCADE_ANGLE",
		38: "IA_FIRST_SACCADE_START_TIME",
		39: "IA_FIRST_SACCADE_END_TIME",
	}

	f2, err := os.Create(path + "RastrOS_Corpus_Eyetracking_Data.tsv")
	if err != nil {
		log.Println("ERRO", err)
	}
	defer f2.Close()

	_, err = f2.WriteString("RECORDING_SESSION_LABEL\tWord_Unique_ID\tText_ID\tGenre\tWord_Number\tSentence_Number\tWord_In_Sentence_Number\tWord_Place_In_Sent\tWord\tWord_Cleaned\tWord_Length\tTotal_Response_Count\tUnique_Count\tOrthographicMatch\tIsModalResponse\tModalResponse\tModalResponseCount\tCertainty\tPoS\tWord_Content_Or_Function\tWord_PoS\tPOSMatch\tWord_Inflection\tInflectionMatch\tSemantic_Word_Context_Score\tSemantic_Response_Match_Score\tSemantic_Response_Context_Score\tFreq_brWaC_fpm\tFreq_Brasileiro_fpm\tFreq_brWaC_log\tFreq_Brasileiro_log\tTime_to_Start\tTyping_Time\tTotal_time\tIA_ID\tIA_LABEL\tTRIAL_INDEX\tIA_LEFT\tIA_RIGHT\tIA_TOP\tIA_BOTTOM\tIA_AREA\tIA_FIRST_FIXATION_DURATION\tIA_FIRST_FIXATION_INDEX\tIA_FIRST_FIXATION_VISITED_IA_COUNT\tIA_FIRST_FIXATION_X\tIA_FIRST_FIXATION_Y\tIA_FIRST_FIX_PROGRESSIVE\tIA_FIRST_FIXATION_RUN_INDEX\tIA_FIRST_FIXATION_TIME\tIA_FIRST_RUN_DWELL_TIME\tIA_FIRST_RUN_FIXATION_COUNT\tIA_FIRST_RUN_START_TIME\tIA_FIRST_RUN_END_TIME\tIA_FIRST_RUN_FIXATION_%\tIA_DWELL_TIME\tIA_FIXATION_COUNT\tIA_RUN_COUNT\tIA_SKIP\tIA_REGRESSION_IN\tIA_REGRESSION_IN_COUNT\tIA_REGRESSION_OUT\tIA_REGRESSION_OUT_COUNT\tIA_REGRESSION_OUT_FULL\tIA_REGRESSION_OUT_FULL_COUNT\tIA_REGRESSION_PATH_DURATION\tIA_FIRST_SACCADE_AMPLITUDE\tIA_FIRST_SACCADE_ANGLE\tIA_FIRST_SACCADE_START_TIME\tIA_FIRST_SACCADE_END_TIME\n")
	if err != nil {
		log.Println("ERRO", err)
	}

	data := readFile(clozeFile)

	lines := strings.Split(data, "\n")
	for i, line := range lines {
		if line == "" {
			continue
		}

		cols := strings.Split(line, "\t")
		if i == 0 {
			log.Println(cols)
			continue
		}

		for k, v := range clozeFields {
			if _, found := clozeMap[cols[0]]; !found {
				clozeMap[cols[0]] = map[string]string{}
			}
			clozeMap[cols[0]][v] = cols[k]
		}

	}

	data = readFileISO(eyeFile)

	lines = strings.Split(data, "\n")
	for i, line := range lines {
		if line == "" {
			continue
		}

		cols := strings.Split(line, "\t")
		if i == 0 {
			// log.Println(cols)
			continue
		}

		// log.Println(cols)

		for k, v := range eyeFields {
			eyeItemMap[v] = strings.ReplaceAll(cols[k], ",", ".")
		}
		wordUniqueID := "UID_" + eyeItemMap["trial"] + "_" + eyeItemMap["IA_ID"]
		log.Println("==============>", wordUniqueID)

		if clozeInfo, found := clozeMap[wordUniqueID]; found {
			log.Println("---->", clozeInfo)
			log.Println("-->", eyeItemMap)

			str := ""
			str += eyeItemMap["RECORDING_SESSION_LABEL"] + "\t"
			str += clozeInfo["Word_Unique_ID"] + "\t"
			str += clozeInfo["Text_ID"] + "\t"
			str += clozeInfo["Genre"] + "\t"
			str += clozeInfo["Word_Number"] + "\t"
			str += clozeInfo["Sentence_Number"] + "\t"
			str += clozeInfo["Word_In_Sentence_Number"] + "\t"
			str += clozeInfo["Word_Place_In_Sent"] + "\t"
			str += clozeInfo["Word"] + "\t"
			str += clozeInfo["Word_Cleaned"] + "\t"
			str += clozeInfo["Word_Length"] + "\t"
			str += clozeInfo["Total_Response_Count"] + "\t"
			str += clozeInfo["Unique_Count"] + "\t"
			str += clozeInfo["OrthographicMatch"] + "\t"
			str += clozeInfo["IsModalResponse"] + "\t"
			str += clozeInfo["ModalResponse"] + "\t"
			str += clozeInfo["ModalResponseCount"] + "\t"
			str += clozeInfo["Certainty"] + "\t"
			str += clozeInfo["PoS"] + "\t"
			str += clozeInfo["Word_Content_Or_Function"] + "\t"
			str += clozeInfo["Word_PoS"] + "\t"
			str += clozeInfo["POSMatch"] + "\t"
			str += clozeInfo["Word_Inflection"] + "\t"
			str += clozeInfo["InflectionMatch"] + "\t"
			str += clozeInfo["Semantic_Word_Context_Score"] + "\t"
			str += clozeInfo["Semantic_Response_Match_Score"] + "\t"
			str += clozeInfo["Semantic_Response_Context_Score"] + "\t"
			str += clozeInfo["Freq_brWaC_fpm"] + "\t"
			str += clozeInfo["Freq_Brasileiro_fpm"] + "\t"
			str += clozeInfo["Freq_brWaC_log"] + "\t"
			str += clozeInfo["Freq_Brasileiro_log"] + "\t"
			str += clozeInfo["Time_to_Start"] + "\t"
			str += clozeInfo["Typing_Time"] + "\t"
			str += clozeInfo["Total_time"] + "\t"
			str += eyeItemMap["IA_ID"] + "\t"
			str += eyeItemMap["IA_LABEL"] + "\t"
			str += eyeItemMap["TRIAL_INDEX"] + "\t"
			str += eyeItemMap["IA_LEFT"] + "\t"
			str += eyeItemMap["IA_RIGHT"] + "\t"
			str += eyeItemMap["IA_TOP"] + "\t"
			str += eyeItemMap["IA_BOTTOM"] + "\t"
			str += eyeItemMap["IA_AREA"] + "\t"
			str += eyeItemMap["IA_FIRST_FIXATION_DURATION"] + "\t"
			str += eyeItemMap["IA_FIRST_FIXATION_INDEX"] + "\t"
			str += eyeItemMap["IA_FIRST_FIXATION_VISITED_IA_COUNT"] + "\t"
			str += eyeItemMap["IA_FIRST_FIXATION_X"] + "\t"
			str += eyeItemMap["IA_FIRST_FIXATION_Y"] + "\t"
			str += eyeItemMap["IA_FIRST_FIX_PROGRESSIVE"] + "\t"
			str += eyeItemMap["IA_FIRST_FIXATION_RUN_INDEX"] + "\t"
			str += eyeItemMap["IA_FIRST_FIXATION_TIME"] + "\t"
			str += eyeItemMap["IA_FIRST_RUN_DWELL_TIME"] + "\t"
			str += eyeItemMap["IA_FIRST_RUN_FIXATION_COUNT"] + "\t"
			str += eyeItemMap["IA_FIRST_RUN_START_TIME"] + "\t"
			str += eyeItemMap["IA_FIRST_RUN_END_TIME"] + "\t"
			str += eyeItemMap["IA_FIRST_RUN_FIXATION_%"] + "\t"
			str += eyeItemMap["IA_DWELL_TIME"] + "\t"
			str += eyeItemMap["IA_FIXATION_COUNT"] + "\t"
			str += eyeItemMap["IA_RUN_COUNT"] + "\t"
			str += eyeItemMap["IA_SKIP"] + "\t"
			str += eyeItemMap["IA_REGRESSION_IN"] + "\t"
			str += eyeItemMap["IA_REGRESSION_IN_COUNT"] + "\t"
			str += eyeItemMap["IA_REGRESSION_OUT"] + "\t"
			str += eyeItemMap["IA_REGRESSION_OUT_COUNT"] + "\t"
			str += eyeItemMap["IA_REGRESSION_OUT_FULL"] + "\t"
			str += eyeItemMap["IA_REGRESSION_OUT_FULL_COUNT"] + "\t"
			str += eyeItemMap["IA_REGRESSION_PATH_DURATION"] + "\t"
			str += eyeItemMap["IA_FIRST_SACCADE_AMPLITUDE"] + "\t"
			str += eyeItemMap["IA_FIRST_SACCADE_ANGLE"] + "\t"
			str += eyeItemMap["IA_FIRST_SACCADE_START_TIME"] + "\t"
			str += eyeItemMap["IA_FIRST_SACCADE_END_TIME"] + "\n"

			_, err = f2.WriteString(str)
			if err != nil {
				log.Println("ERRO", err)
			}

		}

	}

	// for k, v := range clozeMap {
	// 	log.Println(k, v)
	// }

}
