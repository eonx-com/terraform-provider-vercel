package env

type Env struct {
	Type            string      `json:"type"`
	ID              string      `json:"id"`
	Key             string      `json:"key"`
	Value           string      `json:"value"`
	Target          []string    `json:"target"`
	GitBranch       string      `json:"gitBranch"`
	ConfigurationID interface{} `json:"configurationId"`
	UpdatedAt       int64       `json:"updatedAt"`
	CreatedAt       int64       `json:"createdAt"`
}

type ReadEnvResponse struct {
	Envs []Env `json:"envs"`
}

type CreateOrUpdateEnv struct {
	// The type can be `plain`, `secret`, or `system`.
	Type string `json:"type"`

	// The name of the environment variable.
	Key string `json:"key"`

	// If the type is `plain`, a string representing the value of the environment variable.
	// If the type is `secret`, the secret ID of the secret attached to the environment variable.
	// If the type is `system`, the name of the System Environment Variable.
	Value string `json:"value"`

	// 	The target can be a list of `development`, `preview`, or `production`.
	Target []string `json:"target"`

	// The Git branch for this variable, only accepted when the target is exclusively preview.
	GitBranch *string `json:"gitBranch"`
}
