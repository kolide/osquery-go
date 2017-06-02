package server

import (
	"context"
	"encoding/json"

	"github.com/kolide/osquery-golang/gen/osquery"
)

// TablePlugin is the minimum interface required to implement an osquery table
// as a plugin. Any value that implements this interface can be passed to
// NewTablePlugin to satisfy the OsqueryPlugin interface.
type TablePlugin interface {
	// TableName returns the name of the table the plugin implements.
	TableName() string

	// Columns returns the column definition of the table.
	Columns() []ColumnDefinition

	// Generate returns the rows generated by the table. The ctx argument
	// should be checked for cancellation if the generation performs a
	// substantial amount of work. The queryContext argument provides the
	// deserialized JSON query context from osquery.
	Generate(ctx context.Context, queryContext interface{}) ([]map[string]string, error)
}

// NewTablePlugin takes a value that implements TablePlugin and wraps it with
// the appropriate methods to satisfy the OsqueryPlugin interface. Use this to
// easily create plugins implementing osquery tables.
func NewTablePlugin(plugin TablePlugin) OsqueryPlugin {
	return &tablePluginImpl{plugin}
}

type tablePluginImpl struct {
	plugin TablePlugin
}

var _ OsqueryPlugin = (*tablePluginImpl)(nil)

func (t *tablePluginImpl) Name() string {
	return t.plugin.TableName()
}

func (t *tablePluginImpl) RegistryName() string {
	return "table"
}

func (t *tablePluginImpl) Routes() osquery.ExtensionPluginResponse {
	routes := []map[string]string{}
	for _, col := range t.plugin.Columns() {
		routes = append(routes, map[string]string{
			"id":   "column",
			"name": col.Name,
			"type": string(col.Type),
			"op":   "0",
		})
	}
	return routes
}

func (t *tablePluginImpl) Ping() osquery.ExtensionStatus {
	return StatusOK
}

func (t *tablePluginImpl) Call(ctx context.Context, request osquery.ExtensionPluginRequest) osquery.ExtensionResponse {
	switch request["action"] {
	case "generate":
		var queryContext interface{}
		if ctxJSON, ok := request["context"]; ok {
			err := json.Unmarshal([]byte(ctxJSON), &queryContext)
			if err != nil {
				return osquery.ExtensionResponse{
					Status: &osquery.ExtensionStatus{
						Code:    1,
						Message: "error parsing context JSON: " + err.Error(),
					},
				}
			}
		}

		rows, err := t.plugin.Generate(ctx, queryContext)

		if err != nil {
			return osquery.ExtensionResponse{
				Status: &osquery.ExtensionStatus{
					Code:    1,
					Message: "error generating table: " + err.Error(),
				},
			}
		}

		return osquery.ExtensionResponse{
			Status:   &StatusOK,
			Response: rows,
		}

	case "columns":
		return osquery.ExtensionResponse{
			Status:   &StatusOK,
			Response: t.Routes(),
		}

	default:
		return osquery.ExtensionResponse{
			Status: &osquery.ExtensionStatus{
				Code:    1,
				Message: "unknown action: " + request["action"],
			},
		}
	}

}

func (t *tablePluginImpl) Shutdown() {}

// ColumnDefinition defines the relevant information for a column in a table
// plugin. Both values are mandatory.
type ColumnDefinition struct {
	Name string
	Type ColumnType
}

// StringColumn is a helper for defining columns containing strings.
func StringColumn(name string) ColumnDefinition {
	return ColumnDefinition{
		Name: name,
		Type: ColumnTypeString,
	}
}

// StringColumn is a helper for defining columns containing integers.
func IntegerColumn(name string) ColumnDefinition {
	return ColumnDefinition{
		Name: name,
		Type: ColumnTypeInteger,
	}
}

// ColumnType is a strongly typed representation of the data type string for a
// column definition.
type ColumnType string

// ColumnTypeString is used for columns containing strings.
const ColumnTypeString ColumnType = "TEXT"

// ColumnTypeInteger is used for columns containing integers.
const ColumnTypeInteger ColumnType = "INTEGER"
