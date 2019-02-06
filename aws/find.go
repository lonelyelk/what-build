package aws

import "errors"

// FindBuild looks for a build in SSM config by name
func FindBuild(name string) (build *Build, err error) {
	for _, b := range RemoteConfig.Builds {
		if b.Name == name {
			build = &b
			break
		}
	}
	if build == nil {
		err = errors.New("Build config not found")
	}
	return
}

// FindProject looks for a project in SSM config by name
func FindProject(name string) (project *Project, err error) {
	for _, p := range RemoteConfig.Projects {
		if p.Name == name {
			project = &p
			break
		}
	}
	if project == nil {
		err = errors.New("Project config not found")
	}
	return
}
