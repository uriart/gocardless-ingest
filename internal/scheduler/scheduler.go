package scheduler

import (
	"log"
	"os"

	"github.com/gocardless-ingest/internal/clients"
	"github.com/robfig/cron/v3"
)

// StartCron inicia el cron y ejecuta la función de fetch periódicamente
func StartCron(gocardless_client *clients.GoCardlessClient) {
	c := cron.New()
	c.AddFunc("0 8,12 * * *", func() {
		log.Println("INICIO: Consultando transacciones...")

		account := os.Getenv("ACCOUNT_ID")
		transactions, err := gocardless_client.GetTransactions(account)
		if err != nil {
			log.Println("Error al obtener transacciones:", err)
			return
		}

		dsn := os.Getenv("MYSQL_DSN")
		sql_client, err := clients.NewMySQLClient(dsn)
		if err != nil {
			log.Println("Error comprobando transacciones en BD:", err)
			return
		}

		amqpConnStr := os.Getenv("RABBITMQ_CONN_STR")
		rabbitmq_client := clients.NewRabbitMQClient(amqpConnStr)

		var counter = 0
		for _, tx := range transactions {
			already_exists, err := sql_client.TransactionExists(tx)
			if err == nil && !already_exists {
				err := rabbitmq_client.SendTransaction(tx)
				if err != nil {
					log.Printf("Error enviando a la cola: %v\n", err)
				} else {
					counter++
				}
			}
		}

		defer sql_client.Close()
		log.Printf("FIN: Enviadas %v nuevas transacciones de %v procesadas.", counter, len(transactions))
	})
	c.Start()

	select {}
}
