package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"reflect"
	"testing"
)

func TestWithContext(t *testing.T) {
	type (
		args struct {
			ctx   context.Context
			key   FieldKey
			value any
		}

		field struct {
			Username any `json:"username"`
		}
	)

	const username = "unemil"

	var (
		buf bytes.Buffer
		f   field

		test = struct {
			args args
			want any
		}{
			args: args{
				ctx:   context.Background(),
				key:   "username",
				value: username,
			},
			want: username,
		}
	)

	r, w, _ := os.Pipe()
	os.Stdout = w

	logger = newLogger()

	ctx := WithContext(test.args.ctx, test.args.key, test.args.value)
	Info(ctx, "test")

	w.Close()
	defer r.Close()

	io.Copy(&buf, r)
	json.Unmarshal(buf.Bytes(), &f)

	if !reflect.DeepEqual(f.Username, test.want) {
		t.Errorf("got: %v, want: %v", f.Username, test.want)
	}
}
