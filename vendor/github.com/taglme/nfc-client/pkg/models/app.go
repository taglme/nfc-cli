package models

type AppInfo struct {
	Name           string `json:"name"`
	Version        string `json:"version"`
	Commit         string `json:"commit"`
	SDKInfo        string `json:"sdk_info"`
	Platform       string `json:"platform"`
	BuildTime      string `json:"build_time"`
	CheckSuccess   bool   `json:"check_success"`
	Supported      bool   `json:"suported"`
	HaveUpdate     bool   `json:"have_update"`
	UpdateVersion  string `json:"update_version"`
	UpdateDownload string `json:"update_download"`
	StartedAt      string `json:"started_at"`
}
