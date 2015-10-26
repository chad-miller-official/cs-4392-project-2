package main

import (
    "fmt"
    "os"
    "strconv"
)

var Size, NumHoles int
var Neighbors [][]int
var best Board

func solveBoard(b Board) {
    if b.NumPegs() > best.NumPegs() {
        if len(b.Moves) > 0 {
            for _, m := range b.Moves {
                solveBoard(b.ExecuteMove(m))
            }
        } else {
            best = b.Clone()
        }
    }
}

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
    
    startHoles := make([]int, NumHoles)
    
    for i, _ := range startHoles {
        startHoles[i] = i
    }
    
    best = ConstructBoard(startHoles, []Move{})
    
    for i, _ := range startHoles {
        solveBoard(ConstructBoard([]int{i}, []Move{}))
    }
    
    fmt.Println((best.History[0].End + 1), ",", len(best.History))
    
    for _, m := range best.History {
        fmt.Println((m.Start + 1), ",", (m.End + 1))
    }
}

