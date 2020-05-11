package blizzauth

import (
	"reflect"
	"testing"
)

func TestGetAuth(t *testing.T) {
	type args struct {
		apiName string
	}
	tests := []struct {
		name  string
		args  args
		wantT *Auth
	}{
		{
			"test",
			args{
				"test_api",
			},
			&Auth{
				clientName: "test_api",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotT := GetAuth(tt.args.apiName); !reflect.DeepEqual(gotT, tt.wantT) {
				t.Errorf("GetAuth() = %v, want %v", gotT, tt.wantT)
			}
		})
	}
}
