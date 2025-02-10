package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	log "github.com/obalunenko/logger"
)

type service struct {
	cache  map[PackRequest]PackResponse
	mu     sync.RWMutex
	packer Packer
}

func newService(packer Packer) *service {
	return &service{
		cache:  make(map[PackRequest]PackResponse),
		mu:     sync.RWMutex{},
		packer: packer,
	}
}

// packHandler - handler for /pack endpoint.
//
//	@Summary		Get the number of packs needed to ship to a customer
//	@Tags			pack
//	@Description	Calculates the number of packs needed to ship to a customer
//	@ID				shipping-pack-optimizer-pack	post
//	@Accept			json
//	@Produce		json
//	@Param			data	body		PackRequest				true	"Request data"
//	@Success		200		{object}	PackResponse			"Successful response with packs data"
//	@Failure		400		{object}	badRequestError			"Invalid request data
//	@Failure		405		{object}	methodNotAllowedError	"Method not allowed"
//	@Failure		500		{object}	internalServerError		"Internal server error"
//	@Router			/api/v1/pack [post]
func (s *service) packHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if r.Method != http.MethodPost {
			makeResponse(
				r.Context(),
				w,
				http.StatusMethodNotAllowed,
				PackResponse{},
				errors.New(http.StatusText(http.StatusMethodNotAllowed)),
			)

			return
		}

		b, err := io.ReadAll(r.Body)
		if err != nil {
			makeResponse(
				r.Context(),
				w,
				http.StatusBadRequest,
				PackResponse{},
				fmt.Errorf("failed to read request body: %w", err),
			)

			return
		}

		defer func() {
			if err = r.Body.Close(); err != nil {
				log.WithError(ctx, err).Error("Error closing request body")
			}
		}()

		var req PackRequest

		if err = json.Unmarshal(b, &req); err != nil {
			makeResponse(ctx, w, http.StatusBadRequest, PackResponse{}, fmt.Errorf("failed to unmarshal request: %w", err))

			return
		}

		var resp PackResponse

		s.mu.RLock()
		resp, ok := s.cache[req]
		s.mu.RUnlock()

		if ok {
			log.Debug(ctx, "Response found in cache")
			makeResponse(ctx, w, http.StatusOK, resp, nil)
			return
		}

		log.Debug(ctx, "Response not found in cache, calculating")

		items, err := fromAPIRequest(req)
		if err != nil {
			makeResponse(
				ctx,
				w,
				http.StatusBadRequest,
				PackResponse{},
				fmt.Errorf("invalid request: %w", err),
			)

			return
		}

		order := s.packer.PackOrder(ctx, items)

		resp = toAPIResponse(order)

		log.Debug(ctx, "Saving response to cache")

		s.mu.Lock()
		s.cache[req] = resp
		s.mu.Unlock()

		makeResponse(ctx, w, http.StatusOK, resp, nil)
	}
}
