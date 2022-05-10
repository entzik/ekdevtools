package gradleparser

import (
	"io/ioutil"
	"regexp"
	"strings"
)

var dependencyExtractorRegex = regexp.MustCompile(`(?s)(implementation|testImplementation|api)\s*\(?\s*[\',\"]([\w\.-]+):([\w\.-]+):([\w\.\-\+]+)[\',\"]\s*\)?.*`)

type MavenDependency struct {
	Scope    string
	Group    string
	Artefact string
	Version  string
}

func NewMavenDependencyFromLine(line string) *MavenDependency {
	components := dependencyExtractorRegex.FindStringSubmatch(line)
	if components != nil {
		return &MavenDependency{
			Scope:    components[1],
			Group:    components[2],
			Artefact: components[3],
			Version:  components[4],
		}
	} else {
		return nil
	}

}

func (mavenDependency MavenDependency) ToString() string {
	return mavenDependency.Scope + "('" + mavenDependency.Group + ":" + mavenDependency.Artefact + ":" + mavenDependency.Version + "')"
}

func (mavenDependency MavenDependency) IsSameArtefact(otherDependency *MavenDependency) bool {
	return mavenDependency.Group == otherDependency.Group && mavenDependency.Artefact == otherDependency.Artefact
}

func ContainsDependency(path string, targetDependency *MavenDependency) (bool, error) {
	//log.Default().Printf("parsing gradle file: %s", path)
	input, err := ioutil.ReadFile(path)
	if err != nil {
		return false, err
	} else {
		lines := strings.Split(string(input), "\n")

		for _, line := range lines {
			dependency := NewMavenDependencyFromLine(line)
			if dependency != nil {
				if dependency.IsSameArtefact(targetDependency) {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

func isDependency() bool {
	return false
}
