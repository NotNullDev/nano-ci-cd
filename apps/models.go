package apps

import "gorm.io/gorm"

type BuildArguments struct {
	RepoUrl string `json:"repoUrl"`
	AppName string `json:"appName"`
	// DockerfilePath           string   `json:"dockerfilePath"`
	// PreDeployScriptLocation  string   `json:"preDeployScriptLocation"`
	// DeployScriptLocation     string   `json:"deployScriptLocation"`
	// PostDeployScriptLocation string   `json:"postDeployScriptLocation"`
	// EnvFileNames             []string `json:"envFileNames"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type NanoApp struct {
	gorm.Model
	RepoUrl                     string `json:"repoUrl"`
	AppName                     string `json:"appName"`
	EnvFile                     string `json:"envFile"`
	EnvFileLocation             string `json:"envFileLocation"`
	ComposeRepoUrl              string `json:"composeRepoUrl"`
	ComposeFileRelativeLocation string `json:"composeFileRelativeLocation"`
}
