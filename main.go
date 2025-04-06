package main

import (
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
    fmt.Println("  microblog genrss <microblog file> [output file] [--url <base url>]")
    fmt.Println("    Generate an RSS feed using the microblog file.")
}

func cmdMicroblog() {

}

func main() {
    if len(os.Args) < 1 {
        printUsage()
        return
    }

    switch os.Args[0] {
    case "microblog":
        cmdMicroblog()
        return
    default:
        printUsage()
    }
}

