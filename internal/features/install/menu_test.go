package install

import (
	"strings"
	"testing"

	"snailproxy/internal/terminal"
)

func TestSelectMapping(t *testing.T) {
	tests := []struct {
		input string
		want  Action
	}{
		{input: "1\n", want: ActionLocal},
		{input: "2\n", want: ActionOnline},
		{input: "3\n", want: ActionService},
		{input: "0\n", want: ActionReturn},
	}

	for _, tt := range tests {
		t.Run(strings.TrimSpace(tt.input), func(t *testing.T) {
			withInput(t, tt.input)

			action, err := Select()
			if err != nil {
				t.Fatalf("Select() error = %v", err)
			}
			if action != tt.want {
				t.Fatalf("action = %d, want %d", action, tt.want)
			}
		})
	}
}

func TestFormatAssetSize(t *testing.T) {
	tests := []struct {
		size int64
		want string
	}{
		{size: 512, want: "512 B"},
		{size: 1024, want: "1.0 KiB"},
		{size: 5 * 1024 * 1024, want: "5.0 MiB"},
	}
	for _, tt := range tests {
		if got := formatAssetSize(tt.size); got != tt.want {
			t.Errorf("formatAssetSize(%d) = %q, want %q", tt.size, got, tt.want)
		}
	}
}

func withInput(t *testing.T, input string) {
	t.Helper()
	restore := terminal.SetInput(strings.NewReader(input))
	t.Cleanup(restore)
}
