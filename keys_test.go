package blizzauth

import (
	"reflect"
	"testing"
)

func TestNewKeys(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want *keys
	}{
		{
			"nominal",
			args{
				"test",
			},
			&keys{
				id:     "bob'syourun cle roy's yourneighbor",
				secret: "maryjosephandroy",
				name:   "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newKeys(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}
