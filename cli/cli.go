package cli

import (
    "os"
)

type CLI struct {
    Command string
    Subcommand string
    Flags map[string]string
    Extras []string
}

func Parse() *CLI {
    command := ""
    subcommand := ""
    var flags = make(map[string]string)
    var extras []string
    var flagAwaitingValue = ""

    for i := 1; i < len(os.Args); i++ {
        a := os.Args[i]
        if a[0] == '-' && len(a) > 1 {
            if flagAwaitingValue != "" {
                flags[flagAwaitingValue] = ""
            }
            if a[1] == '-' && len(a) > 2 {
                flagAwaitingValue = a[2:]
            } else if len(a[1:]) == 1 {
                flagAwaitingValue = a[1:]
            } else {
                extras = append(extras, a)
            }
        } else if flagAwaitingValue != "" {
            flags[flagAwaitingValue] = a
            flagAwaitingValue = ""
        } else if command == "" {
            command = a 
        } else if subcommand == "" {
            subcommand = a
        } else {
            extras = append(extras, a)
        }
    }
    if flagAwaitingValue != "" {
        flags[flagAwaitingValue] = ""
    }

    return &CLI{
        Command: command,
        Subcommand: subcommand,
        Flags: flags,
        Extras: extras,
    }
}

