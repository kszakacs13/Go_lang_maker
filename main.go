package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var InputString = "6+6+5+5"

var Rules = map[string]string{
	"ADD": "NUM '\\+' ADD || NUM",
	"NUM": "'[0-9]+'",
}

func main() {

	parserTreeInterface := parser(InputString, "ADD", Rules) // Instead "ADD" we will put the foremost element

	parserTreeToDisplay := displayParserTree(parserTreeInterface)

	for _, row := range parserTreeToDisplay {

		fmt.Println(row)

	}

}

func parser(expressionString string, ruleRowName string, rules map[string]string) []interface{} {

	var parserTree []interface{}

	ruleRow := rules[ruleRowName]

	rulesInRow := strings.Split(ruleRow, " || ")

	// Going through all rules in a row and extract the operators belonging to them
	ruleNum := 0
	ruleNumForActualRow := ""

	for ruleNum < len(rulesInRow) {

		ruleElements := strings.Split(rulesInRow[ruleNum], " ")

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

		ruleNumForActualRow = strconv.Itoa(ruleNum) // To save the actual rule number for later, when we want to save it to the parser tree

		// If we did not find a regex specified in the rule, we continue to the next rule in the row, however otherwise we do not want to get the next rule

		if someRegexNotFound {
			ruleNum++
			continue
		} else {
			ruleNum = len(rulesInRow) // Terminate the loop - we did not use break keyword, cause then i would need to include all the code below into this brach of the if - this would not appeal to me
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

			parserTree = append(parserTree, ruleRowName) // Rule name for the lexer to work

			if allElements[i][1] == "non-regex" {
				subTree := parser(extractedPositions[nonRegexElemCount], allElements[i][0], rules) // we ought to cut off the string we wanna work with -- but how to find the part which is representative of it? - go back to where we extract the inbetweens

				parserTree = append(parserTree, subTree[0]) // Add the number of the rule (again for the lexer)

				parserTree = append(parserTree, subTree[1])

				nonRegexElemCount++
			} else {

				// If there are non-regex elements present as well, we mark the element with an operator flag - this will be used in the lexer only when there are no rules

				if foundNonRegex {
					parserTree = append(parserTree, allElements[i][0])
					parserTree = append(parserTree, "OP-REGEX")
				} else {
					parserTree = append(parserTree, allElements[i][0])
				}

			}

		}

	}

	// Return the value with the number of rule in the row

	var returnInterface []interface{}

	returnInterface = append(returnInterface, ruleNumForActualRow)
	returnInterface = append(returnInterface, parserTree)

	return returnInterface

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

// Function to display parserTree

func displayParserTree(parserTree []interface{}) []string {

	displaySlice := []string{}

	// If it contains an interface we pass it down one level lower into the recursion

	for _, elem := range parserTree {

		switch v := elem.(type) {

		case []interface{}:

			subInterfaceBreak := displayParserTree(v)

			for _, subElem := range subInterfaceBreak {

				// We also add tabs to all slice elements we got in return

				displaySlice = append(displaySlice, "\t"+subElem)

			}

		case string:

			// If it does not contain any interfaces, only strings we return those strings with no tab

			displaySlice = append(displaySlice, v)

		}

	}

	return displaySlice

}

// LEXER: goes through the interface, gets into the deepest object recursively, and if the elements are not interfaces anymore, applies the appropriate function to the rule

// DO NOT NEED TO ADD THE OP-REGEX ELEMENT IF THERE IS A RULE PARSED TO THE GRAMMER

// Also: check which rule applies in the row - func ruleForGramm(row string, ruleNum int, inputVars []string) - and count the non regex elemnts from parsing tree
