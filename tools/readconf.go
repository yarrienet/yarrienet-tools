package main

import (
    "yarrienet/config"
    "os"
    "fmt"
)

func main() {
    f, err := os.Open("yarrienet.conf")
    if err != nil {
        panic(err)
    } 

    conf, err := config.ReadFile(f)
    if err != nil {
        panic(err)
    }

    fmt.Println(conf)
}

