package main

import (
    "yarrienet/cli"
    "yarrienet/config"
    "yarrienet/microblog"
    "fmt"
    "os"
    "path/filepath"
    "time"
)

const defaultConfigPath = "~/.config/yarrienet.conf"

const usageInformation string = `USAGE
  yarrienet [command] [subcommand]

DESCRIPTION
  Tool

COMMANDS
  microblog new <microblog file> [-d / --date <yyyy-mm-ddThh-mm-ss>]
    Insert an empty post into the microblog HTML source code in place.

  microblog genrss <microblog file> <rss file> [output file] [--url <base url>]
    Generate an RSS feed using the microblog file.

  help
    Print usage information.`

func printUsage() {
    fmt.Println(usageInformation)
}

//
// TODO support - for stdin not just stdout
//

// Returns the status of the command.
func cmdMicroblogNew() int {
    // check if extra arguments were provided, and error
    // TODO should extraneous flags be treated in a similar way?
    if len(c.Extras) > 1 {
        fmt.Fprintf(os.Stderr, "[error] unrecognized arguments provided\n")
        return 1
    }
    var htmlPath string
    if conf != nil {
        htmlPath = conf.MicroblogHtmlFile
    }
    if len(c.Extras) == 1 {
        htmlPath = c.Extras[0]
    } else if htmlPath == "" {
        fmt.Fprintf(os.Stderr, "[error] missing html path\n")
        return 1
    }
    htmlPath = resolvePath(htmlPath)

    f, err := os.OpenFile(htmlPath, os.O_RDWR, 0644)
    if err != nil {
        fmt.Fprintf(os.Stderr, "[error] failed to open html file: %s\n", htmlPath)
        return 1
    }
    defer f.Close()

    // format datetime (if present)
    var datetime time.Time
    var datetimeStr string
    if v, ok := c.Flags["d"]; ok {
        datetimeStr = v    
    } else if v, ok := c.Flags["date"]; ok {
        datetimeStr = v    
    }
    if datetimeStr != "" {
        datetime, err = time.Parse("2006-01-02-15:04:05", datetimeStr)
        if err != nil {
            fmt.Fprintf(os.Stderr, "[error] invalid date provided: %s\n", err)
            return 1
        }
    } else {
        datetime = time.Now()
    }
    datetime = datetime.In(time.Local)

    err = microblog.InsertNewPostFile(f, datetime)
    if err != nil {
        fmt.Fprintf(os.Stderr, "[error] failed to insert new post: %s\n", err)
        return 1
    }
    return 0
}

// Returns the command status.
func cmdMicroblogGenrss() int {
    if len(c.Extras) > 2 {
        fmt.Fprintf(os.Stderr, "[error] unrecognized arguments provided\n")
        return 1
    }

    // use the config paths provided
    var htmlPath string
    var outputPath string
    if conf != nil {
        htmlPath = conf.MicroblogHtmlFile
        outputPath = conf.MicroblogRssFile
    }

    // override the config paths with provided arguments
    if len(c.Extras) >= 1 {
        htmlPath = c.Extras[0]
    } else if htmlPath == "" {
        fmt.Fprintf(os.Stderr, "[error] missing microblog html file\n")
        return 1
    }
    htmlPath = resolvePath(htmlPath)

    if len(c.Extras) == 2 {
        outputPath = c.Extras[1]
    }
    outputPath = resolvePath(outputPath)
    // - signifies stdout
    if outputPath == "-" {
        outputPath = ""
    }

    // open the html file
    f, err := os.Open(htmlPath)
    if err != nil {
        fmt.Fprintf(os.Stderr, "[error] failed to open microblog file: %s\n", err)
        return 1
    }
    defer f.Close()

    // generate the final rss
    s, err := microblog.GenRssFromFile(f)
    if err != nil {
        fmt.Fprintf(os.Stderr, "[error] failed to generate rss: %s\n", err)
        return 1
    }

    if outputPath == "" {
        // if no output file provided then output to cli
        fmt.Println(s)
        return 0
    }

    outputFile, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
    if err != nil {
        fmt.Fprintf(os.Stderr, "[error] failed to open output file: %s\n", err)
        return 1
    }
    defer outputFile.Close()

    _, err = outputFile.WriteString(s)
    if err != nil {
        fmt.Fprintf(os.Stderr, "[error] failed to write generated rss to output file: %s\n", err)
        return 1
    }
    return 0
}

func resolvePath(path string) string {
    if path[0] == '~' {
        home, err := os.UserHomeDir()
        if err != nil {
            return path
        }
        return filepath.Join(home, path[1:])
    }
    return path
}

var c *cli.CLI
var conf *config.Config
func main() {
    // parse cli
    c = cli.Parse()
    if c == nil {
        fmt.Fprintf(os.Stderr, "[error] failed to parse command line arguments\n")
        return
    }
    if c.Command == "" || c.Command == "help" {
        printUsage()
        return
    }

    var configPath = defaultConfigPath
    configFlag, ok := c.Flags["c"]
    if ok {
        configPath = configFlag
    } else {
        configFlag, ok = c.Flags["config"]
        if ok {
            configPath = configFlag
        }
    }
    configPath = resolvePath(configPath)

    configFile, err := os.Open(configPath)
    if err != nil {
        if configFlag != "" {
            fmt.Fprintf(os.Stderr, "[error] failed to open config file: %s\n", configPath)
            os.Exit(1)
            return
        }
        // else ignore...
    } else {
        conf, err = config.ReadFile(configFile) 
        if err != nil {
            configFile.Close()
            fmt.Fprintf(os.Stderr, "[error] failed to parse config file: %s\n", err)
            os.Exit(1)
            return
        }
        configFile.Close() 
    }

    switch c.Command {
    case "microblog":
        if c.Subcommand == "" {
            fmt.Fprintf(os.Stderr, "[error] microblog requires a subcommand\n")
            return
        }
        switch c.Subcommand {
            case "new":
                s := cmdMicroblogNew() 
                os.Exit(s)
            case "genrss":
                s := cmdMicroblogGenrss()
                os.Exit(s)
            default:
                fmt.Fprintf(os.Stderr, "[error] unknown microblog subcommand '%s'\n", c.Subcommand)
                os.Exit(1)
        }
    default:
        fmt.Fprintf(os.Stderr, "[error] unknown command '%s'\n", c.Command)
        os.Exit(1)
    }
}

