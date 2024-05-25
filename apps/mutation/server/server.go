package server

import (
	"context"
	"fmt"

	"github.com/kurdilesmana/go-saving-api/pkg/logging"
	"github.com/redis/go-redis/v9"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
)

func Redis(handler Handler, redisClient *redis.Client, log *logging.Logger) (err error) {
	ctx := context.Background()
	uniqueID := xid.New().String()

	subject := "transactions"
	consumersGroup := "transactions-consumer-group"
	err = redisClient.XGroupCreateMkStream(ctx, subject, consumersGroup, "0").Err()
	if err != nil {
		log.Warn(logrus.Fields{
			"subject":         subject,
			"consumers_group": consumersGroup,
		}, nil, err.Error())
	}

	for {
		entries, err := redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    consumersGroup,
			Consumer: uniqueID,
			Streams:  []string{subject, ">"},
			Count:    2,
			Block:    0,
			NoAck:    false,
		}).Result()
		if err != nil {
			log.Fatal(logrus.Fields{
				"subject":         subject,
				"consumers_group": consumersGroup,
			}, entries, err.Error())
		}

		log.Info(logrus.Fields{"subject": subject, "consumers_group": consumersGroup}, entries, "entries data...")
		for i := 0; i < len(entries[0].Messages); i++ {
			messageID := entries[0].Messages[i].ID
			values := entries[0].Messages[i].Values
			eventDescription := fmt.Sprintf("%v", values["whatHappened"])
			eventID := fmt.Sprintf("%v", values["eventID"])
			eventData := fmt.Sprintf("%v", values["eventData"])

			if eventDescription == "savings received" {
				log.Info(logrus.Fields{
					"subject":         subject,
					"consumers_group": consumersGroup,
					"event_id":        eventID,
					"event_data":      eventData,
				}, nil, "event data...")

				err = handler.transactionHandler.SavingHandler(ctx, eventID, eventData)
				if err != nil {
					log.Error(logrus.Fields{
						"subject":         subject,
						"consumers_group": consumersGroup,
						"event_id":        eventID,
						"event_data":      eventData,
					}, nil, err.Error())
				}

				redisClient.XAck(ctx, subject, consumersGroup, messageID)
			} else if eventDescription == "cashwithdrawl received" {
				log.Info(logrus.Fields{
					"subject":         subject,
					"consumers_group": consumersGroup,
					"event_id":        eventID,
					"event_data":      eventData,
				}, nil, "event data...")

				err = handler.transactionHandler.CashWithdrawlHandler(ctx, eventID, eventData)
				if err != nil {
					log.Error(logrus.Fields{
						"subject":         subject,
						"consumers_group": consumersGroup,
						"event_id":        eventID,
						"event_data":      eventData,
					}, nil, err.Error())
				}

				redisClient.XAck(ctx, subject, consumersGroup, messageID)
			}
		}
	}
}
