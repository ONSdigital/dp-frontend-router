package router_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-frontend-router/middleware/allRoutes/allroutestest"
	"github.com/ONSdigital/dp-frontend-router/middleware/datasetType/mocks"
	"github.com/ONSdigital/dp-frontend-router/router"
	"github.com/ONSdigital/dp-frontend-router/router/routertest"
	. "github.com/smartystreets/goconvey/convey"
)

func NewHandlerMock() *routertest.HandlerMock {
	return &routertest.HandlerMock{
		ServeHTTPFunc: func(in1 http.ResponseWriter, in2 *http.Request) {},
	}
}

func TestRouter(t *testing.T) {

	Convey("Given a configured router", t, func() {

		healthCheckHandler := NewHandlerMock()
		analyticsHandler := NewHandlerMock()
		searchHandler := NewHandlerMock()
		downloadHandler := NewHandlerMock()
		cookieHandler := NewHandlerMock()
		datasetHandler := NewHandlerMock()
		filterHandler := NewHandlerMock()
		filterFlexHandler := NewHandlerMock()
		feedbackHandler := NewHandlerMock()
		babbageHandler := NewHandlerMock()
		geographyHandler := NewHandlerMock()
		homepageHandler := NewHandlerMock()
		interactivesHandler := NewHandlerMock()
		censusAtlasHandler := NewHandlerMock()

		zebedeeClient := &allroutestest.ZebedeeClientMock{
			GetWithHeadersFunc: func(ctx context.Context, userAccessToken string, path string) ([]byte, http.Header, error) {
				return make([]byte, 0), http.Header{}, nil
			},
		}

		filterClient := &mocks.FilterClientMock{
			GetJobStateFunc: func(ctx context.Context, userAuthToken, serviceAuthToken, downloadServiceToken, collectionID, filterID string) (filter.Model, string, error) {
				return filter.Model{}, "", nil
			},
		}

		datasetClient := &mocks.DatasetClientMock{
			GetFunc: func(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, datasetID string) (dataset.DatasetDetails, error) {
				return dataset.DatasetDetails{}, nil
			},
		}

		config := router.Config{
			AnalyticsHandler:    analyticsHandler,
			SearchHandler:       searchHandler,
			DownloadHandler:     downloadHandler,
			HealthCheckHandler:  healthCheckHandler.ServeHTTPFunc,
			CookieHandler:       cookieHandler,
			DatasetHandler:      datasetHandler,
			DatasetClient:       datasetClient,
			FilterHandler:       filterHandler,
			FilterClient:        filterClient,
			FilterFlexHandler:   filterFlexHandler,
			FeedbackHandler:     feedbackHandler,
			ZebedeeClient:       zebedeeClient,
			BabbageHandler:      babbageHandler,
			GeographyHandler:    geographyHandler,
			HomepageHandler:     homepageHandler,
			InteractivesHandler: interactivesHandler,
			CensusAtlasHandler:  censusAtlasHandler,
		}

		Convey("When a analytics request is made", func() {

			url := "/redir/123"
			req := httptest.NewRequest("GET", url, nil)
			res := httptest.NewRecorder()

			router := router.New(config)
			router.ServeHTTP(res, req)

			Convey("Then no requests are sent to Zebedee", func() {
				So(len(zebedeeClient.GetWithHeadersCalls()), ShouldEqual, 0)
			})

			Convey("Then the request is sent to the search handler", func() {
				So(len(analyticsHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(analyticsHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
		})

		Convey("When a download request is made", func() {

			url := "/download/123"
			req := httptest.NewRequest("GET", url, nil)
			res := httptest.NewRecorder()

			router := router.New(config)
			router.ServeHTTP(res, req)

			Convey("Then no requests are sent to Zebedee", func() {
				So(len(zebedeeClient.GetWithHeadersCalls()), ShouldEqual, 0)
			})

			Convey("Then the request is sent to the download handler", func() {
				So(len(downloadHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(downloadHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
		})

		Convey("When a cookie request is made", func() {

			url := "/cookies/123/345"
			req := httptest.NewRequest("GET", url, nil)
			res := httptest.NewRecorder()

			router := router.New(config)
			router.ServeHTTP(res, req)

			Convey("Then no requests are sent to Zebedee", func() {
				So(len(zebedeeClient.GetWithHeadersCalls()), ShouldEqual, 0)
			})

			Convey("Then the request is sent to the cookie handler", func() {
				So(len(cookieHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(cookieHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
		})

		Convey("When a dataset request is made", func() {

			url := "/datasets/cpih"
			req := httptest.NewRequest("GET", url, nil)
			res := httptest.NewRecorder()

			router := router.New(config)
			router.ServeHTTP(res, req)

			Convey("Then no requests are sent to Zebedee", func() {
				So(len(zebedeeClient.GetWithHeadersCalls()), ShouldEqual, 0)
			})

			Convey("Then the request is sent to the dataset handler", func() {
				So(len(datasetHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(datasetHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
		})

		Convey("When a filter request is made, but the filter/flex handler is not enabled", func() {

			url := "/filters/123"
			req := httptest.NewRequest("GET", url, nil)
			res := httptest.NewRecorder()

			router := router.New(config)
			router.ServeHTTP(res, req)

			Convey("Then no requests are sent to Zebedee", func() {
				So(len(zebedeeClient.GetWithHeadersCalls()), ShouldEqual, 0)
			})

			Convey("Then the request is sent to the filter handler", func() {
				So(len(filterHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(filterHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})

			Convey("Then no requests are sent to the filter/flex handler", func() {
				So(len(filterFlexHandler.ServeHTTPCalls()), ShouldEqual, 0)
			})
		})

		Convey("When a filter request is made and the filter/flex route is enabled", func() {

			url := "/filters/123"
			req := httptest.NewRequest("GET", url, nil)
			res := httptest.NewRecorder()
			config.FilterFlexEnabled = true

			filterableDataset := &mocks.DatasetClientMock{
				GetFunc: func(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, datasetID string) (dataset.DatasetDetails, error) {
					return dataset.DatasetDetails{
						Type: "filterable",
					}, nil
				},
			}
			config.DatasetClient = filterableDataset

			router := router.New(config)
			router.ServeHTTP(res, req)

			Convey("Then no requests are sent to Zebedee", func() {
				So(len(zebedeeClient.GetWithHeadersCalls()), ShouldEqual, 0)
			})

			Convey("Then the request is sent to the filter handler", func() {
				So(len(filterHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(filterHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})

			Convey("Then no requests are sent to the filter/flex handler", func() {
				So(len(filterFlexHandler.ServeHTTPCalls()), ShouldEqual, 0)
			})
		})

		Convey("When a filter request is made for a valid flexible dataset", func() {

			url := "/filters/123"
			req := httptest.NewRequest("GET", url, nil)
			res := httptest.NewRecorder()
			config.FilterFlexEnabled = true

			flexDataset := &mocks.DatasetClientMock{
				GetFunc: func(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, datasetID string) (dataset.DatasetDetails, error) {
					return dataset.DatasetDetails{
						Type: "canatabular_flexible_table",
					}, nil
				},
			}
			config.DatasetClient = flexDataset

			router := router.New(config)
			router.ServeHTTP(res, req)

			Convey("Then no requests are sent to Zebedee", func() {
				So(len(zebedeeClient.GetWithHeadersCalls()), ShouldEqual, 0)
			})

			Convey("Then the request is sent to the filter/flex handler", func() {
				So(len(filterFlexHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(filterFlexHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})

			Convey("Then no requests are sent to the filter handler", func() {
				So(len(filterHandler.ServeHTTPCalls()), ShouldEqual, 0)
			})
		})

		Convey("When a filter-output request is made", func() {

			url := "/filter-outputs/321"
			req := httptest.NewRequest("GET", url, nil)
			res := httptest.NewRecorder()

			router := router.New(config)
			router.ServeHTTP(res, req)

			Convey("Then no requests are sent to Zebedee", func() {
				So(len(zebedeeClient.GetWithHeadersCalls()), ShouldEqual, 0)
			})

			Convey("Then the request is sent to the filter handler", func() {
				So(len(filterHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(filterHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
		})

		Convey("When a feedback request is made", func() {

			url := "/feedback/homepage"
			req := httptest.NewRequest("GET", url, nil)
			res := httptest.NewRecorder()

			router := router.New(config)
			router.ServeHTTP(res, req)

			Convey("Then no requests are sent to Zebedee", func() {
				So(len(zebedeeClient.GetWithHeadersCalls()), ShouldEqual, 0)
			})

			Convey("Then the request is sent to the feedback handler", func() {
				So(len(feedbackHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(feedbackHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
		})

		Convey("When a geography request is made, but the geography handler is not enabled", func() {

			url := "/geography/newport"
			req := httptest.NewRequest("GET", url, nil)
			res := httptest.NewRecorder()

			config.GeographyEnabled = false
			router := router.New(config)
			router.ServeHTTP(res, req)

			Convey("Then a request is sent to Zebedee to check the page type", func() {
				So(len(zebedeeClient.GetWithHeadersCalls()), ShouldEqual, 1)
			})

			Convey("Then no request is sent to the geography handler", func() {
				So(len(geographyHandler.ServeHTTPCalls()), ShouldEqual, 0)
			})

			Convey("Then the request is sent to Babbage", func() {
				So(len(babbageHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(babbageHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
		})

		Convey("When a geography request is made, and the geography handler is enabled", func() {

			url := "/geography/newport"
			req := httptest.NewRequest("GET", url, nil)
			res := httptest.NewRecorder()

			config.GeographyEnabled = true
			router := router.New(config)
			router.ServeHTTP(res, req)

			Convey("Then no requests are sent to Zebedee", func() {
				So(len(zebedeeClient.GetWithHeadersCalls()), ShouldEqual, 0)
			})

			Convey("Then the request is sent to the geography handler", func() {
				So(len(geographyHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(geographyHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})

			Convey("Then no request is sent to Babbage", func() {
				So(len(babbageHandler.ServeHTTPCalls()), ShouldEqual, 0)
			})
		})

		Convey("When a search request is made, but the search handler is not enabled", func() {

			url := "/search"
			req := httptest.NewRequest("GET", url, nil)
			res := httptest.NewRecorder()

			config.SearchRoutesEnabled = false
			router := router.New(config)
			router.ServeHTTP(res, req)

			Convey("Then no request is sent to the search handler", func() {
				So(len(searchHandler.ServeHTTPCalls()), ShouldEqual, 0)
			})

			Convey("Then the request is sent to Babbage", func() {
				So(len(babbageHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(babbageHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
		})

		Convey("When a search request is made, and the search handler is enabled", func() {

			url := "/search"
			req := httptest.NewRequest("GET", url, nil)
			res := httptest.NewRecorder()

			config.SearchRoutesEnabled = true
			router := router.New(config)
			router.ServeHTTP(res, req)

			Convey("Then no requests are sent to Zebedee", func() {
				So(len(zebedeeClient.GetWithHeadersCalls()), ShouldEqual, 0)
			})

			Convey("Then the request is sent to the search handler", func() {
				So(len(searchHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(searchHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})

			Convey("Then no request is sent to Babbage", func() {
				So(len(babbageHandler.ServeHTTPCalls()), ShouldEqual, 0)
			})
		})

		Convey("When a homepage request is made", func() {

			url := "/"
			req := httptest.NewRequest("GET", url, nil)
			res := httptest.NewRecorder()

			router := router.New(config)
			router.ServeHTTP(res, req)

			Convey("Then no requests are sent to Zebedee", func() {
				So(len(zebedeeClient.GetWithHeadersCalls()), ShouldEqual, 0)
			})

			Convey("Then the request is sent to the homepage handler", func() {
				So(len(homepageHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(homepageHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
		})

		Convey("When a legacy page request is made", func() {

			url := "/economy"
			expectedZebedeeURL := "/data?uri=" + url
			req := httptest.NewRequest("GET", url, nil)
			res := httptest.NewRecorder()

			router := router.New(config)
			router.ServeHTTP(res, req)

			Convey("Then a request is sent to Zebedee to check the page type", func() {
				So(len(zebedeeClient.GetWithHeadersCalls()), ShouldEqual, 1)
				So(zebedeeClient.GetWithHeadersCalls()[0].Path, ShouldEqual, expectedZebedeeURL)
			})

			Convey("Then the request is sent to the babbage", func() {
				So(len(babbageHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(babbageHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
		})

		Convey("When a data.json request is made", func() {

			url := "/somepage/data.json"
			req := httptest.NewRequest("GET", url, nil)
			res := httptest.NewRecorder()

			router := router.New(config)
			router.ServeHTTP(res, req)

			Convey("Then a request is not sent to Zebedee to check the page type", func() {
				So(len(zebedeeClient.GetWithHeadersCalls()), ShouldEqual, 0)
			})

			Convey("Then the request is sent to Babbage", func() {
				So(len(babbageHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(babbageHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
		})

		Convey("When a request with a file extension is made", func() {

			url := "/website/main.css"
			req := httptest.NewRequest("GET", url, nil)
			res := httptest.NewRecorder()

			router := router.New(config)
			router.ServeHTTP(res, req)

			Convey("Then a request is not sent to Zebedee to check the page type", func() {
				So(len(zebedeeClient.GetWithHeadersCalls()), ShouldEqual, 0)
			})

			Convey("Then the request is sent to the babbage", func() {
				So(len(babbageHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(babbageHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
		})

		Convey("When a request for a visualisation endpoint is made", func() {

			url := "/visualisations/dvc1119"
			req := httptest.NewRequest("GET", url, nil)
			res := httptest.NewRecorder()

			router := router.New(config)
			router.ServeHTTP(res, req)

			Convey("Then a request is not sent to Zebedee to check the page type", func() {
				So(len(zebedeeClient.GetWithHeadersCalls()), ShouldEqual, 0)
			})

			Convey("Then the request is sent to the babbage", func() {
				So(len(babbageHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(babbageHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
		})

		Convey("When a request for a interactives endpoint is made, but the interactives handler is not enabled", func() {

			url := "/interactives/an_identifier"
			req := httptest.NewRequest("GET", url, nil)
			res := httptest.NewRecorder()

			router := router.New(config)
			router.ServeHTTP(res, req)

			Convey("Then no request is sent to the interactive visualisation handler", func() {
				So(len(interactivesHandler.ServeHTTPCalls()), ShouldEqual, 0)
			})

			Convey("Then the request is sent to Zebedee to check the page type", func() {
				So(len(zebedeeClient.GetWithHeadersCalls()), ShouldEqual, 1)
			})

			Convey("Then the request is sent to Babbage", func() {
				So(len(babbageHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(babbageHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
		})

		Convey("When a request for a interactives endpoint is made, and the interactives handler is enabled", func() {

			url := "/interactives/an_identifier"
			req := httptest.NewRequest("GET", url, nil)
			res := httptest.NewRecorder()

			config.InteractivesEnabled = true
			router := router.New(config)
			router.ServeHTTP(res, req)

			Convey("Then the request is sent to the interactive visualisation handler", func() {
				So(len(interactivesHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(interactivesHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})

			Convey("Then no requests are sent to Zebedee", func() {
				So(len(zebedeeClient.GetWithHeadersCalls()), ShouldEqual, 0)
			})

			Convey("Then no request is sent to Babbage", func() {
				So(len(babbageHandler.ServeHTTPCalls()), ShouldEqual, 0)
			})
		})

		Convey("When a request for a legacy ons redirect endpoint is made", func() {

			url := "/ons/some/old/page"
			req := httptest.NewRequest("GET", url, nil)
			res := httptest.NewRecorder()

			router := router.New(config)
			router.ServeHTTP(res, req)

			Convey("Then a request is not sent to Zebedee to check the page type", func() {
				So(len(zebedeeClient.GetWithHeadersCalls()), ShouldEqual, 0)
			})

			Convey("Then the request is sent to the babbage", func() {
				So(len(babbageHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(babbageHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
		})

		Convey("When a request for a known babbage endpoint is made", func() {

			url := "/file"
			req := httptest.NewRequest("GET", url, nil)
			res := httptest.NewRecorder()

			router := router.New(config)
			router.ServeHTTP(res, req)

			Convey("Then a request is not sent to Zebedee to check the page type", func() {
				So(len(zebedeeClient.GetWithHeadersCalls()), ShouldEqual, 0)
			})

			Convey("Then the request is sent to the babbage", func() {
				So(len(babbageHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(babbageHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
		})

		Convey("When a /data request is made", func() {

			url := "/somepage/data"
			req := httptest.NewRequest("GET", url, nil)
			res := httptest.NewRecorder()

			router := router.New(config)
			router.ServeHTTP(res, req)

			Convey("Then a request is not sent to Zebedee to check the page type", func() {
				So(len(zebedeeClient.GetWithHeadersCalls()), ShouldEqual, 0)
			})

			Convey("Then the request is sent to Babbage", func() {
				So(len(babbageHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(babbageHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
		})

		Convey("When a /latest request is made", func() {

			url := "/economy/environmentalaccounts/bulletins/ukenvironmentalaccounts/latest"
			req := httptest.NewRequest("GET", url, nil)
			res := httptest.NewRecorder()

			router := router.New(config)
			router.ServeHTTP(res, req)

			Convey("Then a request is not sent to Zebedee to check the page type", func() {
				So(len(zebedeeClient.GetWithHeadersCalls()), ShouldEqual, 0)
			})

			Convey("Then the request is sent to Babbage", func() {
				So(len(babbageHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(babbageHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
		})

		Convey("When a census atlas request is made, but the census atlas handler is not enabled", func() {

			url := "/census-atlas"
			req := httptest.NewRequest("GET", url, nil)
			res := httptest.NewRecorder()

			config.CensusAtlasEnabled = false
			router := router.New(config)
			router.ServeHTTP(res, req)

			Convey("Then no request is sent to the census atlas handler", func() {
				So(len(censusAtlasHandler.ServeHTTPCalls()), ShouldEqual, 0)
			})
		})

		Convey("When a census atlas request is made, and the census atlas handler is enabled", func() {

			url := "/census-atlas"
			req := httptest.NewRequest("GET", url, nil)
			res := httptest.NewRecorder()

			config.CensusAtlasEnabled = true
			router := router.New(config)
			router.ServeHTTP(res, req)

			Convey("Then the request is sent to the census atlas handler", func() {
				So(len(censusAtlasHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(censusAtlasHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
		})

	})
}
