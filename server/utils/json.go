package utils

import (
	"encoding/json"
	"os"
)

func JsonParse[T any](bodyData []byte) T {
	var data T
	err := json.Unmarshal(bodyData, &data)
	Check(err)
	return data
}

func JsonSerialize[T any](data T) []byte {
	res, err := json.Marshal(data)
	Check(err)
	return res
}

func JsonStringSerialize[T any](data T) string {
	return string(JsonSerialize(data))
}

func JsonReadFile[T any](path string) T {
	jsonFile, err := os.ReadFile(path)
	Check(err)
	return JsonParse[T](jsonFile)
}

func JsonWriteFile[T any](path string, data T) {
	serializedData := JsonSerialize(data)
	err := os.WriteFile(path, serializedData, 0666)
	Check(err)
}
