package aes_test

import (
	"github.com/atong007/kit/crypto/aes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAesGCM_DecryptedCode(t *testing.T) {
	type fields struct {
		encKey string
	}
	type args struct {
		enCode string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantCode string
		wantErr  bool
	}{
		{"success", fields{"1ab3adcf15eeb01bc812aae31b24efb5"}, args{"0cfaeaa3999edf03f2458fdd1b66469c939cda0ed1435d4e5cb359604f8010d6a9"}, "test", false},
		{"decrypt-err", fields{"1ab3adcf15eeb01bc812aae31b24efb5"}, args{"123456"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := aes.NewAesGCM(tt.fields.encKey)
			assert.NoError(t, err)
			gotCode, err := a.DecryptedCode(tt.args.enCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecryptedCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCode != tt.wantCode {
				t.Errorf("DecryptedCode() gotCode = %v, want %v", gotCode, tt.wantCode)
			}
		})
	}
}

func TestAesGCM_EncryptedCode(t *testing.T) {
	type fields struct {
		encKey string
	}
	type args struct {
		code string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantEnCode string
		wantErr    bool
	}{
		{"success", fields{"1ab3adcf15eeb01bc812aae31b24efb5"}, args{"test"}, "test", false},
		{"encKey-err", fields{"123"}, args{""}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := aes.NewAesGCM(tt.fields.encKey)
			if tt.wantErr {
				assert.Error(t, err)
				return
			} else {
				assert.NoError(t, err)
			}
			gotEnCode, err := a.EncryptedCode(tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncryptedCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotEnCode, err = a.DecryptedCode(gotEnCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncryptedCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotEnCode != tt.wantEnCode {
				t.Errorf("EncryptedCode() gotEnCode = %v, want %v", gotEnCode, tt.wantEnCode)
			}
		})
	}
}
