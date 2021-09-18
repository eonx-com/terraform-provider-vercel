package project

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chronark/terraform-provider-vercel/pkg/vercel/api"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/httpApi"
)

type ProjectHandler struct {
	Api   httpApi.API
	Apiv2 *api.Api
}

func (p *ProjectHandler) Create(project CreateProject, teamId string) (*Project, error) {
	url := "/v6/projects"

	if teamId != "" {
		url = fmt.Sprintf("%s/?teamId=%s", url, teamId)
	}

	res, err := p.Api.Request("POST", url, project)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var createdProject Project

	if err := json.NewDecoder(res.Body).Decode(&createdProject); err != nil {
		return nil, nil
	}

	return &createdProject, nil
}

func (p *ProjectHandler) Read(id string, teamId string) (project Project, err error) {
	url := fmt.Sprintf("/v1/projects/%s", id)

	if teamId != "" {
		url = fmt.Sprintf("%s/?teamId=%s", url, teamId)
	}

	res, err := p.Api.Request("GET", url, nil)

	if err != nil {
		return Project{}, fmt.Errorf("Unable to fetch project from vercel: %w", err)
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&project)

	if err != nil {
		return Project{}, fmt.Errorf("Unable to unmarshal project: %w", err)
	}

	return project, nil
}

func (p *ProjectHandler) Update(id string, project UpdateProject, teamId string) error {
	branch := Branch{
		Branch: project.Branch,
	}
	projectInternal := UpdateProjectInternal{
		Framework:                project.Framework,
		PublicSource:             project.PublicSource,
		InstallCommand:           project.InstallCommand,
		BuildCommand:             project.BuildCommand,
		DevCommand:               project.DevCommand,
		OutputDirectory:          project.OutputDirectory,
		ServerlessFunctionRegion: project.ServerlessFunctionRegion,
		RootDirectory:            project.RootDirectory,
		Name:                     project.Name,
		NodeVersion:              project.NodeVersion,
	}

	url := fmt.Sprintf("/v2/projects/%s", id)
	if teamId != "" {
		url = fmt.Sprintf("%s/?teamId=%s", url, teamId)
	}

	res, err := p.Api.Request("PATCH", url, projectInternal)
	if err != nil {
		return fmt.Errorf("Unable to update project: %w", err)
	}
	defer res.Body.Close()

	if branch.Branch != "" {
		branchUrl := fmt.Sprintf("/v4/projects/%s/branch", id)
		if teamId != "" {
			branchUrl = fmt.Sprintf("%s/?teamId=%s", branchUrl, teamId)
		}

		resBranch, errBranch := p.Api.Request("PATCH", branchUrl, branch)
		if errBranch != nil {
			return fmt.Errorf("Unable to update project branch: %w", errBranch)
		}
		defer resBranch.Body.Close()
	}

	return nil
}

func (p *ProjectHandler) UpdateProductionBranch(projectID, teamID, branch string) *api.VercelError {
	url := fmt.Sprintf("/v4/projects/%s/branch", projectID)

	if teamID != "" {
		url = fmt.Sprintf("%s?teamId=%s", url, teamID)
	}

	_, err := p.Apiv2.Request(http.MethodPatch, url, map[string]string{
		"branch": branch,
	}, nil)

	return err
}

func (p *ProjectHandler) Delete(id string, teamId string) error {
	url := fmt.Sprintf("/v1/projects/%s", id)
	if teamId != "" {
		url = fmt.Sprintf("%s/?teamId=%s", url, teamId)
	}

	res, err := p.Api.Request("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("Unable to delete project: %w", err)
	}
	defer res.Body.Close()
	return nil
}
