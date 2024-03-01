package utils

import (
	"time"
	"unicode"

	"github.com/nbd-wtf/go-nostr"
)

func ValidEventId(id string) bool {
	if len(id) != 64 {
		return false
	}
	for _, c := range id {
		if unicode.IsLetter(c) || unicode.IsNumber(c) {
			if unicode.IsLetter(c) && !unicode.IsLower(c) {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func ValidKind(kind int) bool {
	validKinds := []int{0, 1}

	for _, k := range validKinds {
		if kind == k {
			return true
		}
	}

	return false
}

func ValidContent(content string, kind int) bool {
	return !(kind == 1 && len(content) == 0)
}

func ValidCreatedAt(createdAt int64) bool {
	now := time.Now()
	ts := time.Unix(createdAt, 0)
	return now.Before(ts) || now.Equal(ts)
}

func ValidTags(tags nostr.Tags) bool {
	for _, tag := range tags {
		if len(tag) < 2 || len(tag[0]) == 0 || len(tag[1]) == 0 {
			return false
		}
	}
	return true
}
