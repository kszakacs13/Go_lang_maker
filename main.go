package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// InputString

var InputString = "6+4+5+5"

// grammer

var Grammer = map[string]string{
	"ADD": "NUM '\\+' ADD || NUM",
	"NUM": "'[0-9]+'",
}

// functions for rules

func addition(ruleNum int, allValuesToPass []string) string {

	returnString := ""

	switch ruleNum {
	case 0:

		numOne, _ := strconv.Atoi(allValuesToPass[0])
		numTwo, _ := strconv.Atoi(allValuesToPass[1])

		additionVal := numOne + numTwo

		returnString = strconv.Itoa(additionVal)

	case 1:
		returnString = allValuesToPass[0]

	}

	return returnString
}

// rules assigned to grammer

var Rules = map[string]func(ruleNum int, allValuesToPass []string) string{
	"ADD": addition,
}

// test

func main() {

	parserTreeInterface := parser(InputString, "ADD", Grammer) // Instead "ADD" we will put the foremost element

	/* DISPLAY PASER TREE

	parserTreeToDisplay := displayParserTree(parserTreeInterface)

	for _, row := range parserTreeToDisplay {

		fmt.Println(row)

	}

	*/

	fmt.Println(lexer(Rules, parserTreeInterface))

}

// THE FUNCTIONAL PART

// TODO: be able to omit some characters in the rules by specifying them under the name ~omit~ - func omitChars(Grammer)

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

			var subParserTree []interface{}

			subParserTree = append(subParserTree, ruleRowName) // Rule name for the lexer to work

			subParserTree = append(subParserTree, ruleNumForActualRow) // Add the number of the rule (again for the lexer)

			if allElements[i][1] == "non-regex" {
				subTree := parser(extractedPositions[nonRegexElemCount], allElements[i][0], rules) // we ought to cut off the string we wanna work with -- but how to find the part which is representative of it? - go back to where we extract the inbetweens

				subParserTree = append(subParserTree, subTree)

				nonRegexElemCount++
			} else {

				// If there are non-regex elements present as well, we mark the element with an operator flag - this will be used in the lexer only when there are no rules

				if foundNonRegex {
					subParserTree = append(subParserTree, allElements[i][0])
					subParserTree = append(subParserTree, "OP-REGEX")
				} else {
					subParserTree = append(subParserTree, allElements[i][0])
				}

			}

			parserTree = append(parserTree, subParserTree)

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

func lexer(rules map[string]func(ruleNum int, allValuesToPass []string) string, parserTree []interface{}) string {

	returnedSolution := ""

	ruleName := ""

	ruleNum := -1

	// Check for the grammar name in each elment of the row

	for id, i := range parserTree {

		switch interf := i.(type) {
		case []interface{}:

			if id == 0 {

				ruleName = getInterfaceElement(interf, 0)

				ruleNum, _ = strconv.Atoi(getInterfaceElement(interf, 1))

			}

			// If the third elemnt is an iterface, loop further

			if checkIfInterfaceElem(interf, 2) {
				returnedElem := lexer(rules, getSubInterface(interf, 2))

				parserTree[id] = changeInterfaceElem(interf, 2, returnedElem)
			}

		}

	}

	// Here everything should be assigned a value, not an interface

	// If there is a rule we omit the OP-REGEX types

	containsRule := false

	if _, ok := rules[ruleName]; ok {
		containsRule = true
	} else {
		containsRule = false
	}

	allValuesToPass := []string{}

	for _, i := range parserTree {

		// Return the value of the string - or if there is a function for it return the val of that function with giving in the non OP-REGEX variables

		switch interf := i.(type) {
		case []interface{}:

			if containsRule {

				if lenInterface(interf) < 4 { // Smaller than 4 if not OP-REGEX
					allValuesToPass = append(allValuesToPass, getInterfaceElement(interf, 2))
				}

			} else {

				allValuesToPass = append(allValuesToPass, getInterfaceElement(interf, 2))

			}
		}

	}

	// If we have a funciton to this we return the values with that

	if containsRule {
		returnedSolution = rules[ruleName](ruleNum, allValuesToPass)
	} else {
		returnedSolution = strings.Join(allValuesToPass, "")
	}

	return returnedSolution

}

func getInterfaceElement(inputInterface []interface{}, elementNum int) string {

	returnVal := ""

	for index, i := range inputInterface {
		if index == elementNum {
			switch v := i.(type) {
			case string:
				returnVal = v
			}
		}
	}

	return returnVal

}

func getSubInterface(inputInterface []interface{}, elementNum int) []interface{} {

	var returnVal []interface{}

	for index, i := range inputInterface {
		if index == elementNum {
			switch v := i.(type) {
			case []interface{}:
				returnVal = v
			}
		}
	}

	return returnVal

}

func checkIfInterfaceElem(inputInterface []interface{}, elementNum int) bool {

	var returnVal bool

	for index, i := range inputInterface {
		if index == elementNum {
			_, ok := i.([]interface{})
			returnVal = ok
		}
	}

	return returnVal

}

func changeInterfaceElem(inputInterface []interface{}, elementNum int, assignNewVal string) []interface{} {

	var returnVal []interface{}

	for index, i := range inputInterface {
		if index == elementNum {
			_, ok := i.([]interface{})
			if ok {
				returnVal = append(returnVal, assignNewVal)
			} else {
				returnVal = append(returnVal, i)
			}
		} else {
			returnVal = append(returnVal, i)
		}

	}

	return returnVal

}

func lenInterface(inputInterface []interface{}) int {

	returnLen := 0

	for ind, _ := range inputInterface {
		returnLen = ind + 1
	}

	return returnLen

}
