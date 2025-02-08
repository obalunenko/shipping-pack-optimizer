package packer

import (
	"context"
	"fmt"
	"slices"

	log "github.com/obalunenko/logger"
)

type Packer struct {
	boxes []uint
}

var DefaultBoxes = []uint{
	250,
	500,
	1000,
	2000,
	5000,
}

type Option func(*Packer)

func WithBoxes(boxes []uint) Option {
	return func(p *Packer) {
		slices.Sort(boxes)

		unique := make(map[uint]struct{}, len(boxes))

		isUnique := func(x uint) bool {
			if _, exist := unique[x]; !exist {
				unique[x] = struct{}{}

				return true
			}

			return false
		}

		filtered := boxes[:0]

		for _, b := range boxes {
			if isUnique(b) {
				filtered = append(filtered, b)
			}
		}

		p.boxes = filtered
	}
}

func WithDefaultBoxes() Option {
	return func(p *Packer) {
		p.boxes = DefaultBoxes
	}
}

func NewPacker(ctx context.Context, opts ...Option) (*Packer, error) {
	var p Packer

	if len(opts) == 0 {
		opts = []Option{WithDefaultBoxes()}
	}

	for _, opt := range opts {
		opt(&p)
	}

	if err := p.validate(); err != nil {
		return nil, fmt.Errorf("failed to validate packer: %w", err)
	}

	log.WithField(ctx, "boxes", p.boxes).Info("Packer created")

	return &p, nil
}

func (p Packer) validate() error {
	if len(p.boxes) == 0 {
		return fmt.Errorf("boxes list is empty")
	}

	// There should be no box with zero volume.
	for _, box := range p.boxes {
		if box == 0 {
			return fmt.Errorf("box with zero volume")
		}
	}

	return nil

}

func (p Packer) PackOrder(ctx context.Context, items uint) map[uint]uint {
	log.WithFields(ctx, log.Fields{
		"items": items,
		"boxes": p.boxes,
	}).Debug("Packing order")

	if items == 0 {
		return nil
	}

	if p.boxes[0] == 0 {
		// This should never happen, cause we validate boxes on creation.
		panic(fmt.Errorf("packer has box with zero volume: boxes [%v]", p.boxes))
	}

	// Preallocate memory for the result slice.
	// Make a prediction based on the number of items and the smallest box.
	result := make(map[uint]uint, len(p.boxes))

	if len(p.boxes) == 1 {
		box := p.boxes[0]

		if items < box {
			result[box]++

			return result
		}

		n := items / box

		last := items % box
		if last != 0 {
			n++
		}

		result[box] += n

		return result
	}

	for i := len(p.boxes) - 1; i >= 0; i-- {
		box := p.boxes[i]

		if box >= items {
			if i == 0 {
				result[box]++

				break
			}

			continue
		}

		if box < items {
			if i == 0 {
				result[p.boxes[i+1]]++

				break
			}

			n := items / box

			if n == 0 {
				result[box]++

				break
			}

			result[box] += n

			left := items % box

			if left == 0 {
				break
			}

			items = left
		}
	}

	return result
}
