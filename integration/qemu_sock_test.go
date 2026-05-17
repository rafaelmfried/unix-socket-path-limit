//go:build integration

package integration

import (
	"context"
	"io"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/rafaelmfried/unix-socket-path-limit/runtimedir"
	"github.com/testcontainers/testcontainers-go"
	tcexec "github.com/testcontainers/testcontainers-go/exec"
	"github.com/testcontainers/testcontainers-go/wait"
)

func startQEMUContainer(ctx context.Context, t *testing.T) testcontainers.Container {
	t.Helper()

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    path.Join("..", "docker"),
			Dockerfile: "Dockerfile.qemu",
		},
		WaitingFor: wait.ForExec([]string{"qemu-system-x86_64", "--version"}).
			WithStartupTimeout(2 * time.Minute),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start QEMU container: %v", err)
	}

	t.Cleanup(func() { _ = container.Terminate(ctx) })

	return container
}

func runQEMU(ctx context.Context, container testcontainers.Container, sockPath string) (string, error) {
	if code, _, err := container.Exec(ctx, []string{"mkdir", "-p", path.Dir(sockPath)}); err != nil || code != 0 {
		return "", err
	}

	cmd := []string{
		"timeout", "3",
		"qemu-system-x86_64", "-machine", "none", "-display", "none",
		"-qmp", "unix:" + sockPath + ",server=on,wait=off",
	}

	_, reader, err := container.Exec(ctx, cmd, tcexec.Multiplexed())
	if err != nil {
		return "", err
	}

	output, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func TestQEMUFailsWithLongSocketPath(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	container := startQEMUContainer(ctx, t)

	longSockPath := path.Join("/tmp", strings.Repeat("a", 100), "/", strings.Repeat("b", 100), "/.minikube/machines/demo/monitor")

	if len(longSockPath) <= runtimedir.SunPathLinux {
		t.Fatalf("Test setup error: path is only %d characters long, not long enough to trigger the limit", len(longSockPath))
	}

	output, err := runQEMU(ctx, container, longSockPath)
	if err != nil {
		t.Fatalf("Failed to run QEMU: %v", err)
	}

	if !strings.Contains(strings.ToLower(output), "too long") {
		t.Errorf("Expected qemu to reject the over-limit socket path, got: %s", output)
	}
}

func TestQEMUAcceptsRuntimedirPath(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	container := startQEMUContainer(ctx, t)

	machineName := strings.Repeat("m", 64)
	shortPath := runtimedir.Resolve("/run", machineName, "monitor")

	if len(shortPath) >= runtimedir.SunPathLinux {
		t.Fatalf("Test setup error: path is %d characters long, not short enough to be accepted", len(shortPath))
	}

	output, err := runQEMU(ctx, container, shortPath)
	if err != nil {
		t.Fatalf("Failed to run QEMU: %v", err)
	}

	if strings.Contains(strings.ToLower(output), "too long") {
		t.Errorf("Did not expect qemu to reject the runtimedir path, got: %s", output)
	}
}
