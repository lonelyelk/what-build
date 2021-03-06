package aws

func appendBuild(selection *[]Build, name string, builds *[]Build) {
	for _, b := range *builds {
		if b.Name == name {
			*selection = append(*selection, b)
			return
		}
	}
}

// FindBuild looks for a build by name in a provided slice
func FindBuild(name string, builds *[]Build) *Build {
	for _, b := range *builds {
		if b.Name == name {
			return &b
		}
	}
	return nil
}

// FindBuilds looks for builds in SSM config by names
func FindBuilds(names []string, builds *[]Build) []Build {
	var buildCfgs []Build
	if len(names) != 0 {
		buildCfgs = make([]Build, 0)
		for _, name := range names {
			appendBuild(&buildCfgs, name, builds)
		}
	} else {
		buildCfgs = make([]Build, len(*builds))
		for i, b := range *builds {
			buildCfgs[i] = b
		}
	}
	return buildCfgs
}

func appendProject(selection *[]Project, name string, projects *[]Project) {
	for _, p := range *projects {
		if p.Name == name {
			*selection = append(*selection, p)
			return
		}
	}
}

// FindProject looks for a project by name in the provided slice
func FindProject(name string, projects *[]Project) *Project {
	for _, p := range *projects {
		if p.Name == name {
			return &p
		}
	}
	return nil
}

// FindProjects looks for projects in SSM config by names
func FindProjects(names []string, projects *[]Project) []Project {
	var projCfgs []Project
	if len(names) != 0 {
		projCfgs = make([]Project, 0)
		for _, name := range names {
			appendProject(&projCfgs, name, projects)
		}
	} else {
		projCfgs = make([]Project, len(*projects))
		for i, p := range *projects {
			projCfgs[i] = p
		}
	}
	return projCfgs
}
