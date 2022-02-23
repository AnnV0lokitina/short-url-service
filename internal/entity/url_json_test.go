package entity

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestURL_MarshalJSON(t *testing.T) {
	type fields struct {
		checksum string
		full     string
	}

	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "test json marshal positive",
			fields: fields{
				checksum: "checksum",
				full:     "full",
			},
			want:    []byte("{\"checksum\":\"checksum\",\"full_url\":\"full\"}"),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &URL{
				checksum: tt.fields.checksum,
				full:     tt.fields.full,
			}
			got, err := json.Marshal(u)
			if !tt.wantErr(t, err, fmt.Sprintf("MarshalJSON()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "MarshalJSON()")
		})
	}
}

func TestURL_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		want    *URL
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "test json unmarshal positive",
			want: NewURL("full", "checksum"),
			args: args{
				data: []byte("{\"checksum\":\"checksum\",\"full_url\":\"full\"}"),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got URL
			tt.wantErr(t, json.Unmarshal(tt.args.data, &got), fmt.Sprintf("UnmarshalJSON(%v)", tt.args.data))
			assert.ObjectsAreEqual(got, *tt.want)
		})
	}
}
