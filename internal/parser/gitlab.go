package parser

import (
	"strings"

	"github.com/AndreaPallotta/civet/internal/rules"
	"gopkg.in/yaml.v3"
)

var gitlabReservedKeys = map[string]bool{
	"default":       true,
	"include":       true,
	"stages":        true,
	"variables":     true,
	"workflow":      true,
	"cache":         true,
	"services":      true,
	"after_script":  true,
	"before_script": true,
	"image":         true,
	"pages":         true,
}

// gitlabJob struct models the GitLab CI job YAML structure.
type gitlabJob struct {
	Stage         string            `yaml:"stage"`
	Image         string            `yaml:"image"`
	Script        []string          `yaml:"script"`
	Needs         []string          `yaml:"needs"`
	Rules         []rules.JobRule   `yaml:"rules"`
	Only          []string          `yaml:"only"`
	Except        []string          `yaml:"except"`
	Cache         *rules.Cache      `yaml:"cache"`
	Artifacts     *rules.Artifacts  `yaml:"artifacts"`
	Services      []rules.Service   `yaml:"services"`
	Environment   string            `yaml:"environment"`
	ResourceGroup string            `yaml:"resource_group"`
	Interruptible *bool             `yaml:"interruptible"`
	AllowFailure  bool              `yaml:"allow_failure"`
	Timeout       string            `yaml:"timeout"`
	Retry         *rules.Retry      `yaml:"retry"`
	Variables     map[string]string `yaml:"variables"`
}

type gitlabInclude struct {
	Template string `yaml:"template"`
	Project  string `yaml:"project"`
	File     string `yaml:"file"`
	Remote   string `yaml:"remote"`
	Local    string `yaml:"local"`
	Ref      string `yaml:"ref"`
}

// ParseGitLab parses a .gitlab-ci.yml file.
func ParseGitLab(content, filePath string) (*rules.Pipeline, error) {
	var root yaml.Node
	if err := yaml.Unmarshal([]byte(content), &root); err != nil {
		return nil, err
	}

	pipeline := &rules.Pipeline{
		Platform: rules.PlatformGitLab,
		FilePath: filePath,
		Raw:      content,
		Jobs:     make(map[string]*rules.Job),
	}

	if len(root.Content) == 0 {
		return pipeline, nil
	}
	doc := root.Content[0]
	if doc.Kind != yaml.MappingNode {
		return pipeline, nil
	}

	for i := 0; i < len(doc.Content); i += 2 {
		keyNode := doc.Content[i]
		valNode := doc.Content[i+1]
		key := keyNode.Value

		if strings.HasPrefix(key, ".") {
			continue // hidden job/template
		}

		switch key {
		case "stages":
			var stages []string
			valNode.Decode(&stages)
			pipeline.Stages = stages
		case "include":
			var includes []gitlabInclude
			// include can be a single object, an array, or an array of objects
			if valNode.Kind == yaml.SequenceNode {
				valNode.Decode(&includes)
			} else if valNode.Kind == yaml.MappingNode {
				var inc gitlabInclude
				valNode.Decode(&inc)
				includes = append(includes, inc)
			} else if valNode.Kind == yaml.ScalarNode {
				includes = append(includes, gitlabInclude{Local: valNode.Value})
			}
			for _, inc := range includes {
				pipeline.Include = append(pipeline.Include, rules.Include(inc))
			}
		case "workflow":
			var wf struct {
				Rules []rules.JobRule `yaml:"rules"`
			}
			valNode.Decode(&wf)
			pipeline.Workflow = &rules.Workflow{Rules: wf.Rules}
		default:
			if !gitlabReservedKeys[key] {
				var gj gitlabJob
				valNode.Decode(&gj)
				
				job := &rules.Job{
					Name:          key,
					Stage:         gj.Stage,
					Image:         gj.Image,
					Script:        gj.Script,
					Needs:         gj.Needs,
					Rules:         gj.Rules,
					Only:          gj.Only,
					Except:        gj.Except,
					Cache:         gj.Cache,
					Artifacts:     gj.Artifacts,
					Services:      gj.Services,
					Environment:   gj.Environment,
					ResourceGroup: gj.ResourceGroup,
					Interruptible: gj.Interruptible,
					AllowFailure:  gj.AllowFailure,
					Timeout:       gj.Timeout,
					Retry:         gj.Retry,
					Variables:     gj.Variables,
					Line:          keyNode.Line,
				}
				pipeline.Jobs[key] = job
			}
		}
	}

	return pipeline, nil
}
