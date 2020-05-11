package blizzauth

import (
	"testing"
)

func TestGetAuth(t *testing.T) {
	type args struct {
		apiName string
	}
	tests := []struct {
		name    string
		args    args
		wantT   *Auth
		wantErr bool
	}{
		{
			"test",
			args{
				"cap",
			},
			&Auth{
				clientName: "cap",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotT, err := GetAuth(tt.args.apiName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAuth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotT.clientName != tt.wantT.clientName {
				t.Errorf("GetAuth() = %v, want %v", gotT, tt.wantT)
			}
		})
	}
}

func TestAuth_GetAccessToken(t *testing.T) {

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			"normal",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := GetAuth("cap")
			if err != nil {
				t.Error(err)
				return
			}
			got, err := a.GetAccessToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("Auth.GetAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == "" {
				t.Errorf("Auth.GetAccessToken() = no key")
			}
		})
	}
}
