package runtimedir

import (
	"strings"
	"testing"
)

func TestMachineHash(t *testing.T) {
	tests := []struct {
		name        string
		machineName string
	}{
		{"✅ short name", "minikube"},
		{"✅ very long name", strings.Repeat("a", 200)},
		{"✅ name with dots and dashes", "weird-name.with-dots"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := machineHash(tt.machineName)
			if len(got) != 32 {
				t.Errorf("hash length = %d, want 32", len(got))
			}
			if again := machineHash(tt.machineName); got != again {
				t.Errorf("hash not deterministic: %s != %s", got, again)
			}
		})
	}
}

func TestMachineHashNoCollision(t *testing.T) {
	// Distinct machine names must not collide on the truncated hash.
	names := []string{
		"minikube", "minikube2", "profile-a", "profile-b",
		strings.Repeat("x", 64), strings.Repeat("y", 64),
	}
	seen := make(map[string]string, len(names))
	for _, n := range names {
		h := machineHash(n)
		if prev, ok := seen[h]; ok {
			t.Errorf("❌ collision: %q and %q both hash to %s", prev, n, h)
		}
		seen[h] = n
	}
}

// TestResolveBudget is the unit-level proof of the fix: no matter how long
// the machine name is, the resolved path fits within the tightest sun_path
// limit (macOS, 104). Asserting against the tightest keeps it portable.
func TestResolveBudget(t *testing.T) {
	tests := []struct {
		name        string
		baseDir     string
		machineName string
		sockName    string
	}{
		{"✅ typical", "/tmp/501", "minikube", "monitor"},
		{"✅ worst-case machine name", "/tmp/501", strings.Repeat("m", 64), "monitor"},
		{"✅ longest socket name in tree", "/tmp/501", strings.Repeat("m", 64), "vmnet-helper.sock-krun.sock"},
		{"✅ high uid base", "/tmp/4294967295", strings.Repeat("m", 200), "monitor"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Resolve(tt.baseDir, tt.machineName, tt.sockName)
			if len(got) > SunPathDarwin {
				t.Errorf("resolved path is %d bytes, over the tightest limit (%d):\n  %s",
					len(got), SunPathDarwin, got)
			}
		})
	}
}
