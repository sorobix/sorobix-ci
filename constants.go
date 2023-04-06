package main

import "os"

var (
	SCRIPT = "../run1.sh"
	SECRET = getSecret()
	DBNAME = "sorobix-ci-db"
)

func getSecret() string {
	res := os.Getenv("secret")
	if res != "" {
		return res
	}
	return "sorobixci"
}
