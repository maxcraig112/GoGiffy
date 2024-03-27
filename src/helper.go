package main

import "os"

func GetToken(filename string) (token string, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	content, err := os.ReadFile(wd + `\\` + filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
