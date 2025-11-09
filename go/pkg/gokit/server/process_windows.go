package server

import "syscall"

func SetSysProcAttr(p *Process) error {
    return nil
}

// Kill the entire Process group.
func Kill(p *Process) error {
    return nil
}

// Signal sends a signal to the Process
func Signal(p *Process, sig syscall.Signal) error {
    return nil
}
