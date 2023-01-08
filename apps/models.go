package apps

import "gorm.io/gorm"

type NanoConfig struct {
	gorm.Model        `json:"-"`
	GlobalEnvironment string `json:"globalEnvironment"`
	Token             string `json:"token"`
	NanoContextID     uint   `json:"-"`
}

type NanoApp struct {
	gorm.Model
	AppName           string `json:"appName" gorm:"unique"`
	AppStatus         string `json:"appStatus"`
	EnvVal            string `json:"envVal"`
	EnvMountPath      string `json:"envMountPath"`
	BuildVal          string `json:"buildVal"`
	BuildValMountPath string `json:"buildValMountPath"`
	RepoUrl           string `json:"repoUrl"`
	RepoBranch        string `json:"repoBranch"`
	NanoContextID     uint   `json:"-"`

	// ComposeRepoUrl              string `json:"composeRepoUrl"`
	// ComposeFileRelativeLocation string `json:"composeFileRelativeLocation"`
}

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

type NanoContext struct {
	gorm.Model `json:"-"`
	Apps       []NanoApp  `json:"apps"`
	NanoConfig NanoConfig `json:"nanoConfig"`
}
