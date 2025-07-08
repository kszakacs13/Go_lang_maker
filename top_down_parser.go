package main

// Top-down parser

func td_parser(expressionString string, ruleRowName string, rules map[string]string) []interface{} {

	var parsedInterface []interface{}

	// Take all the only regex, and not-just-regex containing grammar rules

	only_regex := map[string]string{}
	other_rules := map[string]string{}

	for key, rule := range rules {
		regexLen := len(betweenSigns(rule, RegexSeparator))
		ruleLen := len(rule)

		if regexLen+2 >= ruleLen {
			only_regex[key] = rule
		} else {
			other_rules[key] = rule
		}

	}

	// Go through the string - parse every just-regex containing element

	return parsedInterface

}
