package main

import (
    "yarrienet/cli"
    "fmt"
    "os"
)

func printUsage() {
    fmt.Println("USAGE")
    fmt.Println("  yarrienet [command] [subcommand]")
    fmt.Println("")
    fmt.Println("DESCRIPTION")
    fmt.Println("  Tool")
    fmt.Println("")
    fmt.Println("COMMANDS")
    fmt.Println("  microblog new <microblog file> [output file] [--date <rfc3339>]")
    fmt.Println("    Insert an empty post into the microblog HTML source code.")
    fmt.Println("")
    fmt.Println("  microblog genrss <microblog file> <rss file> [output file] [--url <base url>]")
    fmt.Println("    Generate an RSS feed using the microblog file.")
    fmt.Println("")
    fmt.Println("  help")
    fmt.Println("    Print usage information.")
}

func cmdMicroblogNew() {
    fmt.Println("cmd: microblog new")
}

func cmdMicroblogGenrss() {
    fmt.Println("cmd: microblog genrss")
}

func main() {
    var c = cli.Parse()
    if c.Command == "" {
        printUsage()
        return
    }

    switch c.Command {
    case "microblog":
        if c.Subcommand == "" {
            fmt.Fprintf(os.Stderr, "[error] microblog requires a subcommand\n")
            return
        }
        switch c.Subcommand {
            case "new":
               cmdMicroblogNew() 
            case "genrss":
                cmdMicroblogGenrss()
            default:
                fmt.Fprintf(os.Stderr, "[error] unknown microblog subcommand '%s'\n", c.Subcommand)
        }
    case "help":
        printUsage()
    default:
        fmt.Fprintf(os.Stderr, "[error] unknown command '%s'\n", c.Command)
    }
}

