package nacos

import (
    source2 "github.com/ncraft-io/ncraft/go/pkg/ncraft/config/source"
)

type nacosSource struct {
}

func (n *nacosSource) Read() (*source2.ChangeSet, error) {

    return nil, nil
}

func (n *nacosSource) Write(*source2.ChangeSet) error {
    return nil
}

func (n *nacosSource) Watch() (source2.Watcher, error) {
    return nil, nil
}

func (n *nacosSource) String() string {
    return ""
}

func NewSource(opts ...source2.Option) source2.Source {
    return &nacosSource{}
}
