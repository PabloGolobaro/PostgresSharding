package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{name: "first", path: ".", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotConfig, err := LoadConfig(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, dsn := range gotConfig.DSNS {
				t.Logf("LoadConfig() config.dsn = %v", dsn)
			}

		})
	}
}
