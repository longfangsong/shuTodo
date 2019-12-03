package tools

import "log"

func CheckErr(err error, errMessage string) {
	if err != nil {
		log.Fatal(err, errMessage)
	}
}
