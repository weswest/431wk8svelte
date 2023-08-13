package spell

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jdkato/regexp"
)

// affixType is either an affix prefix or suffix
type affixType int

// specific Affix types
const (
	Prefix affixType = iota
	Suffix
)

// affix is a rule for affix (adding prefixes or suffixes)
type affix struct {
	Rules        []rule    // -
	Type         affixType // either PFX or SFX
	CrossProduct bool      // -
}

// expand provides all variations of a given word based on this affix rule
func (a affix) expand(word string, out []string) []string {
	for _, r := range a.Rules {
		if r.matcher != nil && !r.matcher.MatchString(word) {
			continue
		}
		if a.Type == Prefix {
			out = append(out, r.AffixText+word)
			// TODO is does Strip apply to prefixes too?
		} else {
			stripWord := word
			if r.Strip != "" && strings.HasSuffix(word, r.Strip) {
				stripWord = word[:len(word)-len(r.Strip)]
			}
			out = append(out, stripWord+r.AffixText)
		}
	}
	return out
}

// rule is a Affix rule
type rule struct {
	Strip     string
	AffixText string         // suffix or prefix text to add
	Pattern   string         // original matching pattern from AFF file
	matcher   *regexp.Regexp // matcher to see if this rule applies or not
}

// dictConfig is a partial representation of a Hunspell AFF (Affix) file.
type dictConfig struct {
	IconvReplacements []string
	Replacements      [][2]string
	CompoundRule      []string
	Flag              string
	TryChars          string
	WordChars         string
	CompoundOnly      string
	AffixMap          map[rune]affix
	CamelCase         int
	CompoundMin       int64
	compoundMap       map[rune][]string
	NoSuggestFlag     rune
}

// expand expands a word/affix using dictionary/affix rules
//  This also supports CompoundRule flags
func (a dictConfig) expand(wordAffix string, out []string) ([]string, error) {
	out = out[:0]
	idx := strings.Index(wordAffix, "/")

	// not found
	if idx == -1 {
		out = append(out, wordAffix)
		return out, nil
	}
	if idx == 0 || idx+1 == len(wordAffix) {
		return nil, fmt.Errorf("Slash char found in first or last position")
	}
	// safe
	word, keyString := wordAffix[:idx], wordAffix[idx+1:]

	// check to see if any of the flags are in the
	// "compound only".  If so then nothing to add
	compoundOnly := false
	for _, key := range keyString {
		if strings.IndexRune(a.CompoundOnly, key) != -1 {
			compoundOnly = true
			continue
		}
		if _, ok := a.compoundMap[key]; !ok {
			// the isn't a compound flag
			continue
		}
		// is a compound flag
		a.compoundMap[key] = append(a.compoundMap[key], word)
	}

	if compoundOnly {
		return out, nil
	}

	out = append(out, word)
	prefixes := make([]affix, 0, 5)
	suffixes := make([]affix, 0, 5)
	for _, key := range keyString {
		// want keyString to []?something?
		// then iterate over that
		af, ok := a.AffixMap[key]
		if !ok {
			// TODO: How should we handle this?
			continue
		}
		if !af.CrossProduct {
			out = af.expand(word, out)
			continue
		}
		if af.Type == Prefix {
			prefixes = append(prefixes, af)
		} else {
			suffixes = append(suffixes, af)
		}
	}

	// expand all suffixes with out any prefixes
	for _, suf := range suffixes {
		out = suf.expand(word, out)
	}
	for _, pre := range prefixes {
		prewords := pre.expand(word, nil)
		out = append(out, prewords...)

		// now do cross product
		for _, suf := range suffixes {
			for _, w := range prewords {
				out = suf.expand(w, out)
			}
		}
	}
	return out, nil
}

func isCrossProduct(val string) (bool, error) {
	switch val {
	case "Y":
		return true, nil
	case "N":
		return false, nil
	}
	return false, fmt.Errorf("CrossProduct is not Y or N: got %q", val)
}

// newDictConfig reads an Hunspell AFF file
func newDictConfig(file io.Reader) (*dictConfig, error) {
	aff := dictConfig{
		Flag:        "ASCII",
		AffixMap:    make(map[rune]affix),
		compoundMap: make(map[rune][]string),
		CompoundMin: 3, // default in Hunspell
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "TRY":
			if len(parts) < 2 {
				return nil, fmt.Errorf("TRY stanza had %d fields, expected 2", len(parts))
			}
			aff.TryChars = parts[1]
		case "ICONV":
			// if only 2 fields, then its the first stanza that just provides a count
			//  we don't care, as we dynamically allocate
			if len(parts) == 2 {
				continue
			} else if len(parts) < 3 {
				return nil, fmt.Errorf("ICONV stanza had %d fields, expected 2", len(parts))
			}
			aff.IconvReplacements = append(aff.IconvReplacements, parts[1], parts[2])
		case "REP":
			if len(parts) == 2 {
				continue
			} else if len(parts) < 3 {
				return nil, fmt.Errorf("REP stanza had %d fields, expected 2", len(parts))
			}
			aff.Replacements = append(aff.Replacements, [2]string{parts[1], parts[2]})
		case "COMPOUNDMIN":
			if len(parts) < 2 {
				return nil, fmt.Errorf("COMPOUNDMIN stanza had %d fields, expected 2", len(parts))
			}
			val, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("COMPOUNDMIN stanza had %q expected number", parts[1])
			}
			aff.CompoundMin = val
		case "ONLYINCOMPOUND":
			if len(parts) < 2 {
				return nil, fmt.Errorf("ONLYINCOMPOUND stanza had %d fields, expected 2", len(parts))
			}
			aff.CompoundOnly = parts[1]
		case "COMPOUNDRULE":
			if len(parts) < 2 {
				return nil, fmt.Errorf("COMPOUNDRULE stanza had %d fields, expected 2", len(parts))
			}
			val, err := strconv.ParseInt(parts[1], 10, 64)
			if err == nil {
				aff.CompoundRule = make([]string, 0, val)
			} else {
				aff.CompoundRule = append(aff.CompoundRule, parts[1])
				for _, char := range parts[1] {
					if _, ok := aff.compoundMap[char]; !ok {
						aff.compoundMap[char] = []string{}
					}
				}
			}
		case "NOSUGGEST":
			if len(parts) < 2 {
				return nil, fmt.Errorf("NOSUGGEST stanza had %d fields, expected 2", len(parts))
			}
			// should use runes or parse correctly
			chars := []rune(parts[1])
			if len(chars) != 1 {
				return nil, fmt.Errorf("NOSUGGEST stanza had more than one flag: %q", parts[1])
			}
			aff.NoSuggestFlag = chars[0]
		case "WORDCHARS":
			if len(parts) < 2 {
				return nil, fmt.Errorf("WORDCHAR stanza had %d fields, expected 2", len(parts))
			}
			aff.WordChars = parts[1]
		case "FLAG":
			if len(parts) < 2 {
				return nil, fmt.Errorf("FLAG stanza had %d, expected 1", len(parts))
			}
			aff.Flag = parts[1]
		case "PFX", "SFX":
			atype := Prefix
			if parts[0] == "SFX" {
				atype = Suffix
			}

			sections := len(parts)
			if sections > 4 {
				// does this need to be split out into suffix and prefix?
				flag := rune(parts[1][0])
				a, ok := aff.AffixMap[flag]
				if !ok {
					return nil, fmt.Errorf("Got rules for flag %q but no definition", flag)
				}

				strip := ""
				if parts[2] != "0" {
					strip = parts[2]
				}

				var matcher *regexp.Regexp
				var err error
				pat := parts[4]
				if pat != "." {
					if a.Type == Prefix {
						pat = "^" + pat
					} else {
						pat = pat + "$"
					}
					matcher, err = regexp.Compile(pat)
					if err != nil {
						return nil, fmt.Errorf("Unable to compile %s", pat)
					}
				}

				// See #499.
				//
				// TODO: Is this safe to do in all cases?
				if parts[3] == "0" {
					parts[3] = ""
				}

				a.Rules = append(a.Rules, rule{
					Strip:     strip,
					AffixText: parts[3],
					Pattern:   parts[4],
					matcher:   matcher,
				})
				aff.AffixMap[flag] = a
			} else if sections > 3 {
				cross, err := isCrossProduct(parts[2])
				if err != nil {
					return nil, err
				}
				// this is a new Affix!
				a := affix{
					Type:         atype,
					CrossProduct: cross,
				}
				flag := rune(parts[1][0])
				aff.AffixMap[flag] = a
			}
		default:
			// Do nothing.
			//
			// Hunspell ignores lines that don't start with a known directive.
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &aff, nil
}
