package main

import (
	"fmt"
	"./poem_model"
	"./morph"
	"time"
	_ "strconv"
)

func main () {
	testPoemModel()
	//testMorph()
}

func testPoemModel() {
	pm := new(poem_model.PoemModel)

	fmt.Println("Loading w2v:")
	pm.LoadW2VModel("C:/data/ruscorpora_1_300_10.bin")
	pm.LoadJsonModel("./data/poems_model.json")
	pm.Vectorize()
	pm.Matricize()
	//fmt.Println(pm.Poems[0])
	//fmt.Println(pm.Bags[0])


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

	queryWords := []string{"вальс"}

	for i := 0; i < 20; i ++ {
		start := time.Now()
		poems := pm.SimilarPoemsMx(queryWords, 1)
		elapsed := time.Since(start)
		fmt.Println(elapsed)
		fmt.Printf("%+v\n", poems)
	}

	//fmt.Println("start...")
	//for i := 0; i < 100; i ++ {
	//	go poemQuery(pm, queryWords, strconv.Itoa(i) + "a")
	//	time.Sleep(time.Microsecond * 200)
	//	//go poemQuery(pm, queryWords, strconv.Itoa(i) + "b")
	//	//time.Sleep(time.Microsecond * 200)
	//	//go poemQuery(pm, queryWords, strconv.Itoa(i) + "c")
	//	//time.Sleep(time.Microsecond * 200)
	//}

	//fmt.Scanln()
}

func poemQuery(pm * poem_model.PoemModel, words []string, id string) {
	start := time.Now()
	_ = pm.SimilarPoems(words, 1)
	elapsed := time.Since(start)
	fmt.Printf("%s : %v\n", id, elapsed)
}

func testMorph () {
	words, norms, tags := morph.Parse("еж")
	for i := range words {
		fmt.Printf("%-4s %-5s %s\n", words[i], norms[i], tags[i])
	}
}

