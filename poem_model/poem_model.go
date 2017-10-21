package poem_model

import (
	"io/ioutil"
	"encoding/json"
	"../morph"
	"strings"
)

type PoemModel struct {
	Poems []string   `json:"poems"`
	Bags  [][]string `json:"bags"`
	W2V   W2VModel
}

func (pm *PoemModel) LoadJsonModel(fileName string) error {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, pm)
	if err != nil {
		return err
	}
	return nil
}


func (pm *PoemModel) LoadW2VModel(fileName string) error {
	pm.W2V.Load(fileName)
	return nil
}

func (pm *PoemModel) tokenizeWords(words []string) []string {
	POS_TAGS := map[string]string {
		"NOUN": "_NOUN",
		"VERB": "_VERB", "INFN": "_VERB", "GRND": "_VERB", "PRTF": "_VERB", "PRTS": "_VERB",
		"ADJF": "_ADJ", "ADJS": "_ADJ",
		"ADVB": "_ADV",
		"PRED": "_ADP",
	}

	STOP_TAGS := map[string]bool {"PREP": true, "CONJ": true, "PRCL": true, "NPRO": true, "NUMR": true}

	result := make([]string, len(words), len(words))

	for _, w := range words {
		_, morphNorms, morphTags := morph.Parse(w)
		if len(morphNorms) == 0 {
			continue
		}

		suffixes := make(map[string]bool)

		for i, tags := range morphTags {
			norm := morphNorms[i]
			tag := strings.Split(tags, ",")[0]
			_, hasStopTag := STOP_TAGS[tag]
			if hasStopTag {
				break
			}

			suffix, hasPosTag := POS_TAGS[tag]
			_, hasSuffix := suffixes[suffix]
			if hasPosTag && ! hasSuffix {
				result = append(result, norm + suffix)
				suffixes[suffix] = true
			}
		}
	}

	return result
}
