package common

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

// UnicodeToUTF8 ...
func UnicodeToUTF8(unicode string) (utf8 string, err error) {
	sUnicodev := strings.Split(unicode, "\\u")
	var context bytes.Buffer
	for _, v := range sUnicodev {
		if len(v) < 1 {
			continue
		}
		temp, err := strconv.ParseInt(v, 16, 32)
		if err != nil {
			glog.Error(err)
			return utf8, err
		}
		context.WriteString(fmt.Sprintf("%c", temp))
	}
	utf8 = context.String()
	return
}
