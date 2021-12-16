package main

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	"log"
	"sync"
	"testing"
	"time"
)

var topic = "test6"

func Test_main(t *testing.T) {
	//go testNSQ("127.0.0.1:4150")
	fmt.Println(producer("127.0.0.1:4150"))
}

func producer(addr string) error {
	//  try to connect
	cfg := nsq.NewConfig()

	producer, err := nsq.NewProducer(addr, cfg)
	if nil != err {
		return err
	}

	if err = producer.Ping(); err != nil {
		producer.Stop()
		return err
	}
	for {
		producer.Publish(topic, []byte("jmz"))
		time.Sleep(time.Second)
	}

}

type NSQHandler struct {
}

func (this *NSQHandler) HandleMessage(message *nsq.Message) error {
	log.Println("recv:", string(message.Body))
	return nil
}

func testNSQ(addr string) {
	waiter := sync.WaitGroup{}
	waiter.Add(1)

	go func() {
		defer waiter.Done()

		consumer, err := nsq.NewConsumer(topic, "ch1", nsq.NewConfig())
		if nil != err {
			log.Println(err)
			return
		}

		consumer.AddHandler(&NSQHandler{})
		err = consumer.ConnectToNSQD(addr)
		if nil != err {
			log.Println(err)
			return
		}

		select {}
	}()

	waiter.Wait()
}
