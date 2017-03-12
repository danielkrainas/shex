package handlers

import (
	"net/http"

	"github.com/danielkrainas/gobag/context"
	"github.com/danielkrainas/gobag/decouple/cqrs"

	"github.com/danielkrainas/shex/api/v1"
	"github.com/danielkrainas/shex/registry/actions"
	"github.com/danielkrainas/shex/registry/queries"
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
	ctx := r.Context()
	mods, err := q.Execute(ctx, &queries.SearchMods{})
	if err != nil {
		acontext.GetLogger(ctx).Error(err)
		acontext.TrackError(ctx, err)
		return
	}

	if err := v1.ServeJSON(w, mods); err != nil {
		acontext.GetLogger(ctx).Errorf("error sending mods json: %v", err)
	}
}
