package service

import (
	"errors"
	"maps"
	"sort"
)

// ErrEmptyItems is returned when items is zero or empty.
var ErrEmptyItems = errors.New("empty items")

func fromAPIRequest(req PackRequest) (uint, error) {
	if req.Items == 0 {
		return 0, ErrEmptyItems
	}

	return req.Items, nil
}

func toAPIResponse(boxes map[uint]uint) PackResponse {
	var resp PackResponse

	for k, v := range maps.All(boxes) {
		resp.Packs = append(resp.Packs, Pack{
			Box:      k,
			Quantity: v,
		})
	}

	sort.Slice(resp.Packs, func(i, j int) bool {
		return resp.Packs[i].Box > resp.Packs[j].Box
	})

	return resp
}
