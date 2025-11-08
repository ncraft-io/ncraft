package http

import (
	"fmt"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/mojo-lang/core/go/pkg/mojo/core"
)

func init() {
	core.RegisterJSONTypeDecoder("http.EnvelopedResponse", &EnvelopedResponseCodec{})
	core.RegisterJSONTypeEncoder("http.EnvelopedResponse", &EnvelopedResponseCodec{})
}

type EnvelopedResponseCodec struct {
}

type BareEnvelopedResponse EnvelopedResponse

func (codec *EnvelopedResponseCodec) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	a := iter.ReadAny()
	resp := (*EnvelopedResponse)(ptr)
	if a.ValueType() == jsoniter.ObjectValue {
		bareResponse := (*BareEnvelopedResponse)(resp)
		a.ToVal(bareResponse)

		err := &core.Error{}
		a.ToVal(err)

		if err.Code != nil || len(err.Message) > 0 {
			resp.Error = err
		}
	}
}

func (codec *EnvelopedResponseCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	resp := (*EnvelopedResponse)(ptr)
	stream.WriteObjectStart()
	first := true
	if resp.Error != nil {
		if resp.Error.Code != nil {
			stream.WriteObjectField("code")
			if c := resp.Error.Code.HttpStatusCode; c > 0 {
				stream.WriteString(fmt.Sprint(c))
			} else {
				stream.WriteString(resp.Error.Code.Format())
			}

			first = false
		}

		if len(resp.Error.Message) > 0 {
			if !first {
				stream.WriteMore()
			}
			stream.WriteObjectField("message")
			stream.WriteString(resp.Error.Message)
		}

		if len(resp.Error.Details) > 0 {
			if !first {
				stream.WriteMore()
			}
			stream.WriteObjectField("details")
			stream.WriteVal(resp.Error.Details)
		}
	}

	if !first {
		stream.WriteMore()
	}
	stream.WriteObjectField("data")
	stream.WriteVal(resp.Data)

	if resp.TotalCount > 0 {
		stream.WriteMore()
		stream.WriteObjectField("totalCount")
		stream.WriteVal(resp.TotalCount)
	}

	if len(resp.NextPageToken) > 0 {
		stream.WriteMore()
		stream.WriteObjectField("nextPageToken")
		stream.WriteString(resp.NextPageToken)
	}

	stream.WriteObjectEnd()
}

func (codec *EnvelopedResponseCodec) IsEmpty(ptr unsafe.Pointer) bool {
	e := (*EnvelopedResponse)(ptr)
	return e == nil
}
