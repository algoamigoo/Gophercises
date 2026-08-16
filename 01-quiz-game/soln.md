# Go Standard Library Reference: File I/O, Text Processing & CLI

This document provides a consolidated reference of the Go standard library packages and methods used in the CSV quiz application, along with closely related alternatives.

---

## 1. `os` (File & System I/O)

| Method / Symbol | Usage Type | Purpose |
| :--- | :--- | :--- |
| **`os.Open(name)`** | **Used** | Opens a file for reading; returns an `*os.File` pointer and an `error`. |
| **`file.Close()`** | **Used** | Closes the open file handle; usually paired with `defer file.Close()`. |
| `os.ReadFile(name)` | *Related* | Reads an entire file directly into memory as a `[]byte` slice. |
| `os.WriteFile(name, data, perm)` | *Related* | Writes raw byte data to a file, creating or overwriting it. |
| `os.Create(name)` | *Related* | Creates or truncates a file for writing. |

---

## 2. `encoding/csv` (CSV Processing)

| Method / Symbol | Usage Type | Purpose |
| :--- | :--- | :--- |
| **`csv.NewReader(r)`** | **Used** | Initializes a CSV reader wrapping any `io.Reader` (e.g., `*os.File`). |
| **`reader.Read()`** | **Used** | Reads a single line/record as a slice of strings (`[]string`). |
| `reader.ReadAll()` | *Related* | Reads all remaining records at once into a 2D slice (`[][]string`). |
| `csv.NewWriter(w)` | *Related* | Initializes a CSV writer to output formatted CSV lines. |

---

## 3. `io` & `bufio` (Stream & Line Input)

| Method / Symbol | Usage Type | Purpose |
| :--- | :--- | :--- |
| **`io.EOF`** | **Used** | Sentinel error value indicating the End Of File was reached during reading. |
| `io.ReadAll(r)` | *Related* | Reads all bytes from an `io.Reader` stream until EOF. |
| `bufio.NewScanner(r)` | *Related* | Creates a scanner for reading plain text files line-by-line (`scanner.Scan()`). |

---

## 4. `fmt` (Formatting & Console I/O)

| Method / Symbol | Usage Type | Purpose |
| :--- | :--- | :--- |
| **`fmt.Printf(format, a...)`** | **Used** | Formats text using specifiers (`%d`, `%s`, `%v`) and writes to stdout. |
| **`fmt.Println(a...)`** | **Used** | Formats using default rules and writes to stdout with a newline. |
| **`fmt.Scan(a...)`** | **Used** | Reads space-delimited text from standard input (`os.Stdin`). |
| `fmt.Scanln(a...)` | *Related* | Reads space-delimited text from standard input, stopping at a newline. |
| `fmt.Sprintf(format, a...)` | *Related* | Formats string data like `Printf`, but returns a string instead of printing. |

---

## 5. `strings` & `strconv` (Text Manipulation & Type Conversion)

| Method / Symbol | Usage Type | Purpose |
| :--- | :--- | :--- |
| **`strings.TrimSpace(s)`** | **Used** | Trims all leading and trailing whitespace from a string. |
| **`strings.ToLower(s)`** | **Used** | Converts all characters in a string to lowercase. |
| `strconv.Atoi(s)` | *Related* | Converts string digits (e.g., `"42"`) into an integer (`int`). |
| `strconv.Itoa(i)` | *Related* | Converts an integer (`int`) into its string representation. |

---

## 6. `log` & `encoding/json` (Logging & Serialization)

| Method / Symbol | Usage Type | Purpose |
| :--- | :--- | :--- |
| **`log.Fatalf(format, v...)`** | **Used** | Prints formatted log messages to stderr and exits the program immediately (`os.Exit(1)`). |
| `json.NewDecoder(r).Decode(v)` | *Related* | Parses structured JSON data directly from a file stream into a struct/map. |

---

## 7. `flag` (Command-Line Flags)

### What is the `flag` package?

The `flag` package is part of Go’s standard library.  
It provides a simple and idiomatic way to define and parse command-line flags (also called options or switches).

Instead of manually inspecting `os.Args`, you declare the flags you need, give them default values and help text, and then call `flag.Parse()`. After parsing, the package automatically fills in the values.

It is the standard way in Go to accept configuration from the command line.

### Basic Usage Pattern

```go
csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
timeLimit  := flag.Int("limit", 15, "the time limit for the quiz in seconds")
shuffle    := flag.Bool("shuffle", false, "shuffle the quiz questions order")

flag.Parse()

// Use the values (they are pointers)
fmt.Println(*csvFilename)
fmt.Println(*timeLimit)
fmt.Println(*shuffle)