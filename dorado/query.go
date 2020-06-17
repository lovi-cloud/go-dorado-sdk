package dorado

import (
	"net/http"
)

// SearchQuery is query struct for search function
type SearchQuery struct {
	Filter         string
	Range          string
	timeConversion TimeConversion

	AssociateObjType string
	AssociateObjID   string
	Type             string
}

// TimeConversion is type of time
type TimeConversion int

// TimeConversion const
const (
	UTC TimeConversion = iota
	LocalTime
)

// String is function compatible for fmt.Stringer
func (tc TimeConversion) String() string {
	switch tc {
	case UTC:
		return "0"
	case LocalTime:
		return "1"
	default:
		return ""
	}
}

// ToFilter convert to REST API's filter
func ToFilter(param, value string) string {
	return param + "::" + value
}

// NewSearchQueryHostname create hostname filter SearchQuery
func NewSearchQueryHostname(hostname string) *SearchQuery {
	return &SearchQuery{
		Filter: ToFilter("NAME", encodeHostName(hostname)),
	}
}

// NewSearchQueryName create name filter SearchQuery
func NewSearchQueryName(name string) *SearchQuery {
	return &SearchQuery{
		Filter: ToFilter("NAME", name),
	}
}

// NewSearchQueryID create ID filter SearchQuery
func NewSearchQueryID(id string) *SearchQuery {
	return &SearchQuery{
		Filter: ToFilter("ID", id),
	}
}

// AddSearchQuery add url parameter by SearchQuery
func AddSearchQuery(req *http.Request, query *SearchQuery) *http.Request {
	if query == nil {
		return req
	}

	q := req.URL.Query()

	if query.Filter != "" {
		q.Add("filter", query.Filter)
	}
	if query.Range != "" {
		q.Add("range", query.Range)
	}
	if query.timeConversion != UTC {
		q.Add("timeConversion", query.timeConversion.String())
	}

	if query.AssociateObjType != "" {
		q.Add("ASSOCIATEOBJTYPE", query.AssociateObjType)
	}
	if query.AssociateObjID != "" {
		q.Add("ASSOCIATEOBJID", query.AssociateObjID)
	}
	if query.Type != "" {
		q.Add("TYPE", query.Type)
	}

	req.URL.RawQuery = q.Encode()

	return req
}
