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

func printUsage() {
    fmt.Println("USAGE")
    fmt.Println("  yarrienet [command] [subcommand]")
    fmt.Println("")
    fmt.Println("DESCRIPTION")
    fmt.Println("  Tool")
    fmt.Println("")
    fmt.Println("COMMANDS")
    fmt.Println("  microblog new <microblog file> [-d / --date <yyyy-mm-ddThh-mm-ss>]")
    fmt.Println("    Insert an empty post into the microblog HTML source code in place.")
    fmt.Println("")
    fmt.Println("  microblog genrss <microblog file> <rss file> [output file] [--url <base url>]")
    fmt.Println("    Generate an RSS feed using the microblog file.")
    fmt.Println("")
    fmt.Println("  help")
    fmt.Println("    Print usage information.")
}

func cmdMicroblogNew() {
    // check if extra arguments were provided, and error
    // TODO should extraneous flags be treated in a similar way?
    if len(c.Extras) > 1 {
        fmt.Fprintf(os.Stderr, "[error] unrecognized arguments provided\n")
        os.Exit(1)
        return
    }
    var htmlPath string
    if conf != nil {
        htmlPath = conf.MicroblogHtmlFile
    }
    if len(c.Extras) == 1 {
        htmlPath = c.Extras[0]
    } else if htmlPath == "" {
        fmt.Fprintf(os.Stderr, "[error] missing html path\n")
        os.Exit(1)
        return
    }
    htmlPath = resolvePath(htmlPath)

    f, err := os.OpenFile(htmlPath, os.O_RDWR, 0644)
    if err != nil {
        fmt.Fprintf(os.Stderr, "[error] failed to open html file: %s\n", htmlPath)
        os.Exit(1)
        return
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
            os.Exit(1)
            return
        }
    } else {
        datetime = time.Now()
    }
    datetime = datetime.In(time.Local)

    // TODO parse the date flag
    err = microblog.InsertNewPostFile(f, datetime)
    if err != nil {
        fmt.Fprintf(os.Stderr, "[error] failed to insert new post: %s\n", err)
        os.Exit(1)
        return
    }
    // done
}

func cmdMicroblogGenrss() {
    fmt.Println("cmd: microblog genrss")
}

func resolvePath(path string) string {
    if path[0] == '~' {
        home, err := os.UserHomeDir()
        if err != nil {
            return path
        }
        path := filepath.Join(home, path[1:])
        return path
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
    if c.Command == "" {
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
        fmt.Println("[debug] success opened config file")
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

