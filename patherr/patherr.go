// Copyright (c) 2019, AT&T Intellectual Property.
// All rights reserved.
//
// Copyright (c) 2014-2016 by Brocade Communications Systems, Inc.
// All rights reserved.
//
// SPDX-License-Identifier: MPL-2.0

package patherr

import (
	"bytes"
	"fmt"
	"github.com/danos/utils/natsort"
	"strings"
	"text/tabwriter"
    "runtime"
    "os"
    "time"
)

const cfgPath = `Configuration path: `
const operCommand = `Command: `

type CommandInval struct {
	Path []string
	Fail string
}

func printStackToFile() {
    // Create a slice to hold the stack trace
    const size = 1024
    buf := make([]byte, size)
    
    // Capture the stack trace
    n := runtime.Stack(buf, true)
    
    // Define the file name with a timestamp to avoid overwriting
    fileName := fmt.Sprintf("/tmp/stacktrace_%d.log", time.Now().Unix())
    
    // Open the file for writing (creates a new file or truncates an existing file)
    file, err := os.Create(fileName)
    if err != nil {
        fmt.Printf("Error creating file: %v\n", err)
        return
    }
    defer file.Close() // Ensure the file is closed when the function exits
    
    // Write the stack trace to the file
    _, err = file.Write(buf[:n])
    if err != nil {
        fmt.Printf("Error writing to file: %v\n", err)
        return
    }
    
    // Optionally print the file path to the console
    fmt.Printf("Stack trace written to: %s\n", fileName)
}

func (e *CommandInval) Error() string {
	if len(e.Path) == 0 {
        printStackToFile()
		return fmt.Sprintf("EZ3: Invalid command: [%s]", e.Fail)
	}
	return fmt.Sprintf("EZ2: Invalid command: %s [%s]", strings.Join(e.Path, " "), e.Fail)
}

type PathInval struct {
	Path        []string
	Fail        string
	Operational bool
}

func (e *PathInval) Error() string {
	prefix := cfgPath
	if e.Operational {
		prefix = operCommand
	}
	if len(e.Path) == 0 {
		return fmt.Sprintf("%s [%s] is not valid", prefix, e.Fail)
	}
	return fmt.Sprintf("%s %s [%s] is not valid", prefix, strings.Join(e.Path, " "), e.Fail)
}

type PathAmbig struct {
	Path        []string
	Fail        string
	Matches     map[string]string
	Operational bool
}

func (e *PathAmbig) Error() string {
	var buf = new(bytes.Buffer)
	twriter := tabwriter.NewWriter(buf, 8, 0, 1, ' ', 0)

	prefix := cfgPath
	if e.Operational {
		prefix = operCommand
	}

	if len(e.Path) == 0 {
		fmt.Fprintf(buf, "%s [%s] is ambiguous\n", prefix, e.Fail)
	} else {
		fmt.Fprintf(buf, "%s %s [%s] is ambiguous\n", prefix, strings.Join(e.Path, " "), e.Fail)
	}
	fmt.Fprintf(buf, "\n  EZ: Possible completions:\n")

	sorted := make([]string, 0, len(e.Matches))
	for n, _ := range e.Matches {
		sorted = append(sorted, n)
	}

	natsort.Sort(sorted)
	for i, name := range sorted {
		fmt.Fprintf(twriter, "    %s\t%s", name, e.Matches[name])
		if i != len(sorted)-1 {
			fmt.Fprintf(twriter, "\n")
		}
	}
	twriter.Flush()
	return buf.String()
}
