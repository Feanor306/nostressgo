package client

import (
	"github.com/feanor306/nostressgo/src/logger"
	"github.com/nbd-wtf/go-nostr"
)

type Client struct {
	privateKey string
	publcKey   string

	testCases []TestCase
	testIndex int
}

func NewClient() *Client {
	privateKey := nostr.GeneratePrivateKey()
	publicKey, err := nostr.GetPublicKey(privateKey)
	if err != nil {
		logger.New().Error().Err(err).Msg("generate public key")
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
	if increment && c.testIndex < len(c.testCases)-1 {
		c.testIndex += 1
	}
	return currentTest
}
