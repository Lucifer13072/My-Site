package main

import (
	"encoding/base64"
)

func auth(user string, pass string) bool {

	return true
}

func register(user string, pass string) bool {

	return true
}

func hashing(str string) string {

	res := base64.StdEncoding.EncodeToString([]byte(str))

	return res
}

func unhasing(str string) string {

	return "тест"
}
