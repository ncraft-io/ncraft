package jwt

import "strings"

type UrlPath struct {
	Segments []Segment
	Trailing bool
}

type Segment struct {
	Value string
	Any   bool
}

func NewUrlPath(path string) *UrlPath {
	ss := strings.Split(path, "/")
	trailing := ss[len(ss)-1] == "*"

	var segments []Segment
	if trailing {
		segments = make([]Segment, len(ss)-1)
	} else {
		segments = make([]Segment, len(ss))
	}

	for i := range segments {
		if strings.HasPrefix(ss[i], "{") && strings.HasSuffix(ss[i], "}") {
			segments[i] = Segment{Value: "*", Any: true}
		} else {
			segments[i] = Segment{Value: ss[i]}
		}
	}

	return &UrlPath{Segments: segments, Trailing: trailing}
}

func (p *UrlPath) Match(s string) bool {
	for index, segment := range p.Segments {
		i := strings.IndexByte(s, '/')
		j := i + 1
		if i == -1 {
			i = len(s)
			j = len(s)

			if index != len(p.Segments)-1 || p.Trailing {
				return false
			}
		} else {
			if index == len(p.Segments)-1 && !p.Trailing {
				return false
			}
		}

		if !segment.Any {
			if s[:i] != segment.Value {
				return false
			}
		}

		s = s[j:]
	}

	return true
}
