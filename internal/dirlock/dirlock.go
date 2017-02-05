// +build !windows

package dirlock

import "os"

type DirLock struct {
	dir string
	f   *os.File
}
