package client

import (
	"fmt"
	"log"

	"github.com/feanor306/nostressgo/src/service"
	"github.com/feanor306/nostressgo/src/types"
	"github.com/gorilla/websocket"
	"github.com/nbd-wtf/go-nostr"
)

type Client struct {
	Conn          *websocket.Conn
	Service       *service.Service
	Subscriptions map[string]*nostr.ReqEnvelope
}

func NewClient(conn *websocket.Conn, svc *service.Service) *Client {
	return &Client{
		Conn:          conn,
		Service:       svc,
		Subscriptions: make(map[string]*nostr.ReqEnvelope),
	}
}

func (c *Client) Read() {
	defer func() {
		c.Service.DB.Close()
		c.Conn.Close()
	}()

	for {
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err) // srv err
			continue
		}

		if len(p) == 0 {
			continue
		}

		chanGroup := types.NewChanGroup()
		requestEnvelope := types.NewEventWrapper(p)

		go c.HandleRequestMessage(requestEnvelope, chanGroup)
		go func() {
			chanGroup.WaitClose()
		}()

		// handle req.filter.limit?
		for responseEnvelope := range chanGroup.Chan {
			c.Respond(responseEnvelope)
		}

		// ["EVENT", <subscription_id>, <event JSON as defined above>], used to send events requested by clients.
		// ["EOSE", <subscription_id>], used to indicate the end of stored events and the beginning of events newly received in real-time.
		// ["CLOSED", <subscription_id>, <message>]

		// log.Println(string(p))
		// var receivedEvent nostr.Event

		// CLOSE MESSAGE
		// c.Conn.WriteJSON("")
		// c.Conn.WriteControl(websocket.CloseMessage, FormatCloseMessage(CloseMessageTooBig, ""), time.Now().Add(writeWait))
	}
}

func (c *Client) Respond(result *types.EnvelopeWrapper) {
	data, err := result.MarshalJSON()
	if err != nil {
		c.SendErrorResponse(err)
	}

	err = c.Conn.WriteMessage(1, data)
	if err != nil {
		log.Println(err) // srv err
	}
}

func (c *Client) SendErrorResponse(err error) {
	log.Println(err)
	ne := nostr.NoticeEnvelope(err.Error())

	data, err := ne.MarshalJSON()
	if err != nil {
		log.Println(err)
	}

	err = c.Conn.WriteMessage(1, data)
	if err != nil {
		log.Println(err)
	}
}

func (c *Client) HandleNewSubscription(ew *types.EnvelopeWrapper, chanGroup *types.ChanGroup) {
	reqEnv, ok := ew.Envelope.(*nostr.ReqEnvelope)

	if !ok {
		chanGroup.Chan <- ew.NoticeResponse(fmt.Errorf("unable to parse req"))
		return
	}

	if _, ok := c.Subscriptions[reqEnv.SubscriptionID]; ok {
		chanGroup.Chan <- ew.ClosedResponse(reqEnv.SubscriptionID, "subscription already exists")
	}

	c.Subscriptions[reqEnv.SubscriptionID] = reqEnv

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
				// return or wait a bit and send the rest?
				return
			}
			chanGroup.Chan <- event.EventResponse(nil)
		}
		chanGroup.Chan <- ew.EoseResponse()
	}
}

func (c *Client) RemoveSubscription(subscriptionId string) error {
	if _, ok := c.Subscriptions[subscriptionId]; ok {
		delete(c.Subscriptions, subscriptionId)
		return nil
	}

	return fmt.Errorf("subscription not found")
}

func (c *Client) HandleRequestMessage(ew *types.EnvelopeWrapper, chanGroup *types.ChanGroup) {
	defer chanGroup.Done()

	// CLIENT => SERVER message handling
	switch ew.Envelope.Label() {
	case types.EVENT:
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

		switch ee.Event.Kind {
		case 0:
			err = c.Service.HandleZeroEvent(&ee.Event)
		case 1:
			err = c.Service.CreateEvent(&ee.Event)
		}

		chanGroup.Chan <- ew.EventResponse(err)
		// Notify all other clients, use hub?
	case types.REQ:
		// defer of chanGroup.Done will close the channel we send prematurely, find a fix?
		go c.HandleNewSubscription(ew, chanGroup)
		return
	case types.CLOSE:
		closeEnv, ok := ew.Envelope.(*nostr.CloseEnvelope)

		if !ok {
			chanGroup.Chan <- ew.NoticeResponse(fmt.Errorf("unable to parse close"))
			return
		}

		c.RemoveSubscription(closeEnv.String())
		chanGroup.Chan <- ew.ClosedResponse(closeEnv.String(), "subscription removed successfully")
		return
	default:
		chanGroup.Chan <- ew.NoticeResponse(fmt.Errorf("invalid payload"))
		return
	}
}
