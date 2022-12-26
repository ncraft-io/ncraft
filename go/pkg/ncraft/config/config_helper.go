package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config/reader"
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config/source"
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config/source/env"
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/config/source/file"
)

var supportedFileSuffixes = map[string]bool{
	"yaml": true,
	"yml":  true,
	"json": true,
	"toml": true,
}

// const (
// 	DevDeployEnvironment   = "dev"
// 	TestDeployEnvironment  = "test"
// 	StageDeployEnvironment = "stage"
// 	ProdDeployEnvironment  = "prod"
// )
//
// const maxConfigsFolderLevel = 3

func init() {
	_ = Load(defaultSources()...)
}

func defaultSources() []source.Source {
	workDir, _ := os.Getwd()
	dirs := []string{
		filepath.Join(workDir, "conf"),
		filepath.Join(workDir, "configs"),
	}

	if configPath := getEnv("CONFIG_PATH"); len(configPath) > 0 {
		dirs = append([]string{configPath}, dirs...)
	}

	// debug mode in the cmd folder
	if strings.Contains(workDir, "/cmd/") || strings.HasSuffix(workDir, "/cmd") {
		dirs = append(dirs, []string{"../configs", "../../configs"}...)
	}

	var sources []source.Source
	for _, dir := range dirs {
		if sources = newFileSources(dir, getEnv("DEPLOY_ENVIRONMENT")); len(sources) > 0 {
			break
		}
	}
	sources = append(sources, env.NewSource())
	return sources
}

func newFileSources(dir string, env string) []source.Source {
	var sources []source.Source

	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		if f.IsDir() {
			ss := newFileSources(f.Name(), env)
			sources = append(sources, ss...)
		} else {
			segments := strings.Split(f.Name(), ".")
			suffix := ""
			if len(segments) >= 2 {
				suffix = segments[len(segments)-1]
			}
			if !supportedFileSuffixes[suffix] {
				continue
			}
			p := path.Join(dir, f.Name())
			if len(env) > 0 {
				name := strings.Join(segments[:len(segments)-1], ".")
				if strings.HasSuffix(name, env) {
					sources = append(sources, file.NewSource(file.WithPath(p)))
				}
			} else {
				sources = append(sources, file.NewSource(file.WithPath(p)))
			}
		}
	}

	return sources
}

func getEnv(key string) string {
	val := os.Getenv("NCRAFT_" + key)
	if len(val) == 0 {
		return os.Getenv(key)
	}
	return ""
}

type watchCloser struct {
	exit chan struct{}
}

func (w watchCloser) Close() error {
	fmt.Println("close")
	w.exit <- struct{}{}
	return nil
}

func WatchFunc(handle func(reader.Value), paths ...string) (io.Closer, error) {
	p := make([]string, 0, len(paths))
	for _, v := range paths {
		p = append(p, strings.Split(v, ".")...)
	}

	exit := make(chan struct{})
	w, err := Watch(p...)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			v, err := w.Next()
			// if err == err_code.WatchStoppedError {
			//	return
			// }
			if err != nil {
				continue
			}

			// if v.Empty() {
			//	continue
			// }

			if handle != nil {
				handle(v)
			}
		}
	}()

	go func() {
		select {
		case <-exit:
			_ = w.Stop()
		}
	}()

	return watchCloser{exit: exit}, nil
}
