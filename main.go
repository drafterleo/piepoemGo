package main

import (
	"fmt"
	//word2vec "./word2vec"
	poem_model "./poem_model"
)

func main() {
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

	seedWords := []string{"принц_NOUN", "нищий_NOUN"}
	fmt.Printf("%+v", seedWords)
	bestWords := pm.W2V.MostSimilar(seedWords)
	fmt.Printf("Best Words %+v\n", bestWords)

}