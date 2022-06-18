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
	if str == "" {
		return info, nil
	}
	for _, part := range strings.Split(str, ",") {
		if part == "inject" {
			info.inject = true
			continue
		}
		if len(part) > 4 && strings.HasPrefix(part, "key=") {
			info.key = part[4:]
			continue
		}
		// unknown tag
		return nil, fmt.Errorf("unknown tag: %s", part)
	}
	return info, nil
}
