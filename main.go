package main

import (
	"./poem_model"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"fmt"
	"strings"
)

var poemModel *poem_model.PoemModel

func main () {
	//testPoemModel()
	//testMorph()

	startPoemModel()
	startRouter()
}

func startPoemModel() {
	poemModel = new(poem_model.PoemModel)

	fmt.Println("Loading w2v:")
	poemModel.LoadW2VModel("C:/data/ruscorpora_1_300_10.bin")
	poemModel.LoadJsonModel("./data/poems_model.json")
	poemModel.Matricize()
}

func startRouter(){
	router := gin.Default()

	router.Use(cors.Default())

	router.LoadHTMLGlob("./site/*.html")

	router.Static("/css", "./site/css")
	router.Static("/img", "./site/img")
	router.Static("/fonts", "./site/fonts")
	router.Static("/dist", "./site/dist")
	router.StaticFile("/favicon.ico", "./resources/favicon.ico")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
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

	fmt.Printf("words: <%s>\n", data.Words)
	queryWords := strings.Fields(data.Words)
	fmt.Println(queryWords)
	poems := poemModel.SimilarPoemsMx(queryWords, 10)
	for idx, poem := range poems {
		poems[idx] = strings.Replace(poem, "\n", "<br>", -1)
	}

	c.JSON(http.StatusOK, poems)
}
