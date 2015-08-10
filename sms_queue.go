package ideamart

/*
	Request throttling SMS queue. Auto-retries when retryable errors are encountered.
*/

import (
	"log"
	"time"
)

type smsMessage struct {
	message        string
	recipients     []string
	chargingAmount float32
	reportDelivery bool
	retries        int
}

// SMS Queue with auto-retrying for retryable errors.
type smsQueue struct {
	channel                 chan smsMessage
	maxRetryCount           int
	messagesPerSecond       int
	started                 bool
	client                  SMSClient
	sentMessageCallbackFunc func(smsMessage, recipient, smsMessageId string)
}

// Enqueues a message in the SMS queue.
func (q *smsQueue) EnqueueMessage(message string, recipients []string, chargingAmount float32, reportDelivery bool) {
	go func() {
		m := smsMessage{message, recipients, chargingAmount, reportDelivery, 0}
		q.enqueueMessage(m)
	}()
}

func (q *smsQueue) enqueueMessage(m smsMessage) {
	addrBlocks := splitAddrSlice(m.recipients, q.messagesPerSecond*q.client.MaxAddressCount)
	for _, block := range addrBlocks {
		nm := m
		nm.recipients = block
		q.channel <- nm
	}
}

func (q *smsQueue) requeueMessage(m smsMessage) {
	if m.retries < q.maxRetryCount {
		m.retries++
		q.enqueueMessage(m)
	}
}

func (q *smsQueue) sendMessage(m smsMessage) {
	responses, failures, err := q.client.SendTextMessage(m.message, m.recipients, m.chargingAmount, m.reportDelivery)
	if err != nil {
		q.requeueMessage(m)
		return
	}
	for i := range responses {
		if responses[i].Error != nil && responses[i].Error.Retryable {
			failures = append(failures, responses[i].Address)
		} else {
			go q.sentMessageCallbackFunc(m.message, responses[i].Address, responses[i].MessageID)
		}
	}
	if len(failures) > 0 {
		newM := smsMessage{m.message, failures, m.chargingAmount, m.reportDelivery, m.retries + 1}
		q.requeueMessage(newM)
	}
}

// Starts the SMS queue. This method should be called only once. Subsequent calls will not do anything.
func (q *smsQueue) Start() {
	if q.started {
		log.Print("SMS queue is already running.")
		return
	}
	q.started = true
	sentCount := 0
	for {
		if sentCount >= q.messagesPerSecond {
			time.Sleep(time.Second)
			sentCount = 0
		}
		m := <-q.channel
		go q.sendMessage(m)
		sentCount += len(m.recipients) / q.client.MaxAddressCount
		if len(m.recipients)%q.client.MaxAddressCount != 0 {
			sentCount++
		}
	}
}

// Initializes and returns a new SMS queue.
// Make sure that an application has only one queue if request throttling should be properly functional.
func NewSMSQueue(client *SMSClient, capacity, messagesPerSecond, maxRetryCount int, sendCallback func(smsMessage, recipient, smsMessageId string)) smsQueue {
	if client == nil {
		panic("SMS client is nil")
	}
	q := smsQueue{
		channel:                 make(chan smsMessage, capacity),
		maxRetryCount:           maxRetryCount,
		messagesPerSecond:       messagesPerSecond,
		client:                  *client,
		sentMessageCallbackFunc: sendCallback,
	}
	return q
}
