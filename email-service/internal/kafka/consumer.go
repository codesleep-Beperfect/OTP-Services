package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"email-service/config"
	"email-service/internal/email"
	"email-service/internal/model"

	"github.com/segmentio/kafka-go"
)

func StartConsumer() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: config.App.Kafka.Brokers,
		Topic:   config.App.Kafka.Topic,
		GroupID: "email-service-group",
	})

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("error:", err)
			continue
		}

		var event model.OTPEvent
		json.Unmarshal(msg.Value, &event)

		//  Expiry check
		if time.Now().Unix() > event.ExpiresAt {
			log.Println("OTP expired, skipping email")
			continue
		}

		err = email.Send(event.Identifier, event.OTP)
		if err != nil {
			log.Println("failed to send email:", err)
			continue
		}

		log.Println("Email sent to:", event.Identifier)
	}
}