package servkit

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type UrlEncodeMarshal struct {
	runtime.Marshaler
}

// ContentType means the content type of the response
func (u UrlEncodeMarshal) ContentType(_ interface{}) string {
	return "application/json"
}

func (u UrlEncodeMarshal) Marshal(v interface{}) ([]byte, error) {
	j := runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			EmitUnpopulated: true,
		},
	}
	return j.Marshal(v)
}

// NewDecoder indicates how to decode the request
func (u UrlEncodeMarshal) NewDecoder(r io.Reader) runtime.Decoder {
	return runtime.DecoderFunc(func(p interface{}) error {
		log.Println(fmt.Sprintf("AAAAA:%T", p))

		formData, err := io.ReadAll(r)

		if err != nil {
			return err
		}

		values, err := url.ParseQuery(string(formData))
		if err != nil {
			return err
		}

		msg, ok := p.(proto.Message)
		if ok {
			filter := &utilities.DoubleArray{}

			err = runtime.PopulateQueryParameters(msg, values, filter)

			if err != nil {
				return err
			}
		} else {
			log.Println("not a proto message")

			// we need to bind values to struct type
			jsonMap := make(map[string]interface{}, len(values))
			for k, v := range values {
				if len(v) == 1 {
					jsonMap[k] = v[0]
				} else {
					jsonMap[k] = v
				}
			}

			jsonBody, err := json.Marshal(jsonMap)
			if err != nil {
				return fmt.Errorf("failed to marshal: %v", err.Error())
			}

			log.Println(string(jsonBody))
			err = json.Unmarshal(jsonBody, &p)
			if err != nil {
				return fmt.Errorf("failed to unmarshal: %v", err.Error())
			}
		}

		return nil
	})
}
