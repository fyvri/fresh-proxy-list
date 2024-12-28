package utils

import (
	"net/url"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		rawURL string
	}
	tests := []struct {
		name      string
		args      args
		want      *url.URL
		wantError error
	}{
		{
			name: "ValidURL",
			args: args{
				rawURL: testRawURL,
			},
			want: &url.URL{
				Scheme: testScheme,
				Host:   testHost,
			},
			wantError: nil,
		},
		{
			name: "ValidURLWithFullURL",
			args: args{
				rawURL: testFullURL,
			},
			want: &url.URL{
				Scheme:   testScheme,
				Host:     testHost,
				Path:     testPath,
				RawQuery: testRawQuery,
			},
			wantError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewURLParser()
			got, err := u.Parse(tt.args.rawURL)

			if tt.wantError != nil {
				t.Errorf(expectedErrorButGotMessage, "Parse()", err, tt.wantError)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(expectedErrorButGotMessage, "Parse()", err, tt.wantError)
			}
		})
	}
}
