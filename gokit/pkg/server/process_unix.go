//go:build aix || darwin || dragonfly || freebsd || (js && wasm) || linux || netbsd || openbsd || solaris || plan9
// +build aix darwin dragonfly freebsd js,wasm linux netbsd openbsd solaris plan9

package server

import (
    "strconv"
    "syscall"
)

// SetSysProcAttr - set Process group ID and owner (run on behalf)
func SetSysProcAttr(p *Process) error {
    sysProcAttr := &syscall.SysProcAttr{
        Setpgid: true, // Set Process group ID to Pgid, or, if Pgid == 0, to new pid.
        Pgid:    0,    // Child's Process group ID if Setpgid.
    }

    // set owner
    if p.user != nil {
        uid, err := strconv.Atoi(p.user.Uid)
        if err != nil {
            return err
        }
        gid, err := strconv.Atoi(p.user.Gid)
        if err != nil {
            return err
        }
        sysProcAttr.Credential = &syscall.Credential{
            Uid: uint32(uid),
            Gid: uint32(gid),
        }
    }

    // set the attributes
    p.Cmd.SysProcAttr = sysProcAttr

    return nil
}

// Kill the entire Process group.
func Kill(p *Process) error {
    ProcessGroup := 0 - p.Cmd.Process.Pid
    return syscall.Kill(ProcessGroup, syscall.SIGKILL)
}

// Signal sends a signal to the Process
func Signal(p *Process, sig syscall.Signal) error {
    return syscall.Kill(p.Cmd.Process.Pid, sig)
}
