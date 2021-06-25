package table

// ColumnDefinition defines the relevant information for a column in a table
// plugin. Both values are mandatory. Prefer using the *Column helpers to
// create ColumnDefinition structs.
type ColumnDefinition struct {
	Name        string     `json:"name,omitempty"`
	Type        ColumnType `json:"type,omitempty"`
	Description string     `json:"description,omitempty"`

	// Options from https://github.com/osquery/osquery/blob/master/osquery/core/sql/column.h#L37
	Index      bool `json:"index"`
	Required   bool `json:"required"`
	Additional bool `json:"additional"`
	Optimized  bool `json:"optimized"`
	Hidden     bool `json:"hidden"`
}

// ColumnType is a strongly typed representation of the data type string for a
// column definition. The named constants should be used.
type ColumnType string

// The following column types are defined in osquery tables.h.
const (
	ColumnTypeUnknown        ColumnType = "UNKNOWN"
	ColumnTypeText                      = "TEXT"
	ColumnTypeInteger                   = "INTEGER"
	ColumnTypeBigInt                    = "BIGINT"
	ColumnTypeUnsignedBigInt            = "UNSIGNED BIGINT"
	ColumnTypeDouble                    = "DOUBLE"
	ColumnTypeBlob                      = "BLOB"
)

type ColumnOpt func(*ColumnDefinition)

// TextColumn is a helper for defining columns containing strings.
func TextColumn(name string, opts ...ColumnOpt) ColumnDefinition {
	return NewColumn(name, ColumnTypeText, opts...)
}

// IntegerColumn is a helper for defining columns containing integers.
func IntegerColumn(name string, opts ...ColumnOpt) ColumnDefinition {
	return NewColumn(name, ColumnTypeInteger, opts...)
}

// BigIntColumn is a helper for defining columns containing big integers.
func BigIntColumn(name string, opts ...ColumnOpt) ColumnDefinition {
	return NewColumn(name, ColumnTypeBigInt, opts...)
}

// DoubleColumn is a helper for defining columns containing floating point
// values.
func DoubleColumn(name string, opts ...ColumnOpt) ColumnDefinition {
	return NewColumn(name, ColumnTypeDouble, opts...)
}

// NewColumn returns a ColumnDefinition for the specified column.
func NewColumn(name string, ctype ColumnType, opts ...ColumnOpt) ColumnDefinition {
	cd := ColumnDefinition{
		Name: name,
		Type: ctype,
	}

	for _, opt := range opts {
		opt(&cd)
	}

	return cd

}

// IndexColumn is a functional argument to declare this as an indexed
// column. Depending on impmelentation, this can significantly change
// performance.  See osquery source code for more information.
func IndexColumn() ColumnOpt {
	return func(cd *ColumnDefinition) {
		cd.Index = true
	}
}

// RequiredColumn is a functional argument that sets this as a
// required column. sqlite will not process queries, if a required
// column is missing. See osquery source code for more information.
func RequiredColumn() ColumnOpt {
	return func(cd *ColumnDefinition) {
		cd.Required = true
	}

}

// AdditionalColumn is a functional argument that sets this as an
// additional column. See osquery source code for more information.
func AdditionalColumn() ColumnOpt {
	return func(cd *ColumnDefinition) {
		cd.Additional = true
	}

}

// OptimizedColumn is a functional argument that sets this as an
// optimized column. See osquery source code for more information.
func OptimizedColumn() ColumnOpt {
	return func(cd *ColumnDefinition) {
		cd.Optimized = true
	}

}

// HiddenColumn is a functional argument that sets this as an
// hidden column. This omits it from `select *` queries. See osquery source code for more information.
func HiddenColumn() ColumnOpt {
	return func(cd *ColumnDefinition) {
		cd.Hidden = true
	}

}

// ColumnDescription sets the column description. This is not
// currently part of the underlying osquery api, it is here for human
// consumption. It may become part of osquery spec generation.
func ColumnDescription(d string) ColumnOpt {
	return func(cd *ColumnDefinition) {
		cd.Description = d
	}
}

// Options returns the bitmask representation of the boolean column
// options. This uses the values as encoded in
// https://github.com/osquery/osquery/blob/master/osquery/core/sql/column.h#L37
func (c *ColumnDefinition) Options() uint8 {
	optionsBitmask := uint8(0)

	optionValues := map[uint8]bool{
		1:  c.Index,
		2:  c.Required,
		4:  c.Additional,
		8:  c.Optimized,
		16: c.Hidden,
	}

	for v, b := range optionValues {
		if b {
			optionsBitmask = optionsBitmask | v
		}
	}
	return optionsBitmask
}
