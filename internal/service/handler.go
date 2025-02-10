package service

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"

	log "github.com/obalunenko/logger"

	"github.com/obalunenko/shipping-pack-optimizer/internal/service/assets"
)

type Packer interface {
	PackOrder(ctx context.Context, items uint) map[uint]uint
}

func NewRouter(p Packer) *http.ServeMux {
	mux := http.NewServeMux()

	mw := []func(http.Handler) http.Handler{
		logRequestMiddleware,
		logResponseMiddleware,
		requestIDMiddleware,
		recoverMiddleware,
		loggerMiddleware,
		corsMiddleware,
	}

	mwApply := func(h http.Handler) http.Handler {
		for i := range mw {
			h = mw[i](h)
		}

		return h
	}

	mux.Handle("/", mwApply(indexHandler()))
	mux.Handle("/favicon.ico", mwApply(faviconHandler()))

	svc := newService(p)

	// Group api/v1 routes.
	mux.Handle("/api/v1/pack", mwApply(svc.packHandler()))

	return mux
}

func indexHandler() http.HandlerFunc {
	homePageHTML := string(assets.MustLoad("index.gohtml"))
	homePageTmpl := template.Must(template.New("index").Parse(homePageHTML))

	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		if err := homePageTmpl.Execute(w, nil); err != nil {
			http.Error(w, "failed to execute template", http.StatusInternalServerError)

			return
		}
	}
}

func faviconHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}
}

func makeResponse(ctx context.Context, w http.ResponseWriter, code int, resp PackResponse, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	var response any

	response = resp

	if err != nil {
		log.WithError(ctx, err).Error("Error processing request")

		response = newHTTPError(ctx, code, err.Error())
	}

	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
