package logsutils

import (
    "fmt"
    "github.com/ncraft-io/ncraft/go/pkg/ncraft/config"
    "github.com/ncraft-io/ncraft/go/pkg/ncraft/logs"
    "os"
    "path"
    "strings"
    "syscall"
)

func GetLogFileName() string {
    logConfig := &logs.Config{}
    if err := config.Get("logs").Scan(logConfig); err != nil {
        return ""
    }

    if len(logConfig.File.Path) == 0 {
        return ""
    }
    return HandleFileName(logConfig.File.Path)
}

func RedirectStderr(f *os.File) {
    err := syscall.Dup2(int(f.Fd()), int(os.Stderr.Fd()))
    if err != nil {
        fmt.Println("Failed to redirect stderr to file: ", err.Error())
    }
}

func HandleFileName(filename string) string {
    filename = path.Clean(filename)
    parts := make([]string, 0)
    var ret string
    paths := strings.Split(filename, string(os.PathSeparator))
    for _, v := range paths {
        val := handleTemplateFileName(v)
        if len(val) > 0 {
            parts = append(parts, val)
        }
    }

    if path.IsAbs(filename) {
        ret = string(os.PathSeparator) + path.Join(parts...)
    } else {
        ret = path.Join(parts...)
    }
    return ret
}

func handleTemplateFileName(template string) string {
    // foo1{hostname}foo2{port}foo3
    lefts := make([]int, 0)
    rights := make([]int, 0)

    size := len(template)
    for i := 0; i < size; i++ {
        if template[i] == '{' {
            lefts = append(lefts, i)
        } else if template[i] == '}' {
            rights = append(rights, i)
        }
    }

    leftSize := len(lefts)
    rightSize := len(rights)
    var minSize int
    if leftSize < rightSize {
        minSize = leftSize
    } else {
        minSize = rightSize
    }

    ret := template
    for i := minSize - 1; i >= 0; i-- {
        variableName := ret[lefts[i]+1 : rights[i]]
        v := os.Getenv(variableName)
        ret = ret[:lefts[i]] + v + ret[rights[i]+1:]
    }
    return ret
}
