package musicbot

import (
	"os"
	"path/filepath"
	"runtime"
)

func applyLibDaveRuntimePath() {
	oncePaths.Do(initPaths)
	libDir := filepath.Join(baseDir, "libdave", "lib")
	if _, err := os.Stat(libDir); err != nil {
		return
	}
	sep := string(os.PathListSeparator)
	switch runtime.GOOS {
	case "darwin":
		current := os.Getenv("DYLD_LIBRARY_PATH")
		os.Setenv("DYLD_LIBRARY_PATH", libDir+sep+current)
	case "linux":
		current := os.Getenv("LD_LIBRARY_PATH")
		os.Setenv("LD_LIBRARY_PATH", libDir+sep+current)
	}
}
