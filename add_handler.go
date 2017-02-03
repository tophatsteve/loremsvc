package loremsvc

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"net/http"
)

type content struct {
	Language  string   `json:"language"`
	Sentences []string `json:"sentences"`
}

func (c content) OK() error {

	if len(c.Language) == 0 {
		return errors.New("language is blank")
	}

	if len(c.Sentences) == 0 {
		return errors.New("no sentences supplied")
	}

	return nil
}

func addContentHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(r)

	var c content
	err := decode(r, &c)
	if err != nil {
		respondErr(ctx, w, r, err, http.StatusBadRequest)
		return
	}

	for _, v := range c.Sentences {
		key := datastore.NewIncompleteKey(ctx, c.Language, nil)

		entity := struct {
			Value string `json:"value"`
		}{Value: v}

		_, err = datastore.Put(ctx, key, &entity)

		if err != nil {
			respondErr(ctx, w, r, err, http.StatusBadRequest)
			return
		}
	}

	respond(ctx, w, r, nil, http.StatusOK)
}
