package loremsvc

import (
	"bytes"
	"encoding/json"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
	"net/http"
)

func decode(r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		return err
	}

	if valid, ok := v.(interface {
		OK() error
	}); ok {
		err = valid.OK()
		if err != nil {
			return err
		}
	}

	return nil
}

func respond(ctx context.Context, w http.ResponseWriter, r *http.Request, v interface{}, code int) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(v)
	if err != nil {
		respondErr(ctx, w, r, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	_, err = buf.WriteTo(w)
	if err != nil {
		log.Errorf(ctx, "respond: %s", err)
	}
}

func respondErr(ctx context.Context, w http.ResponseWriter, r *http.Request, err error, code int) {
	errObj := struct {
		Error string `json:"error"`
	}{Error: err.Error()}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	err = json.NewEncoder(w).Encode(errObj)
	if err != nil {
		log.Errorf(ctx, "respondErr: %s", err)
	}
}
