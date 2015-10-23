module Project1.Game

open System

let triangleNumber n = n * (n + 1) / 2
let avg a b = (a + b) / 2

// We want to be able to modify these two as well, hence the "mutable"
let mutable numHoles = 0
let mutable lookupTable = [[]]

// Lookup index lists to aid in calculating neighbors
let l1a i = [ i - 1; i - 1; i; i; i + 1; i + 1 ]
let l1b j = [ j - 1; j; j - 1; j + 1; j; j + 1 ] 
let l2a i = [ i - 2; i - 2; i; i; i + 2; i + 2 ]
let l2b j = [ j - 2; j; j - 2; j + 2; j; j + 2 ]

// Returns the row and column in the lookup table (as a tuple) of a peg.
let getLookupIndices n =
    // Returns the row in the lookup table that a peg is in.
    let rec getLookupIndicesHelper row n =
        match n with
        | p when p <= (triangleNumber row) - 1 -> row - 1
        | _ -> getLookupIndicesHelper (row + 1) n
    let row = getLookupIndicesHelper 0 n
    let col = (List.filter (fun i -> lookupTable.[row].[i] = n) [0 .. lookupTable.[row].Length - 1]).[0]
    (row, col)

type Move (start0 : int, end0 : int) =
    let start1 = start0
    let end1 = end0
    // Get the peg in between two pegs. Assumes that they are only one away from each ther.
    let middle =
        let (si, sj) = getLookupIndices start0
        let (ei, ej) = getLookupIndices end0
        lookupTable.[(avg si ei)].[(avg sj ej)]
    member this.Start
        with get() = start1
    member this.End
        with get() = end1
    member this.Middle
        with get() = middle

type Board (empty : int list, ?history : Move list) =
    // holes is a boolean array - true at an index means the hole there is empty; false means it has a peg
    let holes = [ for i in 0 .. numHoles - 1 -> Seq.exists ((=) i) empty ]
    // Gets a list that contains all valid moves on a board
    let moves =
        // Gets a list of all valid moves given a starting point
        let getMovesHelper (holes : bool list) n =
            // Gets a valid move's endpoint given a starting point, using an index for the lookup index lists.
            let getMovesHelperHelper (holes : bool list) k n =
                let (i, j) = getLookupIndices n
                // This has the potential to throw an index-out-of-bounds exception, in which case there is no move
                try
                    match (not holes.[lookupTable.[(l1a i).[k]].[(l1b j).[k]]]) && (not holes.[lookupTable.[(l2a i).[k]].[(l2b j).[k]]]) with
                    | true -> lookupTable.[(l2a i).[k]].[(l2b j).[k]]
                    | _ -> -1   // Return -1 if there's no move
                with
                | _ -> -1
            [for m in 
                (List.filter
                    (fun i -> i <> -1)  // Filter out all nonexistant moves
                    [for k in 0 .. 5 -> getMovesHelperHelper holes k n]) ->
                        Move(m, n) ]
        List.fold
            (fun acc elem -> acc @ elem)
            []
            [for i in (List.filter (fun i -> holes.[i]) [0 .. numHoles - 1]) -> getMovesHelper holes i]
    let history = defaultArg history []
    member this.Holes
        with get() = holes
    member this.Moves
        with get() = moves
    member this.History
        with get() = history
    member this.Pegs = (List.filter (fun elem -> not elem) holes).Length
    member this.NextBoard i =
        // Gets the holes array for a board, using a previous board's holes array and a move to calculate it.
        let getNextBoardHoles (move : Move) =
            List.filter
                (fun i -> i <> move.End)    // Filter out the end move
                ((List.filter (fun i -> holes.[i]) [0 .. numHoles - 1]) @ [ move.Start; move.Middle ])
        Board(getNextBoardHoles moves.[i], history @ [moves.[i]])
    override this.Equals(other) =
        match other with
        | :? Board as o -> (this.Holes = o.Holes)
        | _ -> false
    override this.GetHashCode() = hash (this.Holes, this.Moves, this.History)
    interface IComparable with
        member this.CompareTo other =
            match other with
            | :? Board as o -> compare this.Pegs o.Pegs
            | _ -> invalidArg "other" "not a Board"

let canExit = ref false

#nowarn "40"
let boardAgent = MailboxProcessor<Board>.Start(fun inbox ->
    // Prints the stats of the given board.
    let printStats (b : Board) =
        printfn "%i, %i" (b.History.[0].End + 1) b.History.Length
        List.iter (fun (m : Move) -> printfn "%i, %i" (m.Start + 1) (m.End + 1)) b.History
    let rec inboxLoop (allBoards : Board list) = async {
        let! board = inbox.Receive()
        printfn "allBoards.Length = %i" (allBoards.Length + 1)
        match allBoards.Length with
        | i when i < numHoles - 1 -> return! inboxLoop (board :: allBoards)
        | _ -> printStats (List.max (board :: allBoards)); canExit := true
    }
    inboxLoop ([])
)

type BoardSolver (start : int) =
    let first = Board([start])
    let best = ref (Board([0 .. numHoles - 1]))
    member this.Best
        with get() = !best
    member this.GetBest =
        (*
         * Recursively gets the best board given an empty starting hole.
         * Tests all possible boards, updating the 'best' ref with the
         * best one we've found so far.
         *)
        let rec getBestHelper (b : Board) =
            match b.Pegs with
            | s when s > (!best).Pegs -> match b.Moves.Length with
                                         | 0 -> match b.Pegs with
                                                | p when p > (!best).Pegs -> best := b
                                                | _ -> ignore 0
                                         | _ -> List.iter (fun i -> getBestHelper (b.NextBoard i)) [0 .. b.Moves.Length - 1]
            | _ -> ignore 0
        getBestHelper first; !best

let makeSolverTask i = async {
    boardAgent.Post (BoardSolver(i)).GetBest
}

// We need to change all of these at once for consistency
let resetSize size =
    numHoles <- triangleNumber size
    lookupTable <- [for i in 0 .. size - 1 -> [(triangleNumber i) .. (triangleNumber (i + 1)) - 1]]

// Get the size parameter we passed in.
let getBoardSize args =
    match (args : string []).Length with
    | 2 -> match args.[0] with
           | "-s" -> System.Int32.Parse args.[1]
           | _ -> -1
    | _ -> -1

[<EntryPoint>]
let main(args) =
    let n = getBoardSize args
    match n with
    // Given an invalid board size, print usage and exit
    | s when s < 5 -> printfn "Usage: mono Game.exe -s [board size]"; exit -1
    | _ -> resetSize n
           List.map (fun i -> makeSolverTask i) [0 .. numHoles - 1]
           |> Async.Parallel
           |> Async.RunSynchronously
           |> ignore
           while not !canExit do
               ()
    0
