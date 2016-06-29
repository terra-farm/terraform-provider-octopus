package octopus

import (
	"net/url"
	"strconv"
)

// EntitySummary represents summary information for an Octopus entity.
type EntitySummary struct {
	ID   string `json:"Id"`
	Name string `json:"Name"`
}

// PagedResults represents a page of results from the Octopus API.
//
// This structure is designed for its fields to be embedded in paged result structures for specific result types.
type PagedResults struct {
	ItemType     string            `json:"ItemType"`
	TotalResults int               `json:"TotalResults"`
	ItemsPerPage int               `json:"ItemsPerPage"`
	IsStale      bool              `json:"IsStale"`
	Links        map[string]string `json:"Links"`
}

// GetSkipForNextPage gets the number of items to skip for the next page of results.
//
// Returns false if there is no next page.
func (pagedResults *PagedResults) GetSkipForNextPage() (skip int, ok bool) {
	return pagedResults.skipForPageLink("Page.Next")
}

// GetSkipForPreviousPage gets the number of items to skip for the previous page of results.
//
// Returns false if there is no previous page.
func (pagedResults *PagedResults) GetSkipForPreviousPage() (skip int, ok bool) {
	return pagedResults.skipForPageLink("Page.Previous")
}

func (pagedResults *PagedResults) skipForPageLink(pageLinkName string) (skip int, ok bool) {
	var link string
	link, ok = pagedResults.Links[pageLinkName]
	if !ok {
		return
	}

	linkURL, err := url.Parse(link)
	if err != nil {
		ok = false

		return
	}

	skipParameter := linkURL.Query().Get("skip")
	if len(skipParameter) == 0 {
		ok = false

		return
	}

	skip, err = strconv.Atoi(skipParameter)
	ok = err != nil

	return
}
