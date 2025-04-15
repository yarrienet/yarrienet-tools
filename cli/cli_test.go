package cli

import (
    "os"
    "testing"
)

// Testing the Parse() functionality of the cli package which separates user
// provided command line arguments into a command, subcommand, flags, and
// arguments.
//
// Tests both short and long flags, empty and value flags.
func TestParse(t *testing.T) {
    // overriding os.Args
    os.Args = []string{"yarrienet", "command1", "subcommand1", "-xyz", "arg1", "--long", "value", "arg2", "-", "--", "---", "x", "--verbose", "--config", ""}
    // main cli arguments parse function
    cli := Parse()

    // expected values according to the cli args
    expectedCommand := "command1"
    expectedSubcommand := "subcommand1"
    expectedEmptyFlags := []string{"x", "y", "z", "verbose", "config"}
    expectedValueFlags := map[string]string{
        "long": "value",
    }
    expectedArguments := []string{"arg1", "arg2", "-", "--", "---", "x"}

    // testing command and subcommands
    if cli.Command != expectedCommand {
        t.Errorf("expected cli.Command to be '%s' not '%s'", expectedCommand, cli.Command)
    } 
    if cli.Subcommand != expectedSubcommand {
        t.Errorf("expected cli.Subcommand to be '%s' not '%s'", expectedSubcommand, cli.Subcommand)
    }

    // testing flag length
    flagsLen := len(cli.Flags)
    expectedFlagsLen := len(expectedEmptyFlags) + len(expectedValueFlags)
    if flagsLen != expectedFlagsLen {
        t.Fatalf("expected length of cli.Flags to be %d not %d -- expected: %v + %v, actual: %v -- struct: %v", expectedFlagsLen, flagsLen, expectedEmptyFlags, expectedValueFlags, cli.Flags, cli)
    }
    
    // testing empty flags, verbose should be empty
    for _, f := range expectedEmptyFlags {
        if v, ok := cli.Flags[f]; ok {
            if v != "" {
                t.Errorf("expected flag '%s' to contain an empty value, not '%s'", f, v)
            }
        } else {
            t.Errorf("missing flag '%s'", f)
        }
    }

    // testing flags with values
    for fk, fv := range expectedValueFlags {
        if v, ok := cli.Flags[fk]; ok {
            if fv != v {
                t.Errorf("value of flag '%s' expected '%s' not '%s'", fk, fv, v)
            }
        } else {
            t.Errorf("missing flag '%s'", fk)
        }
    }

    // testing arguments
    argsLen := len(cli.Arguments)
    expectedArgsLen := len(expectedArguments)
    if argsLen != expectedArgsLen {
        t.Fatalf("expected length of cli.Arguments to be %d not %d -- expected: %v, actual: %v -- struct: %v", expectedArgsLen, argsLen, expectedArguments, cli.Arguments, cli)
    }
    for i, v := range expectedArguments {
        if cli.Arguments[i] != v {
            t.Errorf("expected cli.Arguments[%d] to be '%s' not '%s'", i, v, cli.Arguments[i])
        }
    }
}

