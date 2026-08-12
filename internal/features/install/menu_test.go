package install

import (
	"errors"
	"strings"
	"testing"

	"snailproxy/internal/infra/github"
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

func TestSelectAssetSupportsReturnAliases(t *testing.T) {
	assets := []github.Asset{{Name: "mihomo-linux-amd64-v1.0.0.gz", Size: 1024}}
	for _, input := range []string{"0\n", "q\n", "Q\n", "exit\n"} {
		t.Run(strings.TrimSpace(input), func(t *testing.T) {
			withInput(t, input)

			_, err := SelectAsset(assets)
			if !errors.Is(err, errReturnToInstallMenu) {
				t.Fatalf("SelectAsset() error = %v, want errReturnToInstallMenu", err)
			}
		})
	}
}

func TestSelectAssetReturnsSelectedAsset(t *testing.T) {
	assets := []github.Asset{
		{Name: "first.gz", Size: 1024},
		{Name: "second.gz", Size: 2048},
	}
	withInput(t, "2\n")

	asset, err := SelectAsset(assets)
	if err != nil {
		t.Fatalf("SelectAsset() error = %v", err)
	}
	if asset.Name != "second.gz" {
		t.Fatalf("asset.Name = %q, want second.gz", asset.Name)
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
