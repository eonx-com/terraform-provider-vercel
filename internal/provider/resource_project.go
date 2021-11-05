package provider

import (
	"context"

	"github.com/chronark/terraform-provider-vercel/pkg/vercel"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/project"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Description: "https://vercel.com/docs/api#endpoints/projects",

		CreateContext: resourceProjectCreate,
		ReadContext:   resourceProjectRead,
		UpdateContext: resourceProjectUpdate,
		DeleteContext: resourceProjectDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Internal id of this project",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The name of the project.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"team_id": {
				Description: "By default, you can access resources contained within your own user account. To access resources owned by a team, you can pass in the team ID",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "",
			},
			"git_repository": {
				Description: "The git repository that will be connected to the project. Any pushes to the specified connected git repository will be automatically deployed.",
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description: "The git provider of the repository. Must be either `github`, `gitlab`, or `bitbucket`.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"repo": {
							Description: "The name of the git repository. For example: `chronark/terraform-provider-vercel`",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
			"branch": {
				Description: "By default, every commit pushed to the main branch will trigger a Production Deployment instead of the usual Preview Deployment.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"account_id": {
				Description: "The unique ID of the user or team the project belongs to.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"created_at": {
				Description: "A number containing the date when the project was created in milliseconds.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"updated_at": {
				Description: "A number containing the date when the project was updated in milliseconds.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"framework": {
				Description: "The framework that is being used for this project. When null is used no framework is selected.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"public_source": {
				Description: " Specifies whether the source code and logs of the deployments for this project should be public or not.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"install_command": {
				Description: "The install command for this project. When null is used this value will be automatically detected.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"build_command": {
				Description: "The build command for this project. When null is used this value will be automatically detected.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"dev_command": {
				Description: "The dev command for this project. When null is used this value will be automatically detected.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"output_directory": {
				Description: "The output directory of the project. When null is used this value will be automatically detected.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
			},
			"serverless_function_region": {
				Description: "The region to deploy Serverless Functions in this project.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"root_directory": {
				Description: "The name of a directory or relative path to the source code of your project. When null is used it will default to the project root.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
			},
			"node_version": {
				Description: "The Node.js Version for this project.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"alias": {
				Description: "A list of production domains for the project.",
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*vercel.Client)
	// Terraform does not have nested objects with different types yet, so I am using a `TypeList`
	// Here we have to typecast to list first and then take the first item and cast again.
	repo := d.Get("git_repository").([]interface{})[0].(map[string]interface{})

	createProject := project.CreateProject{
		Name: d.Get("name").(string),
		GitRepository: struct {
			Type string `json:"type"`
			Repo string `json:"repo"`
		}{
			Type: repo["type"].(string),
			Repo: repo["repo"].(string),
		},
	}

	if framework, frameworkSet := d.GetOk("framework"); frameworkSet {
		createProject.Framework = framework.(string)
	}

	if publicSource, publicSourceSet := d.GetOk("public_source"); publicSourceSet {
		createProject.PublicSource = publicSource.(bool)
	}

	if installCommand, installCommandSet := d.GetOk("install_command"); installCommandSet {
		createProject.InstallCommand = installCommand.(string)
	}

	if buildCommand, buildCommandSet := d.GetOk("build_command"); buildCommandSet {
		createProject.BuildCommand = buildCommand.(string)
	}

	if devCommand, devCommandSet := d.GetOk("dev_command"); devCommandSet {
		createProject.DevCommand = devCommand.(string)
	}

	if outputDirectory, outputDirectorySet := d.GetOk("output_directory"); outputDirectorySet {
		createProject.OutputDirectory = outputDirectory.(string)
	}

	if serverlessFunctionRegion, serverlessFunctionRegionSet := d.GetOk("serverless_function_region"); serverlessFunctionRegionSet {
		createProject.ServerlessFunctionRegion = serverlessFunctionRegion.(string)
	}

	if rootDirectory, rootDirectorySet := d.GetOk("root_directory"); rootDirectorySet {
		createProject.RootDirectory = rootDirectory.(string)
	}

	if nodeVersion, nodeVersionSet := d.GetOk("node_version"); nodeVersionSet {
		createProject.NodeVersion = nodeVersion.(string)
	}

	branch, branchSet := d.GetOk("branch")

	teamID := d.Get("team_id").(string)

	project, err := client.Project.Create(createProject, teamID)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(project.ID)

	if branchSet && branch != "" && project.Link.ProductionBranch != branch {
		client.Project.UpdateProductionBranch(project.ID, teamID, branch.(string))
	}

	return resourceProjectRead(ctx, d, meta)
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*vercel.Client)

	id := d.Id()

	project, err := client.Project.Read(id, d.Get("team_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("name", project.Name); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("account_id", project.AccountID); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("branch", project.Link.ProductionBranch); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("created_at", project.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("updated_at", project.UpdatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("framework", project.Framework); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("public_source", project.PublicSource); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("install_command", project.InstallCommand); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("build_command", project.BuildCommand); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("dev_command", project.DevCommand); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("output_directory", project.OutputDirectory); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("serverless_function_region", project.ServerlessFunctionRegion); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("root_directory", project.RootDirectory); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("node_version", project.NodeVersion); err != nil {
		return diag.FromErr(err)
	}

	aliases := make([]string, 0)
	for i := 0; i < len(project.Aliases); i++ {
		aliases = append(aliases, project.Aliases[i].Domain)
	}

	if err = d.Set("alias", aliases); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*vercel.Client)
	var update project.UpdateProject

	if d.HasChange("name") {
		update.Name = d.Get("name").(string)
	}

	if d.HasChange("framework") {
		update.Framework = d.Get("framework").(string)
	}

	if d.HasChange("public_source") {
		update.PublicSource = d.Get("public_source").(bool)
	}

	if d.HasChange("install_command") {
		update.InstallCommand = d.Get("install_command").(string)
	}

	if d.HasChange("build_command") {
		update.BuildCommand = d.Get("build_command").(string)
	}

	if d.HasChange("dev_command") {
		update.DevCommand = d.Get("dev_command").(string)
	}

	if d.HasChange("output_directory") {
		update.OutputDirectory = d.Get("output_directory").(string)
	}

	if d.HasChange("serverless_function_region") {
		update.ServerlessFunctionRegion = d.Get("serverless_function_region").(string)
	}

	if d.HasChange("root_directory") {
		update.RootDirectory = d.Get("root_directory").(string)
	}

	if d.HasChange("node_version") {
		update.NodeVersion = d.Get("node_version").(string)
	}

	if d.HasChange("branch") {
		update.Branch = d.Get("branch").(string)
	}

	err := client.Project.Update(d.Id(), update, d.Get("team_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceProjectRead(ctx, d, meta)
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*vercel.Client)
	err := client.Project.Delete(d.Id(), d.Get("team_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}
