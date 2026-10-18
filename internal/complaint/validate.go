package complaint

import (
	"fmt"
	"strings"
	"unicode/utf8"
)
const (
	MinMessageLen    = 10
	MinOrgNameLen    = 2
	MinPersonNameLen = 3
)

var schools = map[string]struct{}{
	"ШЭиНГИ": {}, "ШГ": {}, "ШИТиИ": {}, "БШ": {}, "МШЭ": {},
	"КМА": {}, "ШПМ": {}, "ШХИ": {}, "ШМиЗТ": {},
}

func (in Input) normalized() Input {
	in.Organization = strings.TrimSpace(in.Organization)
	in.School = strings.TrimSpace(in.School)
	in.TeacherName = strings.TrimSpace(in.TeacherName)
	in.Subject = strings.TrimSpace(in.Subject)
	in.Message = strings.TrimSpace(in.Message)
	return in
}

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
