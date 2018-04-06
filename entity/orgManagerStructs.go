package entity

type OrgUser struct {
	Metadata CFMetadata
	Entity   OrgUserEntity
}

type OrgUserEntity struct {
	Username string
	Admin    bool
}

type OrgUserResponse struct {
	Total_Results int
	Total_Pages   int
	Prev_Url      string
	Next_Url      string
	Resources     []OrgUser
}
