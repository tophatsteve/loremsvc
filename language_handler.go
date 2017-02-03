package loremsvc

import (
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"net/http"
)

func languageHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(r)
	kinds, err := getLanguages(ctx)

	if err != nil {
		respondErr(ctx, w, r, err, http.StatusBadRequest)
		return
	}

	respond(ctx, w, r, &kinds, http.StatusOK)
}

func getLanguages(ctx context.Context) ([]string, error) {
	return datastore.Kinds(ctx)
}

func isLanguageValid(ctx context.Context, language string) bool {
	l, err := getLanguages(ctx)

	if err != nil {
		return false
	}

	for _, v := range l {
		if v == language {
			return true
		}
	}

	return false
}
