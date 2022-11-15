package datasetType

import (
	"context"
	"net/http"
	"strings"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-frontend-router/helpers"
	"github.com/ONSdigital/log.go/v2/log"
)

//go:generate moq -out mocks/filter.go -pkg mocks . FilterClient
//go:generate moq -out mocks/dataset.go -pkg mocks . DatasetClient

// FilterClient is an interface with the methods required for a filter client
type FilterClient interface {
	GetJobState(ctx context.Context, userAuthToken, serviceAuthToken, downloadServiceToken, collectionID, filterID string) (f filter.Model, eTag string, err error)
}

// DatasetClient is an interface with methods required for a dataset client
type DatasetClient interface {
	Get(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, datasetID string) (m dataset.DatasetDetails, err error)
}

// Handler is middleware that a accepts a filter and dataset client that returns either the filter or filterFlex handler dependent on the type of dataset
func Handler(filterClient FilterClient, datasetClient DatasetClient) func(filter, filterFlex http.Handler) http.Handler {
	return func(filter, filterFlex http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			path := req.URL.Path
			ctx := req.Context()

			filterID, err := helpers.ReturnSecondSegmentFromPath(path)
			if err != nil {
				log.Error(ctx, "failed to extract filter id info from path", err, log.Data{"filter_id": filterID, "path": path})
				return
			}

			// Obtain access_token from cookie
			userAccessToken := ""
			c, err := req.Cookie(`access_token`)
			if err == nil && len(c.Value) > 0 {
				userAccessToken = c.Value
				log.Info(req.Context(), "obtained access_token Cookie")
			}

			f, _, err := filterClient.GetJobState(ctx, userAccessToken, "", "", "", filterID)
			if err != nil {
				log.Warn(ctx, "failed to get job state - falling through to default handling", log.Data{"filter_id": filterID})
				filter.ServeHTTP(w, req)
				return
			}

			d, err := datasetClient.Get(ctx, userAccessToken, "", "", f.Dataset.DatasetID)
			if err != nil {
				log.Warn(ctx, "failed to get dataset id - falling through to default handling", log.Data{"dataset_id": f.Dataset.DatasetID})
				filter.ServeHTTP(w, req)
				return
			}

			if strings.Contains(d.Type, "cantabular") {
				log.Info(ctx, "using flex handler")
				filterFlex.ServeHTTP(w, req)
				return
			}

			log.Info(ctx, "using filter handler")
			filter.ServeHTTP(w, req)

		})
	}
}
