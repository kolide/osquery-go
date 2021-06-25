package table

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type osqueryTableSpec struct {
	Cacheable bool               `json:"cacheable"`
	Evented   bool               `json:"evented"`
	Name      string             `json:"name,omitempty"`
	Url       string             `json:"url,omitempty"`
	Platforms []string           `json:"platforms,omitempty"`
	Columns   []ColumnDefinition `json:"columns,omitempty"`
}

func (t *Plugin) Spec() (string, error) {
	// FIXME: the columndefinition type is upcased, is that an issue?
	tableSpec := osqueryTableSpec{
		Name:    t.name,
		Columns: t.columns,
		//Platforms: []string{"FIXME"},
	}
	specBytes, err := json.MarshalIndent(tableSpec, "", "  ")
	if err != nil {
		return "", errors.Wrap(err, "marshalling")
	}
	return string(specBytes), nil
}
