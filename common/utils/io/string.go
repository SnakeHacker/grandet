package io

import "strings"

// RemovePunctuation ...
func RemovePunctuation(source string) (target string) {
	target = strings.Replace(source, ":", "_", -1)
	target = strings.Replace(target, "\\", "_", -1)
	target = strings.Replace(target, "/", "_", -1)
	target = strings.Replace(target, "?", "_", -1)
	target = strings.Replace(target, "*", "_", -1)
	target = strings.Replace(target, "[", "_", -1)
	target = strings.Replace(target, "]", "_", -1)
	return
}
