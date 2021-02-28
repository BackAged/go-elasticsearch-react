package worker_test

import (
	"context"
	"time"

	"github.com/BackAged/steadyrabbit"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/streadway/amqp"
)

var _ = Describe("Worker", func() {
	Context("Performance Benchmark: ", func() {
		It("publish lots of message", func() {
			cnf := &steadyrabbit.Config{
				URL: "amqp://root:root@localhost:5672",
			}

			eName := "catalogue"
			tName := "catalog.brand.create"

			cnf.Publisher = &steadyrabbit.PublisherConfig{
				Exchange: &steadyrabbit.ExchangeConfig{
					ExchangeName:    eName,
					ExchangeDeclare: true,
					ExchangeType:    amqp.ExchangeTopic,
					ExchangeDurable: true,
				},
			}

			p, err := steadyrabbit.NewPublisher(cnf)
			Expect(err).ToNot(HaveOccurred())

			b := []byte(`[
				{
					"slug": "shahin3",
					"name": "aha",
					"image_url": "aha.coom"
				},
				{
					"slug": "shahin4",
					"name": "aha",
					"image_url": "aha.coom"
				},
				{
					"slug": "shahin5",
					"name": "aha",
					"image_url": "aha.coom"
				},
				{
					"slug": "shahin6",
					"name": "aha",
					"image_url": "aha.coom"
				}
			]`)

			err = p.Publish(context.Background(), tName, b)
			Expect(err).ToNot(HaveOccurred())

			time.Sleep(1 * time.Second)
		})
	})
})
