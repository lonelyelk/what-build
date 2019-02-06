package what

// CIBuildResponse is a JSON extraits for build entity on circleci
type CIBuildResponse struct {
	BuildNum        int                    `json:"build_num"`
	Branch          string                 `json:"branch"`
	VcsRevision     string                 `json:"vcs_revision"`
	Subject         string                 `json:"subject"`
	Why             string                 `json:"why"`
	DontBuild       string                 `json:"dont_build"`
	StopTime        string                 `json:"stop_time"`
	BuildTimeMillis int                    `json:"build_time_millis"`
	Status          string                 `json:"status"`
	BuildParameters map[string]interface{} `json:"build_parameters"`
}
