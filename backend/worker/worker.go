package worker

import (
	"context"
	"fmt"
	"log"

	"github.com/BackAged/steadyrabbit"
	"github.com/streadway/amqp"
)

// Worker defines worker
type Worker struct {
	taskMap     map[string]TaskFunc
	concurrency int
	buffChan    chan struct{}
	consumer    *steadyrabbit.Consumer
}

// NewWorker returns a worker
func NewWorker(concurrency int, consumer *steadyrabbit.Consumer) *Worker {
	w := &Worker{
		taskMap:     make(map[string]TaskFunc, 0),
		concurrency: concurrency,
		consumer:    consumer,
	}

	bc := make(chan struct{}, concurrency)
	for i := 0; i < concurrency; i++ {
		bc <- struct{}{}
	}

	w.buffChan = bc

	return w
}

func (w *Worker) dispatch(taskExecutor TaskFunc, msg amqp.Delivery, errchan chan error) chan error {
	errchan <- taskExecutor(msg.Body)

	return errchan
}

func (w *Worker) reduce(msg amqp.Delivery) error {
	taskExecutor, ok := w.taskMap[msg.RoutingKey]
	if !ok {
		fmt.Println(msg)
		log.Printf("received task-%s not registered!!\n", msg.RoutingKey)
		return msg.Ack(false)
	}

	errChan := make(chan error, 0)

	go w.dispatch(taskExecutor, msg, errChan)

	err := <-errChan
	if err != nil {
		if msg.Redelivered {
			return msg.Nack(false, false)
		}
		return msg.Nack(false, true)
	}

	return msg.Ack(false)
}

// TaskFunc defines task executor
type TaskFunc func(msg []byte) error

// RegisterTask registers a task
func (w *Worker) RegisterTask(task string, taskExecutor TaskFunc) error {
	if _, ok := w.taskMap[task]; ok {
		log.Printf("task-%s already registered!!\n", task)
		return NewErrAlreadyRegisteredTask(task)
	}

	w.taskMap[task] = taskExecutor
	return nil
}

// ConsumeNewMessage ...
func (w *Worker) ConsumeNewMessage() {
	if err := w.consumer.ConsumeOne(context.Background(), w.reduce); err != nil {
		log.Printf("error consuming new message: %s\n", err)
	}
}

// Start starts the worker
func (w *Worker) Start() error {
	log.Println("started worker...")
	for {
		<-w.buffChan
		go w.ConsumeNewMessage()
		w.buffChan <- struct{}{}
	}
}
