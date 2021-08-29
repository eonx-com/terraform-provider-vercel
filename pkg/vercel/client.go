package vercel

import (
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/api"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/env"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/httpApi"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/projectdomain"

	"github.com/chronark/terraform-provider-vercel/pkg/vercel/alias"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/dns"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/domain"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/project"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/secret"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/team"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/user"
)

type Client struct {
	Project       *project.ProjectHandler
	User          *user.UserHandler
	Env           *env.Handler
	Secret        *secret.Handler
	Team          *team.Handler
	Alias         *alias.Handler
	Domain        *domain.Handler
	ProjectDomain *projectdomain.Handler
	DNS           *dns.Handler
}

func New(token string) *Client {
	apiv1 := httpApi.New(token)
	apiv2 := api.New(token)

	return &Client{
		Project: &project.ProjectHandler{
			Api: apiv1,
		},
		User: &user.UserHandler{
			Api: apiv1,
		},
		Env:           &env.Handler{Api: apiv2},
		Secret:        &secret.Handler{Api: apiv1},
		Team:          &team.Handler{Api: apiv1},
		Alias:         &alias.Handler{Api: apiv1},
		Domain:        &domain.Handler{Api: apiv1},
		ProjectDomain: &projectdomain.Handler{Api: apiv2},
		DNS:           &dns.Handler{Api: apiv1},
	}
}
