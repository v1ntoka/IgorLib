package myzip

import "io"

func IsErr(err error) bool {
	return err != nil && err != io.EOF
}

func CriticalErr(err error) {
	if IsErr(err) {
		panic(err)
	}
}
