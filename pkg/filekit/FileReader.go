package filekit

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"demo/pkg/logger"

	"github.com/xuri/excelize/v2"
)

var (
	separatorPattern = regexp.MustCompile(`^-+(\w+)\r\n`)
)

func ReadFile(filePath string, filename string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.Error("close file failed", logger.WithError(err))
		}
	}(file)

	fileType := filepath.Ext(file.Name())
	switch fileType {
	case ".xlsx":
		return readExcel(filePath, filename)
	case ".pdf":
		return readPDF(filePath)
	case ".jpg", ".jpeg", ".png":
		return readImage(filePath, filename)
	}

	return nil, nil
}

func readExcel(filePath string, filename string) ([]byte, error) {
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}

	sheets := file.GetSheetList()
	rows, err := file.GetRows(sheets[0])
	if err != nil {
		return nil, err
	}

	cols := rows[0]
	data := make([]map[string]string, 0)
	for i, row := range rows {
		if i == 0 {
			continue
		}
		// 判斷欄位與資料是否一致
		if len(row) == len(cols) {
			rowData := make(map[string]string)
			for k, col := range cols {
				rowData[col] = row[k]
			}
			data = append(data, rowData)
		}

	}

	jsonBytes, err := json.Marshal(map[string]interface{}{
		"filename": filename,
		"data":     data,
	})
	if err != nil {
		return nil, err
	}

	return jsonBytes, err
}

func readPDF(finePath string) ([]byte, error) {
	return readBinary(finePath)
}

func readImage(finePath string, filename string) ([]byte, error) {
	bytes, err := readBinary(finePath)
	if err != nil {
		return nil, err
	}

	imageBytes, err := extractFirstValue(bytes)
	if err != nil {
		return nil, err
	}

	b64 := base64.StdEncoding.EncodeToString(imageBytes)
	if b64 == "" {
		return nil, err
	}

	data := map[string]interface{}{
		"filename": filename,
		"data":     []string{b64},
	}
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return jsonBytes, err
}

func extractFirstValue(bytes []byte) ([]byte, error) {
	payload := string(bytes)
	separatorMatches := separatorPattern.FindStringSubmatch(payload)
	if len(separatorMatches) <= 0 {
		return nil, errors.New("invalid payload, cannot find data separator")
	}
	separator := separatorMatches[1]

	re := regexp.MustCompile(`-+` + separator + `\r\nContent-Disposition: form-data; name=".+"; filename=".+"\r\nContent-Type: .+\r\n\r\n([\s\S]+)\r\n-+` + separator + `-+\r\n`)
	matches := re.FindStringSubmatch(payload)
	if len(matches) <= 0 {
		return nil, errors.New("invalid payload, cannot extract image data")
	}
	imageBytes := []byte(matches[1])

	return imageBytes, nil
}

func readBinary(finePath string) ([]byte, error) {
	f, err := os.Open(finePath)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return data, nil
}
