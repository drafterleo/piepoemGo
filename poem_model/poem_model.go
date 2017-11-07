package poem_model

import (
	"io/ioutil"
	"encoding/json"
	"../morph"
	"strings"
	"sort"
	"github.com/gonum/matrix/mat64"
)

type PoemModel struct {
	Poems    []string   `json:"poems"`
	Bags     [][]string `json:"bags"`
	W2V      W2VModel
	Vectors  [][][]float32
	Matrices []mat64.Matrix
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


func (pm *PoemModel) Vectorize() {
	pm.Vectors = make([][][]float32, len(pm.Bags))
	for idx, bag := range pm.Bags {
		pm.Vectors[idx] = pm.TokenVectors(bag)
	}
}


func (pm *PoemModel) Matricize() {
	pm.Matrices = make([]mat64.Matrix, len(pm.Bags))
	for idx, bag := range pm.Bags {
		data, rows := pm.TokenVectorsData(bag)
		pm.Matrices[idx] = mat64.NewDense(rows, pm.W2V.Size, data).T()
	}
}


func (pm *PoemModel) TokenizeWords(words []string) []string {
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

		suffixes := make(map[string]bool) // added suffixes

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

func (pm *PoemModel) TokenVectors(tokens []string) [][]float32 {
	vecs := make([][]float32, 0, len(tokens))
	for _, token := range tokens {
		vector, err := pm.W2V.WordVector(token)
		if err != nil {
			continue
		}
		vecs = append(vecs, vector)
	}
	return vecs
}

func (pm *PoemModel) TokenVectorsData(tokens []string) (data []float64, rows int) {
	data = make([]float64, 0, pm.W2V.Size)
	rows = 0
	for _, token := range tokens {
		vector, err := pm.W2V.WordVector(token)
		if err != nil {
			continue
		}
		for _, val := range vector {
			data = append(data, float64(val))
		}
		rows ++
	}

	if rows > 0 {
		return data, rows
	} else {
		return  nil, 0
	}
}


func (pm *PoemModel) SimilarPoems(queryWords []string, topN int) []string {
	simPoems := make([]string, 0, topN)
	tokens := pm.TokenizeWords(queryWords)
	if len(tokens) == 0 || topN <= 0 {
		return simPoems
	}

	queryVecs := pm.TokenVectors(tokens)

	type PoemSimilarity struct {
		Idx	int
		Sim float32
	}

	sims := make([]PoemSimilarity, len(pm.Bags))

	for idx, _ := range pm.Bags {
		//poemVecs := pm.TokenVectors(pm.Bags[idx])
		poemVecs := pm.Vectors[idx]
		var sim float32
		for _, qv := range queryVecs {
			for _, pv := range poemVecs {
				var dist float32
				for i := 0; i < pm.W2V.Size; i ++ {
					// dot production
					dist += qv[i] * pv[i]
				}
				sim += dist
			}
		}

		if len(poemVecs) > 0 {
			sim /= float32(len(poemVecs) * len(queryVecs))
		}

		sims[idx].Idx = idx
		sims[idx].Sim = sim
	}

	sort.Slice(sims, func (i, j int) bool {
		return sims[i].Sim > sims[j].Sim
	})

	for i := 0; i < topN; i ++ {
		simPoems = append(simPoems, pm.Poems[sims[i].Idx])
	}

	return simPoems
}


func (pm *PoemModel) SimilarPoemsMx(queryWords []string, topN int) []string {
	simPoems := make([]string, 0, topN)
	tokens := pm.TokenizeWords(queryWords)
	queryData, queryVecsN := pm.TokenVectorsData(tokens)
	if len(tokens) == 0 || topN <= 0 || queryVecsN == 0{
		return simPoems
	}

	queryMx := mat64.NewDense(queryVecsN, pm.W2V.Size, queryData)

	type PoemSimilarity struct {
		Idx	int
		Sim float64
	}

	sims := make([]PoemSimilarity, len(pm.Bags))

	for idx, _ := range pm.Bags {
		var resMx mat64.Dense
		bagMx := pm.Matrices[idx]
		_, poemVecsN := bagMx.Dims()
		resMx.Mul(queryMx, bagMx)
		sim := mat64.Sum(&resMx)

		if poemVecsN > 0 {
			sim /= float64(poemVecsN + queryVecsN)
		}

		sims[idx].Idx = idx
		sims[idx].Sim = sim
	}

	sort.Slice(sims, func (i, j int) bool {
		return sims[i].Sim > sims[j].Sim
	})

	for i := 0; i < topN; i ++ {
		simPoems = append(simPoems, pm.Poems[sims[i].Idx])
	}

	return simPoems
}
