package env

import (
	"fmt"
	"net/http"

	"github.com/chronark/terraform-provider-vercel/pkg/vercel/api"
)

type Handler struct {
	Api *api.Api
}

// Create a Project Environment Variable
// https://vercel.com/docs/api#endpoints/projects/create-a-project-environment-variable
func (h *Handler) Create(projectID string, payload CreateOrUpdateEnv, teamId string) (*Env, *api.VercelError) {
	url := fmt.Sprintf("/v6/projects/%s/env", projectID)

	if teamId != "" {
		url = fmt.Sprintf("%s/?teamId=%s", url, teamId)
	}

	var env Env

	if _, err := h.Api.Request(http.MethodPost, url, payload, &env); err != nil {
		return nil, err
	}

	return &env, nil
}

// Get Project Environment Variables
// https://vercel.com/docs/api#endpoints/projects/get-project-environment-variables
func (h *Handler) Read(projectID, teamID, envID string) (*Env, *api.VercelError) {
	url := fmt.Sprintf("/v6/projects/%s/env", projectID)

	if teamID != "" {
		url = fmt.Sprintf("%s/?teamId=%s", url, teamID)
	}

	var envResponse ReadEnvResponse

	if _, err := h.Api.Request(http.MethodGet, url, nil, &envResponse); err != nil {
		return nil, err
	}

	var env *Env
	for _, e := range envResponse.Envs {
		if e.ID == envID {
			env = &e
			break
		}
	}

	if env == nil {
		return nil, &api.VercelError{
			Code: api.ErrCodeNotFound,
		}
	}

	return env, nil
}

// Edit a Project Environment Variable
// https://vercel.com/docs/api#endpoints/projects/edit-a-project-environment-variable
func (h *Handler) Update(projectID string, envID string, payload CreateOrUpdateEnv, teamId string) (*Env, *api.VercelError) {
	url := fmt.Sprintf("/v6/projects/%s/env/%s", projectID, envID)

	if teamId != "" {
		url = fmt.Sprintf("%s/?teamId=%s", url, teamId)
	}

	var env Env

	if _, err := h.Api.Request(http.MethodPatch, url, payload, &env); err != nil {
		return nil, err
	}

	return &env, nil
}

// Delete a Specific Environment Variable
// https://vercel.com/docs/api#endpoints/projects/delete-a-specific-environment-variable
func (h *Handler) Delete(projectID, envKey string, teamId string) *api.VercelError {
	url := fmt.Sprintf("/v8/projects/%s/env/%s", projectID, envKey)

	if teamId != "" {
		url = fmt.Sprintf("%s/?teamId=%s", url, teamId)
	}

	if _, err := h.Api.Request(http.MethodDelete, url, nil, nil); err != nil {
		return err
	}

	return nil
}
