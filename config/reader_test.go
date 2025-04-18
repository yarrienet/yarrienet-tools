package config

import (
    "testing"
    "fmt"
    "os"
    "regexp"
    "strings"
)

// Test parsing valid integer and boolean values with parseValue.
func TestValidParseValue(t *testing.T) {
    // testing integer
    var input = "45"
    val, err := parseValue(input)
    if err != nil {
        t.Errorf("failed to parse '%s': %s", input, err)
    } else {
        if i, ok := val.(int); ok {
            if i != 45 {
                t.Errorf("parsed value of '%s' should be 45 not %d", input, i) 
            }
        } else {
            t.Errorf("parsed value of '%s' should be an int not %v", input, val)
        }
    }

    // testing true boolean
    input = "true"
    val, err = parseValue(input)
    if err != nil {
        t.Errorf("failed to parse '%s': %s'", input, err)
    } else {
        if b, ok := val.(bool); ok {
            if b != true {
                t.Errorf("parsed value of '%s' should be true not %t", input, b)
            }
        } else {
            t.Errorf("parsed value of '%s' should be a boolean not %v", input, val)
        }
    }
}

// Test parsing valid and invalid string values with parseValue.
func TestStringParseValue(t *testing.T) {
    // valid string
    var valid = map[string]string{
        `"example"`: "example",
        `"multiple words"`: "multiple words",
        `"a"`: "a",
        `""`: "",
        `"\""`: `"`,
        `"a\"a"`: `a"a`,
        `    "test content" `: "test content",
        `"64"`: "64",
    }
    for input, expected := range valid {
        val, err := parseValue(input)
        if err != nil {
            t.Errorf("failed to parse '%s': %s", input, err)
            continue
        }
        if s, ok := val.(string); ok {
            if s != expected { 
                t.Errorf("parsed value '%s' should be '%s' not '%s'", input, expected, s)
            }
        } else {
            t.Errorf("parsed value of '%s' should be a string not %v", input, val)
        }
    }

    // invalid strings, all should error
    var invalid = []string{ `"incomplete`, `incomplete2"`, `no quotes`, `q\"uotes`, `"`, `"""`, "", "   "}
    for _, input := range invalid {
        val, err := parseValue(input)
        if err == nil {
            t.Errorf("invalid string '%s' should result in error not '%v'", input, val)
        }
        if val != nil {
            t.Errorf("invalid string '%s' results in val '%v'", input, val)
        }
    } 
}

// Test updating the parsed config structure with valid and invalid keys and
// values with updateConfig.
func TestUpdateConfig(t *testing.T) {
    var config = Config{}

    // testing value change of microblog_html_file
    key, value, expected := "microblog_html_file", `"path"`, "path"
    err := updateConfig(&config, key, value)
    if err != nil {
        t.Errorf("updating config key '%s' resulted in an error: %s", key, err)
    } else {
        if config.MicroblogHtmlFile != expected {
            t.Errorf("failed to update config '%s' (config.MicroblogHtmlFile), expected '%s' not '%s'", key, expected, config.MicroblogHtmlFile)
        }
    }

    // testing value change of microblog_rss_file
    key, value, expected = "microblog_rss_file", `"path2"`, "path2"
    err = updateConfig(&config, key, value)
    if err != nil {
        t.Errorf("updating config key '%s' resulted in an error: %s", key, err)
    } else {
        if config.MicroblogRssFile != expected {
            t.Errorf("failed to update config '%s' (config.MicroblogRssFile), expected '%s' not '%s'", key, expected, config.MicroblogRssFile)
        }
    }

    // invalid key
    key, value = "microblog_invalid_file", `"x"`
    err = updateConfig(&config, key, value)
    if err == nil {
        t.Errorf("updating invalid config key '%s' should result in error", key)
    }

    // invalid type for key (45 will be parsed as a integer, as per parseValue)
    key, value = "microblog_html_file", "45"
    err = updateConfig(&config, key, value)
    if err == nil {
        t.Errorf("updating config key '%s' with integer '%s' should result in an error", key, value)
    }

    // invalid string, testing pass through to parseValue
    key, value = "microblog_html_file", `"value`
    err = updateConfig(&config, key, value)
    if err == nil {
        t.Errorf("updating config key '%s' with invalid string '%s' should result in an error", key, value)
    }
}

// Create and open a config file of a random name in the operating system's
// temporary directory. Uses the pattern 'yarrienet-tools-test[rand].conf'.
// Returns file on success, error on failure.
func openTempConfigFile() (*os.File, error) {
    var f, err = os.CreateTemp("", "yarrienet-tools-test*.conf")
    if err != nil {
        return nil, err
    }
    return f, nil
}

// Determine the line number from ReadFile error messages. Returns a string
// containing the line number on success, an empty string when extraction
// fails.
func determineLineNumber(err string) string {
    // regex to extract the line number
    lineRe := regexp.MustCompile(`line (\d+)`)
    matches := lineRe.FindStringSubmatch(err)
    // matches is the structure of [ fullString, matchedSubstring ]
    if len(matches) > 1 {
        return matches[1]
    } else {
        return ""
    }
}

// Test the main function that parses the config file and returns a config
// structure containing user options. Tests both the config structure and the
// errors returned.
func TestConfigFileHelper(t *testing.T) {
    // test if a valid config file succeeds (without a terminating new line)
    var expectedHtmlFile = "~/yarrie.net/microblog/index.html"
    var expectedRssFile = "~/yarrie.net/microblog/feed.xml"
    var configString = fmt.Sprintf(`microblog_html_file "%s"
# test comment
microblog_rss_file "%s"`, expectedHtmlFile, expectedRssFile)

    var validFile, err = openTempConfigFile()
    if err != nil {
        t.Fatalf("failed to open temp testing config file: %s", err)
    }
    defer validFile.Close()

    validFile.WriteString(configString)
    validFile.Seek(0, 0)

    conf, err := ReadFile(validFile)
    if err != nil {
        t.Fatalf("failed to parse config file: %s", err)
    }

    if conf.MicroblogHtmlFile != expectedHtmlFile {
        t.Errorf("expected MicroblogHtmlFile to be '%s' not '%s'", expectedHtmlFile, conf.MicroblogHtmlFile)
    }
    if conf.MicroblogRssFile != expectedRssFile {
        t.Errorf("expected MicroblogRssFile to be '%s' not '%s'", expectedRssFile, conf.MicroblogRssFile)
    }

    // testing if mangled spaces work
    configString = fmt.Sprintf(`# comment 1
microblog_html_file    "%s"     

# comment 2

`, expectedHtmlFile)
    validMangledFile, err := openTempConfigFile() 
    if err != nil {
        t.Fatalf("failed to open temp config file: %s", err)
    }
    defer validMangledFile.Close()

    validMangledFile.WriteString(configString)
    validMangledFile.Seek(0, 0)

    conf, err = ReadFile(validMangledFile)
    if err != nil {
        t.Fatalf("failed to parse config file: %s", err)
    }

    if conf.MicroblogHtmlFile != expectedHtmlFile {
        t.Errorf("expected mangled MicroblogHtmlFile to be '%s' not '%s'", expectedHtmlFile, conf.MicroblogHtmlFile)
    }
    if conf.MicroblogRssFile != "" {
        t.Errorf("expected MicroblogRssFile to be empty not '%s'", conf.MicroblogRssFile)
    }

    // testing if error reporting is valid
    configString = "microblog_html_file 10"
    invalidFile1, err := openTempConfigFile()
    if err != nil {
        t.Fatalf("failed to open temp config file")
    }
    defer invalidFile1.Close()

    invalidFile1.WriteString(configString)
    invalidFile1.Seek(0, 0)

    conf, err = ReadFile(invalidFile1)
    if err != nil {
        // extract line number
        errStr := err.Error()
        lineNumber := determineLineNumber(errStr)
        if lineNumber == "" {
            t.Errorf("missing '(line x)' from error: %s", err)
        } else if lineNumber != "1" {
            t.Errorf("expecting error containing line number 1 not %s", lineNumber)
        }

        // check if error is valid
        if !strings.Contains(errStr, "'microblog_html_file' expects a string") {
            t.Errorf("incorrect error for providing an integer to a string: %s", errStr)
        }
    } else {
        t.Errorf("expected error for invalid config file 'invalidFile1'")
    }

    // testing if line number reporting is valid on error
    // error occurs in configString on line 6
    configString = `


# comment

microblog_html_file 10`

    invalidFile2, err := openTempConfigFile()
    if err != nil {
        t.Fatalf("failed to open temp config file")
    }
    defer invalidFile2.Close()

    invalidFile2.WriteString(configString)
    invalidFile2.Seek(0, 0)

    conf, err = ReadFile(invalidFile2)
    if err != nil {
        errStr := err.Error()
        lineNumber := determineLineNumber(errStr)
        if lineNumber == "" {
            t.Errorf("missing '(line x)' from error: %s", err)
        } else if lineNumber != "6" {
            t.Errorf("expecting error containing line number 6 not %s", lineNumber)
        }
    } else {
        t.Errorf("expected error for invalid config file 'invalidFile2'")
    }
}

