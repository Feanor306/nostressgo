package helpers

import (
	"strconv"
	"time"

	"github.com/nbd-wtf/go-nostr"
)

func GetGtags(event *nostr.Event) []string {
	var tags nostr.Tags = make(nostr.Tags, 0, len(event.Tags))

	for _, tag := range event.Tags {
		key := tag.Key()
		if len(key) == 1 && key != "e" && key != "p" && key != "d" {
			tag[1] = "#" + tag[0] + ":" + tag[1]
			tags.AppendUnique(tag)
		}
	}

	return GetTagValues(tags)
}

func GetEtags(event *nostr.Event) []string {
	return GetTagValues(event.Tags.GetAll([]string{"e"}))
}

func GetPtags(event *nostr.Event) []string {
	return GetTagValues(event.Tags.GetAll([]string{"p"}))
}

func GetExpiration(event *nostr.Event) (int64, error) {
	tags := event.Tags.GetAll([]string{"expiration"})

	if len(tags) == 0 {
		return 0, nil
	}

	i, err := strconv.ParseInt(tags[0].Value(), 10, 64)
	if err != nil {
		return 0, err
	}

	tm := time.Unix(i, 0)
	return tm.Unix(), nil
}

func GetTagValues(tags nostr.Tags) []string {
	result := make([]string, 0, len(tags))

	for _, tag := range tags {
		result = append(result, tag.Value())
	}

	return result
}
