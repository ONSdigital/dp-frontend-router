package analytics

import (
	"context"
	"encoding/json"
	"github.com/ONSdigital/dp-frontend-router/analytics/analyticstest"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"net/http"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSQSBackend(t *testing.T) {

	Convey("SQS backend initialise without error", t, func() {
		backend, err := NewSQSBackend(context.Background(), "https://fake.url")
		So(err, ShouldBeNil)
		So(backend, ShouldNotBeNil)
	})

	Convey("SQS backend should capture the right data", t, func() {
		mockSQSClient := &analyticstest.SQSClientMock{
			SendMessageFunc: func(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
				msgID := "test-message-id"
				return &sqs.SendMessageOutput{
					MessageId: &msgID,
				}, nil
			},
		}

		sqsBackend := &sqsBackend{
			mockSQSClient,
			"https://fake.url",
		}

		fakeReq, err := http.NewRequest("GET", "/", nil)
		So(err, ShouldBeNil)

		sqsBackend.Store(fakeReq, "/some/url", "some term", "list type", "gaID", "gID", 10, 20, 30)
		requestParams := mockSQSClient.SendMessageCalls()[0].Params
		So(requestParams, ShouldNotBeNil)

		var input map[string]interface{}
		err = json.Unmarshal([]byte(*requestParams.MessageBody), &input)
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
