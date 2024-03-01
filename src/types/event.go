package types

import (
	"fmt"
	"strings"

	"github.com/feanor306/nostressgo/src/utils"
	"github.com/nbd-wtf/go-nostr"
)

// Move to models?
type Event struct {
	nostr.Event
	Etags      []string
	Ptags      []string
	Gtags      []string
	Dtag       string
	Expiration nostr.Timestamp
	Json       string
}

func NewEvent(ne *nostr.Event) *Event {
	return &Event{
		Event: *ne,
	}
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

// TODO validation for kind 0 and longform kind 30023?
func (e *Event) Validate() error {
	if !utils.ValidEventId(e.ID) {
		return fmt.Errorf("invalid event id")
	}

	if !utils.ValidEventId(e.PubKey) {
		return fmt.Errorf("invalid pubkey")
	}

	if !utils.ValidKind(e.Kind) {
		return fmt.Errorf("invalid kind")
	}

	if !utils.ValidContent(e.Content, e.Kind) {
		return fmt.Errorf("invalid content")
	}

	if !utils.ValidCreatedAt(int64(e.CreatedAt)) {
		return fmt.Errorf("invalid created_at")
	}

	if !utils.ValidTags(e.Tags) {
		return fmt.Errorf("invalid tag")
	}

	sigOk, err := e.CheckSignature()
	if err != nil {
		return err
	}
	if !sigOk {
		return fmt.Errorf("invalid signature")
	}

	return nil
}

func (e *Event) ToEnvelopeWrapper() *EnvelopeWrapper {
	var envelope nostr.EventEnvelope

	e.SetTags()
	envelope.Event = e.Event

	return &EnvelopeWrapper{
		Envelope: &envelope,
	}
}
