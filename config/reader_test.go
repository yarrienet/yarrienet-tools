package config

import (
    "testing"
)

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

