package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"gosema/sysv"
	"os"
	"time"

	"golang.org/x/sys/unix"
)

/*

	an example of ipc sysv with golang - this can be done with newer apis
	like shm_open(), /dev/shm or mmap() too, sysv is the oldest way to do
	this, and it's commonly used.
	author: sina@snix.ir

*/

const (
	//SHMSEMFLG = 0666 | unix.IPC_CREAT | unix.IPC_EXCL
	SHMSEMFLG = 0666 | unix.IPC_CREAT // open shm and sem with this flag
	SEMSHEMID = 0x12                  // shm and sem id
	NUMSOFSEM = 0x01                  // number of sems
	SHRDMEMRY = 0x20                  // size of shared memory (32 byte)
)

const SEM_A = 0x00

func clntOrServ() bool {
	if os.Args[1] == server {
		return true
	}
	if os.Args[1] != client {
		usage(os.Args[0])
	}
	return false
}

func main() {
	defer defexit()

	if len(os.Args) < 2 {
		usage(os.Args[0])
	}

	serv := clntOrServ()

	semkey, err := sysv.Ftok(os.DevNull, SEMSHEMID)
	if err != nil {
		fatal(err)
	}

	shmkey, err := sysv.Ftok(os.DevNull, SEMSHEMID)
	if err != nil {
		fatal(err)
	}

	shmidd, err := unix.SysvShmGet(int(shmkey), SHRDMEMRY, SHMSEMFLG)
	if err != nil {
		fatal(err)
	}

	defer unix.SysvShmCtl(shmidd, unix.IPC_RMID, &unix.SysvShmDesc{})

	mem, err := unix.SysvShmAttach(shmidd, 0, 0)
	if err != nil {
		fatal(err)
	}

	defer unix.SysvShmDetach(mem)

	if len(mem) != SHRDMEMRY {
		fatal(errNotSame)
	}

	logst("shared memory is now ready to use", SHRDMEMRY)

	sempht, err := sysv.NewSemGet(int(semkey), NUMSOFSEM, SHMSEMFLG)
	if err != nil {
		fatal(err)
	}

	defer sempht.DelSem()

	// set semval to 1, which force client to wait for semaphore

	switch {
	case serv:
		if err := sempht.SetVal(SEM_A, 1); err != nil {
			warng(err)
		}

		for {

			l, _ := sempht.GetVal(SEM_A) // value of semaphore
			s, _ := sempht.GetPID(SEM_A) // pid that holds semaphore
			fmt.Printf("sem locked by process pid: %d, semval: %d\n", s, l)

			// lock sem[0], make semval = semval-1 and if semval is 0 then no one
			// can lock semaphore any more
			if err := sempht.Lock(SEM_A); err != nil {
				warng(err)
				continue
			}

			// l, _ = sempht.GetVal(SEM_A)
			// fmt.Println(l)

			// write some dummy random data to shared mem
			if err := makeRand(mem); err != nil {
				warng(err)
			}

			time.Sleep(time.Second) // every one sec
			fmt.Println("i just put some random jokes in memory.. yeay")
			if err := sempht.Unlock(SEM_A); err != nil {
				warng(err)
			} //unlock sem[0], make semval = semval+1
		}
	}

	for {
		if err := sempht.Lock(SEM_A); err != nil {
			warng(err)
			continue // do not unlock it if you haven't locked it
		}

		fmt.Printf("data: %s\n", hex.EncodeToString(mem))
		if err := sempht.Unlock(SEM_A); err != nil {
			warng(err)
		}
	}
}

func makeRand(b []byte) error {
	_, err := rand.Read(b)
	return err

}
