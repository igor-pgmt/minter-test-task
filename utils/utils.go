package utils

import (
	"fmt"
	"syscall"
)

func MaxOpenFiles() (*syscall.Rlimit, error) {
	var rLimit syscall.Rlimit

	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return nil, fmt.Errorf("Error Getting Rlimit: %v", err)
	}

	if rLimit.Cur < rLimit.Max {
		rLimit.Cur = rLimit.Max
		err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
		if err != nil {
			return nil, fmt.Errorf("Error Setting Rlimit: %v", err)
		}
	}

	return &rLimit, nil
}
