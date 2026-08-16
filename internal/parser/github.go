package parser

import (
	"strings"

	"github.com/AndreaPallotta/civet/internal/rules"
	"gopkg.in/yaml.v3"
)

// githubWorkflow models the root of a GitHub Actions YAML.
type githubWorkflow struct {
	Name        string                       `yaml:"name"`
	On          interface{}                  `yaml:"on"`
	Env         map[string]string            `yaml:"env"`
	Permissions map[string]string            `yaml:"permissions"`
	Concurrency *rules.Concurrency           `yaml:"concurrency"`
	Jobs        map[string]yaml.Node         `yaml:"jobs"`
}

type githubJob struct {
	Name           string            `yaml:"name"`
	RunsOn         string            `yaml:"runs-on"`
	Needs          []string          `yaml:"needs"`
	TimeoutMinutes int               `yaml:"timeout-minutes"`
	Env            map[string]string `yaml:"env"`
	Permissions    map[string]string `yaml:"permissions"`
	Concurrency    *rules.Concurrency`yaml:"concurrency"`
	Steps          []rules.Step      `yaml:"steps"`
	Services       map[string]rules.Service `yaml:"services"`
}

// ParseGitHub parses a GitHub Actions YAML file.
func ParseGitHub(content, filePath string) (*rules.Pipeline, error) {
	var gw githubWorkflow
	if err := yaml.Unmarshal([]byte(content), &gw); err != nil {
		return nil, err
	}

	pipeline := &rules.Pipeline{
		Platform:    rules.PlatformGitHub,
		FilePath:    filePath,
		Raw:         content,
		Jobs:        make(map[string]*rules.Job),
		Permissions: gw.Permissions,
		Concurrency: gw.Concurrency,
	}

	// Parse 'on:'
	if gw.On != nil {
		pipeline.Workflow = &rules.Workflow{
			Triggers: make(map[string]interface{}),
		}
		
		switch v := gw.On.(type) {
		case string:
			pipeline.Workflow.Triggers[v] = nil
		case []interface{}:
			for _, trig := range v {
				if s, ok := trig.(string); ok {
					pipeline.Workflow.Triggers[s] = nil
				}
			}
		case map[string]interface{}:
			pipeline.Workflow.Triggers = v
		}
	}

	// We need line numbers for jobs, so we parse jobs using the yaml.Node mapping we stored in githubWorkflow
	// wait, `yaml:"jobs"` as `map[string]yaml.Node` doesn't give us the key node's line number natively.
	// But it's good enough to decode the value node. The key node isn't exposed directly when decoding into map[string]yaml.Node.
	// Let's decode the whole file as a yaml.Node to find job line numbers.
	
	var root yaml.Node
	yaml.Unmarshal([]byte(content), &root)
	
	jobLines := make(map[string]int)
	if len(root.Content) > 0 && root.Content[0].Kind == yaml.MappingNode {
		doc := root.Content[0]
		for i := 0; i < len(doc.Content); i += 2 {
			if doc.Content[i].Value == "jobs" && doc.Content[i+1].Kind == yaml.MappingNode {
				jobsNode := doc.Content[i+1]
				for j := 0; j < len(jobsNode.Content); j += 2 {
					jobKey := jobsNode.Content[j]
					jobLines[jobKey.Value] = jobKey.Line
				}
				break
			}
		}
	}

	for jobID, jobNode := range gw.Jobs {
		var gj githubJob
		jobNode.Decode(&gj)
		
		var services []rules.Service
		for svcName, svc := range gj.Services {
			svc.Name = svcName
			services = append(services, svc)
		}

		timeoutStr := ""
		if gj.TimeoutMinutes > 0 {
			timeoutStr = "timeout" // placeholder, since github uses ints for minutes
		}

		var cache *rules.Cache
		for _, step := range gj.Steps {
			if strings.Contains(step.Uses, "actions/cache") {
				cache = &rules.Cache{Key: step.With["key"]}
				break
			}
			if strings.HasPrefix(step.Uses, "actions/setup-") {
				if _, ok := step.With["cache"]; ok {
					cache = &rules.Cache{Key: step.Uses}
					break
				}
			}
		}

		job := &rules.Job{
			Name:          jobID,
			RunsOn:        gj.RunsOn,
			Needs:         gj.Needs,
			Steps:         gj.Steps,
			Services:      services,
			Variables:     gj.Env,
			Timeout:       timeoutStr,
			Cache:         cache,
			Line:          jobLines[jobID],
		}
		pipeline.Jobs[jobID] = job
	}

	return pipeline, nil
}
