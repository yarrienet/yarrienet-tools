package cli

import (
    "os"
)

// Structure containing the parsed result of the given command line arguments.
// Separated by the first two arguments (command and subcommand respectively),
// extra arguments in an array, then flags and their optional value.
type CLI struct {
    // The first argument.
    Command string
    // The second argument.
    Subcommand string
    // Map containing each included flag and its associated value. Both single
    // and double dash flags are supported (e.g. -v vs. --verbose). Multiple
    // characters following a single dash will each appear in the map.
    //
    // A flag can be present without containing a value, this is represented
    // with an empty string. When using a flag value do not only check the
    // presence of the flag, you must always confirm the string value is not
    // empty.
    //
    // Single or double dashes with no characters following will be parsed as
    // arguments and can be safely expected, this is a common pattern for
    // representing stdout or stdin.
    Flags map[string]string
    // Extra arguments after the command and subcommand. Commonly used for
    // files.
    Arguments []string
}

// Parse the given command line arguments into a CLI structure which contains
// each parsed element. There is no schema logic and error handling for invalid
// flags, commands and arguments should be handled after the parse by the code
// that called it.
func Parse() *CLI {
    var command string
    var subcommand string
    var flags = make(map[string]string)
    var arguments []string
    // flag name awaiting value, used to track until found value and can be
    // added to map.
    var flagAwaitingValue = ""

    // loop each word (space separated)
    for i := 1; i < len(os.Args); i++ {
        a := os.Args[i]
        if a[0] == '-' && len(a) > 1 {
            // if word begins with a dash, most likely a flag

            // one flag following another means that the previous does not
            // contain a value, add to map without a value and continue parsing
            if flagAwaitingValue != "" {
                flags[flagAwaitingValue] = ""
            }

            if a[1] == '-' && len(a) > 2 {
                // if first dash is followed by another, seek the value
                // after flag
                flagAwaitingValue = a[2:]
            } else if len(a) > 1 {
                // if only a single dash with characters following, loop each
                // character after dash and add to map of present flags
                for _, f := range a[1:] {
                    flags[string(f)] = ""    
                }
            } else {
                // if no characters follow dash, then add to extra arguments.
                // a single - is a credible argument.
                arguments = append(arguments, a)
            }
        } else if flagAwaitingValue != "" {
            
            flags[flagAwaitingValue] = a
            flagAwaitingValue = ""
        } else if command == "" {
            // first non-flag word encountered is command
            command = a 
        } else if subcommand == "" {
            // second non-flag word encountered is the subcommand, determined
            // when command is not empty
            subcommand = a
        } else {
            // non-flag word when command and subcommand is determined is an
            // argument
            arguments = append(arguments, a)
        }
    }
    // word loop ended before flag found its value, add to map as present but
    // without value
    if flagAwaitingValue != "" {
        flags[flagAwaitingValue] = ""
    }

    return &CLI{
        Command: command,
        Subcommand: subcommand,
        Flags: flags,
        Arguments: arguments,
    }
}

