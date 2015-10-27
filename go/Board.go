package main

func TriangleNum(n int) int {
    return ( n * ( n + 1 ) ) / 2
}

func getXY(n int) (int, int) {
    var x, y int
    
    for i := 0; i <= len(Neighbors); i++ {
        if n <= TriangleNum(i) - 1 {
            x = i - 1
            break
        }
    }
    
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

func (m Move) getMiddle() int {
    x1, y1 := getXY(m.Start)
    x2, y2 := getXY(m.End)
    retval := Neighbors[(x1 + x2) / 2][(y1 + y2) / 2]
    return retval
}

func NewMove(start, end int) *Move {
    retval := new(Move)
    retval.Start = start
    retval.End = end
    retval.Middle = retval.getMiddle()
    return retval
}

func (m Move) Clone() *Move {
    retval := new(Move)
    retval.Start = m.Start
    retval.End = m.End
    retval.Middle = m.Middle
    return retval
}

type Board struct {
    Holes []bool
    Moves, History []*Move
}

func (b Board) getNeighborSafe(x1, y1, x2, y2 int) int {
    if x1 < 0 || x2 < 0 || y1 < 0 || y2 < 0 {
        return -1
    }
    
    if x1 >= len(Neighbors) || x2 >= len(Neighbors) {
        return -1
    }
    
    if y1 >= len(Neighbors[x1]) || y2 >= len(Neighbors[x2]) {
        return -1
    }
    
    if b.Holes[Neighbors[x1][y1]] || b.Holes[Neighbors[x2][y2]] {
        return -1
    }
    
    return Neighbors[x2][y2]
}

func (b Board) getTwoAwayNeighbors(center int) []int {
    var retval []int
    
    x, y := getXY(center)
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
        
        if neighbor != -1 {
            retval = append(retval, neighbor)
        }
    }
    
    return retval
}

func (b Board) getMoves() []*Move {
    var retval []*Move
    
    for i := 0; i < len(b.Holes); i++ {
        if b.Holes[i] {
            for _, j := range b.getTwoAwayNeighbors(i) {
                retval = append(retval, NewMove(j, i))
            }
        }
    }
    
    return retval
}

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

func (b Board) NumPegs() int {
    sum := 0
    
    for _, h := range b.Holes {
        if !h {
            sum++
        }
    }
    
    return sum
}

func (b Board) ExecuteMove(m *Move) *Board {
    retval := new(Board)
    
    retval.Holes = make([]bool, len(b.Holes))
    copy(retval.Holes, b.Holes)
    retval.Holes[m.Start] = true
    retval.Holes[m.Middle] = true
    retval.Holes[m.End] = false
    
    retval.Moves = retval.getMoves()
    
    retval.History = make([]*Move, len(b.History) + 1)
    
    for i, _ := range b.History {
        retval.History[i] = b.History[i].Clone()
    }
    
    retval.History[len(b.History)] = m.Clone()
    
    return retval
}

func (b Board) Clone() *Board {
    retval := new(Board)
    
    retval.Holes = make([]bool, len(b.Holes))
    copy(retval.Holes, b.Holes)
    
    retval.Moves = make([]*Move, len(b.Moves))
    
    for i, _ := range b.Moves {
        retval.Moves[i] = b.Moves[i].Clone()
    }
    
    retval.History = make([]*Move, len(b.History))
    
    for i, _ := range b.History {
        retval.History[i] = b.History[i].Clone()
    }
    
    return retval
}

