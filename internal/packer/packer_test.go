package packer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/obalunenko/shipping-pack-optimizer/internal/testlogger"
)

func TestPacker_PackOrder(t *testing.T) {
	ctx := testlogger.New(context.Background())

	type fields struct {
		boxes []uint
	}

	type args struct {
		items uint
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[uint]uint
	}{
		{
			name: "default. 1 - 1x250",
			fields: fields{
				boxes: DefaultBoxes(),
			},
			args: args{
				items: 1,
			},
			want: map[uint]uint{
				250: 1,
			},
		},
		{
			name: "default. 251 - 1x500",
			fields: fields{
				boxes: DefaultBoxes(),
			},
			args: args{
				items: 251,
			},
			want: map[uint]uint{
				500: 1,
			},
		},
		{
			name: "default. 501 - 1x500, 1x250",
			fields: fields{
				boxes: DefaultBoxes(),
			},
			args: args{
				items: 501,
			},
			want: map[uint]uint{
				500: 1,
				250: 1,
			},
		},
		{
			name: "default. 12001  - 2x5000, 1x2000, 1x250",
			fields: fields{
				boxes: DefaultBoxes(),
			},
			args: args{
				items: 12001,
			},
			want: map[uint]uint{
				5000: 2,
				2000: 1,
				250:  1,
			},
		},
		{
			name: "default. 29292929292929  - ?",
			fields: fields{
				boxes: DefaultBoxes(),
			},
			args: args{
				items: 29292929292929,
			},
			want: map[uint]uint{
				5000: 5858585858,
				2000: 1,
				500:  2,
			},
		},

		{
			name: "custom[1, 2, 4, 8]. 1 - 1",
			fields: fields{
				boxes: []uint{1, 2, 4, 8},
			},
			args: args{
				items: 1,
			},
			want: map[uint]uint{
				1: 1,
			},
		},

		{
			name: "custom[3]. 7 - 3 3 3",
			fields: fields{
				boxes: []uint{3},
			},
			args: args{
				items: 7,
			},
			want: map[uint]uint{
				3: 3,
			},
		},

		{
			name: "custom[1,2,4]. 7 - 4, 2, 1?",
			fields: fields{
				boxes: []uint{1, 2, 4},
			},
			args: args{
				items: 7,
			},
			want: map[uint]uint{
				4: 1,
				2: 1,
				1: 1,
			},
		},
		{
			name: "custom edge cases[23, 31, 53]. 500_000 -23: 2, 31: 7, 53: 9429",
			fields: fields{
				boxes: []uint{23, 31, 53},
			},
			args: args{
				items: 500_000,
			},
			want: map[uint]uint{
				53: 9429,
				31: 7,
				23: 2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewPacker(ctx, WithBoxes(tt.fields.boxes))
			require.NoError(t, err)

			got := p.PackOrder(ctx, tt.args.items)

			assert.EqualValues(t, tt.want, got)
		})
	}
}

func TestNewPacker(t *testing.T) {
	ctx := testlogger.New(context.Background())

	type args struct {
		opts []Option
	}

	tests := []struct {
		name    string
		args    args
		want    *Packer
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "default boxes",
			args: args{
				opts: []Option{},
			},
			want: &Packer{
				boxes: DefaultBoxes(),
			},
			wantErr: assert.NoError,
		},
		{
			name: "custom boxes",
			args: args{
				opts: []Option{
					WithBoxes([]uint{32, 1, 2, 2, 4, 16, 8, 16}),
				},
			},
			want: &Packer{
				boxes: []uint{1, 2, 4, 8, 16, 32},
			},
			wantErr: assert.NoError,
		},
		{
			name: "custom boxes empty - error",
			args: args{
				opts: []Option{
					WithBoxes([]uint{}),
				},
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "custom boxes contains zero - error",
			args: args{
				opts: []Option{
					WithBoxes([]uint{9, 0, 2}),
				},
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPacker(ctx, tt.args.opts...)
			if !tt.wantErr(t, err) {
				return
			}

			assert.EqualValues(t, tt.want, got)
		})
	}
}
