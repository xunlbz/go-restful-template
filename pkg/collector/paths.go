package collector

import (
	"path/filepath"
)

const (
	DefaultProcMountPoint = "/proc"
)

func procFilePath(name string) string {
	return filepath.Join(DefaultProcMountPoint, name)
}
