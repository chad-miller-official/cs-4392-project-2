package main

func TriangleNum(n int) int {
    return ( n * ( n + 1 ) ) / 2
}

func getXY(n int) (int, int) {
    var x, y int
    
    for i := 0; i <= len(Neighbors); i++ {
        if n <= TriangleNum(i) {
            x = i - 1
            break
        }
    }
    
    for _, i := range Neighbors[x] {
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
    return Neighbors[(x1 + x2) / 2][(y1 + y2) / 2]
}
