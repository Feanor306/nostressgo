package types

import (
	"fmt"
	"strings"

	"github.com/feanor306/nostressgo/src/utils"
	"github.com/nbd-wtf/go-nostr"
)

type Event struct {
	nostr.Event
	Etags      []string
	Ptags      []string
	Gtags      []string
	Dtag       string
	Subject    string
	Expiration nostr.Timestamp
	Json       string
}

func NewEvent(ne *nostr.Event) *Event {
	return &Event{
		Event: *ne,
	}
}

func (e *Event) SetTags() {
	tags := make(nostr.Tags, 0, len(e.Etags)+len(e.Ptags)+len(e.Gtags)+3)

	if len(e.Etags) > 0 {
		for _, etag := range e.Etags {
			tags = tags.AppendUnique(nostr.Tag{"e", etag})
		}
	}
	if len(e.Ptags) > 0 {
		for _, ptag := range e.Ptags {
			tags = tags.AppendUnique(nostr.Tag{"p", ptag})
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

	if len(e.Subject) > 0 {
		tags = tags.AppendUnique(nostr.Tag{"subject", e.Subject})
	}

	if e.Expiration > 0 {
		tags = tags.AppendUnique(nostr.Tag{"expiration", fmt.Sprint(e.Expiration)})
	}

	e.Tags = tags
}

func (e *Event) Validate() error {
	switch e.Kind {
	case 0:
		return e.ValidateEvent0()
	case 1:
		return e.ValidateEvent1()
	case 5:
		return e.ValidateEvent5()
	default:
		return fmt.Errorf("unsupported event kind")
	}
}

func (e *Event) ValidateEvent0() error {
	if !utils.ValidEventId(e.ID) {
		return fmt.Errorf("invalid id")
	}
	if !utils.ValidEventId(e.PubKey) {
		return fmt.Errorf("invalid pubkey")
	}
	if !utils.ValidContent(e.Content, e.Kind) {
		return fmt.Errorf("invalid content")
	}
	return nil
}

func (e *Event) ValidateEvent1() error {
	if !utils.ValidEventId(e.ID) {
		return fmt.Errorf("invalid event id")
	}

	if !utils.ValidEventId(e.PubKey) {
		return fmt.Errorf("invalid pubkey")
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

func (e *Event) ValidateEvent5() error {
	if !utils.ValidEventId(e.PubKey) {
		return fmt.Errorf("invalid pubkey")
	}
	if !utils.ValidTags(e.Tags) {
		return fmt.Errorf("invalid tag")
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
