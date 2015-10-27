package main

import (
    "fmt"
    "os"
    "runtime"
    "strconv"
)

var Size, NumHoles int

/*
 * This is a jagged array used to calculate hole locations. It approximates
 * the game board. At the end of the static block, for a number like s=5,
 * neighbors looks like:
 * 
 * [[0]
 *  [1, 2]
 *  [3, 4, 5]
 *  [6, 7, 8, 9]
 *  [10, 11, 12, 13, 14]]
 */
var Neighbors [][]int
var allBest *Board

/*
 * This class is basically an instance of a function
 * that solves a board given a starting hole.
 */
type BoardSolver struct {
    // first is the initial board we're starting with
    first, Best *Board
}

/*
 * Basically a BoardSolver constructor.
 */
func NewBoardSolver(start int) *BoardSolver {
    retval := new(BoardSolver)
    retval.first = NewBoard([]int{start})
    
    // Create an array that's just [0, 1, ..., NumHoles]
    startHoles := make([]int, NumHoles)
    
    for i, _ := range startHoles {
        startHoles[i] = i
    }
    
    // And use it to create the "best" board so far, which has no pegs
    retval.Best = NewBoard(startHoles)
    return retval
}

/*
 * This is a recursive algorithm. It is very memory-intensive.
 * Basically, we start with a board, and test every possible subsequent board
 * we can get to from this board, recursively.
 */
func (bs BoardSolver) solveBoardHelper(b *Board) {
    // If we're on a board worth testing that has moves...
    if b.NumPegs() > bs.Best.NumPegs() {
        if len(b.Moves) > 0 {
            // Recursively test each next board
            for _, m := range b.Moves {
                bs.solveBoardHelper(b.ExecuteMove(m))
            }
        } else {
            // If we have no moves, see if this is the best we've done so far
            *bs.Best = *b
        }
    }
}

// Basically a wrapper for solveBoardHelper()
func (bs BoardSolver) SolveBoard() {
    bs.solveBoardHelper(bs.first)
}

func main() {
    // Take advantage of every core we can work with
    runtime.GOMAXPROCS(runtime.NumCPU())
    
    // Make sure we get the correct arguments
	args := os.Args[1:]
	
	if len(args) != 2 || args[0] != "-s" {
	    fmt.Println("Usage: ./Game -s [board size]")
	    os.Exit(-1)
    }
    
    Size, _ = strconv.Atoi(args[1])
    NumHoles = TriangleNum(Size)
    Neighbors = make([][]int, Size)
    
    // Create the Neighbors lookup table
    counter := 0
    
    for i := 0; i < Size; i++ {
        Neighbors[i] = make([]int, i + 1)
        
        for j := 0; j < i + 1; j++ {
            Neighbors[i][j] = counter
            counter++
        }
    }
    
    // Create an array that's just [0, 1, ..., NumHoles]
    startHoles := make([]int, NumHoles)
    
    for i, _ := range startHoles {
        startHoles[i] = i
    }
    
    // And use it to create the "best" board so far, which has no pegs
    allBest = NewBoard(startHoles)
    
    // Initialize the channels array we're going to work with
    channels := make([]chan *Board, NumHoles)
    
    for i, _ := range startHoles {
        // Initialize each individual channel
        channels[i] = make(chan *Board)
        
        // And start a new goroutine to calculate the solution for a board with a given starting point
        go func (channel chan *Board, start int) {
            // Solve the board...
            solver := NewBoardSolver(start)
            solver.SolveBoard()
            
            // ...And pass the result to the channel
            channel <- solver.Best
        }(channels[i], i)
    }
    
    /*
     * Get the result from each channel. We're iterating over all of them
     * starting at index 0 rather than waiting for each one as they
     * return. This is because we need to wait for the last one before
     * we can continue anyway.
     */
    for i, _ := range startHoles {
        best := <-channels[i]
        
        /*
		 * If the board we get is better than our current
		 * best, make our current best the one we just got
		 */
        if best.NumPegs() > allBest.NumPegs() {
            allBest = best
        }
    }
    
    // Print out the solution
    fmt.Printf("%d, %d\n", (allBest.History[0].End + 1), len(allBest.History))
    
    for _, m := range allBest.History {
        fmt.Printf("%d, %d\n", (m.Start + 1), (m.End + 1))
    }
}

