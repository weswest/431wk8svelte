package check

import (
	"strings"

	"github.com/errata-ai/regexp2"
	"github.com/errata-ai/vale/v2/internal/core"
	"github.com/errata-ai/vale/v2/internal/nlp"
)

// Repetition looks for repeated uses of Tokens.
type Repetition struct {
	Definition `mapstructure:",squash"`
	Tokens     []string
	Max        int
	pattern    *regexp2.Regexp
	Ignorecase bool
	Alpha      bool
}

// NewRepetition creates a new `repetition`-based rule.
func NewRepetition(cfg *core.Config, generic baseCheck, path string) (Repetition, error) {
	rule := Repetition{}

	err := decodeRule(generic, &rule)
	if err != nil {
		return rule, readStructureError(err, path)
	}

	err = checkScopes(rule.Scope, path)
	if err != nil {
		return rule, err
	}

	regex := ""
	if rule.Ignorecase {
		regex += ignoreCase
	}

	regex += `(` + strings.Join(rule.Tokens, "|") + `)`
	re, err := regexp2.CompileStd(regex)
	if err != nil {
		return rule, core.NewE201FromPosition(err.Error(), path, 1)
	}

	rule.pattern = re
	return rule, nil
}

// Run executes the the `repetition`-based rule.
//
// The rule looks for repeated matches of its regex -- such as "this this".
func (o Repetition) Run(blk nlp.Block, f *core.File) ([]core.Alert, error) {
	var curr, prev string
	var hit bool
	var ploc []int
	var count int
	var alerts []core.Alert

	txt := blk.Text
	for _, loc := range o.pattern.FindAllStringIndex(txt, -1) {
		converted, err := re2Loc(txt, loc)
		if err != nil {
			return alerts, err
		}
		curr = strings.TrimSpace(converted)

		if o.Ignorecase {
			hit = strings.EqualFold(curr, prev) && curr != ""
		} else {
			hit = curr == prev && curr != ""
		}

		hit = hit && (!o.Alpha || core.IsLetter(curr))
		if hit {
			count++
		}

		if hit && count > o.Max {
			pos := []int{ploc[0], loc[1]}

			converted, err = re2Loc(txt, pos)
			if err != nil {
				return alerts, err
			}

			if !strings.Contains(converted, "\n") {
				floc := []int{ploc[0], loc[1]}

				a, err := makeAlert(o.Definition, floc, txt)
				if err != nil {
					return alerts, err
				}

				a.Message, a.Description = formatMessages(o.Message,
					o.Description, curr)
				alerts = append(alerts, a)
				count = 0
			}
		}
		ploc = loc
		prev = curr
	}

	return alerts, nil
}

// Fields provides access to the internal rule definition.
func (o Repetition) Fields() Definition {
	return o.Definition
}

// Pattern is the internal regex pattern used by this rule.
func (o Repetition) Pattern() string {
	return o.pattern.String()
}
