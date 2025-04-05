package config

import (
    "bufio"
    "strings"
    "io"
    "fmt"
    "os"
    "strconv"
)

type ReadState int
const (
    ReadStateKey ReadState = iota
    ReadStateValue
)

type Config struct {
    MicroblogHtmlFile string 
}

func parseValue(s string) (interface{}, error) {
    if val, err := strconv.Atoi(s); err == nil {
        return val, nil
    }
    if val, err := strconv.ParseBool(s); err == nil {
        return val, nil
    }

    var stringTerminated = false
    var stringEscape = false

    var sb strings.Builder
    for _, r := range s {
        if stringTerminated {
            return nil, fmt.Errorf("string terminated?")
        }
    
        if r == '"' && !stringEscape {
            if sb.Len() != 0 {
                // string terminated
                stringTerminated = true 
            }
        } else if r == '\\' {
            stringEscape = true
        } else {
            sb.WriteRune(r)
        }
    } 
    return sb.String(), nil
}

func updateConfig(config *Config, key string, value string) error {
    parsedValue, err := parseValue(value)
    if err != nil {
        return err
    }

    switch key {
    case "microblog_html_file":
        if s, ok := parsedValue.(string); ok {
            config.MicroblogHtmlFile = s
        } else {
            return fmt.Errorf("not expecting string")
        }
    }
    return nil
}

func ReadFile(f *os.File) (*Config, error) {
    r := bufio.NewReader(f)

    config := &Config{}
    
    var key string
    var value string

    state := ReadStateKey
    var sb strings.Builder
    for {
        r, _, err := r.ReadRune()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }

        switch state {
        case ReadStateKey:
            if r == ' ' {
                key = sb.String()
                state = ReadStateValue
                sb.Reset()
            } else if r == '\n' {
                if sb.Len() != 0 {
                    err = fmt.Errorf("expected value after key")
                    return nil, err
                }
            } else {
                sb.WriteRune(r) 
            }
        case ReadStateValue:
            if r == '\n' {
                value = sb.String()
                updateConfig(config, key, value)

                state = ReadStateKey
                sb.Reset()
            } else {
                sb.WriteRune(r)
            } 
        }
    }
    // result := sb.String()
    return config, nil
}

