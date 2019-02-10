package aws_test

import (
	"testing"

	"github.com/lonelyelk/what-build/aws"
)

func TestFindProjects(t *testing.T) {
	projects := []aws.Project{aws.Project{Name: "name1"}, aws.Project{Name: "name2"}}
	found := aws.FindProjects([]string{"name2", "name3"}, &projects)
	if len(found) != 1 {
		t.Errorf("Expected %v to have length of 1", found)
	}
	if found[0].Name != "name2" {
		t.Errorf("Expected to find 'name2' in %v", found)
	}
}

func TestFindProjects_Empty(t *testing.T) {
	projects := []aws.Project{aws.Project{Name: "name1"}, aws.Project{Name: "name2"}}
	found := aws.FindProjects([]string{}, &projects)
	if len(found) != len(projects) {
		t.Errorf("Expected %v to have length of %d", found, len(projects))
	}
	if found[0].Name != projects[0].Name || found[1].Name != projects[1].Name {
		t.Errorf("Expected to find all values from %v in %v", projects, found)
	}
}

func TestFindBuilds(t *testing.T) {
	builds := []aws.Build{aws.Build{Name: "name1"}, aws.Build{Name: "name2"}}
	found := aws.FindBuilds([]string{"name2", "name3"}, &builds)
	if len(found) != 1 {
		t.Errorf("Expected %v to have length of 1", found)
	}
	if found[0].Name != "name2" {
		t.Errorf("Expected to find 'name2' in %v", found)
	}
}

func TestFindBuilds_Empty(t *testing.T) {
	builds := []aws.Build{aws.Build{Name: "name1"}, aws.Build{Name: "name2"}}
	found := aws.FindBuilds([]string{}, &builds)
	if len(found) != len(builds) {
		t.Errorf("Expected %v to have length of %d", found, len(builds))
	}
	if found[0].Name != builds[0].Name || found[1].Name != builds[1].Name {
		t.Errorf("Expected to find all values from %v in %v", builds, found)
	}
}
