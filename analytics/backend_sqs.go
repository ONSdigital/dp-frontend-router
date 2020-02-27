package analytics

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ONSdigital/log.go/log"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/sqsiface"
)

var _ ServiceBackend = &sqsBackend{}

type sqsBackend struct {
	sqsService  sqsiface.SQSAPI // *sqs.SQS
	queueURL    string
	sendMessage func(sqs.SendMessageRequest) (*sqs.SendMessageOutput, error)
}

// NewSQSBackend creates a new SQS backend for storing analytics data
func NewSQSBackend(queueURL string) (ServiceBackend, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, err
	}

	return &sqsBackend{sqs.New(cfg), queueURL, func(s sqs.SendMessageRequest) (*sqs.SendMessageOutput, error) {
		return s.Send()
	}}, nil
}

func (b *sqsBackend) Store(req *http.Request, url, term, listType, gaID string, gID string, pageIndex, linkIndex, pageSize float64) {
	var data = map[string]interface{}{
		"created":   time.Now().Format(time.RFC3339),
		"url":       url,
		"term":      term,
		"listType":  listType,
		"gaID":      gaID, // 2 year expiration cookie (_ga)
		"gID":       gID,  // 24 hour expiration cookie (_gid)
		"pageIndex": pageIndex,
		"linkIndex": linkIndex,
		"pageSize":  pageSize,
	}

	jb, err := json.Marshal(&data)
	if err != nil {
		log.Event(req.Context(), "error marshaling json", log.ERROR, log.Error(err))
		return
	}

	json := string(jb)
	smr := b.sqsService.SendMessageRequest(&sqs.SendMessageInput{
		MessageBody: &json,
		QueueUrl:    &b.queueURL,
	})

	smo, err := b.sendMessage(smr)
	if err != nil {
		log.Event(req.Context(), "error sending sqs message", log.ERROR, log.Error(err))
		return
	}

	log.Event(req.Context(), "stored analytics data in SQS", log.INFO, log.Data{"message_id": *smo.MessageId})
}
