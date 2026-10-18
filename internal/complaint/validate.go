package complaint

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Validation limits. They are measured in runes (characters), not bytes:
// a Cyrillic letter takes 2 bytes in UTF-8, so len(s) would be wrong here.
const (
	MinMessageLen    = 10
	MinOrgNameLen    = 2
	MinPersonNameLen = 3
)

// schools are the schools of KBTU a staff complaint can refer to.
var schools = map[string]struct{}{
	"ШЭиНГИ": {}, "ШГ": {}, "ШИТиИ": {}, "БШ": {}, "МШЭ": {},
	"КМА": {}, "ШПМ": {}, "ШХИ": {}, "ШМиЗТ": {},
}

// normalized returns a copy of the input with surrounding spaces removed from
// every text field. Value receiver: the original input is not changed.
func (in Input) normalized() Input {
	in.Organization = strings.TrimSpace(in.Organization)
	in.School = strings.TrimSpace(in.School)
	in.TeacherName = strings.TrimSpace(in.TeacherName)
	in.Subject = strings.TrimSpace(in.Subject)
	in.Message = strings.TrimSpace(in.Message)
	return in
}

// Validate checks the input against the rules of its category and returns an
// error wrapping ErrInvalidInput for the first rule that is broken.
func (in Input) Validate() error {
	in = in.normalized()

	switch in.Category {
	case CategoryOrg:
		if runeLen(in.Organization) < MinOrgNameLen {
			return invalidf("organization name must be at least %d characters", MinOrgNameLen)
		}
	case CategoryStaff:
		if _, ok := schools[in.School]; !ok {
			return invalidf("unknown school %q", in.School)
		}
		if runeLen(in.TeacherName) < MinPersonNameLen {
			return invalidf("teacher name must be at least %d characters", MinPersonNameLen)
		}
	case CategoryCafeteria, CategoryDormitory, CategoryOther:
		// only the message is required for these categories
	default:
		return invalidf("unknown category %q", in.Category)
	}

	if runeLen(in.Message) < MinMessageLen {
		return invalidf("message must be at least %d characters", MinMessageLen)
	}
	return nil
}

func runeLen(s string) int { return utf8.RuneCountInString(s) }

func invalidf(format string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrInvalidInput, fmt.Sprintf(format, args...))
}
