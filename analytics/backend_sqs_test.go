package analytics

import (
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/sqsiface"

	. "github.com/smartystreets/goconvey/convey"
)

type fakeSQS struct {
	sqsiface.SQSAPI
	input *sqs.SendMessageInput
}

func (s *fakeSQS) SendMessageRequest(i *sqs.SendMessageInput) sqs.SendMessageRequest {
	s.input = i
	return sqs.SendMessageRequest{
		Request: &aws.Request{},
	}
}

func TestSQSBackend(t *testing.T) {
	Convey("SQS backend should capture the right data", t, func() {
		backend, err := NewSQSBackend("https://fake.url")
		So(err, ShouldBeNil)
		So(backend, ShouldNotBeNil)

		fake := &fakeSQS{}
		backend.(*sqsBackend).sqsService = fake

		fakeReq, err := http.NewRequest("GET", "/", nil)
		So(err, ShouldBeNil)
		backend.Store(fakeReq, "/some/url", "some term", "list type", "gaID", "gID", 10, 20, 30)

		So(fake.input, ShouldNotBeNil)
	})
}
