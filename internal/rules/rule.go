package rules

import (
	"fmt"
	"strings"
)

// Severity represents the severity level of a finding.
type Severity int

const (
	SeverityInfo Severity = iota
	SeverityWarning
	SeverityError
)

func (s Severity) String() string {
	switch s {
	case SeverityInfo:
		return "INFO"
	case SeverityWarning:
		return "WARNING"
	case SeverityError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

func (s Severity) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", strings.ToLower(s.String()))), nil
}

// Category represents a scoring category for pipeline analysis.
type Category int

const (
	CategoryPerformance Category = iota
	CategorySecurity
	CategoryReliability
	CategoryMaintainability
	CategoryObservability
	CategoryCompliance
)

var categoryNames = map[Category]string{
	CategoryPerformance:     "Performance",
	CategorySecurity:        "Security",
	CategoryReliability:     "Reliability",
	CategoryMaintainability: "Maintainability",
	CategoryObservability:   "Observability",
	CategoryCompliance:      "Compliance",
}

func (c Category) String() string {
	if name, ok := categoryNames[c]; ok {
		return name
	}
	return "Unknown"
}

func (c Category) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", strings.ToLower(c.String()))), nil
}

// AllCategories returns all scoring categories in display order.
func AllCategories() []Category {
	return []Category{
		CategoryPerformance,
		CategorySecurity,
		CategoryReliability,
		CategoryMaintainability,
		CategoryObservability,
		CategoryCompliance,
	}
}

// Platform represents a CI/CD platform.
type Platform int

const (
	PlatformGitLab Platform = iota
	PlatformGitHub
)

func (p Platform) String() string {
	switch p {
	case PlatformGitLab:
		return "GitLab CI"
	case PlatformGitHub:
		return "GitHub Actions"
	default:
		return "Unknown"
	}
}

func (p Platform) MarshalJSON() ([]byte, error) {
	switch p {
	case PlatformGitLab:
		return []byte(`"gitlab"`), nil
	case PlatformGitHub:
		return []byte(`"github"`), nil
	default:
		return []byte(`"unknown"`), nil
	}
}

// Location identifies where a finding was detected.
type Location struct {
	File string `json:"file"`
	Line int    `json:"line,omitempty"`
	Job  string `json:"job,omitempty"`
}

func (l Location) String() string {
	s := l.File
	if l.Line > 0 {
		s = fmt.Sprintf("%s:%d", s, l.Line)
	}
	if l.Job != "" {
		s = fmt.Sprintf("%s (job: %s)", s, l.Job)
	}
	return s
}

// Finding represents a single issue detected by a rule.
type Finding struct {
	RuleID     string   `json:"rule_id"`
	RuleName   string   `json:"rule_name"`
	Severity   Severity `json:"severity"`
	Category   Category `json:"category"`
	Message    string   `json:"message"`
	Location   Location `json:"location"`
	Suggestion string   `json:"suggestion,omitempty"`
	DocURL     string   `json:"doc_url,omitempty"`
}

// Rule is the interface that all pipeline analysis rules must implement.
type Rule interface {
	// ID returns the unique rule identifier (e.g., "UNI-001", "GL-003").
	ID() string
	// Name returns a human-readable name for the rule.
	Name() string
	// Category returns the scoring category this rule belongs to.
	Category() Category
	// DefaultSeverity returns the default severity level.
	DefaultSeverity() Severity
	// Platforms returns which platforms this rule applies to.
	Platforms() []Platform
	// Check runs the rule against a parsed pipeline and returns any findings.
	Check(pipeline *Pipeline) []Finding
}

// Pipeline is the normalized representation of a CI/CD pipeline configuration
// that rules operate on. Parsers convert platform-specific YAML into this struct.
type Pipeline struct {
	Platform Platform
	FilePath string
	Raw      string // raw YAML content

	// Parsed structure
	Stages   []string
	Jobs     map[string]*Job
	Workflow *Workflow

	// GitLab-specific
	Include []Include
	// GitHub-specific
	Permissions map[string]string
	Concurrency *Concurrency
}

// Job represents a single job/step in a pipeline.
type Job struct {
	Name          string
	Stage         string            // GitLab stage name
	Image         string
	Tags          []string          // GitLab runner tags
	RunsOn        string            // GitHub runs-on
	Script        []string          // GitLab script lines
	Steps         []Step            // GitHub steps
	Needs         []string          // explicit dependencies
	Rules         []JobRule         // GitLab rules: / GitHub if:
	Cache         *Cache
	Artifacts     *Artifacts
	Services      []Service
	Environment   string
	ResourceGroup string            // GitLab resource_group
	Interruptible *bool
	AllowFailure  bool
	Timeout       string
	Retry         *Retry
	Variables     map[string]string
	Only          []string          // deprecated GitLab only:
	Except        []string          // deprecated GitLab except:
	Line          int               // line number in source file
}

// Step represents a GitHub Actions step.
type Step struct {
	Name    string
	Uses    string
	Run     string
	With    map[string]string
	Env     map[string]string
	If      string
	Timeout int
}

// JobRule represents a conditional rule for job execution.
type JobRule struct {
	If      string
	When    string
	Changes []string
}

// Cache represents cache configuration.
type Cache struct {
	Key   string
	Paths []string
}

// Artifacts represents artifact configuration.
type Artifacts struct {
	Paths    []string
	ExpireIn string
	Reports  map[string]string
}

// Service represents a service container.
type Service struct {
	Name  string
	Alias string
}

// Retry represents retry configuration.
type Retry struct {
	Max  int
	When []string
}

// Include represents a GitLab CI include directive.
type Include struct {
	Template string
	Project  string
	File     string
	Remote   string
	Local    string
	Ref      string
}

// Workflow represents top-level workflow configuration.
type Workflow struct {
	// GitLab: workflow:rules
	Rules []JobRule
	// GitHub: on: triggers
	Triggers map[string]interface{}
}

// Concurrency represents GitHub Actions concurrency configuration.
type Concurrency struct {
	Group            string
	CancelInProgress bool
}
