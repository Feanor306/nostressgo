package types

import "github.com/nbd-wtf/go-nostr"

type EnvelopeWrapper struct {
	Envelope nostr.Envelope
}

func NewEnvelopeWrapper(data []byte) *EnvelopeWrapper {
	return &EnvelopeWrapper{
		Envelope: nostr.ParseMessage(data),
	}
}

func (ew *EnvelopeWrapper) MarshalJSON() ([]byte, error) {
	return ew.Envelope.MarshalJSON()
}

func (ew *EnvelopeWrapper) IsEvent() bool {
	_, ok := ew.Envelope.(*nostr.EventEnvelope)
	return ok
}

func (ew *EnvelopeWrapper) IsSubscriptionReq() bool {
	_, ok := ew.Envelope.(*nostr.ReqEnvelope)
	return ok
}

func (ew *EnvelopeWrapper) IsClose() bool {
	_, ok := ew.Envelope.(*nostr.CloseEnvelope)
	return ok
}

func (ew *EnvelopeWrapper) ToEvent() *nostr.EventEnvelope {
	return ew.Envelope.(*nostr.EventEnvelope)
}

func (ew *EnvelopeWrapper) ToReq() *nostr.ReqEnvelope {
	return ew.Envelope.(*nostr.ReqEnvelope)
}

func (ew *EnvelopeWrapper) ToClose() *nostr.CloseEnvelope {
	return ew.Envelope.(*nostr.CloseEnvelope)
}

func (ew *EnvelopeWrapper) EventResponse(err error) *EnvelopeWrapper {
	ee := ew.Envelope.(*nostr.EventEnvelope)
	reason := "event saved successfully"

	if err != nil {
		reason = err.Error()
	}

	return &EnvelopeWrapper{
		Envelope: &nostr.OKEnvelope{
			EventID: ee.ID,
			OK:      err == nil,
			Reason:  reason,
		},
	}
}

func (ew *EnvelopeWrapper) NoticeResponse(err error) *EnvelopeWrapper {
	var reason nostr.NoticeEnvelope = "event saved successfully"

	if err != nil {
		reason = nostr.NoticeEnvelope(err.Error())
	}

	return &EnvelopeWrapper{
		Envelope: &reason,
	}
}

func (ew *EnvelopeWrapper) ClosedResponse(id string, reason string) *EnvelopeWrapper {
	return &EnvelopeWrapper{
		Envelope: &nostr.ClosedEnvelope{
			SubscriptionID: id,
			Reason:         reason,
		},
	}
}

func (ew *EnvelopeWrapper) EoseResponse() *EnvelopeWrapper {
	var reason nostr.EOSEEnvelope = "end of stored events"
	return &EnvelopeWrapper{
		Envelope: &reason,
	}
}
