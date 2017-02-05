// +build !windows

package dirlock

import (
	"fmt"
	"os"
	"syscall"
)

// DirLock 应该是一个目录锁
type DirLock struct {
	dir string
	f   *os.File
}

// New 生成一个目录锁
func New(dir string) *DirLock {
	return &DirLock{
		dir: dir,
	}
}

func (l *DirLock) Lock() error {
	f, err := os.Open(l.dir)
	if err != nil {
		return err
	}
	l.f = f
	// 需要认真解读，目录无法理解
	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		return fmt.Errorf("cannot flock directory %s - %s", l.dir, err)
	}
	return nil
}

func (l *DirLock) Unlock() error {
	defer l.f.Close()
	return syscall.Flock(int(l.f.Fd()), syscall.LOCK_UN)
}
