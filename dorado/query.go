package dorado

import "net/http"

type SearchQuery struct {
	Filter         string
	Range          string
	timeConversion TimeConversion
}

type TimeConversion int

const (
	UTC TimeConversion = iota
	LocalTime
)

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

func ToFilter(param, value string) string {
	return param + "::" + value
}

func CreateSearchName(name string) *SearchQuery {
	return &SearchQuery{
		Filter: ToFilter("NAME", name),
	}
}

func CreateSearchId(id string) *SearchQuery {
	return &SearchQuery{
		Filter: ToFilter("ID", id),
	}
}

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

	req.URL.RawQuery = q.Encode()

	return req
}
