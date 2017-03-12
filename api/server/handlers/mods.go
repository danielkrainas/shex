package handlers

import (
	"net/http"

	"github.com/danielkrainas/gobag/decouple/cqrs"

	"github.com/danielkrainas/shex/registry/actions"
)

func Mods(actionPack actions.Pack) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			CreateMod(actionPack, w, r)
		case http.MethodGet:
			SearchMods(actionPack, w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

func CreateMod(c cqrs.CommandHandler, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func SearchMods(q cqrs.QueryExecutor, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
