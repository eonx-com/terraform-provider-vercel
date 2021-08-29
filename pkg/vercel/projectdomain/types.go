package projectdomain

type CreateOrUpdateProjectDomain struct {
	Name               string  `json:"name"`
	Redirect           *string `json:"redirect"`
	RedirectStatusCode *int    `json:"redirectStatusCode"`
	GitBranch          *string `json:"gitBranch"`
}

type ProjectDomain struct {
	Name               string `json:"name"`
	GitBranch          string `json:"gitBranch"`
	Redirect           string `json:"redirect"`
	RedirectStatusCode int    `json:"redirectStatusCode"`
	ProjectID          string `json:"projectId"`
	CreatedAt          int64  `json:"createdAt"`
	UpdatedAt          int64  `json:"updatedAt"`
}
