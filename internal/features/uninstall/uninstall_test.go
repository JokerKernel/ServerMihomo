package uninstall

import (
	"context"
	"errors"
	"strings"
	"testing"

	"snailproxy/internal/domain/mihomo"
	"snailproxy/internal/feature"
	"snailproxy/internal/infra/platform"
	"snailproxy/internal/terminal"
)

func TestRunCancelDoesNotCreatePlatformManager(t *testing.T) {
	restore := terminal.SetInput(strings.NewReader("n\n"))
	t.Cleanup(restore)

	runtime := uninstallRuntime{managerErr: errors.New("platform manager should not be created")}
	err := (Feature{}).Run(context.Background(), runtime)
	if !errors.Is(err, feature.ErrReturn) {
		t.Fatalf("Run() error = %v, want ErrReturn", err)
	}
}

func TestRunConfirmedUninstalls(t *testing.T) {
	restore := terminal.SetInput(strings.NewReader("y\n"))
	t.Cleanup(restore)

	manager := &uninstallManager{}
	if err := (Feature{}).Run(context.Background(), uninstallRuntime{manager: manager}); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if !manager.uninstalled {
		t.Fatal("Uninstall() was not called")
	}
}

type uninstallRuntime struct {
	manager    platform.Manager
	managerErr error
}

func (uninstallRuntime) Terminal() terminal.Terminal {
	return terminal.Default()
}

func (uninstallRuntime) NewMihomoStore() mihomo.Store {
	return mihomo.Store{}
}

func (r uninstallRuntime) NewPlatformManager() (platform.Manager, error) {
	return r.manager, r.managerErr
}

type uninstallManager struct {
	uninstalled bool
}

func (*uninstallManager) PrepareBinary(context.Context, string, string, bool) error { return nil }
func (*uninstallManager) InstallService(context.Context) error                      { return nil }
func (*uninstallManager) StartService(context.Context) error                        { return nil }
func (*uninstallManager) RestartService(context.Context) error                      { return nil }
func (*uninstallManager) StopService(context.Context) error                         { return nil }
func (*uninstallManager) WriteProxyEnvironment(context.Context) error               { return nil }
func (*uninstallManager) ClearProxyEnvironment(context.Context) error               { return nil }
func (m *uninstallManager) Uninstall(context.Context) error {
	m.uninstalled = true
	return nil
}
