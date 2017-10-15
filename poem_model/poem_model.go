package poem_model

import (
	"io/ioutil"
	"encoding/json"
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
