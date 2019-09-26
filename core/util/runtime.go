package util

import "os"

func GetCommandOfExecution() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	}
	return ""
}
