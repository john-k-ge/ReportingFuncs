package entity

type OrgResponse struct {
	Total_Results int
	Total_Pages   int
	Prev_Url      string
	Next_Url      string
	Resources     []OrgResource
}

type OrgResource struct {
	Metadata CFMetadata
	Entity   OrgEntity
}

type CFMetadata struct {
	Guid       string
	Url        string
	Created_at string
	Updated_at string
}

type OrgEntity struct {
	Name                  string `json:"name"`
	Status                string `json:"status"`
	Quota_definition_url  string `json:"quota_definition_url"`
	Spaces_url            string `json:"spaces_url"`
	Users_url             string `json:"users_url"`
	Managers_url          string `json:"managers_url"`
	Quota_definition_guid string `json:"quota_definition_guid"`
	//Auditors_url         string `json:"auditors_url"`
	//Billing_enabled             bool
	//Domains_url                 string
	//Private_domains_url         string
	//Billing_managers_url        string
	//App_events_url              string
	//Space_quota_definitions_url string
}

type OrgInfo struct {
	Name         string
	Guid         string
	Status       string
	Mem_util     int
	Mem_quota    int
	Quota_name   string
	Quota_url    string
	Quota_guid   string
	Spaces_url   string
	Managers_url string
	Users_url    string
	Managers     map[string]string
	Users        map[string]string
	Services     []string
	Created      string
}

type OrgMemResponse struct {
	Memory_usage_in_mb int
}

type OrgQuota struct {
	Metadata CFMetadata
	Entity   OrgQuotaEntity
}

type OrgQuotaMetadata struct {
	Name string
	Guid string
}

type OrgQuotaEntity struct {
	Memory_limit int
	Name         string
}

type OrgServiceInstance struct {
	Metadata CFMetadata
	Entity   OrgServiceEntity
}

type OrgServiceEntity struct {
	Name             string
	Service_plan_url string
	PlanName         string
	Service_url      string
	ServiceName      string
}

type QuotaResponse struct {
	Total_Results int
	Total_Pages   int
	Prev_Url      string
	Next_Url      string
	Resources     []QuotaResource
}

type QuotaResource struct {
	Metadata CFMetadata
	Entity   QuotaEntity
}

type QuotaEntity struct {
	Name string
}
