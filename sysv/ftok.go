package sysv

import (
	"golang.org/x/sys/unix"
)

func Ftok(p string, id uint) (uint, error) {
	fss := &unix.Stat_t{}
	if err := unix.Stat(p, fss); err != nil {
		return 0, err
	}

	return uint((uint(fss.Ino) & 0xffff) |
		uint((fss.Dev&0xff)<<16) |
		((id & 0xff) << 24)), nil
}
