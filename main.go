package main

import (
	"fmt"
	"regexp"
	"strings"
)

var InputString = "6+6+5"

var Rules = map[string]string{
	"ADD": "NUM '+' NUM || NUM",
	"NUM": "'[0-9]+'",
}

func main() {

	fmt.Println(parser(InputString, "ADD", Rules))

}

func parser(expressionString string, ruleRowName string, rules map[string]string) []interface{} {

	var parserTree []interface{}

	ruleRow := rules[ruleRowName]

	rulesInRow := strings.Split(ruleRow, "||")

	// Going through all rules in a row and extract the operators belonging to them

	for j := 0; j < len(rulesInRow); j++ {

		ruleElements := strings.Split(rulesInRow[j], " ")

		someRegexNotFound := false

		allElements := [][]string{}

		for i := 0; i < len(ruleElements); i++ {

			element := ruleElements[i]

			if strings.Contains(element, "'") {
				regex := betweenSigns(element, "'")
				// regex = "`" + regex + "`"

				pattern := regexp.MustCompile(regex)

				// Searching if the REGEX occurs in the expressionString

				positionOfRegex := pattern.FindStringIndex(expressionString)

				if positionOfRegex == nil {

					someRegexNotFound = true // Continue to the next rule in the row if no matches found

				} else {

					regexElement := pattern.FindString(expressionString)

					allElements = append(allElements, []string{regexElement, "regex"})
				}

			} else {

				// Here we extract every non REGEX element - if we have just one and no regex elements, we ought to go recursive with that one deeper in the Rules map

				allElements = append(allElements, []string{element, "non-regex"})

			}

		}

		if someRegexNotFound {
			continue // If we did not find a regex specified in the rule, we continue to the next rule in the row
		}

		// Go recursive with all the non-regex elements, if there is none then we have to return the regex match

		for i := 0; i < len(allElements); i++ {

			if allElements[i][1] == "non-regex" {
				subTree := parser(expressionString, allElements[i][0], rules)
				parserTree = append(parserTree, ruleRowName)
				parserTree = append(parserTree, subTree)
			} else {
				parserTree = append(parserTree, ruleRowName)
				parserTree = append(parserTree, allElements[i][0])
			}

		}

	}

	return parserTree

}

func betweenSigns(stringToSlice string, signs string) string {

	insideSigns := false

	outputString := ""

	for character := 0; character < len(stringToSlice); character++ {

		actualChar := stringToSlice[character : character+1]

		if actualChar == signs {
			// If we bump into a separator char, we will track if we are inside or outside the string enclosed by them
			insideSigns = !insideSigns

		} else {
			// If not a separator and we are not inbetween separators, we add to the output string
			if insideSigns {

				outputString += actualChar

			}
		}

	}

	return outputString

}
