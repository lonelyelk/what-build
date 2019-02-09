package aws

func appendBuild(name string, builds []Build, selection *[]Build) {
	for _, b := range builds {
		if b.Name == name {
			*selection = append(*selection, b)
			return
		}
	}
}

// FindBuilds looks for builds in SSM config by names
func FindBuilds(names []string) []Build {
	config := GetRemoteConfig()
	var buildCfgs []Build
	if len(names) != 0 {
		buildCfgs = make([]Build, 0)
		for _, name := range names {
			appendBuild(name, config.Builds, &buildCfgs)
		}
	} else {
		buildCfgs = make([]Build, len(config.Builds))
		for i, b := range config.Builds {
			buildCfgs[i] = b
		}
	}
	return buildCfgs
}

func appendProject(name string, projects []Project, selection *[]Project) {
	for _, p := range projects {
		if p.Name == name {
			*selection = append(*selection, p)
			return
		}
	}
}

// FindProjects looks for projects in SSM config by names
func FindProjects(names []string) []Project {
	config := GetRemoteConfig()
	var projCfgs []Project
	if len(names) != 0 {
		projCfgs = make([]Project, 0)
		for _, name := range names {
			appendProject(name, config.Projects, &projCfgs)
		}
	} else {
		projCfgs = make([]Project, len(config.Projects))
		for i, p := range config.Projects {
			projCfgs[i] = p
		}
	}
	return projCfgs
}
