package main

import "fmt"

func main() {
    arr := []int{64, 34, 25, 12, 22, 11, 90}
    for i, _ := range arr {
        fmt.Print(arr[i], " ")
    }
    fmt.Println()
}