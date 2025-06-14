package model

import (
	"fmt"
	"regexp"
)

// VarValue は変数の値
type VarValue string

func NewVarValue(value any) VarValue {
	v := ""
	switch value.(type) {
	case string:
		v = value.(string)
	case *string:
		v = *value.(*string)
	default:
		fmt.Printf("[🚨WARNING!!] unknown var value type: type = %T. value = %+v\n", value, value)
	}

	return VarValue(v)
}

// Default はデフォルト値を返す
func (v VarValue) Default() string {
	pipePattern := regexp.MustCompile(`\{\{\s*\..+\s*\|\s*default\s+"(.+)"\s*\}\}`)
	pipeMatches := pipePattern.FindStringSubmatch(string(v))
	if len(pipeMatches) == 2 {
		return pipeMatches[1]
	}

	prefixPattern := regexp.MustCompile(`\{\{\s*default\s+"(.+)"\s+\..+\s*\}\}`)
	prefixMatches := prefixPattern.FindStringSubmatch(string(v))
	if len(prefixMatches) == 2 {
		return prefixMatches[1]
	}

	return ""
}

func (v VarValue) IsOptional(name string) bool {
	escapedName := regexp.QuoteMeta(name)
	pattern := fmt.Sprintf(`\{\{\s*\.%s\s*\}\}`, escapedName)

	regex := regexp.MustCompile(pattern)
	return regex.MatchString(string(v))
}

func (v VarValue) IsOptionalWithDefault(name string) bool {
	escapedName := regexp.QuoteMeta(name)

	pipePattern := fmt.Sprintf(`\{\{\s*\.%s\s*\|\s*default\s+.+?\s*\}\}`, escapedName)
	prefixPattern := fmt.Sprintf(`\{\{\s*default\s+.+?\s+\.%s\s*\}\}`, escapedName)
	fullPattern := fmt.Sprintf(`^(?:%s|%s)$`, pipePattern, prefixPattern)

	regex := regexp.MustCompile(fullPattern)
	return regex.MatchString(string(v))
}
