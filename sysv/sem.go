package sysv

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	SYS_GETPID int = 0x0B + iota
	SYS_GETVAL
	SYS_SETVAL int = 0x010
)

const (
	non = iota - one
	one = 0x01
	zro = 0x00
)

type oprations struct {
	num uint16
	opt int16
	flg int16
}

type Semaphore struct {
	nums int
	smid int
}

func NewSemGet(key, nums, flags int) (*Semaphore, error) {
	id, _, errno := unix.Syscall(
		unix.SYS_SEMGET, uintptr(key),
		uintptr(nums), uintptr(flags),
	)
	if errno != 0 {
		return nil, error(errno)
	}
	return &Semaphore{nums: nums, smid: int(id)}, nil
}

func (sem *Semaphore) GetVal(semnum int) (int, error) {
	vl, _, errno := unix.Syscall(
		unix.SYS_SEMCTL, uintptr(sem.smid),
		uintptr(semnum), uintptr(SYS_GETVAL),
	)
	if errno != 0 {
		return int(vl), error(errno)
	}
	return int(vl), nil
}

func (sem *Semaphore) SetVal(semnum, v int) error {
	_, _, errno := unix.Syscall6(
		unix.SYS_SEMCTL, uintptr(sem.smid),
		uintptr(semnum), uintptr(SYS_SETVAL), uintptr(v), 0, 0,
	)
	if errno != 0 {
		return error(errno)
	}
	return nil
}

func (sem *Semaphore) GetPID(semnum int) (int, error) {
	pid, _, errno := unix.Syscall(
		unix.SYS_SEMCTL, uintptr(sem.smid),
		uintptr(semnum), uintptr(SYS_GETPID),
	)
	if errno != 0 {
		return int(pid), error(errno)
	}
	return int(pid), nil
}

func (sem *Semaphore) setops(n uint16, o, f int16) error {
	opt := oprations{num: n, opt: o, flg: f}
	_, _, errno := unix.Syscall(
		unix.SYS_SEMOP, uintptr(sem.smid),
		uintptr(unsafe.Pointer(&opt)), uintptr(sem.nums),
	)

	if errno != 0 {
		return error(errno)
	}
	return nil
}

func (sem *Semaphore) Lock(semnum int) error {
	return sem.setops(uint16(semnum), non, zro)
}

func (sem *Semaphore) Unlock(semnum int) error {
	return sem.setops(uint16(semnum), one, zro)
}

func (sem *Semaphore) DelSem() error {
	_, _, errno := unix.Syscall(
		unix.SYS_SEMCTL, uintptr(sem.smid),
		uintptr(0), uintptr(unix.IPC_RMID),
	)
	if errno != 0 {
		return error(errno)
	}
	return nil
}
