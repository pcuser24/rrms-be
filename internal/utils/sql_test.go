package utils

import "testing"

func TestSequelizePlaceholders(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SequelizePlaceholders(tt.args.query); got != tt.want {
				t.Errorf("SequelizePlaceholders() = %v, want %v", got, tt.want)
			}
		})
	}
}
