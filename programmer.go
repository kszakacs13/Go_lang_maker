package main

import (
	"regexp"
	"strconv"
	"strings"
)

// THE FUNCTIONAL PART

var RegexSeparator = "'" // Changeable

// TODO: be able to omit some characters in the rules by specifying them under the name ~omit~ - func omitChars(Grammer)

func Parser(expressionString string, ruleRowName string, rules map[string]string, regexSeparator string) []interface{} {

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

			if strings.Contains(element, regexSeparator) {
				regex := betweenSigns(element, regexSeparator)

				pattern := regexp.MustCompile(regex)

				// Searching if the REGEX occurs in the expressionString

				positionOfRegex := pattern.FindAllStringIndex(expressionString, -1)

				if positionOfRegex == nil {

					someRegexNotFound = true // Continue to the next rule in the row if no matches found

				} else {

					// Loop through all the occurances of the regex pattern, loop through all the downstream elements, check if any applies to the pattern better, go till we find a pattern that best matches ours, if we don't then continue to next loop by setting someRegexNotFound to true

					downStreamElems := getDownStreamElements(rules, ruleRowName, regexSeparator)

					found_best_match := false
					best_match := []int{}

					for _, regexPos := range positionOfRegex {

						regexUsable := true

						for _, dsElem := range downStreamElems {

							// Check if any applies to the pattern

							dsElemRow := rules[dsElem]

							if strings.Contains(dsElemRow, regex) && !found_best_match {

								allBlocks := strings.Split(dsElemRow, " ")

								for _, block := range allBlocks {

									if strings.Contains(block, regexSeparator) && !found_best_match {

										sub_regex := betweenSigns(block, regexSeparator)
										sub_pattern := regexp.MustCompile(sub_regex)

										positionOfSubRegex := sub_pattern.FindAllStringIndex(expressionString, -1)

										for _, sub_regex_pos := range positionOfSubRegex {
											if (regexPos[0] >= sub_regex_pos[0] && regexPos[0] <= sub_regex_pos[1]) ||
												(regexPos[1] >= sub_regex_pos[0] && regexPos[1] <= sub_regex_pos[1]) {
												regexUsable = false
											}
										}

									}

								}

							}

						}

						if regexUsable {
							found_best_match = true
							best_match = regexPos
						}

						if found_best_match {
							break
						}

					}

					if found_best_match {
						allElements = append(allElements, []string{expressionString[best_match[0]:best_match[1]], "regex", strconv.Itoa(best_match[0]), strconv.Itoa(best_match[1])}) // We save the positions to extract the important elements later
					} else {
						someRegexNotFound = true
					}

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
				subTree := Parser(extractedPositions[nonRegexElemCount], allElements[i][0], rules, regexSeparator) // we ought to cut off the string we wanna work with -- but how to find the part which is representative of it? - go back to where we extract the inbetweens

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

func getDownStreamElements(rules map[string]string, element string, separatorChar string) []string {

	return_arr := []string{}

	all_appeared := false

	this_rows_elements := []string{element}

	for !all_appeared {

		last_rows_elements := this_rows_elements

		this_rows_elements = []string{}

		for _, row := range last_rows_elements {

			ruleSides := strings.Split(rules[row], " || ")

			for _, ruleSide := range ruleSides {

				blocks := strings.Split(ruleSide, " ")

				for _, block := range blocks {

					// If the actual block is not yet in return_arr, we add it

					if !arr_contains(return_arr, block) && !strings.Contains(block, separatorChar) {
						return_arr = append(return_arr, block)
						this_rows_elements = append(this_rows_elements, block) // At the same time we add it to the elements we wanna find next loop
					}

				}

			}

		}

		if len(this_rows_elements) == 0 { // This means we did not find any unchecked element in this loop
			all_appeared = true

			return_arr = arr_remove(return_arr, element)
		}

	}

	return return_arr

}

func arr_contains(arr []string, value string) bool {

	contains := false

	for _, i := range arr {
		if i == value {
			contains = true
		}
	}

	return contains

}

func arr_remove(arr []string, value string) []string {

	return_arr := []string{}

	for _, i := range arr {
		if i != value {
			return_arr = append(return_arr, i)
		}
	}

	return return_arr
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

func Lexer(rules map[string]func(ruleNum int, allValuesToPass []string) string, parserTree []interface{}) string {

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
				returnedElem := Lexer(rules, getSubInterface(interf, 2))

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

// Ignore special

func IgnoreParts(grammar map[string]string, inputString string) string {
	re := regexp.MustCompile(grammar["~ignore~"])
	outputExpression := re.ReplaceAllString(inputString, "")

	return outputExpression
}

// Decorate the string as wanted

func ChangeDecorator(inputString string, regexToFind string, regexToKeep string, changeTo []string) string { // We slice the found string by the regexToKeep, and add the elements of the changeTo inbetween, if there is too few however, we will replicate the last element (or recurse if the user want so)

	decoratedString := ""

	cursorLeftAt := 0

	// Go through all the found regexes

	re := regexp.MustCompile(regexToFind)
	allFirstLevelRegexes := re.FindAllStringIndex(inputString, -1)

	for _, regexPatt := range allFirstLevelRegexes {

		actualSubStr := inputString[regexPatt[0]:regexPatt[1]]

		// First we add all the unrelated strings to the output

		decoratedString += inputString[cursorLeftAt:regexPatt[0]]
		cursorLeftAt = regexPatt[1]

		// In one part we check for all the regexes we should keep

		re2 := regexp.MustCompile(regexToKeep)
		allKeepRegexs := re2.FindAllStringIndex(actualSubStr, -1)

		decoratedSubStr := ""

		subCursor := 0

		for i, shouldKeep := range allKeepRegexs {

			indChangeTo := i

			if i > len(changeTo)-1 {
				indChangeTo = len(changeTo) - 1
			}

			// Where we are not inbetween we change the str
			if shouldKeep[0] > subCursor {
				decoratedSubStr += changeTo[indChangeTo] + actualSubStr[shouldKeep[0]:shouldKeep[1]]
			} else {
				decoratedSubStr += actualSubStr[shouldKeep[0]:shouldKeep[1]]
			}

			subCursor = shouldKeep[1]

		}

		lastRegexIndex := allKeepRegexs[len(allKeepRegexs)-1][1]

		if lastRegexIndex < len(actualSubStr) {
			decoratedSubStr += actualSubStr[lastRegexIndex:]
		}

		decoratedString += decoratedSubStr

	}

	return decoratedString
}

func InsertDecorator(inputString string, regexToFind [][]string, insertValues []string) string { // We add as many regexes we want in pairs. The decorator goes through them and wherever it finds them, inserts the same index value from the insert values

	decoratedString := ""

	for ind, regexPair := range regexToFind {

		decoratedString = ""

		re1 := regexp.MustCompile(regexPair[0])
		re2 := regexp.MustCompile(regexPair[1])

		allMatch1 := re1.FindAllStringIndex(inputString, -1)
		allMatch2 := re2.FindAllStringIndex(inputString, -1)

		// Find all the pairs (which start and end at the same direction)

		var insertIndexes []int

		for _, match1 := range allMatch1 {
			for _, match2 := range allMatch2 {
				if match1[1] == match2[0] {
					insertIndexes = append(insertIndexes, match1[1])
				}
			}
		}

		cursor := 0

		for _, insertInd := range insertIndexes {
			decoratedString += inputString[cursor:insertInd] + insertValues[ind]
			cursor = insertInd
		}

		if cursor < len(inputString) {
			decoratedString += inputString[cursor:]
		}

		inputString = decoratedString

	}

	return decoratedString

}
