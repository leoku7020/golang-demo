package servkit

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

type TextPlainMarshal struct {
	runtime.Marshaler
}

func (m TextPlainMarshal) ContentType(_ interface{}) string {
	return "text/plain"
}

func (m TextPlainMarshal) Marshal(v interface{}) ([]byte, error) {
	if msg, ok := v.(proto.Message); ok {
		return prototext.Marshal(msg)
	}
	return nil, status.Errorf(http.StatusInternalServerError, "failed to marshal message")
}

func (m TextPlainMarshal) NewDecoder(r io.Reader) runtime.Decoder {
	return runtime.DecoderFunc(func(v interface{}) error {
		text, err := io.ReadAll(r)
		if err != nil {
			return err
		}

		data := map[string]interface{}{
			"data": string(text),
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}

		err = json.Unmarshal(jsonData, &v)
		if err != nil {
			return err
		}

		return nil
	})
}

func (m TextPlainMarshal) NewEncoder(w io.Writer) runtime.Encoder {
	return runtime.EncoderFunc(func(v interface{}) error {
		data, err := m.Marshal(v)
		if err != nil {
			return err
		}
		_, err = w.Write(data)
		return err
	})
}

func (m TextPlainMarshal) Delimiter() []byte {
	return []byte("\n")
}
