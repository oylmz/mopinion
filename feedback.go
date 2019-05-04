package mopinion

import (
	"context"
	"fmt"
	"strings"
)

// FeedbackInterface holds two methods that return feedback by a given dataset or report.
// FeedbackInterface accepts pagination options and filters.
type FeedbackInterface interface {
	GetByDataset(ctx context.Context, datasetID int, options *PaginationOptions, filters *FilterCollection) (*Feedback, *Response, error)
	GetByReport(ctx context.Context, reportID int, options *PaginationOptions, filters *FilterCollection) (*Feedback, *Response, error)
}

// FeedbackService implements FeedbackInterface.
type FeedbackService struct {
	service
}

func (s *FeedbackService) get(ctx context.Context, url string) (*Feedback, *Response, error) {
	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	feedback := new(Feedback)
	resp, err := s.client.Do(ctx, req, feedback)
	if err != nil {
		return nil, resp, err
	}

	return feedback, resp, nil
}

// FilterCollection holds the filters to be applied on feedback data.
type FilterCollection struct {
	Filters []Filter
}

// Filter is a struct composed of key, modifier and its value.
// For more information, look at the link below.
// https://developer.mopinion.com/api/#section/Requests-and-Responses/Filters
type Filter struct {
	Key      FilterKey
	Modifier FilterModifier
	Value    string
}

// Build returns the string to be appended to the URL.
func (f Filter) Build() string {
	return fmt.Sprintf("filter[%s%s]=%s", f.Modifier, f.Key, f.Value)
}

// FilterKey type represents the possible filters
type FilterKey string

const (
	// Date refers to the time when a record is created.
	Date FilterKey = "date"
	// Rating is a geenral numeric rating, its value should be numeric.
	Rating FilterKey = "rating"
	// Nps is Net Promotor Score, its value should be between 0 and 10.
	Nps FilterKey = "nps"
	// Ces is Customer Effort Score, its value should be between 1 and 5.
	Ces FilterKey = "ces"
	// CesInverse is Customer Effort Score
	CesInverse FilterKey = "ces_inverse"
	// Gcr is Goal Completion Rate. Option are no, partly, yes.
	Gcr FilterKey = "gcr"
	// Tags can be assigned to the feedback items.
	Tags FilterKey = "tags"
)

// FilterModifier type represents the possible modifiers
type FilterModifier string

const (
	// Not is used for logical not
	Not FilterModifier = "!"
	// Lt means less than
	Lt FilterModifier = "<"
	// Lte means less than or equal
	Lte FilterModifier = "<<"
	// Gt means greater than
	Gt FilterModifier = ">"
	// Gte means greater than or equal
	Gte FilterModifier = ">>"
)

func addFilters(s string, filters *FilterCollection) string {
	var pair []string
	if filters == nil {
		return s
	}
	for _, f := range filters.Filters {
		pair = append(pair, f.Build())
	}
	if strings.Contains(s, "?") {
		return fmt.Sprintf("%s&%s", s, strings.Join(pair, "&"))
	}
	return fmt.Sprintf("%s?%s", s, strings.Join(pair, "&"))
}

// GetByDataset returns feedback by given dataset id, pagination options and filters.
func (s *FeedbackService) GetByDataset(ctx context.Context, datasetID int, options *PaginationOptions, filters *FilterCollection) (*Feedback, *Response, error) {
	u := fmt.Sprintf("datasets/%d/feedback", datasetID)
	u, err := addPaginationOptions(u, options)
	if err != nil {
		return nil, nil, err
	}
	u = addFilters(u, filters)
	return s.get(ctx, u)
}

// GetByReport returns feedback by given report id, pagination options and filters.
func (s *FeedbackService) GetByReport(ctx context.Context, reportID int, options *PaginationOptions, filters *FilterCollection) (*Feedback, *Response, error) {
	u := fmt.Sprintf("reports/%d/feedback", reportID)
	u, err := addPaginationOptions(u, options)
	if err != nil {
		return nil, nil, err
	}
	u = addFilters(u, filters)
	return s.get(ctx, u)
}
