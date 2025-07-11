package main

// run "go run main.go programmer.go" in your cmd to run this example

import (
	"fmt"
	"strconv"
)

// A test InputString

var InputString = "6 + 4 + 5 * 5 * 5 + 5 + 6"

// grammar - we specify grammar rules, which will be tested on the string

var Grammar = map[string]string{
	"ADD":      "ADD '\\+' ADD || MULT", // This is the root rule in our case; we can add more sub rule to the root rule, the parser will try them out and use the first suitable; subroots are separated by || characters; all elements of a root is separated by a single space
	"MULT":     "MULT '\\*' MULT || NUM",
	"NUM":      "'[0-9]+'",
	"~ignore~": "[ \t\n]+",
}

// functions for rules - we make functions, which we will add to the rules later on; these functions always get string inputs, as many as non-regex elements were specified in the rule (e.g. in NUM '\\+' ADD its two, NUM and ADD), and should always return a string

func addition(ruleNum int, allValuesToPass []string) string { // ruleNum property is just for the interpreterm allValuesToPass will hold all the values from the current rule

	returnString := ""

	switch ruleNum {
	case 0: // In case of the first sub rule, this part runs (NUM '\\+' ADD)

		numOne, _ := strconv.Atoi(allValuesToPass[0]) // value of NUM
		numTwo, _ := strconv.Atoi(allValuesToPass[1]) // value of ADD

		additionVal := numOne + numTwo

		returnString = strconv.Itoa(additionVal)

	case 1: // In case of the second subrule (NUM)
		returnString = allValuesToPass[0] // value of NUM

	}

	return returnString
}

func multiplication(ruleNum int, allValuesToPass []string) string { // ruleNum property is just for the interpreterm allValuesToPass will hold all the values from the current rule

	returnString := ""

	switch ruleNum {
	case 0: // In case of the first sub rule, this part runs (NUM '\\+' ADD)

		numOne, _ := strconv.Atoi(allValuesToPass[0]) // value of NUM
		numTwo, _ := strconv.Atoi(allValuesToPass[1]) // value of ADD

		additionVal := numOne * numTwo

		returnString = strconv.Itoa(additionVal)

	case 1: // In case of the second subrule (NUM)
		returnString = allValuesToPass[0] // value of NUM

	}

	return returnString
}

// rules assigned to grammar - just add the function name to the grammar rule

var Rules = map[string]func(ruleNum int, allValuesToPass []string) string{
	"ADD":  addition,
	"MULT": multiplication,
}

// test

func main() {

	// If we put a "~ignore~" part into our grammar rule, we should run the ignoreParts function to throw away unwanted characters - this takes our grammar rules and the input string as its two input and will come back with a polished string
	ignoredPartsInput := IgnoreParts(Grammar, InputString) // If we have an ~ignore~ property in the grammer, we can ignore specific parts from the inputstring

	// The parser function will make a parserTree, this we don't need for anything, but to put it into the interpreter later on - this function takes 4 inputs: the polished input string (if you do not want to ignore anything in the inputstring you can just put the raw string here), the root element of our grammar rules (this is the direct or indirect parent of every other rule), the grammar rules, and the regex separator (which specifies the separator character we write the regex in the rules between, by default it is single quotes)
	parserTreeInterface := Parser(ignoredPartsInput, "ADD", Grammar, RegexSeparator) // Instead "ADD" we will put the foremost element, as this is the root by which the parser reads the string

	// lexer, although the name is misleading, assigns the interpreter logic to the grammar rules - it takes the grammar rules and the parser tree as its two arguments
	fmt.Println(Lexer(Rules, parserTreeInterface))

	// The two functions below serve a decorator role, that is good for the preprocessing of the inputString before we give it to the parser

	fmt.Println(ChangeDecorator("AddFGdggDGGggfregg", ".gg", "gg", []string{"_"})) // This function changes a substr that matches a certain regex to the given string or strings, we give the inputstring, the pattern we wanna change, the pattern we wanna keep in the given substr, and the values we want to change the parts of the substring we don't want to keep to. If we add multiple values to the string array, every unkeeped part of the substring will change to the corresponding element (according to the number of the keeped parts). If there are more unkeeped parts than elements in the array, the function will use the last element for all the unkeeped parts that are over the limit.

	fmt.Println(InsertDecorator("AddFGdggDGGggfregg", [][]string{{".g", "g"}, {".", "G"}}, []string{"_", "-"})) // This function puts strings between certain substrings that match the given patterns. First we add the input string, then we specify as many pattern pairs as we want, and then all the string to be inserted inbetween the corresponding patterns.

}
