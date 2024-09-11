package servkit

import (
	"encoding/json"
	"io"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/protobuf/encoding/protojson"
)

type CustomJsonMarshal struct {
	runtime.Marshaler
}

func (u CustomJsonMarshal) ContentType(v interface{}) string {
	if message, ok := v.(*httpbody.HttpBody); ok {
		return message.ContentType
	}

	return "application/json"
}

func (u CustomJsonMarshal) Marshal(v interface{}) ([]byte, error) {
	if httpBody, ok := v.(*httpbody.HttpBody); ok {
		return httpBody.Data, nil
	}
	j := runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			EmitUnpopulated: true,
		},
	}
	return j.Marshal(v)
}

func (u CustomJsonMarshal) NewDecoder(r io.Reader) runtime.Decoder {
	d := json.NewDecoder(r)
	return runtime.DecoderWrapper{
		Decoder: d,
		UnmarshalOptions: protojson.UnmarshalOptions{
			// If DiscardUnknown is set, unknown fields are ignored.
			DiscardUnknown: true,
		},
	}
}
