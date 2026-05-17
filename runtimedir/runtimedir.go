package runtimedir

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

const (
	SunPathLinux  = 108
	SunPathDarwin = 104
)

func machineHash(machineName string) string {
	sum := sha256.Sum256([]byte(machineName))
	return hex.EncodeToString(sum[:])[:32]
}

func Resolve(baseDir, machineName, socketName string) string {
	return filepath.Join(baseDir, "minikube", machineHash(machineName), socketName)
}

func runtimeBaseDir() string {
	uid := strconv.Itoa(os.Getuid())
	switch runtime.GOOS {
	case "linux":
		if xdg := os.Getenv("XDG_RUNTIME_DIR"); xdg != "" {
			if fi, err := os.Stat(xdg); err == nil && fi.IsDir() {
				return xdg
			}
		}
		return filepath.Join("/tmp", uid)
	case "darwin":
		return filepath.Join("/tmp", uid)
	default:
		return os.TempDir()
	}
}

func SocketPath(machineName, socketName string) string {
	return Resolve(runtimeBaseDir(), machineName, socketName)
}

func EnsureDir(socketPath string) error {
	dir := filepath.Dir(socketPath)
	return os.MkdirAll(dir, 0700)
}
