package main

import (
	"fmt"
	"./poem_model"
	"./morph"
	"strings"
)

func main () {
	testPoemModel()
	//testMorph()
}

func testPoemModel() {
	pm := new(poem_model.PoemModel)
	pm.LoadJsonModel("./data/poems_model.json")
	fmt.Println(pm.Poems[0])
	fmt.Println(pm.Bags[0])
	fmt.Println("Loading w2v:")
	pm.LoadW2VModel("C:/data/ruscorpora_1_300_10.bin")
	//for i := 1000; i < 2000; i ++ {
	//  fmt.Printf("%v ", pm.W2V.Vocab[i])
	//}
	//fmt.Println(pm.W2V.Vec[0])

	seedWords := []string{"принц", "нищий"}
	tokens := tokenizeWords(seedWords)
	fmt.Printf("%+v", tokens)
	bestWords := pm.W2V.MostSimilar(tokens)
	fmt.Printf("Best Words %+v\n", bestWords)
}

func testMorph () {
	words, norms, tags := morph.Parse("еж")
	for i := range words {
		fmt.Printf("%-4s %-5s %s\n", words[i], norms[i], tags[i])
	}

	inWords := []string{"принц", "нищий"}
	tokens := tokenizeWords(inWords)
	fmt.Printf("Tokenized words %v\n", tokens)
}

func tokenizeWords(words []string) []string {
	POS_TAGS := map[string]string {
		"NOUN": "_NOUN",
		"VERB": "_VERB", "INFN": "_VERB", "GRND": "_VERB", "PRTF": "_VERB", "PRTS": "_VERB",
		"ADJF": "_ADJ", "ADJS": "_ADJ",
		"ADVB": "_ADV",
		"PRED": "_ADP",
	}

	STOP_TAGS := map[string]bool {"PREP": true, "CONJ": true, "PRCL": true, "NPRO": true, "NUMR": true}

	result := make([]string, 0, len(words))

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