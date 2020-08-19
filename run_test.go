package mlc

import (
	"context"
	"testing"

	"github.com/google/go-github/v32/github"
)

func TestRun(t *testing.T) {
	type args struct {
		cf *config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "single mode",
			args: args{
				cf: &config{
					ctx:      context.TODO(),
					client:   github.NewClient(nil),
					owner:    "d-tsuji",
					repo:     "flower",
					branch:   "master",
					allMode:  false,
					filePath: "README.md",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Run(tt.args.cf); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
