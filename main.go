package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"syscall"
)

const memoryMount = "/sys/fs/cgroup/memory"

func main() {
	if os.Args[0] == "/proc/self/exe" {
		fmt.Printf("now pid is %d", syscall.Getpid())
		cmd := exec.Command("sh", "-c", `stress --vm-bytes 200m --vm-keep -m 1`)
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}
	}

	cmd := exec.Command("/proc/self/exe")
	//cmd := exec.Command("sh")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUSER | syscall.CLONE_NEWNET,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      syscall.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      syscall.Getgid(),
				Size:        1,
			},
		},
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%v", cmd.Process.Pid)
		os.Mkdir(path.Join(memoryMount, "test_memory", "tasks"), 0755)
		ioutil.WriteFile(path.Join(memoryMount, "test_memory", "memory.limit_in_bytes"),
			[]byte("100m"), 0664)
		cmd.Process.Wait()
	}
}
