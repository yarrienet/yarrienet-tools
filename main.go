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
const defaultBaseUrl = "http://yarrie.net/microblog"

const usageInformation string = `USAGE
  yarrienet <command> [<subcommand>] [-h | --help] [-c | --config <config>]

DESCRIPTION
  These tools are built to achieve semantic publishing where writing occurs directly in the
  webpage's HTML source and then all auto-generated elements are built against that. In practice a
  single webpage containing each post can be used to produce the associated RSS feed without the
  backing of an abstract file tree or database.
  
  HTML does not need to be abstracted to produce working RSS unlike other static site generators.
  
  The tool relies heavily on a set schema utilizing semantic elements and as such tools are
  included to insert the HTML source of schema friendly elements directly into the webpage to be
  modified.
  
COMMANDS
  microblog new <microblog file> [-d | --date <YYYY-MM-DD-hh-mm-ss>]
    Insert an empty post into the microblog HTML source code in place.

  microblog genrss <microblog file> [<output rss>] [--url <base url>] 
    Generate an RSS feed using the microblog file. Omitting output or using '-' will print the
    generated RSS feed to stdout.

  help
    Print usage information.`

func printUsage() {
    fmt.Println(usageInformation)
}

//
// TODO finishing commenting other packages
//

// Microblog new item command. Insert the source code of a new microblog item
// at the top of the microblog HTML page. Will parse additional CLI flags
// and extras as part of the command. Returns a status code, success is 0.
func cmdMicroblogNew() int {
    // check if extra arguments were provided, and error
    if len(c.Arguments) > 1 {
        fmt.Fprintf(os.Stderr, "[error] more than one argument provided\n")
        return 1
    }

    // determine html path using one defined in config file or extra flag
    // extra flag should supersede config file entry
    var htmlPath string
    // config file is already parsed
    if conf != nil {
        htmlPath = conf.MicroblogHtmlFile
    }
    if len(c.Arguments) == 1 {
        htmlPath = c.Arguments[0]
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
        // YYYY-MM-DD-hh:mm:ss
        datetime, err = time.Parse("2006-01-02-15:04:05", datetimeStr)
        if err != nil {
            fmt.Fprintf(os.Stderr, "[error] invalid date provided: %s\n", err)
            return 1
        }
    } else {
        // if date not provided then use current. has the side effect of a date
        // flag with a missing value will just produce current date.
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
    if len(c.Arguments) > 2 {
        fmt.Fprintf(os.Stderr, "[error] more than two arguments provided\n")
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
    if len(c.Arguments) >= 1 {
        htmlPath = c.Arguments[0]
    } else if htmlPath == "" {
        // if no entry in config file or not provided then error as parsing is
        // required to generate RSS file
        fmt.Fprintf(os.Stderr, "[error] missing microblog html file\n")
        return 1
    }
    htmlPath = resolvePath(htmlPath)

    // output path is optional, default behavior on unprovided output path is
    // printing to stdout
    if len(c.Arguments) == 2 {
        outputPath = c.Arguments[1]
    }
    outputPath = resolvePath(outputPath)
    // '-' signifies stdout
    if outputPath == "-" {
        // no provided output path produces the generated rss feed being
        // printed to stdout
        outputPath = ""
    }

    // get the base url
    var baseUrl string = defaultBaseUrl
    if baseUrlFlag, ok := c.Flags["url"]; ok {
        if len(baseUrlFlag) > 0 {
            baseUrl = baseUrlFlag
        } else {
            fmt.Fprintf(os.Stderr, "[error] base url flag missing value\n")
            return 1
        }
    }
    metadata := &microblog.RSSMetadata{
        Title: "yarrie",
        Author: "yarrie",
        Description: "yarrie's microblog",
        BaseUrl: baseUrl, 
    }

    // open the html file
    f, err := os.Open(htmlPath)
    if err != nil {
        fmt.Fprintf(os.Stderr, "[error] failed to open microblog file: %s\n", err)
        return 1
    }
    defer f.Close()

    // generate the final rss feed, returns a string containing feed
    s, err := microblog.GenRssFromFile(f, metadata)
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
    if len(path) == 0 {
        return path
    }
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
    //
    // currently all erroring for extraneous arguments and flags are left to
    // the command branches to handle. no extraneous flag checks are completed
    // as there is no command logic beyond the switch statement below.
    c = cli.Parse()
    if c == nil {
        fmt.Fprintf(os.Stderr, "[error] failed to parse command line arguments\n")
        return
    }
    // print usage when missing a command or help argument/flag
    _, hFlag := c.Flags["h"]
    _, helpFlag := c.Flags["help"] 
    if c.Command == "" || c.Command == "help" || (hFlag || helpFlag) {
        printUsage()
        return
    }

    // parsing the config file for command functions to utilize
    // use the default config path unless -c / --config flag provided
    var configPath = defaultConfigPath
    // check both -c and --config flags when determining custom config file
    // path
    var configFlagUsed = false
    var v string
    if v, configFlagUsed = c.Flags["c"]; configFlagUsed {
        configPath = v
    } else if v, configFlagUsed = c.Flags["config"]; configFlagUsed {
        configPath = v
    }
    // check if missing flag value for config path
    if configFlagUsed && len(configPath) == 0 {
        fmt.Fprintf(os.Stderr, "[error] config path flag missing value\n")
        os.Exit(1)
    }
    configPath = resolvePath(configPath)

    // open the config file path
    configFile, err := os.Open(configPath)
    if err != nil {
        // config file failed to open

        if configFlagUsed {
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

