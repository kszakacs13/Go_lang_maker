## Generic Recursive Parser & Interpreter (Go)

This repository contains a lightweight program by which you can make your **own parser and interpreter** for *any grammar you define*.

It’s designed so you can: Define grammar rules\
Attach interpreter functions to them\
Parse and evaluate complex expressions (arithmetic, logical, or domain-specific)

---

## Quick start

Clone the repo and run:

```bash
go run main.go programmer.go
```

*(**`programmer.go`** contains all the functions needed, so you can just import it to your project)*

---

## Overview

This parser lets you define:

- **Grammar**: which describes the structure of valid expressions.
- **Interpreter functions**: which describe *how to process* the parsed parts.

Together, they let you parse and process anything you can express as a recursive grammar:

- Arithmetic (`6 + 5`)
- Boolean expressions (`true && false`)
- Domain-specific languages
- Custom mini-languages

---

## How it works

### Define your grammar

In `Grammar` map:

```go
var Grammar = map[string]string{
	"ADD":      "ADD '\+' ADD || MULT",
	"MULT":     "MULT '\*' MULT || NUM",
	"NUM":      "'[0-9]+'",
	"~ignore~": "[ \t\n]+",
}
```

*Subrules* are separated by `||`.\
Regex parts are written between quotes (e.g., `'\+'`) or what else you define as the separator.

You can define completely different grammars to parse anything, not just arithmetic.

---

### Attach interpreter functions

Create functions that describe **what to do** when each rule matches.

Example:

```go
func addition(ruleNum int, allValuesToPass []string) string {
	// ruleNum is which subrule matched
	// allValuesToPass contains the parsed child values
}
```

Then map them to your rules:

```go
var Rules = map[string]func(ruleNum int, allValuesToPass []string) string{
	"ADD":  addition,
	"MULT": multiplication,
}
```

---

### Run the parser & interpreter

In `main()`:

```go
ignoredPartsInput := IgnoreParts(Grammar, InputString)
parserTree := Parser(ignoredPartsInput, "ADD", Grammar, RegexSeparator)
result := Lexer(Rules, parserTree)

fmt.Println(result) // Output: evaluated result
```

---

## Example

Input string:

```go
var InputString = "6 + 4 + 5 * 5 * 5 + 5 + 6"
```

- Parser builds a tree based on the grammar.
- Interpreter functions compute the result.

> By changing the grammar & functions, you can parse and evaluate any language.

---

## Extra utilities

This project includes decorator functions:

- `ChangeDecorator(...)`: replace matched substrings.
- `InsertDecorator(...)`: insert strings between substrings matching patterns.

They’re useful for preprocessing, formatting, or transforming input strings.

---

## Extend it to your language

You can define:

- New grammar rules for things like `AND`, `OR`, keywords, functions, etc.
- Interpreter functions that do anything: evaluate, transform, compile.

For example:

```go
"LOGIC": "LOGIC '&&' LOGIC || BOOL",
"BOOL": "'true' || 'false'",
```

---

## Running summary

```bash
go run main.go programmer.go
```

Output includes:

- Evaluated result of your expression
- Examples of decorator utilities

---

## Why this is useful

- Build interpreters for small DSLs
- Teach recursive parsing
- Prototype simple compilers
- Learn grammar-driven design

---

## License

Feel free to use, extend, and share!\


