package uidgen_test

import (
	"github.com/atong007/kit/uidgen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewUidGenerator(t *testing.T) {
	type args struct {
		nodeId int64
	}
	tests := []struct {
		name       string
		args       args
		wantNotNil bool
		wantErr    bool
	}{
		{
			"success",
			args{111},
			true,
			false,
		},
		{
			"error with negative node id",
			args{-1},
			false,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ug, err := uidgen.NewUidGenerator(tt.args.nodeId)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantNotNil, ug != nil)
		})
	}
}
func TestUidGenerator_NewId(t *testing.T) {
	ug, err := uidgen.NewUidGenerator(111)
	require.NoError(t, err)
	tests := []struct {
		name        string
		wantNotZero bool
	}{
		{
			"success",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ug.NewId()
			assert.Equal(t, tt.wantNotZero, got != 0)
		})
	}
}
