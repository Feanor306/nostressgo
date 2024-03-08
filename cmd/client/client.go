package main

import (
	"log"

	"github.com/nbd-wtf/go-nostr"
)

type Client struct {
	privateKey string
	publcKey   string

	testCases []TestCase
	testIndex int
}

func newClient() *Client {
	privateKey := nostr.GeneratePrivateKey()
	publicKey, err := nostr.GetPublicKey(privateKey)
	if err != nil {
		log.Println("generate public key:", err)
	}
	return &Client{
		privateKey: privateKey,
		publcKey:   publicKey,
		testCases:  GetTestCases(publicKey, privateKey),
		testIndex:  0,
	}
}

func (c *Client) GetTestCase(increment bool) *TestCase {
	currentTest := &c.testCases[c.testIndex]
	if increment {
		c.testIndex += 1
	}
	return currentTest
}
