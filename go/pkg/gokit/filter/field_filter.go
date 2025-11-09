package filter

import (
    "bytes"
    "fmt"
)

type FieldFilter struct {
    Name   string
    Filter interface{}
}

type FieldFilters []FieldFilter

func (f *FieldFilters) AddFilter(field string, filter interface{}) *FieldFilters {
    *f = append(*f, FieldFilter{
        Name:   field,
        Filter: filter,
    })
    return f
}

func (f FieldFilters) Compile(filter string) Filter {
    buffer := bytes.NewBuffer(nil)
    for i, field := range f {
        if i > 0 {
            buffer.WriteString(" && ")
        }
        buffer.WriteString(field.Name)
        buffer.WriteString(" == ")
        buffer.WriteString(fmt.Sprint(field.Filter))
    }

    if len(filter) > 0 {
        buffer.WriteString(" && ")
        buffer.WriteString("(")
        buffer.WriteString(filter)
        buffer.WriteString(")")
    }
    return Filter(buffer.String())
}
