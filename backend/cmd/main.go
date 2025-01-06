package main

import (
    "fmt"
    "os"
)

func main() {
    d := os.Getenv("DB")
    u := os.Getenv("DB_U")
    p := os.Getenv("DB_P")

    fmt.Printf(d, u, p)
}
