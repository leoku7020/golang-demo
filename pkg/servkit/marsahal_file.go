package servkit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"demo/pkg/filekit"
	"demo/pkg/logger"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"
)

type FileMarshal struct {
	runtime.Marshaler
}

// ContentType means the content type of the response
func (u FileMarshal) ContentType(_ interface{}) string {
	return "application/json"
}

func (u FileMarshal) Marshal(v interface{}) ([]byte, error) {
	j := runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			EmitUnpopulated: true,
		},
	}
	return j.Marshal(v)
}

// NewDecoder indicates how to decode the request
func (u FileMarshal) NewDecoder(r io.Reader) runtime.Decoder {
	return runtime.DecoderFunc(func(p interface{}) error {
		fileData, err := io.ReadAll(r)
		if err != nil {
			return err
		}
		re := regexp.MustCompile(`filename="([^"]+)"`)
		matches := re.FindSubmatch(fileData)
		if len(matches) >= 2 {
			extension := filepath.Ext(string(matches[1]))
			filename := fmt.Sprintf("tmp_*.%s", extension)
			tmpFile, err := os.CreateTemp("", filename)
			if err != nil {
				return err
			}

			if err = os.WriteFile(tmpFile.Name(), fileData, 0777); err != nil {
				return err
			}

			defer func(file *os.File) {
				err := file.Close()
				if err != nil {
					logger.Error("close temp file failed", logger.WithError(err))
				}
			}(tmpFile)

			data, err := filekit.ReadFile(tmpFile.Name(), string(matches[1]))
			if err != nil {
				return err
			}

			err = json.Unmarshal(data, &p)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
