package domain

type Config struct {
	Images map[string]Image `yaml:"images"`
}

type Image struct {
	LabelConstraints map[string]*LabelConstraint `yaml:"labelConstraints,omitempty"`
}

type LabelConstraint struct {
	CompareSemver              *LabelConstraintCompareSemver `yaml:"compareSemver,omitempty"`
	CompareSemverResult        *bool                         `yaml:"compareSemverResult,omitempty"`
	CompareSemverResultMessage string                        `yaml:"compareSemverResultMessage,omitempty"`

	CompareString              *LabelConstraintCompareString `yaml:"compareString,omitempty"`
	CompareStringResult        *bool                         `yaml:"compareStringResult,omitempty"`
	CompareStringResultMessage bool                          `yaml:"compareStringResultMessage,omitempty"`

	Count              *LabelConstraintCounter `yaml:"count,omitempty"`
	CountResult        *bool                   `yaml:"countResult,omitempty"`
	CountResultMessage string                  `yaml:"countMessage,omitempty"`

	Exists              *bool  `yaml:"exists,omitempty"`
	ExistsResult        bool   `yaml:"existsResult,omitempty"`
	ExistsResultMessage string `yaml:"existsResultMessage,omitempty"`
}

type LabelConstraintCompareSemver struct {
	Equal              string `yaml:"equal,omitempty"`
	EqualResult        *bool  `yaml:"equalResult,omitempty"`
	EqualResultMessage string `yaml:"equalResultMessage,omitempty"`

	GreaterThan              string `yaml:"greaterThan,omitempty"`
	GreaterThanResult        *bool  `yaml:"greaterThanResult,omitempty"`
	GreaterThanResultMessage string `yaml:"greaterThanResultMessage,omitempty"`

	LessThan              string `yaml:"lessThan,omitempty"`
	LessThanResult        *bool  `yaml:"lessThanResult,omitempty"`
	LessThanResultMessage string `yaml:"lessThanResultMessage,omitempty"`
}

type LabelConstraintCompareString struct {
	Equal              string `yaml:"equal,omitempty"`
	EqualResult        *bool  `yaml:"equalResult,omitempty"`
	EqualResultMessage string `yaml:"equalResultMessage,omitempty"`

	HasPrefix              string `yaml:"hasPrefix,omitempty"`
	HasPrefixResult        *bool  `yaml:"hasPrefixResult,omitempty"`
	HasPrefixResultMessage string `yaml:"hasPrefixResultMessage,omitempty"`

	HasSuffix              string `yaml:"hasSuffix,omitempty"`
	HasSuffixResult        *bool  `yaml:"hasSuffixResult,omitempty"`
	HasSuffixResultMessage string `yaml:"hasSuffixResultMessage,omitempty"`
}

type LabelConstraintCounter struct {
	Equal              *uint  `yaml:"equal,omitempty"`
	EqualResult        *bool  `yaml:"equalResult,omitempty"`
	EqualResultMessage string `yaml:"equalResultMessage,omitempty"`

	GreaterThan              *uint  `yaml:"greaterThan,omitempty"`
	GreaterThanResult        *bool  `yaml:"greaterThanResult,omitempty"`
	GreaterThanResultMessage string `yaml:"greaterThanResultMessage,omitempty"`

	LessThan              *uint  `yaml:"lessThan,omitempty"`
	LessThanResult        *bool  `yaml:"lessThanResult,omitempty"`
	LessThanResultMessage string `yaml:"lessThanResultMessage,omitempty"`
}
