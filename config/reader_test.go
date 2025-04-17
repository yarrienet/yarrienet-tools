package config

import (
    "testing"
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

// Write
func configFileHelper() {

}

