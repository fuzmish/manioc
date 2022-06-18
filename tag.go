package manioc

import (
	"fmt"
	"reflect"
	"strings"
)

type tagInfo struct {
	inject bool
	key    any
}

func parseTag(tag reflect.StructTag) (*tagInfo, error) {
	// example; manioc:"inject,key=foo"
	info := &tagInfo{
		inject: false,
		key:    nil,
	}
	str := tag.Get("manioc")
	for _, part := range strings.Split(str, ",") {
		if part == "" {
			continue
		}
		if part == "inject" {
			info.inject = true
			continue
		}
		if strings.HasPrefix(part, "key=") {
			// if the value part is empty, remain key as nil
			if len(part) > len("key=") {
				info.key = part[4:]
			}
			continue
		}
		// unknown tag
		return nil, fmt.Errorf("unknown tag: %s", part)
	}
	return info, nil
}
