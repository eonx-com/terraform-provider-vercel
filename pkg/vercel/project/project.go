package project

import (
	"encoding/json"
	"fmt"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/httpApi"
)

type ProjectHandler struct {
	Api httpApi.API
}

func (p *ProjectHandler) Create(project CreateProject, teamId string) (string, error) {
	url := "/v6/projects"
	if teamId != "" {
		url = fmt.Sprintf("%s/?teamId=%s", url, teamId)
	}

	res, err := p.Api.Request("POST", url, project)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var createdProject Project
	err = json.NewDecoder(res.Body).Decode(&createdProject)
	if err != nil {
		return "", nil
	}

	return createdProject.ID, nil
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
