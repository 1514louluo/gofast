package response

import "GF_PROJECT_NAME/config"

type SysConfigResponse struct {
	Config config.Server `json:"config"`
}
