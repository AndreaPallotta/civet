package report

import (
	"encoding/json"
	"io"

	"github.com/AndreaPallotta/civet/internal/engine"
)

type JSONReporter struct{}

func NewJSONReporter() *JSONReporter {
	return &JSONReporter{}
}

func (r *JSONReporter) Write(w io.Writer, rep *engine.Report) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(rep)
}
