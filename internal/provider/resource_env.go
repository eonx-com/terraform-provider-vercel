package provider

import (
	"context"

	"github.com/chronark/terraform-provider-vercel/pkg/vercel"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/api"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/env"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEnv() *schema.Resource {
	return &schema.Resource{
		Description: "https://vercel.com/docs/api#endpoints/projects/get-project-environment-variables",

		CreateContext: resourceEnvCreate,
		ReadContext:   resourceEnvRead,
		UpdateContext: resourceEnvUpdate,
		DeleteContext: resourceEnvDelete,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Description: "The unique project identifier.",
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
			"type": {
				Description: "The type can be `plain`, `secret`, or `system`.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"id": {
				Description: "Unique id for this variable.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"key": {
				Description: "The name of the environment variable.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"value": {
				Description: "If the type is `plain`, a string representing the value of the environment variable. If the type is `secret`, the secret ID of the secret attached to the environment variable. If the type is `system`, the name of the System Environment Variable.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"target": {
				Description: "The target can be a list of `development`, `preview`, or `production`.",
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    3,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"git_branch": {
				Description: "The Git branch for this variable, only accepted when the target is exclusively preview.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"created_at": {
				Description: "A number containing the date when the variable was created in milliseconds.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"updated_at": {
				Description: "A number containing the date when the variable was updated in milliseconds.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

func toCreateOrUpdateEnv(d *schema.ResourceData) env.CreateOrUpdateEnv {
	dto := env.CreateOrUpdateEnv{
		Type:  d.Get("type").(string),
		Key:   d.Get("key").(string),
		Value: d.Get("value").(string),
	}

	// Casting each target because go does not allow typecasting from interface{} to []string
	targetList := d.Get("target").([]interface{})
	dto.Target = make([]string, len(targetList))

	for i := 0; i < len(dto.Target); i++ {
		dto.Target[i] = targetList[i].(string)
	}

	if b := d.Get("git_branch").(string); b != "" {
		dto.GitBranch = &b
	}

	return dto
}

func resourceEnvCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*vercel.Client)

	payload := toCreateOrUpdateEnv(d)

	env, err := client.Env.Create(d.Get("project_id").(string), payload, d.Get("team_id").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(env.ID)

	return resourceEnvRead(ctx, d, meta)
}

func resourceEnvRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*vercel.Client)

	env, err := client.Env.Read(d.Get("project_id").(string), d.Get("team_id").(string), d.Id())

	if err != nil {
		if err.Is(api.ErrCodeNotFound) {
			d.SetId("")
			return diag.Diagnostics{}
		}

		return diag.FromErr(err)
	}

	if err := d.Set("type", env.Type); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("key", env.Key); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("value", env.Value); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("target", env.Target); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("git_branch", env.GitBranch); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("updated_at", env.UpdatedAt); err != nil {
		return diag.FromErr(err)

	}

	if err := d.Set("created_at", env.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func resourceEnvUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*vercel.Client)

	if d.HasChanges("type", "key", "value", "target", "git_branch") {
		payload := toCreateOrUpdateEnv(d)

		if _, err := client.Env.Update(
			d.Get("project_id").(string),
			d.Id(),
			payload,
			d.Get("team_id").(string),
		); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceEnvRead(ctx, d, meta)
}

func resourceEnvDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*vercel.Client)

	if err := client.Env.Delete(
		d.Get("project_id").(string),
		d.Get("id").(string),
		d.Get("team_id").(string),
	); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diag.Diagnostics{}
}
