/*
    Copyright (C) 2020 Accurics, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

package utils

import (
	"regexp"
	"strings"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
)

var (
	skipRulesPattern = regexp.MustCompile(`(#ts:skip=[ \t]*(([A-Za-z0-9]+[.-]{1}){3,5}([\d]+)){1}([ \t]+.*){0,1})`)
	skipRulesPrefix  = "#ts:skip="
)

// GetSkipRules returns a list of rules to be skipped. The rules to be skipped
// can be set in terraform resource config with the following pattern:
// #ts:skip=AWS.S3Bucket.DS.High.1043
// $ts:skip=AWS.S3Bucket.DS.High.1044 reason to skip the rule
// each rule and its optional comment must be in a new line
func GetSkipRules(body string) []output.SkipRule {
	var skipRules []output.SkipRule

	// check if any rules comments are present in body
	if !skipRulesPattern.MatchString(body) {
		return skipRules
	}

	// get all skip rule comments
	comments := skipRulesPattern.FindAllString(body, -1)

	// extract rule ids from comments
	for _, c := range comments {
		c = strings.TrimPrefix(c, skipRulesPrefix)
		skipRule := getSkipRuleObject(c)
		if skipRule != nil {
			skipRules = append(skipRules, *skipRule)
		}
	}
	return skipRules
}

func getSkipRuleObject(s string) *output.SkipRule {
	if s == "" {
		return nil
	}
	var skipRule output.SkipRule
	ruleComment := strings.Fields(s)

	skipRule.Rule = strings.TrimSpace(ruleComment[0])
	if len(ruleComment) > 1 {
		comment := strings.Join(ruleComment[1:], " ")
		skipRule.Comment = strings.TrimSpace(comment)
	}
	return &skipRule
}
