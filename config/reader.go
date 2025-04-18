// Package config parses the custom config file format.
package config

import (
    "bufio"
    "strings"
    "io"
    "fmt"
    "os"
    "strconv"
)

// Represents the read states that the parser can be in.
type readState int
const (
    // Writing the key to the buffer.
    readStateKey readState = iota
    // Writing the value the buffer.
    readStateValue
)

// Structure that holds known config options.
type Config struct {
    // The path of the microblog HTML file. Represented by
    // "microblog_html_file" in the config file, expects a string.
    MicroblogHtmlFile string 
    // The path of the microblog RSS file. Represented by "microblog_rss_file"
    // in the config file, expects a string.
    MicroblogRssFile string
}

// Represents the states that the string reader within parseValue uses.
type stringState int
const (
    // String is uninitialized (not init with double quotes) and not suitable
    // for content reading.
    stringStateUninit stringState = iota
    // String is initialized and is suitable for writing the contents to the
    // buffer.
    stringStateWrite
    // String is escaped and the next reserved character (e.g. double quotes)
    // should not be treated for their intended purpose.
    stringStateEscape
    // String is terminated (final double quotes have been encountered) and is
    // no longer suitable for content reading.
    stringStateTerminated
)

// Parse the raw string value from the config file to a native type. Supports
// integers, booleans, and strings encapsulated in double quotes, e.g. "...".
// Returns an interface of the parsed value on success, an error on failure.
//
// The string parser is permissive and ignores spaces surrounding the string,
// e.g. '   "string content" ' will result in 'string content'.
func parseValue(s string) (interface{}, error) {
    // try parse as int, if success then return
    if val, err := strconv.Atoi(s); err == nil {
        return val, nil
    }
    // try to parse as boolean, if success then return
    if val, err := strconv.ParseBool(s); err == nil {
        return val, nil
    }

    // not integer to boolean, begin string parsing

    var state = stringStateUninit

    // used to build the new string
    var sb strings.Builder

    for _, r := range s {
        // loop over string, rune by rune

        switch state {
        case stringStateUninit:
            // has not encountered the start of the string (a double quote).
            if r == '"' {
                // double quote means the beginning of the string
                state = stringStateWrite
            } else if r != ' ' {
                // error when encountering a character outside of a string,
                // expect a space which the parser allows e.g.
                // '    "string content"  ' will become 'string content'.
                return nil, fmt.Errorf("character '%c' encountered outside of string", r)
            }
        case stringStateWrite:
            // write state allows for string content to be written into the
            // buffer
            if r == '\\' {
                // if encountering escape character then switch state to allow
                // for it
                state = stringStateEscape
            } else if r == '"' {
                // an (unescaped) double quote denotes the end of the string
                state = stringStateTerminated
            } else {
                // else just write
                sb.WriteRune(r)
            }
        case stringStateEscape:
            // write the escaped character and back out
            sb.WriteRune(r)
            state = stringStateWrite
        case stringStateTerminated:
            // don't allow for any characters following an escaped string
            // excluding space characters
            if r != ' ' {
                return nil, fmt.Errorf("character '%c' encountered after string", r)
            }
        }
    } 

    if state != stringStateTerminated {
        // if loop ended without the string terminated, something is wrong with
        // the given string, tailor the error message depending on error
        if sb.Len() > 0 {
            return nil, fmt.Errorf("string not terminated")
        } else {
            // it's possible to find no string
            return nil, fmt.Errorf("no string encountered")
        }
    }
    return sb.String(), nil
}

// Parse the config key and value pair and update the config pointer with
// result. Will validate that key is supported, and that value is valid.
// Passed config structure will be modified upon successful parsing of key
// value pair. Returns an error on parsing failure.
func updateConfig(config *Config, key string, value string) error {
    // parse the value and return an interface
    parsedValue, err := parseValue(value)
    if err != nil {
        return err
    }

    // check if provided key is supported
    switch key {
    case "microblog_html_file":
        // confirm and set value as string
        if s, ok := parsedValue.(string); ok {
            config.MicroblogHtmlFile = s
        } else {
            return fmt.Errorf("'%s' expects a string value", key)
        }
    case "microblog_rss_file":
        // confirm and set value as string
        if s, ok := parsedValue.(string); ok {
            config.MicroblogRssFile = s
        } else {
            return fmt.Errorf("'%s' expects a string value", key)
        }
    default:
        // invalid key is provided
        return fmt.Errorf("'%s' is not a valid key", key)
    }
    return nil
}

// Read and parse a config file and return a parsed Config structure. Can
// an error.
//
// See Config structure for supported config options and their associated key
// within a config file. Lines prefixed with a '#' character will be ignored.
// Supports Unicode.
func ReadFile(f *os.File) (*Config, error) {
    // reader for the file contents, will loop char by char
    r := bufio.NewReader(f)
    // config structure to be modified and returned
    config := &Config{}
    
    // current iteration key
    var key string
    // current iteration value
    var value string
    // is currently a line comment, should skip over
    var lineComment = false

    state := readStateKey
    // string buffer for writing key / value
    var line = 1
    var sb strings.Builder
    for {
        // read each character (rune) in each iteration
        r, _, err := r.ReadRune()
        if err == io.EOF {
            // eof, end loop
            break
        }
        if err != nil {
            return nil, err
        }

        // line comment handling
        if r == '#' && sb.Len() == 0 {
            // line comment found and represented by a #, set state and
            // skip current iteration until \n found (below)
            lineComment = true
            continue
        } else if r == '\n' && lineComment {
            // line comment ends when a new line is found
            lineComment = false
            line++
            continue
        } else if lineComment {
            // line comment active and not a new line should be ignored
            continue
        }

        switch state {
        case readStateKey:
            // writing the key to the string buffer
            if r == '\n' {
                // key must be followed a value before a new line, invalid
                // config syntax
                if sb.Len() != 0 {
                    // unless the buffer is empty in which case represents an
                    // empty line and should be ignored
                    err = fmt.Errorf("expected value after key (line %d)", line)
                    return nil, err
                }
            } else if r == ' ' {
                // key and value is separated by a space, on which key should
                // extracted from buffer and state should be set to value

                // UNICODE NOTE: fine to not convert to rune when only
                // comparing for single byte space
                if sb.Len() == 0 {
                    // check that buffer is not empty on space
                    err = fmt.Errorf("line must not start with a space (line %d)", line)
                    return nil, err
                }
                // extract key string from buffer and write to state
                key = sb.String()
                // reset state
                state = readStateValue
                sb.Reset()
            } else {
                // write character to buffer if not new line or space
                sb.WriteRune(r) 
            }
        case readStateValue:
            // writing the value to the string buffer
            if r == '\n' {
                // each key value pair is terminated by a space, extract value
                // from buffer and write to state
                value = sb.String()
                // parse key value and update config
                err = updateConfig(config, key, value)
                if err != nil {
                    return nil, fmt.Errorf("%s (line %d)", err, line)
                }

                // reset state
                state = readStateKey
                sb.Reset()
            } else {
                // writing character to buffer if not new line
                sb.WriteRune(r)
            } 
        }
        // increment line number
        if r == '\n' {
            line++
        }
    }
    // handling lines which aren't terminated by a new line
    if state == readStateValue {
        // write the value from the buffer
        value = sb.String()
        // update the config
        err := updateConfig(config, key, value)
        if err != nil {
            return nil, fmt.Errorf("%s (line %d)", err, line)
        }
    }
    return config, nil
}

