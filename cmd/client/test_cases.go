package main

import (
	"github.com/nbd-wtf/go-nostr"
)

const P_TAG = "f7234bd4c1394dda46d09f35bd384dd30cc552ad5541990f98844fb06676e9ca"
const E_TAG = "5c83da77af1dec6d7289834998ad7aafbd9e2191396d75ec3cc27f5a77226f36"

type TestCase struct {
	Request  nostr.Envelope
	Response nostr.Envelope
}

func (tc *TestCase) SerializeRequest() ([]byte, error) {
	return tc.Request.MarshalJSON()
}

func (tc *TestCase) SerializeResponse() ([]byte, error) {
	return tc.Response.MarshalJSON()
}

func GetTestCases(publicKey, privateKey string) []TestCase {
	cases := make([]TestCase, 0, 100)

	// Valid event kind 1
	createEvent := &nostr.Event{
		PubKey:    publicKey,
		CreatedAt: nostr.Now(),
		Kind:      nostr.KindTextNote,
		Tags: nostr.Tags{
			{"e", E_TAG, "wss://nostr.example.com"},
			{"p", P_TAG}},
		Content: "Hello Worlddasdsdf!",
	}
	createEvent.Sign(privateKey)

	createValid1 := TestCase{
		Request: &nostr.EventEnvelope{
			Event: *createEvent,
		},
		Response: &nostr.OKEnvelope{
			EventID: createEvent.ID,
			OK:      true,
			Reason:  "event saved successfully",
		},
	}

	cases = append(cases, createValid1)

	return cases
}
