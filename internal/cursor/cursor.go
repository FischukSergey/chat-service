package cursor

import (
	"encoding/base64"
	"encoding/json"
)

func Encode(data any) (string, error) {
	// Фикс: data -> JSON -> base64URL -> result
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	base64URL := base64.URLEncoding.EncodeToString(jsonData)
	return base64URL, nil
}

func Decode(in string, to any) error {
	// Фикс: in -> base64URL -> JSON -> to
	jsonData, err := base64.URLEncoding.DecodeString(in)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonData, to)
	if err != nil {
		return err
	}
	return nil
}
