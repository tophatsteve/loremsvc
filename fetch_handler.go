package loremsvc

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/tophatsteve/margo"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"net/http"
	"strconv"
	"strings"
)

type text struct {
	Language       string   `json:"language"`
	ParagraphCount int      `json:"number_of_paragraphs"`
	SentenceCount  int      `json:"sentences_per_paragraph"`
	Paragraphs     []string `json:"paragraphs"`
}

func rootHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(r)
	respondErr(ctx, w, r, errors.New("please call a supported operation"), http.StatusBadRequest)
}

func fetchHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := appengine.NewContext(r)
	language := p.ByName("language")

	if !isLanguageValid(ctx, language) {
		respondErr(ctx, w, r, fmt.Errorf("%s is not a supported langauge", language), http.StatusBadRequest)
		return
	}

	numberOfParagraphs, err := strconv.Atoi(p.ByName("paragraphs"))
	if err != nil {
		respondErr(ctx, w, r, fmt.Errorf("number of paragraphs, %s, is not a valid number", p.ByName("paragraphs")), http.StatusBadRequest)
		return
	}

	numberOfSentences, err := strconv.Atoi(p.ByName("sentences"))
	if err != nil {
		respondErr(ctx, w, r, fmt.Errorf("number of sentences, %s, is not a valid number", p.ByName("sentences")), http.StatusBadRequest)
		return
	}

	l, err := loadLanguage(ctx, language)
	if err != nil {
		respondErr(ctx, w, r, err, http.StatusBadRequest)
		return
	}

	m := margo.NewMargo(l, 2)

	var paragraphs []string

	for x := 0; x < numberOfParagraphs; x++ {
		var paragraph []string
		for y := 0; y < numberOfSentences; y++ {
			paragraph = append(paragraph, m.GenerateSentence(0, true))
		}
		paragraphs = append(paragraphs, strings.Join(paragraph, " "))
	}

	t := text{
		Language:       language,
		Paragraphs:     paragraphs,
		SentenceCount:  numberOfSentences,
		ParagraphCount: numberOfParagraphs,
	}
	respond(ctx, w, r, &t, http.StatusOK)
}

func loadLanguage(ctx context.Context, language string) ([]string, error) {
	var entities []struct {
		Value string `json:"value"`
	}
	var lines []string
	_, err := datastore.NewQuery(language).GetAll(ctx, &entities)

	if err != nil {
		return lines, err
	}

	for _, e := range entities {
		lines = append(lines, e.Value)
	}

	return lines, nil
}
