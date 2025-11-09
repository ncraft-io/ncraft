package filter

import "net/url"

type HttpQueryFilter struct {
}

func (c HttpQueryFilter) Compile(query url.Values) Filter {
    filter := ""
    filters := &FieldFilters{}
    for k, v := range query {
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
