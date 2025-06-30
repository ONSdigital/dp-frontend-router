//nolint:revive,stylecheck // ignore var-naming Package name "datasetType" is kept for compatibility.
package datasetType

import (
	"context"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
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
