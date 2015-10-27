package main

import (
    "fmt"
    "os"
    "runtime"
    "strconv"
)

var Size, NumHoles int
var Neighbors [][]int
var allBest *Board

type BoardSolver struct {
    first, Best *Board
}

func NewBoardSolver(start int) *BoardSolver {
    retval := new(BoardSolver)
    retval.first = NewBoard([]int{start})
    
    startHoles := make([]int, NumHoles)
    
    for i, _ := range startHoles {
        startHoles[i] = i
    }
    
    retval.Best = NewBoard(startHoles)
    return retval
}

func (bs BoardSolver) solveBoardHelper(b *Board) {
    if b.NumPegs() > bs.Best.NumPegs() {
        if len(b.Moves) > 0 {
            for _, m := range b.Moves {
                bs.solveBoardHelper(b.ExecuteMove(m))
            }
        } else {
            *bs.Best = *b
        }
    }
}

func (bs BoardSolver) SolveBoard() {
    bs.solveBoardHelper(bs.first)
}

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    
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
    
    allBest = NewBoard(startHoles)
    
    channels := make([]chan *Board, NumHoles)
    
    for i, _ := range startHoles {
        channels[i] = make(chan *Board)
        
        go func (channel chan *Board, start int) {
            solver := NewBoardSolver(start)
            solver.SolveBoard()
            channel <- solver.Best
        }(channels[i], i)
    }
    
    for i, _ := range startHoles {
        best := <-channels[i]
        
        if best.NumPegs() > allBest.NumPegs() {
            allBest = best
        }
    }
    
    fmt.Printf("%d, %d\n", (allBest.History[0].End + 1), len(allBest.History))
    
    for _, m := range allBest.History {
        fmt.Printf("%d, %d\n", (m.Start + 1), (m.End + 1))
    }
}

