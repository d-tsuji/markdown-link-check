package markdownlink

import (
	"testing"

	"go.uber.org/goleak"
)

func TestCheck(t *testing.T) {
	type args struct {
		config Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "from url", args: args{config: Config{args: []string{"https://raw.githubusercontent.com/d-tsuji/flower/master/README.md"}}}, wantErr: false},
		{name: "from local file", args: args{config: Config{args: []string{"testdata/README.md"}}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Check(tt.args.config); (err != nil) != tt.wantErr {
				defer goleak.VerifyNone(t)
				t.Errorf("Check() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
