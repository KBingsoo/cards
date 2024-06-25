package cmd

import (
	"os"

	"github.com/KBingsoo/cards/internal/domain/cards"
	"github.com/KBingsoo/cards/internal/gateways/database"
	"github.com/KBingsoo/cards/internal/gateways/pubsub"
	"github.com/KBingsoo/cards/internal/gateways/web"
	"github.com/KBingsoo/entities/pkg/models"
	"github.com/joho/godotenv"
	"github.com/literalog/go-wise/wise"
	"github.com/streadway/amqp"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start webserver",
	RunE: func(cmd *cobra.Command, args []string) error {

		err := godotenv.Load(".env")
		if err != nil {
			return err
		}

		col, err := database.GetCollection("cards")
		if err != nil {
			return err
		}

		repository, err := wise.NewMongoSimpleRepository[models.Card](col)
		if err != nil {
			return err
		}

		connection, err := amqp.Dial(os.Getenv("RABBIT_URL"))
		if err != nil {
			return err
		}
		defer connection.Close()

		consumer, err := pubsub.NewConsumer(connection)
		if err != nil {
			return err
		}

		producer, err := pubsub.NewProducer(connection)
		if err != nil {
			return err
		}

		service := cards.NewManager(repository, producer)

		handler := cards.NewHandler(service)

		server := web.NewServer(handler)

		errCh := make(chan error)

		go func() {
			errCh <- server.Run(8080)
		}()

		go func() {
			errCh <- consumer.Consume()
		}()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
