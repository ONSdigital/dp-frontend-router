package analytics

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ONSdigital/go-ns/log"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

var _ ServiceBackend = &sqsBackend{}

type sqsBackend struct {
	sqsService *sqs.SQS
	queueURL   string
}

// NewSQSBackend creates a new SQS backend for storing analytics data
func NewSQSBackend(queueURL string) (ServiceBackend, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, err
	}

	return &sqsBackend{sqs.New(cfg), queueURL}, nil
}

func (b *sqsBackend) Store(req *http.Request, url, term, listType, gaID string, pageIndex, linkIndex, pageSize float64) {
	var data = map[string]interface{}{
		"created":   time.Now().Format(time.RFC3339),
		"url":       url,
		"term":      term,
		"listType":  listType,
		"gaID":      gaID,
		"pageIndex": pageIndex,
		"linkIndex": linkIndex,
		"pageSize":  pageSize,
	}

	jb, err := json.Marshal(&data)
	if err != nil {
		log.ErrorR(req, err, nil)
		return
	}

	json := string(jb)
	smr := b.sqsService.SendMessageRequest(&sqs.SendMessageInput{
		MessageBody: &json,
		QueueUrl:    &b.queueURL,
	})

	smo, err := smr.Send()
	if err != nil {
		log.ErrorR(req, err, nil)
		return
	}

	log.DebugR(req, "stored analytics data in SQS", log.Data{"message_id": *smo.MessageId})
}
