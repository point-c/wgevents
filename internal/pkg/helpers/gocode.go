package helpers

import (
	"go/format"
)

func GoFmt(src []byte) (dst []byte, err error) {
	dst, err = format.Source(src)
	if err != nil {
		dst = src
	}
	return
}
