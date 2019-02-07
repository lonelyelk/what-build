package aws

func findBuild(name string) (build *Build) {
	for _, b := range RemoteConfig.Builds {
		if b.Name == name {
			return &b
		}
	}
	return nil
}

// FindBuilds looks for builds in SSM config by names
func FindBuilds(names []string) []*Build {
	var buildCfgs []*Build
	if len(names) != 0 {
		buildCfgs = make([]*Build, 0)
		for _, name := range names {
			b := findBuild(name)
			if b != nil {
				buildCfgs = append(buildCfgs, b)
			}
		}
	} else {
		buildCfgs = make([]*Build, len(RemoteConfig.Builds))
		for i, b := range RemoteConfig.Builds {
			buildCfgs[i] = &b
		}
	}
	return buildCfgs
}

func findProject(name string) (project *Project) {
	for _, p := range RemoteConfig.Projects {
		if p.Name == name {
			return &p
		}
	}
	return nil
}

// FindProjects looks for projects in SSM config by names
func FindProjects(names []string) []*Project {

	var projCfgs []*Project
	if len(names) != 0 {
		projCfgs = make([]*Project, 0)
		for _, name := range names {
			p := findProject(name)
			if p != nil {
				projCfgs = append(projCfgs, p)
			}
		}
	} else {
		projCfgs = make([]*Project, len(RemoteConfig.Projects))
		for i, p := range RemoteConfig.Projects {
			projCfgs[i] = &p
		}
	}
	return projCfgs
}
