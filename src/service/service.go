package service

import (
	"fmt"

	"github.com/feanor306/nostressgo/src/config"
	"github.com/feanor306/nostressgo/src/database"
	"github.com/feanor306/nostressgo/src/handlers"
	"github.com/feanor306/nostressgo/src/logger"
	"github.com/feanor306/nostressgo/src/utils"
	"github.com/nbd-wtf/go-nostr"
)

type Service struct {
	Conf     *config.Config
	DB       *database.DB
	Log      *logger.Log
	Upgrader *handlers.Upgrader
}

func NewService(conf *config.Config, db *database.DB, up *handlers.Upgrader) *Service {
	return &Service{
		Conf:     conf,
		DB:       db,
		Log:      logger.New(),
		Upgrader: up,
	}
}

func (s *Service) CreateEvent(event *nostr.Event) error {
	if len(s.Conf.ExclusivePubKey) > 0 && event.PubKey != s.Conf.ExclusivePubKey {
		return fmt.Errorf("only accepting events from exclusive key")
	}

	return s.DB.CreateEvent(event)
}

func (s *Service) HandleZeroEvent(event *nostr.Event) error {
	id, err := s.DB.EventZeroExists(event)
	if err != nil {
		return err
	}
	if len(id) > 0 {
		return s.DB.UpdateEventZero(id, event)
	} else {
		return s.DB.CreateEvent(event)
	}
}

func (s *Service) HandleExpiration(event *nostr.Event) error {
	etags := utils.GetEtags(event)

	if len(etags) > 0 {
		return s.DB.ExpireEvents(etags)
	}

	return nil
}
