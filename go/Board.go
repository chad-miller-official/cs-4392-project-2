package main

func TriangleNum(n int) int {
    return ( n * ( n + 1 ) ) / 2
}

func getXY(n int) (int, int) {
    var x, y int
    
    // X: Go through each row and check to see if we can possibly be in that row
    for i := 0; i <= len(Neighbors); i++ {
        if n <= TriangleNum(i) - 1 {
            x = i - 1
            break
        }
    }
    
    // Y: Go through each element in that row and find our exact position
    for i, _ := range Neighbors[x] {
        if Neighbors[x][i] == n {
            y = i
            break
        }
    }
    
    return x, y
}

type Move struct {
    Start, End, Middle int
}

/*
 * The X and Y coordinaets of the in-between of two pegs that are two away
 * from each other is the average of their X and Y coordinates.
 */
func (m Move) getMiddle() int {
    x1, y1 := getXY(m.Start)
    x2, y2 := getXY(m.End)
    retval := Neighbors[(x1 + x2) / 2][(y1 + y2) / 2]
    return retval
}

/*
 * Basically a Move constructor.
 * Takes in a start and end point. Assumes they form a valid move.
 */
func NewMove(start, end int) *Move {
    retval := new(Move)
    retval.Start = start
    retval.End = end
    retval.Middle = retval.getMiddle()
    return retval
}

type Board struct {
    Holes []bool                // true = does not have peg; false = has peg
    Moves, History []*Move
}

func (b Board) getNeighborSafe(x1, y1, x2, y2 int) int {
    // Bounds check 1. Make sure we're all positive indices.
    if x1 < 0 || x2 < 0 || y1 < 0 || y2 < 0 {
        return -1
    }
    
    // Bounds check 2. Make sure the first index isn't out of bounds.
    if x1 >= len(Neighbors) || x2 >= len(Neighbors) {
        return -1
    }
    
    // Bounds check 3. Make sure the second index isn't out of bounds.
    if y1 >= len(Neighbors[x1]) || y2 >= len(Neighbors[x2]) {
        return -1
    }
    
    // Hole check - make sure the holes are empty.
    if b.Holes[Neighbors[x1][y1]] || b.Holes[Neighbors[x2][y2]] {
        return -1
    }
    
    return Neighbors[x2][y2]
}

/*
 * How we calculate two-away neighbors:
 * Use the "center" peg to index into the neighbors array.
 * Then, for each of the 6 possible pegs around it, test these two conditions:
 * 1. Is the neighbor not empty?
 * 2. Is the two-away neighbor on the same trajectory empty?
 * If both of those conditions are met, we add that to the array of two-away neighbors.
 * This method is used to calculate the possible moves on the board.
 */
func (b Board) getTwoAwayNeighbors(center int) []int {
    var retval []int
    
    x, y := getXY(center)
    
    // Helper table to determine which points are our neighbors
    neighborPoints := [6][4]int{
        { x - 1, y - 1, x - 2, y - 2 },
        { x - 1, y, x - 2, y },
        { x, y - 1, x, y - 2 },
        { x, y + 1, x, y + 2 },
        { x + 1, y, x + 2, y },
        { x + 1, y + 1, x + 2, y + 2 },
    }

    for _, a := range neighborPoints {
        neighbor := b.getNeighborSafe(a[0], a[1], a[2], a[3])
        
        // If the two-away neighbor peg exists, add it to the list
        if neighbor != -1 {
            retval = append(retval, neighbor)
        }
    }
    
    return retval
}

func (b Board) getMoves() []*Move {
    var retval []*Move
    
    // For each hole...
    for i := 0; i < len(b.Holes); i++ {
        // If it isn't empty...
        if b.Holes[i] {
            /*
			 * Get all its two-away neighbors that meet the following conditions:
			 * 1. The neighbor directly adjacent to it is not empty, and
			 * 2. The two-away neighbor in the same trajectory is empty.
			 */
            for _, j := range b.getTwoAwayNeighbors(i) {
                retval = append(retval, NewMove(j, i))
            }
        }
    }
    
    return retval
}

/*
 * Basically a Board constructor.
 * Takes in an array of holes that are initially empty.
 */
func NewBoard(holes []int) *Board {
    retval := new(Board)
    retval.Holes = make([]bool, NumHoles)
    
    for _, i := range holes {
        retval.Holes[i] = true
    }
    
    retval.Moves = retval.getMoves()
    retval.History = make([]*Move, 0)
    return retval
}

// Get the number of non-empty pegs on the board
func (b Board) NumPegs() int {
    sum := 0
    
    for _, h := range b.Holes {
        if !h {
            sum++
        }
    }
    
    return sum
}

/*
 * This method doesn't actually affect the current board.
 * This method actually constructs a new board based on
 * the current board and the move that was executed.
 */
func (b Board) ExecuteMove(m *Move) *Board {
    retval := new(Board)
    
    /*
     * The holes in the next board are the current board's holes, plus
     * the move's start and middle holes, and without the move's end hole,
     * because a peg moved from the start, eliminated a peg in the middle,
     * and ended up at the end.
     */
    retval.Holes = make([]bool, len(b.Holes))
    copy(retval.Holes, b.Holes)
    retval.Holes[m.Start] = true
    retval.Holes[m.Middle] = true
    retval.Holes[m.End] = false
    
    // Calculate next board's moves based on the new holes
    retval.Moves = retval.getMoves()
    
    // Next board's history is identical to the current board's history...
    retval.History = make([]*Move, len(b.History) + 1)
    
    for i, _ := range b.History {
        retval.History[i] = b.History[i]
    }
    
    // ...Plus the new move
    retval.History[len(b.History)] = m
    
    return retval
}

