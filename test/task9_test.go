package test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

var ErrConcrete error = errors.New("this is a concrete error")

func Operate(mode bool) error {
	if !mode {
		return errors.New("this is a concrete error")
	}
	return ErrConcrete
}

func TestOperateConcrete(t *testing.T) {
	tests := []struct {
		name string
		mode bool
		err  error
	}{
		{"Interesting#1", false, errors.New("this is a concrete error")},

		{"Interesting#2", true, ErrConcrete},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calcErr := Operate(tt.mode)
			if tt.mode {
				require.ErrorIs(t, tt.err, calcErr) // errorIs сравнивает по указателю
			} else {
				require.NotErrorIs(t, tt.err, calcErr) // поэтому две errors.New() с одинаковым msg не равны
			}
		})
	}
}
