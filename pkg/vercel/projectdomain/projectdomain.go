package projectdomain

import (
	"fmt"
	"net/http"

	"github.com/chronark/terraform-provider-vercel/pkg/vercel/api"
)

type Handler struct {
	Api *api.Api
}

// Get a Single Project Domain
// https://vercel.com/docs/api#endpoints/projects/get-a-single-project-domain
func (h *Handler) Read(projectID, teamID, domainName string) (*ProjectDomain, *api.VercelError) {
	url := fmt.Sprintf("/v8/projects/%s/domains/%s", projectID, domainName)

	if teamID != "" {
		url = fmt.Sprintf("%s?teamId=%s", url, teamID)
	}

	var domain ProjectDomain

	_, err := h.Api.Request(http.MethodGet, url, nil, &domain)

	return &domain, err
}

// Add a Domain to a Project
// https://vercel.com/docs/api#endpoints/projects/add-a-domain-to-a-project
func (h *Handler) Create(projectID, teamID string, dto CreateOrUpdateProjectDomain) (*ProjectDomain, *api.VercelError) {
	url := fmt.Sprintf("/v8/projects/%s/domains", projectID)

	if teamID != "" {
		url = fmt.Sprintf("%s?teamId=%s", url, teamID)
	}

	var domain ProjectDomain

	_, err := h.Api.Request(http.MethodPost, url, dto, &domain)

	if err != nil {
		return nil, err
	}

	return &domain, nil
}

// Update a Domain of a Project
func (h *Handler) Update(projectID, teamID, domainID string, dto CreateOrUpdateProjectDomain) (*ProjectDomain, *api.VercelError) {
	url := fmt.Sprintf("/v1/projects/%s/domains/%s", projectID, domainID)

	if teamID != "" {
		url = fmt.Sprintf("%s?teamId=%s", url, teamID)
	}

	var domain ProjectDomain

	if _, err := h.Api.Request(http.MethodPatch, url, dto, &domain); err != nil {
		return nil, err
	}

	return &domain, nil
}

// Delete a Specific Production Domain of a Project
// https://vercel.com/docs/api#endpoints/projects/delete-a-specific-production-domain-of-a-project
func (h *Handler) Delete(projectID, teamID, domainName string) *api.VercelError {
	url := fmt.Sprintf("/v8/projects/%s/domains/%s", projectID, domainName)

	if teamID != "" {
		url = fmt.Sprintf("%s?teamId=%s", url, teamID)
	}

	if _, err := h.Api.Request(http.MethodDelete, url, nil, nil); err != nil {
		return err
	}

	return nil
}
