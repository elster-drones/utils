package audit

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

// ErrLoginuidUnset is returned when the loginuid is not set (value 4294967295)
var ErrLoginuidUnset = errors.New("loginuid not set")

func GetPidLoginuid(pid int32) (uint32, error) {
	procfile := fmt.Sprintf("/proc/%d/loginuid", pid)
	f, e := os.Open(procfile)
	if e != nil {
		return 0, e
	}
	defer f.Close()

	var buf []byte = make([]byte, 10)
	n, e := f.Read(buf)
	if e != nil && e != io.EOF && e != io.ErrUnexpectedEOF {
		return 0, e
	}

	uids := string(buf[:n])
	uid, e := strconv.Atoi(uids)
	if e != nil {
		return 0, e
	}

	// 4294967295 (0xFFFFFFFF) means loginuid is not set
	if uint32(uid) == 0xFFFFFFFF {
		return 0, ErrLoginuidUnset
	}

	return uint32(uid), nil
}
