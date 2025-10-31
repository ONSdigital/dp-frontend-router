package router_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-frontend-router/middleware/allRoutes"
	"github.com/ONSdigital/dp-frontend-router/middleware/allRoutes/allroutestest"
	"github.com/ONSdigital/dp-frontend-router/router"
	"github.com/ONSdigital/dp-frontend-router/router/routertest"
	. "github.com/smartystreets/goconvey/convey"
)

func NewHandlerMock() *routertest.HandlerMock {
	return &routertest.HandlerMock{
		ServeHTTPFunc: func(_ http.ResponseWriter, _ *http.Request) {},
	}
}

func TestSecurityHandler(t *testing.T) {
	Convey("Given a security handler", t, func() {
		res := httptest.NewRecorder()
		handler := NewHandlerMock()
		securityHandler := router.SecurityHandler(handler)

		Convey("When a default request is made", func() {
			url := "/"
			req := httptest.NewRequest(http.MethodGet, url, http.NoBody)
			securityHandler.ServeHTTP(res, req)

			Convey("Then xframe-options header is SAMEORIGIN", func() {
				So(res.Header().Get(router.HTTPHeaderKeyXFrameOptions), ShouldEqual, "SAMEORIGIN")
			})

			Convey("And the request is sent to the underlying handler", func() {
				So(len(handler.ServeHTTPCalls()), ShouldEqual, 1)
				So(handler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
		})

		Convey("When a handled request is made we want no xframe-options header", func() {
			urls := []string{
				"/embed",
				"/visualisations/path",
				"/census/maps/path",
			}
			for i, url := range urls {
				req := httptest.NewRequest(http.MethodGet, url, http.NoBody)
				securityHandler.ServeHTTP(res, req)

				Convey("Then no xframe-options header is set: "+url, func() {
					So(res.Header().Get(router.HTTPHeaderKeyXFrameOptions), ShouldBeEmpty)
				})

				Convey("And the request is sent to the underlying handler: "+url, func() {
					So(len(handler.ServeHTTPCalls()), ShouldEqual, i+1)
					So(handler.ServeHTTPCalls()[i].In2.URL.Path, ShouldResemble, url)
				})
			}
		})
	})
}

func TestRouter(t *testing.T) {
	Convey("Given a configured router", t, func() {
		healthCheckHandler := NewHandlerMock()
		searchHandler := NewHandlerMock()
		downloadHandler := NewHandlerMock()
		cookieHandler := NewHandlerMock()
		datasetHandler := NewHandlerMock()
		filterHandler := NewHandlerMock()
		filterFlexHandler := NewHandlerMock()
		feedbackHandler := NewHandlerMock()
		babbageHandler := NewHandlerMock()
		homepageHandler := NewHandlerMock()
		censusAtlasHandler := NewHandlerMock()
		releaseCalendarHandler := NewHandlerMock()
		prefixDatasetHandler := NewHandlerMock()
		proxyHandler := NewHandlerMock()

		zebedeeClient := &allroutestest.ZebedeeClientMock{
			GetWithHeadersFunc: func(_ context.Context, _ string, _ string) ([]byte, http.Header, error) {
				return make([]byte, 0), http.Header{}, nil
			},
		}

		config := router.Config{
			SearchHandler:        searchHandler,
			DownloadHandler:      downloadHandler,
			HealthCheckHandler:   healthCheckHandler.ServeHTTPFunc,
			CookieHandler:        cookieHandler,
			DatasetHandler:       datasetHandler,
			PrefixDatasetHandler: prefixDatasetHandler,
			FilterHandler:        filterHandler,
			FilterFlexHandler:    filterFlexHandler,
			FeedbackHandler:      feedbackHandler,
			ZebedeeClient:        zebedeeClient,
			BabbageHandler:       babbageHandler,
			HomepageHandler:      homepageHandler,
			CensusAtlasHandler:   censusAtlasHandler,
			RelCalHandler:        releaseCalendarHandler,
			ProxyHandler:         proxyHandler,
		}

		// Common URLs
		const (
			alladhocsURL             = "/alladhocs"
			cpiURL                   = "/inflationandpriceindices/datasets/consumerpriceinflation/current"
			datalistURL              = "/datalist"
			economyURL               = "/economy"
			environmentalaccountsURL = "/environmentalaccounts"
			findADatasetURL          = "/census/find-a-dataset"
			filterURL                = "/filters/123"
			previousreleasesURL      = "/previousreleases"
			publicsectorfinanceURL   = "/governmentpublicsectorandtaxes/publicsectorfinance"
			publicationsURL          = "/publications"
			relateddataURL           = "/relateddata"
			searchURL                = "/search"
		)

		Convey("When a download request is made", func() {
			url := "/download/123"
			req := httptest.NewRequest("GET", url, http.NoBody)
			res := httptest.NewRecorder()

			r := router.New(config)
			r.ServeHTTP(res, req)

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
			req := httptest.NewRequest("GET", url, http.NoBody)
			res := httptest.NewRecorder()

			r := router.New(config)
			r.ServeHTTP(res, req)
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
			req := httptest.NewRequest("GET", url, http.NoBody)
			res := httptest.NewRecorder()

			r := router.New(config)
			r.ServeHTTP(res, req)

			Convey("Then no requests are sent to Zebedee", func() {
				So(len(zebedeeClient.GetWithHeadersCalls()), ShouldEqual, 0)
			})

			Convey("Then the request is sent to the dataset handler", func() {
				So(len(datasetHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(datasetHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
		})

		Convey("When a filter-output request is made", func() {
			url := "/filter-outputs/321"
			req := httptest.NewRequest("GET", url, http.NoBody)
			res := httptest.NewRecorder()

			r := router.New(config)
			r.ServeHTTP(res, req)
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
			req := httptest.NewRequest("GET", url, http.NoBody)
			res := httptest.NewRecorder()
			r := router.New(config)
			r.ServeHTTP(res, req)

			Convey("Then no requests are sent to Zebedee", func() {
				So(len(zebedeeClient.GetWithHeadersCalls()), ShouldEqual, 0)
			})
			Convey("Then the request is sent to the feedback handler", func() {
				So(len(feedbackHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(feedbackHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
		})
		Convey("When a search request is made, but the search handler is not enabled", func() {
			url := searchURL
			req := httptest.NewRequest("GET", url, http.NoBody)
			res := httptest.NewRecorder()

			config.SearchRoutesEnabled = false
			r := router.New(config)
			r.ServeHTTP(res, req)
			Convey("Then no request is sent to the search handler", func() {
				So(len(searchHandler.ServeHTTPCalls()), ShouldEqual, 0)
			})
			Convey("Then the request is sent to Babbage", func() {
				So(len(babbageHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(babbageHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
		})
		Convey("When a search request is made, and the search handler is enabled", func() {
			url := searchURL
			req := httptest.NewRequest("GET", url, http.NoBody)
			res := httptest.NewRecorder()

			config.SearchRoutesEnabled = true
			r := router.New(config)
			r.ServeHTTP(res, req)
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
		Convey("When a data aggregation page request is made, but the data aggregation pages are not enabled", func() {
			url := alladhocsURL
			req := httptest.NewRequest("GET", url, http.NoBody)
			res := httptest.NewRecorder()

			config.SearchRoutesEnabled = true
			config.DataAggregationPagesEnabled = false
			r := router.New(config)
			r.ServeHTTP(res, req)
			Convey("Then no request is sent to the search handler", func() {
				So(len(searchHandler.ServeHTTPCalls()), ShouldEqual, 0)
			})
			Convey("Then the request is sent to Babbage", func() {
				So(len(babbageHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(babbageHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
		})
		Convey("When a data aggregation page request is made, but the data aggregation pages are enabled", func() {
			url := alladhocsURL
			req := httptest.NewRequest("GET", url, http.NoBody)
			res := httptest.NewRecorder()

			config.SearchRoutesEnabled = true
			config.DataAggregationPagesEnabled = true
			r := router.New(config)
			r.ServeHTTP(res, req)
			Convey("Then the request is sent to the search handler", func() {
				So(len(searchHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(searchHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
			Convey("Then no request is sent to Babbage", func() {
				So(len(babbageHandler.ServeHTTPCalls()), ShouldEqual, 0)
			})
		})
		Convey("When a legacy data aggregation page request is made, but the data aggregation pages are not enabled", func() {
			url := economyURL + alladhocsURL
			req := httptest.NewRequest("GET", url, http.NoBody)
			res := httptest.NewRecorder()

			config.SearchRoutesEnabled = true
			config.DataAggregationPagesEnabled = false
			r := router.New(config)
			r.ServeHTTP(res, req)
			Convey("Then no request is sent to the search handler", func() {
				So(len(searchHandler.ServeHTTPCalls()), ShouldEqual, 0)
			})
			Convey("Then the request is sent to Babbage", func() {
				So(len(babbageHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(babbageHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
		})
		Convey("When a legacy data aggregation page request is made, but the data aggregation pages are enabled", func() {
			url := economyURL + alladhocsURL
			req := httptest.NewRequest("GET", url, http.NoBody)
			res := httptest.NewRecorder()

			config.SearchRoutesEnabled = true
			config.DataAggregationPagesEnabled = true
			r := router.New(config)
			r.ServeHTTP(res, req)
			Convey("Then the request is not sent to the search handler", func() {
				So(len(searchHandler.ServeHTTPCalls()), ShouldEqual, 0)
			})
			Convey("And the page request has been redirected to the root", func() {
				So(res.Result().StatusCode, ShouldEqual, 301)
				So(res.Result().Header.Get("Location"), ShouldEqual, "/alladhocs")
			})
			Convey("And no request is sent to Babbage", func() {
				So(len(babbageHandler.ServeHTTPCalls()), ShouldEqual, 0)
			})
		})

		Convey("When a topic aggregation page request is made, but the topic aggregation pages are not enabled", func() {
			url := economyURL + publicationsURL
			req := httptest.NewRequest("GET", url, http.NoBody)
			res := httptest.NewRecorder()

			config.SearchRoutesEnabled = true
			config.TopicAggregationPagesEnabled = false

			r := router.New(config)
			r.ServeHTTP(res, req)
			Convey("Then no request is sent to the search handler", func() {
				So(len(searchHandler.ServeHTTPCalls()), ShouldEqual, 0)
			})
			Convey("Then the request is sent to Babbage", func() {
				So(len(babbageHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(babbageHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
		})
		Convey("When a topic aggregation page request is made, but the topic aggregation pages are enabled", func() {
			url := economyURL + publicationsURL
			req := httptest.NewRequest("GET", url, http.NoBody)
			res := httptest.NewRecorder()

			config.SearchRoutesEnabled = true
			config.TopicAggregationPagesEnabled = true

			r := router.New(config)
			r.ServeHTTP(res, req)
			Convey("Then the request is sent to the search handler", func() {
				So(len(searchHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(searchHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
			Convey("Then no request is sent to Babbage", func() {
				So(len(babbageHandler.ServeHTTPCalls()), ShouldEqual, 0)
			})
		})
		Convey("When a subtopic aggregation page request is made, but the subtopic aggregation pages are not enabled", func() {
			url := economyURL + environmentalaccountsURL + datalistURL
			req := httptest.NewRequest("GET", url, http.NoBody)
			res := httptest.NewRecorder()

			config.SearchRoutesEnabled = true
			config.TopicAggregationPagesEnabled = false
			r := router.New(config)
			r.ServeHTTP(res, req)
			Convey("Then no request is sent to the search handler", func() {
				So(len(searchHandler.ServeHTTPCalls()), ShouldEqual, 0)
			})
			Convey("Then the request is sent to Babbage", func() {
				So(len(babbageHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(babbageHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
		})
		Convey("When a subtopic aggregation page request is made, but the subtopic aggregation pages are enabled", func() {
			url := economyURL + environmentalaccountsURL + datalistURL
			req := httptest.NewRequest("GET", url, http.NoBody)
			res := httptest.NewRecorder()

			config.SearchRoutesEnabled = true
			config.TopicAggregationPagesEnabled = true
			r := router.New(config)
			r.ServeHTTP(res, req)
			Convey("Then the request is sent to the search handler", func() {
				So(len(searchHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(searchHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
			Convey("Then no request is sent to Babbage", func() {
				So(len(babbageHandler.ServeHTTPCalls()), ShouldEqual, 0)
			})
		})
		Convey("When a 3rd level topic aggregation page request is made, but the topic aggregation pages are not enabled", func() {
			url := economyURL + publicsectorfinanceURL + datalistURL
			req := httptest.NewRequest("GET", url, http.NoBody)
			res := httptest.NewRecorder()

			config.SearchRoutesEnabled = true
			config.TopicAggregationPagesEnabled = false
			r := router.New(config)
			r.ServeHTTP(res, req)
			Convey("Then no request is sent to the search handler", func() {
				So(len(searchHandler.ServeHTTPCalls()), ShouldEqual, 0)
			})
			Convey("Then the request is sent to Babbage", func() {
				So(len(babbageHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(babbageHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
		})
		Convey("When a 3rd level topic aggregation page request is made, but the topic aggregation pages are enabled", func() {
			url := economyURL + publicsectorfinanceURL + datalistURL
			req := httptest.NewRequest("GET", url, http.NoBody)
			res := httptest.NewRecorder()

			config.SearchRoutesEnabled = true
			config.TopicAggregationPagesEnabled = true
			r := router.New(config)
			r.ServeHTTP(res, req)
			Convey("Then the request is sent to the search handler", func() {
				So(len(searchHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(searchHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
			Convey("Then no request is sent to Babbage", func() {
				So(len(babbageHandler.ServeHTTPCalls()), ShouldEqual, 0)
			})
		})
		Convey("When a related data request is made, the RelatedDataRouteEnabled is enabled and the Legacy Cache Proxy is disabled", func() {
			url := economyURL + relateddataURL
			req := httptest.NewRequest("GET", url, http.NoBody)
			res := httptest.NewRecorder()

			config.RelatedDataRouteEnabled = true
			config.SearchRoutesEnabled = true
			config.LegacyCacheProxyEnabled = false

			r := router.New(config)
			r.ServeHTTP(res, req)
			Convey("Then the request is sent to the search handler", func() {
				So(searchHandler.ServeHTTPCalls(), ShouldHaveLength, 1)
				So(searchHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
			Convey("Then no request is sent to Babbage", func() {
				So(babbageHandler.ServeHTTPCalls(), ShouldHaveLength, 0)
			})
		})
		Convey("When a related data request is made, the RelatedDataRouteEnabled is enabled and the Legacy Cache Proxy is enabled", func() {
			url := economyURL + relateddataURL
			req := httptest.NewRequest("GET", url, http.NoBody)
			res := httptest.NewRecorder()

			config.RelatedDataRouteEnabled = true
			config.SearchRoutesEnabled = true
			config.LegacyCacheProxyEnabled = true
			r := router.New(config)
			r.ServeHTTP(res, req)
			Convey("Then the request is sent to the legacy cache proxy", func() {
				So(proxyHandler.ServeHTTPCalls(), ShouldHaveLength, 1)
				So(proxyHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
			Convey("Then no request is sent to Babbage", func() {
				So(babbageHandler.ServeHTTPCalls(), ShouldHaveLength, 0)
			})
			Convey("Then no request is sent to the Search Handler", func() {
				So(searchHandler.ServeHTTPCalls(), ShouldHaveLength, 0)
			})
		})

		Convey("When a previous releases request is made, the PreviousReleasesRoute is enabled and the Legacy Cache Proxy is disabled", func() {
			url := economyURL + previousreleasesURL
			req := httptest.NewRequest("GET", url, http.NoBody)
			res := httptest.NewRecorder()

			config.PreviousReleasesRouteEnabled = true
			config.SearchRoutesEnabled = true
			config.LegacyCacheProxyEnabled = false

			r := router.New(config)
			r.ServeHTTP(res, req)
			Convey("Then the request is sent to the search handler", func() {
				So(searchHandler.ServeHTTPCalls(), ShouldHaveLength, 1)
				So(searchHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
			Convey("Then no request is sent to Babbage", func() {
				So(babbageHandler.ServeHTTPCalls(), ShouldHaveLength, 0)
			})
		})
		Convey("When a previous releases request is made, the PreviousReleasesRoute is enabled and the Legacy Cache Proxy is enabled", func() {
			url := economyURL + previousreleasesURL
			req := httptest.NewRequest("GET", url, http.NoBody)
			res := httptest.NewRecorder()

			config.PreviousReleasesRouteEnabled = true
			config.SearchRoutesEnabled = true
			config.LegacyCacheProxyEnabled = true
			r := router.New(config)
			r.ServeHTTP(res, req)
			Convey("Then the request is sent to the legacy cache proxy", func() {
				So(proxyHandler.ServeHTTPCalls(), ShouldHaveLength, 1)
				So(proxyHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
			Convey("Then no request is sent to Babbage", func() {
				So(babbageHandler.ServeHTTPCalls(), ShouldHaveLength, 0)
			})
			Convey("Then no request is sent to the Search Handler", func() {
				So(searchHandler.ServeHTTPCalls(), ShouldHaveLength, 0)
			})
		})

		Convey("When a dataset finder request is made, but the Dataset Finder is not enabled", func() {
			url := findADatasetURL
			req := httptest.NewRequest("GET", url, http.NoBody)
			res := httptest.NewRecorder()

			config.DatasetFinderEnabled = false
			r := router.New(config)
			r.ServeHTTP(res, req)
			Convey("Then no request is sent to the search handler", func() {
				So(len(searchHandler.ServeHTTPCalls()), ShouldEqual, 0)
			})
			Convey("Then the request is sent to Babbage", func() {
				So(len(babbageHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(babbageHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
		})
		Convey("When a dataset finder request is made, and the Dataset Finder is enabled", func() {
			url := findADatasetURL
			req := httptest.NewRequest("GET", url, http.NoBody)
			res := httptest.NewRecorder()

			config.DatasetFinderEnabled = true
			r := router.New(config)
			r.ServeHTTP(res, req)
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
			req := httptest.NewRequest("GET", url, http.NoBody)
			res := httptest.NewRecorder()

			r := router.New(config)
			r.ServeHTTP(res, req)
			Convey("Then no requests are sent to Zebedee", func() {
				So(len(zebedeeClient.GetWithHeadersCalls()), ShouldEqual, 0)
			})
			Convey("Then the request is sent to the homepage handler", func() {
				So(len(homepageHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(homepageHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
		})

		Convey("When a legacy page request is made", func() {
			url := economyURL
			expectedZebedeeURL := "/data?uri=" + url
			req := httptest.NewRequest("GET", url, http.NoBody)
			res := httptest.NewRecorder()

			r := router.New(config)
			r.ServeHTTP(res, req)
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
			req := httptest.NewRequest("GET", url, http.NoBody)
			res := httptest.NewRecorder()

			r := router.New(config)
			r.ServeHTTP(res, req)
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
			req := httptest.NewRequest("GET", url, http.NoBody)
			res := httptest.NewRecorder()

			r := router.New(config)
			r.ServeHTTP(res, req)
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
			req := httptest.NewRequest("GET", url, http.NoBody)
			res := httptest.NewRecorder()

			r := router.New(config)
			r.ServeHTTP(res, req)
			Convey("Then a request is not sent to Zebedee to check the page type", func() {
				So(len(zebedeeClient.GetWithHeadersCalls()), ShouldEqual, 0)
			})
			Convey("Then the request is sent to the babbage", func() {
				So(len(babbageHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(babbageHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
		})

		Convey("When a request for a legacy ons redirect endpoint is made", func() {
			url := "/ons/some/old/page"
			req := httptest.NewRequest("GET", url, http.NoBody)
			res := httptest.NewRecorder()

			r := router.New(config)
			r.ServeHTTP(res, req)
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
			req := httptest.NewRequest("GET", url, http.NoBody)
			res := httptest.NewRecorder()

			r := router.New(config)
			r.ServeHTTP(res, req)

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
			req := httptest.NewRequest("GET", url, http.NoBody)
			res := httptest.NewRecorder()

			r := router.New(config)
			r.ServeHTTP(res, req)

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
			req := httptest.NewRequest("GET", url, http.NoBody)
			res := httptest.NewRecorder()

			r := router.New(config)
			r.ServeHTTP(res, req)
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
			req := httptest.NewRequest("GET", url, http.NoBody)
			res := httptest.NewRecorder()

			config.CensusAtlasEnabled = false
			r := router.New(config)
			r.ServeHTTP(res, req)
			Convey("Then no request is sent to the census atlas handler", func() {
				So(len(censusAtlasHandler.ServeHTTPCalls()), ShouldEqual, 0)
			})
		})

		Convey("When a census atlas request is made, and the census atlas handler is enabled", func() {
			url := "/census/maps"
			req := httptest.NewRequest("GET", url, http.NoBody)
			res := httptest.NewRecorder()

			config.CensusAtlasEnabled = true
			r := router.New(config)
			r.ServeHTTP(res, req)

			Convey("Then the request is sent to the census atlas handler", func() {
				So(len(censusAtlasHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(censusAtlasHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
		})

		Convey("When a release calendar request is made", func() {
			url := "/releasecalendar"
			res := httptest.NewRecorder()
			r := router.New(config)
			req := httptest.NewRequest("GET", url, http.NoBody)

			r.ServeHTTP(res, req)

			Convey("Then the request is sent to the release calendar handler", func() {
				So(releaseCalendarHandler.ServeHTTPCalls(), ShouldHaveLength, 1)
			})

			Convey("And no request is sent to the babbage handler", func() {
				So(babbageHandler.ServeHTTPCalls(), ShouldBeEmpty)
			})
		})

		Convey("When a malicious URL with a redirect attempt is made", func() {
			url := "//%5cexample.com"
			req := httptest.NewRequest("GET", url, http.NoBody)
			w := httptest.NewRecorder()

			r := router.New(config)
			r.ServeHTTP(w, req)

			Convey("Then the request is redirected but with the path properly escaped", func() {
				res := w.Result()
				So(res.StatusCode, ShouldEqual, http.StatusMovedPermanently)
				So(res.Header.Get("Location"), ShouldResemble, "/%5Cexample.com")
			})
		})

		Convey("When a legacy dataset request is made, but the dataset handler is not enabled", func() {
			url := economyURL + cpiURL
			req := httptest.NewRequest("GET", url, http.NoBody)
			res := httptest.NewRecorder()

			config.NewDatasetRoutingEnabled = false
			newRouter := router.New(config)
			newRouter.ServeHTTP(res, req)

			Convey("Then a request is sent to Zebedee to check the page type", func() {
				So(len(zebedeeClient.GetWithHeadersCalls()), ShouldEqual, 1)
			})

			Convey("Then no request is sent to the dataset handler", func() {
				So(len(prefixDatasetHandler.ServeHTTPCalls()), ShouldEqual, 0)
			})

			Convey("Then the request is sent to Babbage", func() {
				So(len(babbageHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(babbageHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})
		})

		Convey("When a legacy dataset request is made, and the dataset handler is enabled", func() {
			url := economyURL + cpiURL
			req := httptest.NewRequest("GET", url, http.NoBody)
			res := httptest.NewRecorder()

			expectedZebedeeURL := "/data?uri=" + url

			// mock allRouteHandler's zebedee response to return dataset page type
			zebedeeResponseBody := json.RawMessage(`{"type":"dataset","apiDatasetId":""}`)
			zebedeeClient = &allroutestest.ZebedeeClientMock{
				GetWithHeadersFunc: func(_ context.Context, _ string, _ string) ([]byte, http.Header, error) {
					h := http.Header{}
					h.Add(allRoutes.HeaderOnsPageType, "dataset")
					return zebedeeResponseBody, h, nil
				},
			}
			config.ZebedeeClient = zebedeeClient

			config.NewDatasetRoutingEnabled = true
			config.ContentTypeByteLimit = 5 * 1000 * 1000

			newRouter := router.New(config)
			newRouter.ServeHTTP(res, req)

			Convey("Then a request is sent to Zebedee to check the page type", func() {
				So(len(zebedeeClient.GetWithHeadersCalls()), ShouldEqual, 1)
				So(zebedeeClient.GetWithHeadersCalls()[0].Path, ShouldEqual, expectedZebedeeURL)
			})

			Convey("Then the request is sent to the prefix dataset handler", func() {
				So(len(prefixDatasetHandler.ServeHTTPCalls()), ShouldEqual, 1)
				So(prefixDatasetHandler.ServeHTTPCalls()[0].In2.URL.Path, ShouldResemble, url)
			})

			Convey("Then no request is sent to Babbage", func() {
				So(len(babbageHandler.ServeHTTPCalls()), ShouldEqual, 0)
			})
		})
	})
}

func TestLegacyCacheFeatureFlag(t *testing.T) {
	Convey("When a /releases/ request is made", t, func() {
		babbageHandler := NewHandlerMock()
		releaseCalendarHandler := NewHandlerMock()
		legacyCacheProxyHandler := NewHandlerMock()

		config := router.Config{
			BabbageHandler: babbageHandler,
			RelCalHandler:  releaseCalendarHandler,
			ProxyHandler:   legacyCacheProxyHandler,
		}

		url := "/releases/"
		res := httptest.NewRecorder()

		Convey("And the legacy cache proxy is enabled", func() {
			config.LegacyCacheProxyEnabled = true
			req := httptest.NewRequest("GET", url, http.NoBody)
			r := router.New(config)
			r.ServeHTTP(res, req)

			Convey("Then the request is sent to legacy cache proxy", func() {
				So(legacyCacheProxyHandler.ServeHTTPCalls(), ShouldHaveLength, 1)
			})
			Convey("No request is sent to release calendar", func() {
				So(releaseCalendarHandler.ServeHTTPCalls(), ShouldBeEmpty)
			})
			Convey("No request is sent to Babbage", func() {
				So(babbageHandler.ServeHTTPCalls(), ShouldBeEmpty)
			})
		})

		Convey("And the legacy cache proxy is not enabled", func() {
			config.LegacyCacheProxyEnabled = false
			req := httptest.NewRequest("GET", url, http.NoBody)
			r := router.New(config)
			r.ServeHTTP(res, req)

			Convey("Then the request is sent to release calendar", func() {
				So(releaseCalendarHandler.ServeHTTPCalls(), ShouldHaveLength, 1)
			})
			Convey("No request is sent to Babbage", func() {
				So(babbageHandler.ServeHTTPCalls(), ShouldBeEmpty)
			})
		})
	})
}
