package main

import (
	"fmt"
	"./poem_model"
	"./morph"
	"time"
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

	//seedWords := []string{"принц", "нищий"}
	//tokens := pm.TokenizeWords(seedWords)
	//fmt.Printf("%+v", tokens)
	//bestWords, err := pm.W2V.MostSimilar(tokens, 10)
	//if err == nil {
	//	fmt.Printf("Best Words %+v\n", bestWords)
	//}

	queryWords := []string{"вальс", "гарнизон", "астролябия"}

	for i := 0; i < 10; i ++ {
		start := time.Now()
		poems := pm.SimilarPoems(queryWords, 1)
		elapsed := time.Since(start)
		fmt.Println(elapsed)
		fmt.Printf("%+v\n", poems)
	}
}

func testMorph () {
	words, norms, tags := morph.Parse("еж")
	for i := range words {
		fmt.Printf("%-4s %-5s %s\n", words[i], norms[i], tags[i])
	}
}

