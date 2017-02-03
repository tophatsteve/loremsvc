package loremsvc

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func init() {
	router := httprouter.New()
	router.GET("/", rootHandler)
	router.GET("/languages", languageHandler)
	router.POST("/add", addContentHandler)
	router.GET("/generate/:language/:paragraphs/:sentences", fetchHandler)
	http.Handle("/", router)
}
