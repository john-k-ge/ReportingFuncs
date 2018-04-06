package urls

import "strconv"

const (
	resultsPerPage    = 50
	QuotaNamePath     = "/v2/quota_definitions?q=name:%v"
	QuotaGuidPath     = "/v2/quota_definitions/%v"
	OrgMemUtil        = "/v2/organizations/%v/memory_usage"
	searchServicePath = "/v2/services?q=label:%v"
	allPlansPath      = "/v2/services/%v/service_plans?order-direction=asc&page=%v&results-per-page=50"
	genericPagingPath = "?order-direction=asc&page=%v&results-per-page=50"
	spacePath         = "/v2/spaces/%v"
	orgPath           = "/v2/organizations/%v"
	instancePath      = "/v2/service_instances/%v?accepts_incomplete=true"
)

var (
	OrgApiPath = "/v2/organizations?order-direction=asc&page=%d&results-per-page=" + strconv.Itoa(resultsPerPage)
)
