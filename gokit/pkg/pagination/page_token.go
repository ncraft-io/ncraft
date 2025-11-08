package pagination

import (
    "fmt"
)

type PageToken interface {
    Format() string
    Parse(token string) error
    Value() interface{}

    Create(value interface{}) PageToken
}

func CreatePositionPageToken(value interface{}) PageToken {
    return (&PositionPageToken{}).Create(value)
}

func ParsePageToken(token string) (PageToken, error) {
    if len(token) > 0 {
        switch token[0] {
        case 'p':
            pageToken := &PositionPageToken{}
            if err := pageToken.Parse(token); err != nil {
                return nil, err
            } else {
                return pageToken, nil
            }
        }
    }
    return nil, fmt.Errorf("malformat token: %s", token)
}
