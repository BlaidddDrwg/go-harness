package main

import (
  "errors"
  "fmt"
  "math/rand"
  "os"
  "strconv"
)

func returnError(err error) {
  fmt.Println("You need to call this program with 2 arguments: ")
  fmt.Println("Usage: go run randomizer.go 1 10")
}

func main() {
  // Find the command line arguments that the randomizer was called with:
  if len(os.Args) != 3 {
    returnError(errors.New("Wrong arguments error"))
  }

  // Make sure that the 2 arguments are both integers. Error out if not.
  first, err1 := strconv.Atoi(os.Args[1])
  second, err2 := strconv.Atoi(os.Args[2])

  if err1 != nil {
    returnError(err1)
  }
  if err2 != nil {
    returnError(err2)
  }

  fmt.Println(first + rand.Intn(second - first))
}
