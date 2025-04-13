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
// TODO should extraneous flags error like extras?
// TODO update printUsage to mention -c config flag
// TODO Update DESCRIPTION in printUsage
// TODO finishing commenting other packages
//

// Microblog new item command. Insert the source code of a new microblog item
// at the top of the microblog HTML page. Will parse additional CLI flags
// and extras as part of the command. Returns a status code, success is 0.
func cmdMicroblogNew() int {
    // check if extra arguments were provided, and error
    if len(c.Extras) > 1 {
        fmt.Fprintf(os.Stderr, "[error] unrecognized arguments provided\n")
        return 1
    }
    // determine html path using one defined in config file or extra flag
    // extra flag should supersede config file entry
    var htmlPath string
    // config file is already parsed
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

    // open file for writing
    f, err := os.OpenFile(htmlPath, os.O_RDWR, 0644)
    if err != nil {
        fmt.Fprintf(os.Stderr, "[error] failed to open html file: %s\n", htmlPath)
        return 1
    }
    defer f.Close()

    // format datetime (if present)
    var datetime time.Time
    var datetimeStr string
    // parse date defined in -d or --date
    if v, ok := c.Flags["d"]; ok {
        datetimeStr = v    
    } else if v, ok := c.Flags["date"]; ok {
        datetimeStr = v    
    }
    if datetimeStr != "" {
        // YYYY-MM-DD-hh-mm-ss datetime, err = time.Parse("2006-01-02-15:04:05", datetimeStr)
        if err != nil {
            fmt.Fprintf(os.Stderr, "[error] invalid date provided: %s\n", err)
            return 1
        }
    } else {
        // if date not provided then use current
        datetime = time.Now()
    }
    // when parsing the date convert to current timezone to achieve +0100
    // in <time datetime> -- though this does not work? time.Now() does this
    // automatically
    //
    // datetime = datetime.In(time.Local)

    err = microblog.InsertNewPostFile(f, datetime)
    if err != nil {
        fmt.Fprintf(os.Stderr, "[error] failed to insert new post: %s\n", err)
        return 1
    }
    return 0
}

// Microblog generate RSS feed command. Using a provided microblog HTML file,
// generate an RSS feed from the semantic elements of each post and write to a
// file or stdout. Will parse additional CLI flags and extras as part of the
// command. Returns a status code, success is 0.
func cmdMicroblogGenrss() int {
    // check for extraneous extras
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

    // if provided override the config paths with provided arguments as flags
    // supersede the config file paths
    if len(c.Extras) >= 1 {
        htmlPath = c.Extras[0]
    } else if htmlPath == "" {
        // if no entry in config file or not provided then error as parsing is
        // required to generate RSS file
        fmt.Fprintf(os.Stderr, "[error] missing microblog html file\n")
        return 1
    }
    htmlPath = resolvePath(htmlPath)

    // output path is optional, default behavior on unprovided output path is
    // printing to stdout
    if len(c.Extras) == 2 {
        outputPath = c.Extras[1]
    }
    outputPath = resolvePath(outputPath)
    // '-' signifies stdout
    if outputPath == "-" {
        // no provided output path produces the generated rss feed being
        // printed to stdout
        outputPath = ""
    }

    // open the html file
    f, err := os.Open(htmlPath)
    if err != nil {
        fmt.Fprintf(os.Stderr, "[error] failed to open microblog file: %s\n", err)
        return 1
    }
    defer f.Close()

    // generate the final rss feed, returns a string containing feed
    s, err := microblog.GenRssFromFile(f)
    if err != nil {
        fmt.Fprintf(os.Stderr, "[error] failed to generate rss: %s\n", err)
        return 1
    }

    // default behavior for missing output path is print to stdout
    if outputPath == "" {
        // print and exit with success
        fmt.Println(s)
        return 0
    }

    // open or create the provided output path
    outputFile, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
    if err != nil {
        fmt.Fprintf(os.Stderr, "[error] failed to open output file: %s\n", err)
        return 1
    }
    defer outputFile.Close()

    // write the string to the file
    _, err = outputFile.WriteString(s)
    if err != nil {
        fmt.Fprintf(os.Stderr, "[error] failed to write generated rss to output file: %s\n", err)
        return 1
    }
    return 0
}

// Takes an absolute path and resolves it by replacing any `~` character at
// the start of the path with the user's home directory. Safe to pass an empty
// string to return an empty string. Returns the resolved path.
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

// CLI and config are parsed before command branching.
var c *cli.CLI
var conf *config.Config
func main() {
    // parse cli using helper function which breaks cli args into commmand,
    // subcommand, flags, and extra strings
    //
    // e.g. [command] [subcommand] [--flag] [extra] [extra] [extra]
    c = cli.Parse()
    if c == nil {
        fmt.Fprintf(os.Stderr, "[error] failed to parse command line arguments\n")
        return
    }
    // print usage when missing a command or help provided
    if c.Command == "" || c.Command == "help" {
        printUsage()
        return
    }

    // parsing the config file for command functions to utilize
    // use the default config path unless -c / --config flag provided
    var configPath = defaultConfigPath
    // check both -c and --config flags when determining custom config file
    // path
    if configFlag, ok := c.Flags["c"]; ok {
        configPath = configFlag
    } else if configFlag, ok = c.Flags["config"]; ok {
        configPath = configFlag
    }
    configPath = resolvePath(configPath)

    // open the config file path
    configFile, err := os.Open(configPath)
    if err != nil {
        // config file failed to open

        // if config FLAG, not configPath was provided then error...
        if configFlag != "" {
            fmt.Fprintf(os.Stderr, "[error] failed to open config file: %s\n", configPath)
            os.Exit(1)
            return
        }
        // ... else if no custom config was provided then config file does
        // not exist and that is fine to ignore
    } else {
        // config file open

        // read and parse using helper function
        conf, err = config.ReadFile(configFile) 
        if err != nil {
            configFile.Close()
            fmt.Fprintf(os.Stderr, "[error] failed to parse config file: %s\n", err)
            os.Exit(1)
            return
        }
        configFile.Close() 
    }

    // command branching
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

