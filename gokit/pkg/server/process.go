package server

import (
    "fmt"
    "github.com/ncraft-io/ncraft/go/pkg/ncraft/logs"
    "os"
    "os/exec"
    "os/user"
    "syscall"
    "time"
)

// Pid struct run.yml
type Pid struct {
    Follow string `yaml:",omitempty"`
    Parent string `yaml:",omitempty"`
    Child  string `yaml:",omitempty"`
}

type ProcessConfig struct {
    command    []string
    Cwd        string            `yaml:",omitempty" json:",omitempty"`
    Env        map[string]string `yaml:",omitempty" json:",omitempty"`
    Log        logs.Config       `yaml:",omitempty" json:",omitempty"`
    Stderr     logs.Config       `yaml:",omitempty" json:",omitempty"`
    Logger     string            `yaml:",omitempty" json:",omitempty"`
    Require    []string          `yaml:",omitempty"`
    RequireCmd string            `yaml:"require_cmd,omitempty"`
    PostExit   string            `yaml:"post_exit,omitempty"`
    User       string            `yaml:",omitempty" json:",omitempty"`
    Wait       uint              `yaml:",omitempty"`
    Retries    int               `yaml:",omitempty"`
    Pid        `yaml:",omitempty" json:",omitempty"`
    cli        bool

    configFile string
    ctl        string
    log        bool
    user       *user.User
}

type Process struct {
    *ProcessConfig
    //Logger
    //LoggerStderr Logger
    Cmd *exec.Cmd

    ErrorChan chan error
    QuitChan  chan struct{}

    StartTime time.Time
    EndTime   time.Time
}

// SetEnv set environment variables - If the Cmd.Env contains duplicate
// environment keys, only the last value in the slice for each duplicate
// key is used.
func (p *Process) SetEnv(env []string) {
    if p.Env != nil {
        for k, v := range p.Env {
            env = append(env, fmt.Sprintf("%s=%s", k, v))
        }
        p.Cmd.Env = env
    }
}

// Start runs the command
func (p *Process) Start() (*Process, error) {
    // command obtained from Config parent
    p.Cmd = exec.Command(p.command[0], p.command[1:]...)

    // change working directory
    if p.Cwd != "" {
        p.Cmd.Dir = p.Cwd
    }

    // set environment variables
    p.SetEnv(os.Environ())

    // set sysProcAttr
    if err := SetSysProcAttr(p); err != nil {
        return nil, err
    }

    var (
        //prStdout, prStderr *os.File
        pwStdout, pwStderr *os.File
        //e                   error
    )
    // log only if are available loggers
    //if p.Logger.IsLogging() && p.LoggerStderr.IsLogging() {
    //	// create the pipes for Stdout
    //	prStdout, pwStdout, e = os.Pipe()
    //	if e == nil {
    //		p.Cmd.Stdout = pwStdout
    //		go p.Logger.Log(prStdout)
    //	}
    //	prStderr, pwStderr, e = os.Pipe()
    //	if e == nil {
    //		p.Cmd.Stderr = pwStderr
    //		go p.LoggerStderr.Log(prStderr)
    //	}
    //} else if p.Logger.IsLogging() {
    //	// create the pipes for Stdout
    //	prStdout, pwStdout, e = os.Pipe()
    //	if e == nil {
    //		p.Cmd.Stdout = pwStdout
    //		p.Cmd.Stderr = pwStdout
    //		go p.Logger.Log(prStdout)
    //	}
    //} else if p.LoggerStderr.IsLogging() {
    //	// create the pipes for Stdout
    //	prStderr, pwStderr, e = os.Pipe()
    //	if e == nil {
    //		p.Cmd.Stderr = pwStderr
    //		go p.LoggerStderr.Log(prStderr)
    //	}
    //}

    // Start the Process
    if err := p.Cmd.Start(); err != nil {
        return nil, err
    }

    // set start time
    p.StartTime = time.Now()

    // wait Process to finish in a goroutine
    go p.Wait(pwStdout, pwStderr)

    return p, nil
}

// Wait - wait Process to finish
func (p *Process) Wait(stdout, stderr *os.File) {
    err := p.Cmd.Wait()
    if stdout != nil {
        stdout.Close()
        close(p.QuitChan)
    }
    if stderr != nil {
        stderr.Close()
    }
    p.ErrorChan <- err
}

// Kill the entire Process group.
func (p *Process) Kill() error {
    return Kill(p)
}

// Pid return Process PID
func (p *Process) Pid() int {
    if p.Cmd == nil || p.Cmd.Process == nil {
        return 0
    }
    return p.Cmd.Process.Pid
}

// Signal sends a signal to the Process
func (p *Process) Signal(sig syscall.Signal) error {
    return Signal(p, sig)
}

// NewProcess return Process instance
func NewProcess(cfg *Config) *Process {
    qch := make(chan struct{})
    return &Process{
        //Config: cfg,
        //Logger: &LogWriter{
        //	logger: NewLogger(cfg, qch),
        //},
        //LoggerStderr: &LogWriter{
        //	logger: NewStderrLogger(cfg),
        //},
        ErrorChan: make(chan error, 1),
        QuitChan:  qch,
        StartTime: time.Now(),
    }
}
