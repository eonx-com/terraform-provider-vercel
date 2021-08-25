package pdomain

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chronark/terraform-provider-vercel/pkg/vercel/httpApi"
)

type Handler struct {
	Api httpApi.API
}

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

func (h *Handler) Read(projectID, teamID, domainName string) (ProjectDomain, error) {
	url := fmt.Sprintf("/v8/projects/%s/domains/%s", projectID, domainName)

	if teamID != "" {
		url = fmt.Sprintf("%s?teamId=%s", url, teamID)
	}

	res, err := h.Api.Request("GET", url, nil)
	if err != nil {
		return ProjectDomain{}, fmt.Errorf("unable to fetch project domain from vercel: %w", err)
	}
	defer res.Body.Close()

	var domain ProjectDomain

	err = json.NewDecoder(res.Body).Decode(&domain)
	if err != nil {
		return domain, fmt.Errorf("unable to unmarshal domain response: %w", err)
	}

	return domain, nil
}

func (h *Handler) Create(projectID, teamID string, dto CreateOrUpdateProjectDomain) (*ProjectDomain, error) {
	url := fmt.Sprintf("/v8/projects/%s/domains", projectID)

	if teamID != "" {
		url = fmt.Sprintf("%s?teamId=%s", url, teamID)
	}

	res, err := h.Api.Request(http.MethodPost, url, dto)

	if err != nil {
		return nil, fmt.Errorf("unable to create project domain from vercel: %w", err)
	}

	defer res.Body.Close()

	var domain ProjectDomain

	err = json.NewDecoder(res.Body).Decode(&domain)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal project domain response: %w", err)
	}

	return &domain, nil
}

func (h *Handler) Update(projectID, teamID, domainID string, dto CreateOrUpdateProjectDomain) (*ProjectDomain, error) {
	url := fmt.Sprintf("/v1/projects/%s/domains/%s", projectID, domainID)

	if teamID != "" {
		url = fmt.Sprintf("%s?teamId=%s", url, teamID)
	}

	res, err := h.Api.Request(http.MethodPatch, url, dto)

	if err != nil {
		return nil, fmt.Errorf("unable to update project domain from vercel: %w", err)
	}

	defer res.Body.Close()

	var domain ProjectDomain

	err = json.NewDecoder(res.Body).Decode(&domain)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal project domain response: %w", err)
	}

	return &domain, nil
}

func (h *Handler) Delete(projectID, teamID, domainName string) error {
	url := fmt.Sprintf("/v8/projects/%s/domains/%s", projectID, domainName)

	if teamID != "" {
		url = fmt.Sprintf("%s?teamId=%s", url, teamID)
	}

	res, err := h.Api.Request(http.MethodDelete, url, nil)

	if err != nil {
		return fmt.Errorf("unable to delete project domain from vercel: %w", err)
	}

	defer res.Body.Close()

	return nil
}
