package client

import (
	"time"

	"github.com/nbd-wtf/go-nostr"
)

const P_TAG = "f7234bd4c1394dda46d09f35bd384dd30cc552ad5541990f98844fb06676e9ca"
const E_TAG = "5c83da77af1dec6d7289834998ad7aafbd9e2191396d75ec3cc27f5a77226f36"

type TestCase struct {
	Request  nostr.Envelope
	Response nostr.Envelope
	SubId    string
	Done     bool
}

func (tc *TestCase) SerializeRequest() ([]byte, error) {
	return tc.Request.MarshalJSON()
}

func (tc *TestCase) SerializeResponse() ([]byte, error) {
	return tc.Response.MarshalJSON()
}

func (tc *TestCase) GetEoseResponse() ([]byte, error) {
	var eose nostr.EOSEEnvelope = nostr.EOSEEnvelope(tc.SubId)
	return eose.MarshalJSON()
}

func (tc *TestCase) GetEventResponse() ([]byte, error) {
	ee := tc.Response.(*nostr.EventEnvelope)
	ee.Event.Sig = ""
	return ee.MarshalJSON()
}

func GetTestCases(publicKey, privateKey string) []TestCase {
	cases := make([]TestCase, 0, 100)

	cases = append(cases, GetEvent1TestCases(publicKey, privateKey)...)
	cases = append(cases, GetEvent0TestCases(publicKey, privateKey)...)
	cases = append(cases, GetEvent5TestCases(publicKey, privateKey)...)
	cases = append(cases, GetReqTestCases(publicKey, privateKey)...)

	return cases
}

func GetEvent1TestCases(publicKey, privateKey string) []TestCase {
	cases := make([]TestCase, 0, 10)

	// Valid event kind 1
	event1 := &nostr.Event{
		PubKey:    publicKey,
		CreatedAt: nostr.Now(),
		Kind:      nostr.KindTextNote,
		Tags: nostr.Tags{
			{"e", E_TAG, "wss://nostr.example.com"},
			{"p", P_TAG}},
		Content: "Hello Worlddasdsdf!",
	}
	event1.Sign(privateKey)
	createValid1 := TestCase{
		Request: &nostr.EventEnvelope{
			Event: *event1,
		},
		Response: &nostr.OKEnvelope{
			EventID: event1.ID,
			OK:      true,
			Reason:  "event saved successfully",
		},
	}
	cases = append(cases, createValid1)

	// invalid date
	event1Invalid2 := &nostr.Event{
		PubKey:    publicKey,
		CreatedAt: nostr.Timestamp(time.Now().Add(time.Hour).Unix()),
		Kind:      nostr.KindTextNote,
		Tags: nostr.Tags{
			{"e", E_TAG, "wss://nostr.example.com"},
			{"p", P_TAG}},
		Content: "Hello Worlddasdsdf!",
	}
	event1Invalid2.Sign(privateKey)
	createInvalid2 := TestCase{
		Request: &nostr.EventEnvelope{
			Event: *event1Invalid2,
		},
		Response: &nostr.OKEnvelope{
			EventID: event1Invalid2.ID,
			OK:      false,
			Reason:  "invalid created_at",
		},
	}
	cases = append(cases, createInvalid2)
	return cases
}

func GetEvent0TestCases(publicKey, privateKey string) []TestCase {
	cases := make([]TestCase, 0, 10)

	// Valid event kind 0
	event0valid1 := &nostr.Event{
		PubKey:    publicKey,
		Kind:      nostr.KindProfileMetadata,
		CreatedAt: nostr.Timestamp(time.Now().Add(time.Hour).Unix()),
		Content:   "{\"name\":\"Bob\", \"about\":\"normal dude\", \"picture\":\"face.jpg\"}",
	}
	event0valid1.Sign(privateKey)
	create0Valid1 := TestCase{
		Request: &nostr.EventEnvelope{
			Event: *event0valid1,
		},
		Response: &nostr.OKEnvelope{
			EventID: event0valid1.ID,
			OK:      true,
			Reason:  "event saved successfully",
		},
	}
	cases = append(cases, create0Valid1)

	// Valid event kind 0 should replace previous
	event0valid2 := &nostr.Event{
		PubKey:    publicKey,
		Kind:      nostr.KindProfileMetadata,
		CreatedAt: nostr.Timestamp(time.Now().Add(time.Hour).Unix()),
		Content:   "{\"name\":\"Robert\", \"about\":\"normal friend\", \"picture\":\"head.jpg\"}",
	}
	event0valid2.Sign(privateKey)
	create0Valid2 := TestCase{
		Request: &nostr.EventEnvelope{
			Event: *event0valid2,
		},
		Response: &nostr.OKEnvelope{
			EventID: event0valid2.ID,
			OK:      true,
			Reason:  "event saved successfully",
		},
	}
	cases = append(cases, create0Valid2)

	// Invalid event kind 0
	event0invalid1 := &nostr.Event{
		PubKey:    publicKey,
		Kind:      nostr.KindProfileMetadata,
		CreatedAt: nostr.Timestamp(time.Now().Add(time.Hour).Unix()),
		Content:   "{\"name\":\"Bob\", \"about\":\"normal dude\", \"picture\":\"face.jpg\"}}",
	}
	event0invalid1.Sign(privateKey)
	create0Invalid1 := TestCase{
		Request: &nostr.EventEnvelope{
			Event: *event0invalid1,
		},
		Response: &nostr.OKEnvelope{
			EventID: event0invalid1.ID,
			OK:      false,
			Reason:  "invalid content",
		},
	}
	cases = append(cases, create0Invalid1)

	return cases
}

func GetEvent5TestCases(publicKey, privateKey string) []TestCase {
	cases := make([]TestCase, 0, 10)

	// Valid event kind 1
	event1 := &nostr.Event{
		PubKey:    publicKey,
		CreatedAt: nostr.Now(),
		Kind:      nostr.KindTextNote,
		Tags: nostr.Tags{
			{"e", E_TAG, "wss://nostr.example.com"},
			{"p", P_TAG}},
		Content: "To be expired",
	}
	event1.Sign(privateKey)
	createValid1 := TestCase{
		Request: &nostr.EventEnvelope{
			Event: *event1,
		},
		Response: &nostr.OKEnvelope{
			EventID: event1.ID,
			OK:      true,
			Reason:  "event saved successfully",
		},
	}
	cases = append(cases, createValid1)

	// Valid event kind 1
	event5 := &nostr.Event{
		PubKey:    publicKey,
		CreatedAt: nostr.Now(),
		Kind:      nostr.KindDeletion,
		Tags: nostr.Tags{
			{"e", event1.ID, "wss://nostr.example.com"},
		},
		Content: "Delete this please",
	}
	event5.Sign(privateKey)
	expireEvent5 := TestCase{
		Request: &nostr.EventEnvelope{
			Event: *event5,
		},
		Response: &nostr.OKEnvelope{
			EventID: event5.ID,
			OK:      true,
			Reason:  "event saved successfully",
		},
	}
	cases = append(cases, expireEvent5)

	return cases
}

func GetReqTestCases(publicKey, privateKey string) []TestCase {
	cases := make([]TestCase, 0, 10)

	// Valid event kind 1
	event1 := &nostr.Event{
		PubKey:    publicKey,
		CreatedAt: nostr.Now(),
		Kind:      nostr.KindTextNote,
		Tags: nostr.Tags{
			{"e", E_TAG},
			{"p", P_TAG}},
		Content: "Example event customtext",
	}
	event1.Sign(privateKey)
	createValid1 := TestCase{
		Request: &nostr.EventEnvelope{
			Event: *event1,
		},
		Response: &nostr.OKEnvelope{
			EventID: event1.ID,
			OK:      true,
			Reason:  "event saved successfully",
		},
	}
	cases = append(cases, createValid1)

	sub1 := "sub1"
	req1 := TestCase{
		Request: &nostr.ReqEnvelope{
			SubscriptionID: sub1,
			Filters:        nostr.Filters{nostr.Filter{IDs: []string{event1.ID}}},
		},
		Response: &nostr.EventEnvelope{
			Event: *event1,
		},
		SubId: sub1,
		Done:  false,
	}
	cases = append(cases, req1)

	sub2 := "sub2"
	req2 := TestCase{
		Request: &nostr.ReqEnvelope{
			SubscriptionID: sub2,
			Filters:        nostr.Filters{nostr.Filter{Authors: []string{event1.PubKey}}},
		},
		Response: &nostr.EventEnvelope{
			Event: *event1,
		},
		SubId: sub2,
		Done:  false,
	}
	cases = append(cases, req2)

	return cases
}
