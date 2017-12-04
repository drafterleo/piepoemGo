package main

import (
	"gopkg.in/gcfg.v1"

	"./poem_model"

	"net/http"
	"github.com/gin-gonic/gin"
	//"github.com/gin-contrib/cors"

	"fmt"
	"strings"
)

var config struct {
	Models struct {
		W2V string
		Poems string
	}
}

const defaultConfig = `
[models]
w2v="c:/data/ruscorpora_1_300_10.bin"
poems="./data/poems_model.json"
`
const MAX_SEARCH_LENGTH = 100

var poemModel *poem_model.PoemModel

func main () {
	//testPoemModel()
	//testMorph()

	loadConfig("./config.cfg")
	startPoemModel()
	startRouter()
}

func loadConfig(cfgFile string) error {
	var err error

	if cfgFile != "" {
		err = gcfg.ReadFileInto(&config, cfgFile)
	} else {
		err = gcfg.ReadStringInto(&config, defaultConfig)
	}

	//fmt.Println(config)

	return err
}

func startPoemModel() {
	poemModel = new(poem_model.PoemModel)

	fmt.Println("Loading w2v:")
	poemModel.LoadW2VModel(config.Models.W2V)
	poemModel.LoadJsonModel(config.Models.Poems)
	poemModel.Matricize()
}

func startRouter(){
	router := gin.Default()

	//router.Use(cors.Default())

	router.LoadHTMLGlob("./site/*.html")

	router.Static("/css", "./site/css")
	router.Static("/img", "./site/img")
	router.Static("/fonts", "./site/fonts")
	router.Static("/dist", "./site/dist")
	router.StaticFile("/favicon.ico", "./resources/favicon.ico")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
		// fmt.Println(c.Request.URL.Query())
	})

	router.POST("/poems", postPoems)

	router.Run(":8085")
}


func postPoems (c *gin.Context) {
	var data struct {
		Words string `json:"words"`
	}

	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(data.Words) > MAX_SEARCH_LENGTH {
		data.Words = data.Words[0:MAX_SEARCH_LENGTH]
	}
	fmt.Printf("words: <%s>\n", data.Words)
	queryWords := strings.Fields(strings.ToLower(data.Words))
	fmt.Println(queryWords)
	poems := poemModel.SimilarPoemsMx(queryWords, 10)
	for idx, poem := range poems {
		poems[idx] = strings.Replace(poem, "\n", "<br>", -1)
	}

	c.JSON(http.StatusOK, poems)
}
