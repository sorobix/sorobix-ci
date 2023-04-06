package main

import "os"

var (
	SCRIPT = "../run1.sh"
	SECRET = os.Getenv("secret")
	DBNAME = "sorobix-ci-db"
)
