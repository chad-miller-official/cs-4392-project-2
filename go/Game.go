package main

import (
    "fmt"
    "os"
    "strconv"
)

var Size, NumHoles int
var Neighbors [][]int

func main() {
	args := os.Args[1:]
	
	if len(args) != 2 || args[0] != "-s" {
	    fmt.Println("Usage: ./Game -s [board size]")
	    os.Exit(-1)
    }
    
    Size, _ = strconv.Atoi(args[1])
    NumHoles = TriangleNum(Size)
    Neighbors = make([][]int, Size)
    
    counter := 0
    
    for i := 0; i < Size; i++ {
        Neighbors[i] = make([]int, i + 1)
        
        for j := 0; j < i + 1; j++ {
            Neighbors[i][j] = counter
            counter++
        }
    }
    
    fmt.Println("Size:", Size)
    fmt.Println("NumHoles:", NumHoles)
    fmt.Println("Neighbors:", Neighbors)
}
