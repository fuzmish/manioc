package manioc

import (
	"reflect"
	"strings"
)

type tagInfo struct {
	inject bool
	key    any
}

func parseTag(tag reflect.StructTag) *tagInfo {
	// example; manioc:"inject,key=foo"
	info := &tagInfo{
		inject: false,
		key:    nil,
	}
	str := tag.Get("manioc")
	for _, part := range strings.Split(str, ",") {
		if part == "inject" {
			info.inject = true
			continue
		}
		if len(part) > 4 && strings.HasPrefix(part, "key=") {
			info.key = part[4:]
			continue
		}
	}
	return info
}
