package client

import (
	"fmt"
	"sync"
	"time"

	"github.com/feanor306/nostressgo/src/service"
	"github.com/feanor306/nostressgo/src/types"
	"github.com/gorilla/websocket"
	"github.com/nbd-wtf/go-nostr"
)

type Client struct {
	PubKey  string
	Hub     *Hub
	Conn    *websocket.Conn
	Service *service.Service

	SubMutex      sync.RWMutex
	Subscriptions map[string]*nostr.ReqEnvelope
}

func NewClient(conn *websocket.Conn, svc *service.Service, hub *Hub) *Client {
	return &Client{
		Conn:          conn,
		Service:       svc,
		Hub:           hub,
		Subscriptions: make(map[string]*nostr.ReqEnvelope),
	}
}

// Read handles client websocket read loop
func (c *Client) Read() {
	defer func() {
		c.Service.DB.Close()
		c.Conn.Close()
	}()

	for {
		msgType, p, err := c.Conn.ReadMessage()
		if err != nil {
			c.Service.Log.Error().Err(err).Send()
			continue
		}

		if msgType == websocket.CloseMessage {
			c.Hub.Unregister <- c
			return
		}

		if len(p) == 0 {
			continue
		}

		chanGroup := types.NewChanGroup()
		requestEnvelope := types.NewEnvelopeWrapper(p)

		go c.HandleRequestMessage(requestEnvelope, chanGroup)
		go func() {
			chanGroup.WaitClose()
		}()

		for {
			responseEnvelope, ok := <-chanGroup.Chan
			if !ok {
				break
			}
			c.Respond(responseEnvelope)
		}
	}
}

func (c *Client) HandleRequestMessage(ew *types.EnvelopeWrapper, chanGroup *types.ChanGroup) {
	switch ew.Envelope.Label() {
	case types.EVENT:
		go c.HandleEvent(ew, chanGroup)
	case types.REQ:
		go c.HandleRequestSubscription(ew, chanGroup)
	case types.CLOSE:
		go c.HandleCloseSubscription(ew, chanGroup)
	default:
		chanGroup.Chan <- ew.NoticeResponse(fmt.Errorf("invalid payload"))
		defer chanGroup.Done()
		return
	}
}

func (c *Client) HandleEvent(ew *types.EnvelopeWrapper, chanGroup *types.ChanGroup) {
	defer chanGroup.Done()
	var err error
	ee, ok := ew.Envelope.(*nostr.EventEnvelope)

	if !ok {
		chanGroup.Chan <- ew.EventResponse(fmt.Errorf("unable to parse event"))
		return
	}

	event := types.NewEvent(&ee.Event)

	if err := event.Validate(); err != nil {
		chanGroup.Chan <- ew.EventResponse(err)
		return
	}

	if len(c.PubKey) > 0 {
		c.PubKey = event.PubKey
	}

	switch ee.Event.Kind {
	case 0:
		err = c.Service.HandleZeroEvent(&ee.Event)
	case 1:
		err = c.Service.CreateEvent(&ee.Event)
	case 5:
		err = c.Service.HandleExpiration(&ee.Event)
	}

	chanGroup.Chan <- ew.EventResponse(err)
	c.Hub.Broadcast <- &ee.Event
}

func (c *Client) HandleRequestSubscription(ew *types.EnvelopeWrapper, chanGroup *types.ChanGroup) {
	defer chanGroup.Done()
	reqEnv, ok := ew.Envelope.(*nostr.ReqEnvelope)

	if !ok {
		chanGroup.Chan <- ew.NoticeResponse(fmt.Errorf("unable to parse req"))
		return
	}

	// if subscription exists, overwrite it
	c.SubMutex.Lock()
	c.Subscriptions[reqEnv.SubscriptionID] = reqEnv
	c.SubMutex.Unlock()

	for _, filter := range reqEnv.Filters {
		cg := types.NewChanGroup()
		go func() {
			cg.WaitClose()
		}()
		if err := c.Service.DB.GetEventsByFilter(&filter, cg); err != nil {
			chanGroup.Chan <- ew.ClosedResponse(reqEnv.SubscriptionID, err.Error())
			return
		}

		count := 0
		for event := range cg.Chan {
			count++
			if filter.Limit > 0 && count >= filter.Limit {
				// pause sending because of limit
				// all events matching subscription should be sent eventually
				time.Sleep(time.Second * 5)
			}
			chanGroup.Chan <- event.EventResponse(nil)
		}
		chanGroup.Chan <- ew.EoseResponse()
	}
}

func (c *Client) HandleCloseSubscription(ew *types.EnvelopeWrapper, chanGroup *types.ChanGroup) {
	defer chanGroup.Done()
	closeEnv, ok := ew.Envelope.(*nostr.CloseEnvelope)

	if !ok {
		chanGroup.Chan <- ew.NoticeResponse(fmt.Errorf("unable to parse close"))
		return
	}

	c.SubMutex.Lock()
	defer c.SubMutex.Unlock()
	if _, ok := c.Subscriptions[closeEnv.String()]; ok {
		delete(c.Subscriptions, closeEnv.String())
		chanGroup.Chan <- ew.ClosedResponse(closeEnv.String(), "subscription removed successfully")
		return
	}

	chanGroup.Chan <- ew.NoticeResponse(fmt.Errorf("subscription not found"))
}

func (c *Client) Respond(result *types.EnvelopeWrapper) {
	data, err := result.MarshalJSON()
	if err != nil {
		c.SendErrorResponse(err)
	}

	err = c.Conn.WriteMessage(1, data)
	if err != nil {
		c.Service.Log.Error().Err(err).Send()
	}
}

func (c *Client) SendErrorResponse(err error) {
	c.Service.Log.Error().Err(err).Send()
	ne := nostr.NoticeEnvelope(err.Error())

	data, err := ne.MarshalJSON()
	if err != nil {
		c.Service.Log.Error().Err(err).Send()
	}

	err = c.Conn.WriteMessage(1, data)
	if err != nil {
		c.Service.Log.Error().Err(err).Send()
	}
}
