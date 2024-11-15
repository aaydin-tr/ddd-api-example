package valueobject

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAllocation(t *testing.T) {
	tests := []struct {
		name    string
		value   int
		want    *Allocation
		wantErr error
	}{
		{
			name:    "valid allocation",
			value:   100,
			want:    &Allocation{value: 100},
			wantErr: nil,
		},
		{
			name:    "negative allocation",
			value:   -1,
			want:    nil,
			wantErr: ErrInvalidAllocation,
		},
		{
			name:    "zero allocation",
			value:   0,
			want:    nil,
			wantErr: ErrInvalidAllocation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAllocation(tt.value)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAllocation_GetValue(t *testing.T) {
	allocation := &Allocation{value: 100}
	assert.Equal(t, 100, allocation.GetValue())
}

func TestAllocation_Equals(t *testing.T) {
	tests := []struct {
		name  string
		a     *Allocation
		value ValueObject
		want  bool
	}{
		{
			name:  "equal allocations",
			a:     &Allocation{value: 100},
			value: &Allocation{value: 100},
			want:  true,
		},
		{
			name:  "different allocations",
			a:     &Allocation{value: 100},
			value: &Allocation{value: 200},
			want:  false,
		},
		{
			name:  "nil value",
			a:     &Allocation{value: 100},
			value: nil,
			want:  false,
		},
		{
			name:  "different type",
			a:     &Allocation{value: 100},
			value: &Name{value: "100"},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Equals(tt.value)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAllocation_Scan(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  int
	}{
		{
			name:  "int64 value",
			value: int64(100),
			want:  100,
		},
		{
			name:  "int32 value",
			value: int32(100),
			want:  100,
		},
		{
			name:  "int value",
			value: 100,
			want:  100,
		},
		{
			name:  "nil value",
			value: nil,
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Allocation{}
			err := a.Scan(tt.value)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, a.value)
		})
	}
}

func TestAllocation_Value(t *testing.T) {
	a := &Allocation{value: 100}
	got, err := a.Value()
	assert.NoError(t, err)
	assert.Equal(t, 100, got)
}
