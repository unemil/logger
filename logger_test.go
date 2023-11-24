package logger

import (
	"context"
	"reflect"
	"testing"
)

func TestWithContext(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   any
		value any
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "AdditionTest",
			args: args{
				ctx:   context.Background(),
				key:   "key",
				value: "addition",
			},
			want: context.WithValue(context.Background(), "key", "addition"),
		},
		{
			name: "DuplicationTest",
			args: args{
				ctx:   context.WithValue(context.Background(), "key", "duplication"),
				key:   "key",
				value: "duplication",
			},
			want: context.WithValue(context.WithValue(context.Background(), "key", "duplication"), "key", "duplication"),
		},
		{
			name: "OverwritingTest",
			args: args{
				ctx:   context.WithValue(context.Background(), "key", "default"),
				key:   "key",
				value: "overwriting",
			},
			want: context.WithValue(context.WithValue(context.Background(), "key", "default"), "key", "overwriting"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithContext(tt.args.ctx, tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithContext() = %v, want %v", got, tt.want)
			}
		})
	}
}
