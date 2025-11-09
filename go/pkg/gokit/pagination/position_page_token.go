package pagination

import (
    "reflect"
    "time"

    "github.com/mr-tron/base58"
)

type PositionPageToken struct {
    time.Time
    Position int32
}

func (p *PositionPageToken) Format() string {
    seconds := p.Unix()
    bytes := []byte{
        byte(p.Position >> 16),
        byte(p.Position >> 8),
        byte(p.Position),
        byte(seconds >> 32),
        byte(seconds >> 24),
        byte(seconds >> 16),
        byte(seconds >> 8),
        byte(seconds),
    }
    return "p" + base58.Encode(bytes)
}

func (p *PositionPageToken) Parse(token string) error {
    if len(token) > 0 && token[0] == 'p' {
        token = token[1:]

        bytes, err := base58.Decode(token)
        if err != nil {
            return err
        }

        var seconds int64
        seconds |= int64(bytes[3]) << 32
        seconds |= int64(bytes[4]) << 24
        seconds |= int64(bytes[5]) << 16
        seconds |= int64(bytes[6]) << 8
        seconds |= int64(bytes[7]) << 0

        var position int64
        position |= int64(bytes[0]) << 16
        position |= int64(bytes[1]) << 8
        position |= int64(bytes[2]) << 0

        p.Time = time.Unix(seconds, 0)
        p.Position = int32(position)
    }

    return nil
}

func (p *PositionPageToken) Value() interface{} {
    return p.Position
}

func (*PositionPageToken) Create(value interface{}) PageToken {
    if value != nil {
        token := &PositionPageToken{
            Time: time.Now(),
        }

        v := reflect.ValueOf(value)
        switch v.Kind() {
        case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
            token.Position = int32(v.Int())
            return token
        case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
            token.Position = int32(v.Uint())
            return token
        }
    }
    return &PositionPageToken{}
}
