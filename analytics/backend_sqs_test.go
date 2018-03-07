package analytics

import (
	"encoding/json"
	"log"
	"net/http"
	"testing"
	"time"

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
		backend.(*sqsBackend).sendMessage = func(s sqs.SendMessageRequest) (*sqs.SendMessageOutput, error) {
			msgID := "test-message-id"
			return &sqs.SendMessageOutput{
				MessageId: &msgID,
			}, nil
		}

		fakeReq, err := http.NewRequest("GET", "/", nil)
		So(err, ShouldBeNil)
		backend.Store(fakeReq, "/some/url", "some term", "list type", "gaID", "gID", 10, 20, 30)

		So(fake.input, ShouldNotBeNil)
		log.Printf("%+v", fake.input)

		var input map[string]interface{}
		err = json.Unmarshal([]byte(*fake.input.MessageBody), &input)
		So(err, ShouldBeNil)
		So(input, ShouldContainKey, "created")
		So(input, ShouldContainKey, "gID")
		So(input, ShouldContainKey, "gaID")
		So(input, ShouldContainKey, "linkIndex")
		So(input, ShouldContainKey, "listType")
		So(input, ShouldContainKey, "pageIndex")
		So(input, ShouldContainKey, "pageSize")
		So(input, ShouldContainKey, "term")
		So(input, ShouldContainKey, "url")

		So(input["created"], ShouldStartWith, time.Now().Format("2006-01-02T15:04"))
		So(input["gID"], ShouldEqual, "gID")
		So(input["gaID"], ShouldEqual, "gaID")
		So(input["linkIndex"], ShouldEqual, 20)
		So(input["listType"], ShouldEqual, "list type")
		So(input["pageIndex"], ShouldEqual, 10)
		So(input["pageSize"], ShouldEqual, 30)
		So(input["term"], ShouldEqual, "some term")
		So(input["url"], ShouldEqual, "/some/url")
	})
}
