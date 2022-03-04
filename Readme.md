# gosema
gosema is a simple ipc sysv implementation with golang, nothing fancy..   
gosema use golang/x/sys/unix package and use linux kernel [ipc(2)](https://man7.org/linux/man-pages/man2/ipc.2.html) and [shmget(2)](https://man7.org/linux/man-pages/man2/shmget.2.html), [semctl(2)](https://man7.org/linux/man-pages/man2/semctl.2.html) and other linux ipc handlers

autor: sina@snix.ir

### build and run
you need to install golang compiler 
```console
# git clone https://git.snix.ir/gosema
# cd gosema && go build
# ./gosema [server|client]
```
