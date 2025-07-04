package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var InputString = "6+6+5"

var Rules = map[string]string{
	"ADD": "NUM '\\+' ADD || NUM",
	"NUM": "'[0-9]+'",
}

func main() {

	fmt.Println(parser(InputString, "ADD", Rules)) // Instead "ADD" we will put the foremost element

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

				pattern := regexp.MustCompile(regex)

				// Searching if the REGEX occurs in the expressionString

				positionOfRegex := pattern.FindStringIndex(expressionString)

				if positionOfRegex == nil {

					someRegexNotFound = true // Continue to the next rule in the row if no matches found

				} else {

					regexElement := pattern.FindString(expressionString)

					allElements = append(allElements, []string{regexElement, "regex", strconv.Itoa(positionOfRegex[0]), strconv.Itoa(positionOfRegex[1])}) // We save the positions to extract the important elements later
				}

			} else {

				// Here we extract every non REGEX element - if we have just one and no regex elements, we ought to go recursive with that one deeper in the Rules map - also we have to save the part that corresponds to the non-regex element

				allElements = append(allElements, []string{element, "non-regex"})

			}

		}

		if someRegexNotFound {
			continue // If we did not find a regex specified in the rule, we continue to the next rule in the row
		}

		// Extract the positions of the non-regex parts (they are all the parts that are not in the regex)

		extractStart := 0

		extractedPositions := []string{}

		foundNonRegex := false

		for i := 0; i < len(allElements); i++ {

			actualElement := allElements[i]

			if actualElement[1] == "regex" {

				regexStart, _ := strconv.Atoi(actualElement[2])
				regexEnd, _ := strconv.Atoi(actualElement[3])

				if extractStart != regexStart {
					extractedPositions = append(extractedPositions, expressionString[extractStart:regexStart])
					extractStart = regexEnd
				}

			} else {
				foundNonRegex = true
			}

		}

		if extractStart != len(expressionString) && foundNonRegex {
			extractedPositions = append(extractedPositions, expressionString[extractStart:])
		}

		// Go recursive with all the non-regex elements, if there is none then we have to return the regex match

		nonRegexElemCount := 0

		for i := 0; i < len(allElements); i++ {

			if allElements[i][1] == "non-regex" {
				subTree := parser(extractedPositions[nonRegexElemCount], allElements[i][0], rules) // we ought to cut off the string we wanna work with -- but how to find the part which is representative of it? - go back to where we extract the inbetweens
				parserTree = append(parserTree, ruleRowName)
				parserTree = append(parserTree, subTree)

				nonRegexElemCount++
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
