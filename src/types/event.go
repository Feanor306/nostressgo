package types

import (
	"fmt"
	"strings"

	"github.com/nbd-wtf/go-nostr"
)

type Event struct {
	nostr.Event
	Etags      []string
	Ptags      []string
	Gtags      []string
	Dtag       string
	Expiration int32
	Json       string
}

func (e *Event) SetTags() {
	tags := make(nostr.Tags, 0, len(e.Etags)+len(e.Ptags)+len(e.Gtags)+2)

	if len(e.Etags) > 0 {
		for _, etag := range e.Etags {
			tags = tags.AppendUnique(nostr.Tag{"e", etag})
		}
	}
	if len(e.Ptags) > 0 {
		for _, ptag := range e.Ptags {
			tags = tags.AppendUnique(nostr.Tag{"e", ptag})
		}
	}
	if len(e.Gtags) > 0 {
		for _, gtag := range e.Gtags {
			// saved in following format
			// "#" + tag[0] + ":" + tag[1]
			tag := strings.Split(gtag, ":")
			tags = tags.AppendUnique(nostr.Tag{tag[0][1:], tag[1]})
		}
	}

	if len(e.Dtag) > 0 {
		tags = tags.AppendUnique(nostr.Tag{"d", e.Dtag})
	}

	if e.Expiration > 0 {
		tags = tags.AppendUnique(nostr.Tag{"expiration", fmt.Sprint(e.Expiration)})
	}

	e.Tags = tags
	// process rest of JSON, if applicable?
}
