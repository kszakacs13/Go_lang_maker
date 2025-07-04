package main

import (
	"fmt"
	"regexp"
	"strings"
)

var InputString = "6+6+5"

var Rules = map[string]string{
	"ADD": "NUM '+' NUM || NUM",
	"NUM": "'[0-9]'",
}

func main() {

	fmt.Println("Hi")

}

func parser(expressionString string, ruleRow string) struct{} {

	type parsingTree struct {
		OpType   string
		Children []*parsingTree
	}

	// Looping through the actual ruleRow's elements, searching for a regex

	rulesInRow := strings.Split(ruleRow, "||")

	// Going through all rules in a row and extract the operators belonging to them

	for j := 0; j < len(rulesInRow); j++ {

		ruleElements := strings.Split(ruleRow, " ")

		regex := ""

		for i := 0; i < len(ruleElements); i++ {

			element := ruleElements[i]

			if strings.Contains(element, "'") {
				regex = betweenSigns(element, "'")
				regex = "`" + regex + "`"
				pattern := regexp.MustCompile(regex)

				// Searching if the REGEX occurs in the expressionString

				positionOfRegex := pattern.FindStringIndex(expressionString)

				if positionOfRegex != nil {

					// If we found it, we get the substrings of the rule

				}
			}

		}

	}

}

func betweenSigns(stringToSlice string, signs string) []string {

	insideSigns := false

	outputString := []string{}

	allOthers := []string{}
	allOthersIndex := 0

	for character := 0; character < len(stringToSlice); character++ {

		actualChar := stringToSlice[character : character+1]

		if actualChar == signs {
			// If we bump into a separator char, we will track if we are inside or outside the string enclosed by them
			insideSigns = !insideSigns

			if !insideSigns {
				allOthersIndex++
			}

		} else {
			// If not a separator and we are not inbetween separators, we add to the output string
			if insideSigns {

				if len(outputString)-1 < allOthersIndex {
					outputString = append(outputString, actualChar)
				} else {
					outputString[allOthersIndex] += actualChar
				}

			} else {

				if len(allOthers)-1 < allOthersIndex {
					allOthers = append(allOthers, actualChar)
				} else {
					allOthers[allOthersIndex] += actualChar
				}

			}
		}

	}

	return outputString

}
