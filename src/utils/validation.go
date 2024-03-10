package utils

import (
	"encoding/json"
	"strings"
	"time"
	"unicode"

	"github.com/nbd-wtf/go-nostr"
)

type Kind0Content struct {
	Name    string `json:"name"`
	About   string `json:"about"`
	Picture string `json:"picture"`
}

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

func ValidContent(content string, kind int) bool {
	switch kind {
	case 0:
		var res Kind0Content
		err := json.Unmarshal([]byte(content), &res)
		return len(content) > 0 && err == nil
	case 1:
		return len(content) > 0
	}
	return true
}

func ValidCreatedAt(createdAt int64) bool {
	now := time.Now()
	ts := time.Unix(createdAt, 0)
	return ts.Before(now) || now.Equal(ts)
}

func ValidTags(tags nostr.Tags) bool {
	for _, tag := range tags {
		if len(tag) < 2 || len(tag[0]) == 0 || len(tag[1]) == 0 {
			return false
		}
	}
	return true
}

func EventMatchesFilter(event *nostr.Event, filter *nostr.Filter) bool {
	if len(filter.Authors) > 0 {
		for _, author := range filter.Authors {
			if author == event.PubKey {
				return true
			}
		}
	}

	if len(filter.Kinds) > 0 {
		for _, kind := range filter.Kinds {
			if kind == event.Kind {
				return true
			}
		}
	}

	if len(filter.Search) > 0 {
		if strings.Contains(event.Content, filter.Search) {
			return true
		}
	}

	if !filter.Until.Time().IsZero() {
		if filter.Until.Time().Unix() > event.CreatedAt.Time().Unix() {
			return true
		}
	}

	if len(filter.Tags) > 0 {
		if etags, ok := filter.Tags["e"]; ok {
			eventEtags := GetEtags(event)
			for _, v := range etags {
				for _, etag := range eventEtags {
					if v == etag {
						return true
					}
				}
			}

		}

		if ptags, ok := filter.Tags["p"]; ok {
			eventPtags := GetPtags(event)
			for _, v := range ptags {
				for _, ptag := range eventPtags {
					if v == ptag {
						return true
					}
				}
			}

		}
	}
	return true
}
