package check

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/errata-ai/regexp2"
	"github.com/errata-ai/vale/v2/internal/core"
	"github.com/errata-ai/vale/v2/internal/nlp"
	"github.com/jdkato/regexp"
	"github.com/mitchellh/mapstructure"

	"gopkg.in/yaml.v2"
)

var inlineScopes = []string{"code", "link", "strong", "emphasis"}

// FilterEnv is the environment passed to the `--filter` flag.
type FilterEnv struct {
	Rules []Definition
}

// Rule represents in individual writing construct to enforce.
type Rule interface {
	Run(blk nlp.Block, file *core.File) ([]core.Alert, error)
	Fields() Definition
	Pattern() string
}

// Definition holds the common attributes of rule definitions.
type Definition struct {
	Action      core.Action
	Description string
	Extends     string
	Level       string
	Limit       int
	Link        string
	Message     string
	Name        string
	Scope       []string
	Selector    Selector
}

var defaultStyles = []string{"Vale"}
var extensionPoints = []string{
	"capitalization",
	"conditional",
	"consistency",
	"existence",
	"occurrence",
	"repetition",
	"substitution",
	"readability",
	"spelling",
	"sequence",
	"metric",
	"script",
}
var defaultRules = map[string]map[string]interface{}{
	"Avoid": {
		"extends":    "existence",
		"name":       "Vale.Avoid",
		"level":      "error",
		"message":    "Avoid using '%s'.",
		"scope":      "text",
		"ignorecase": false,
		"tokens":     []string{},
		"path":       "internal",
	},
	"Terms": {
		"extends":    "substitution",
		"name":       "Vale.Terms",
		"level":      "error",
		"message":    "Use '%s' instead of '%s'.",
		"scope":      "text",
		"ignorecase": true,
		"swap":       map[string]string{},
		"path":       "internal",
	},
	"Repetition": {
		"extends": "repetition",
		"name":    "Vale.Repetition",
		"level":   "error",
		"message": "'%s' is repeated!",
		"scope":   "text",
		"alpha":   true,
		"action": core.Action{
			Name:   "edit",
			Params: []string{"truncate", " "},
		},
		"tokens": []string{`[^\s]+`},
		"path":   "internal",
	},
	"Spelling": {
		"extends": "spelling",
		"name":    "Vale.Spelling",
		"message": "Did you really mean '%s'?",
		"level":   "error",
		"scope":   "text",
		"action": core.Action{
			Name:   "suggest",
			Params: []string{"spellings"},
		},
		"ignore": []interface{}{"vocab.txt"},
		"path":   "internal",
	},
}

const (
	ignoreCase      = `(?i)`
	wordTemplate    = `(?m)\b(?:%s)\b`
	nonwordTemplate = `(?m)(?:%s)`
	tokenTemplate   = `^(?:%s)$`
)

type baseCheck map[string]interface{}

func buildRule(cfg *core.Config, generic baseCheck) (Rule, error) {
	path, ok := generic["path"].(string)
	if !ok {
		msg := fmt.Errorf("'%v' is not valid .", generic)
		return Existence{}, core.NewE100("buildRule: path", msg)
	}

	name, ok := generic["extends"].(string)
	if !ok {
		name = "unknown"
	}

	delete(generic, "path")
	switch name {
	case "existence":
		return NewExistence(cfg, generic, path)
	case "substitution":
		return NewSubstitution(cfg, generic, path)
	case "capitalization":
		return NewCapitalization(cfg, generic, path)
	case "occurrence":
		return NewOccurrence(cfg, generic, path)
	case "spelling":
		return NewSpelling(cfg, generic, path)
	case "repetition":
		return NewRepetition(cfg, generic, path)
	case "readability":
		return NewReadability(cfg, generic, path)
	case "conditional":
		return NewConditional(cfg, generic, path)
	case "consistency":
		return NewConsistency(cfg, generic, path)
	case "sequence":
		return NewSequence(cfg, generic, path)
	case "metric":
		return NewMetric(cfg, generic, path)
	case "script":
		return NewScript(cfg, generic, path)
	default:
		return Existence{}, core.NewE201FromTarget(
			fmt.Sprintf("'extends' key must be one of %v.", extensionPoints),
			name,
			path)
	}
}

func formatMessages(msg string, desc string, subs ...string) (string, string) {
	return core.FormatMessage(msg, subs...), core.FormatMessage(desc, subs...)
}

// NOTE: We need to do this because regexp2, the library we use for extended
// syntax, returns its locatons in *rune* offsets.
func re2Loc(s string, loc []int) (string, error) {
	converted := []rune(s)

	size := len(converted)
	if loc[0] < 0 || loc[1] > size {
		msg := fmt.Errorf("%d (%d:%d)", size, loc[0], loc[1])
		return "", core.NewE100("re2loc: bounds", msg)
	}

	return string(converted[loc[0]:loc[1]]), nil
}

func makeAlert(chk Definition, loc []int, txt string) (core.Alert, error) {
	match, err := re2Loc(txt, loc)
	if err != nil {
		return core.Alert{}, err
	}

	a := core.Alert{
		Check: chk.Name, Severity: chk.Level, Span: loc, Link: chk.Link,
		Match: match, Action: chk.Action}
	a.Message, a.Description = formatMessages(chk.Message, chk.Description, match)

	return a, nil
}

func parse(file []byte, path string) (map[string]interface{}, error) {
	generic := map[string]interface{}{}

	if err := yaml.Unmarshal(file, &generic); err != nil {
		r := regexp.MustCompile(`yaml: line (\d+): (.+)`)
		if r.MatchString(err.Error()) {
			groups := r.FindStringSubmatch(err.Error())
			i, err := strconv.Atoi(groups[1])
			if err != nil {
				return generic, core.NewE100("addCheck/Atoi", err)
			}
			return generic, core.NewE201FromPosition(groups[2], path, i)
		}
	} else if err := validateDefinition(generic, path); err != nil {
		return generic, err
	}

	return generic, nil
}

func validateDefinition(generic map[string]interface{}, path string) error {
	if point, ok := generic["extends"]; !ok || point == nil {
		return core.NewE201FromPosition(
			"Missing the required 'extends' key.",
			path,
			1)
	} else if !core.StringInSlice(point.(string), extensionPoints) {
		key := point.(string)
		return core.NewE201FromTarget(
			fmt.Sprintf("'extends' key must be one of %v.", extensionPoints),
			key,
			path)
	}

	if _, ok := generic["message"]; !ok {
		return core.NewE201FromPosition(
			"Missing the required 'message' key.",
			path,
			1)
	}

	if level, ok := generic["level"]; ok {
		if level == nil || !core.StringInSlice(level.(string), core.AlertLevels) {
			return core.NewE201FromTarget(
				fmt.Sprintf("'level' must be one of %v", core.AlertLevels),
				"level",
				path)
		}
	}

	if generic["code"] != nil && generic["code"].(bool) {
		return core.NewE201FromTarget(
			"`code` is deprecated; please use `scope: raw` instead.",
			"code",
			path)
	}

	return nil
}

func readStructureError(err error, path string) error {
	r1 := regexp.MustCompile(`\* '(.+)' (.+)`)
	r2 := regexp.MustCompile(`\* '(?:.*)' (.*): (\w+)`)
	if r1.MatchString(err.Error()) {
		groups := r1.FindStringSubmatch(err.Error())
		return core.NewE201FromTarget(
			groups[2],
			strings.ToLower(groups[1]),
			path)
	} else if r2.MatchString(err.Error()) {
		groups := r2.FindStringSubmatch(err.Error())
		return core.NewE201FromTarget(
			fmt.Sprintf("%s: '%s'", groups[1], groups[2]),
			strings.ToLower(groups[2]),
			path)
	}
	return core.NewE201FromPosition(err.Error(), path, 1)
}

func makeRegexp(
	template string,
	noCase bool,
	word func() bool,
	callback func() string,
	append bool,
) string {
	regex := ""

	if word() {
		if template != "" {
			regex += template
		} else {
			regex += wordTemplate
		}
	} else {
		regex += nonwordTemplate
	}

	if append {
		regex += callback()
	} else {
		regex = callback() + regex
	}

	if noCase {
		regex = ignoreCase + regex
	}

	return regex
}

func matchToken(expected, observed string, ignorecase bool) bool {
	p := expected
	if ignorecase {
		p = ignoreCase + p
	}

	r, err := regexp2.CompileStd(fmt.Sprintf(tokenTemplate, p))
	if core.IsPhrase(expected) || err != nil {
		return expected == observed
	}
	return r.MatchStringStd(observed)
}

func updateExceptions(previous []string, current map[string]struct{}) (*regexp2.Regexp, error) {
	for term := range current {
		previous = append(previous, term)
	}

	// NOTE: This is required to ensure that we have greedy alternation.
	sort.Slice(previous, func(p, q int) bool {
		return len(previous[p]) > len(previous[q])
	})

	regex := makeRegexp(
		"",
		false,
		func() bool { return true },
		func() string { return "" },
		true)

	regex = fmt.Sprintf(regex, strings.Join(previous, "|"))
	if len(previous) > 0 {
		return regexp2.CompileStd(regex)
	}

	return nil, nil
}

func decodeRule(input interface{}, output interface{}) error {
	config := mapstructure.DecoderConfig{
		ErrorUnused:      true,
		Squash:           true,
		WeaklyTypedInput: true,
		Result:           output,
	}

	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}

func checkScopes(scopes []string, path string) error {
	for _, scope := range scopes {
		if strings.Contains(scope, "&") {
			// FIXME: multi part ...
			continue
		}

		// Negation ...
		scope = strings.TrimPrefix(scope, "~")

		// Specification ...
		//
		// TODO: check sub-scopes too?
		scope = strings.Split(scope, ".")[0]

		if core.StringInSlice(scope, inlineScopes) {
			return core.NewE201FromTarget(
				fmt.Sprintf("scope '%v' is no longer supported; use 'raw' instead.", scope),
				"scope",
				path)
		} else if !core.StringInSlice(scope, allowedScopes) {
			return core.NewE201FromTarget(
				fmt.Sprintf("'%v' is not a valid scope; must be one of %v", scope, allowedScopes),
				"scope",
				path)
		}
	}

	return nil
}
