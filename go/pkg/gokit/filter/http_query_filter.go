package filter

import (
	"net/url"
	"sort"
)

type HttpQueryFilter struct {
}

func (c HttpQueryFilter) Compile(query url.Values) Filter {
	filter := ""
	filters := &FieldFilters{}
	keys := make([]string, 0, len(query))
	for k := range query {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := query[k]
		switch k {
		case "filter":
			if len(v) > 0 {
				filter = v[0]
			}
		case "page_size", "page_token", "skip", "show_deleted", "parent", "order_by":
			continue
		default:
			if len(v) > 0 {
				filters.AddFilter(k, v[0])
			}
		}
	}

	return filters.Compile(filter)
}
