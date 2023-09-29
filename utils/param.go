package utils

import "net/url"

func DecodeParam(param string) (string, error) {
	decodedValue, err := url.QueryUnescape(param)
	if err != nil {
		return "", err
	}
	return decodedValue, nil
}

func EncodeParam(param string) string {
	return url.QueryEscape(param)
}
