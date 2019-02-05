package what

// Settings is a settings object for the crawler
type Settings struct {
	PerPage   int `json:"per_page"`
	MaxOffset int `json:"max_offset"`
}

// Project contains info to fetch builds from CircleCI
type Project struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Token string `json:"token"`
}

// Build contains search conditions and identification
type Build struct {
	Name            string                 `json:"name"`
	BuildParameters map[string]interface{} `json:"build_parameters"`
}

// Config contains projects and builds along with settings for the crawler
type Config struct {
	Settings Settings  `json:"settings"`
	Projects []Project `json:"projects"`
	Builds   []Build   `json:"builds"`
}

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
