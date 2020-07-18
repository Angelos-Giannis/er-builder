package util

import (
	pluralize "github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
)

// Util describes the utilities package.
type Util struct {
	pluralize *pluralize.Client
}

// New creates and returns a new utilities instance.
func New() *Util {
	return &Util{
		pluralize: pluralize.NewClient(),
	}
}

// GetCaseOfString converts a given string to the desired format based on provided case option.
func (u *Util) GetCaseOfString(initialValue, convertToCase string) string {
	switch convertToCase {
	case "snake_case":
		return strcase.ToSnake(initialValue)
	case "camelCase":
		return strcase.ToLowerCamel(initialValue)
	case "screaming_snake_case":
		return strcase.ToScreamingSnake(initialValue)
	case "kebab_case":
		return strcase.ToKebab(initialValue)
	default:
		return strcase.ToSnake(initialValue)
	}
}

// GetValueCount converts a string to either pluran or singular.
func (u *Util) GetValueCount(toPlural bool, initialValue string) string {
	if toPlural {
		return u.pluralize.Plural(initialValue)
	}

	return u.pluralize.Singular(initialValue)
}
