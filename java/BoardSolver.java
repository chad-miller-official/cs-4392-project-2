import java.util.concurrent.Callable;

public class BoardSolver implements Callable<Board>
{
	private Board first, best;
	private byte bestNumPegs;
	
	public BoardSolver(byte start)
	{
		first = new Board(new byte[]{ start });
		bestNumPegs = 0;
	}
	
	@Override
	public Board call()
	{
		solve(first);
		return best;
	}
	
	/*
	 * This is a recursive algorithm. It is very memory-intensive.
	 * Basically, we start with a board, and test every possible subsequent board
	 * we can get to from this board, recursively.
	 */
	private void solve(Board b)
	{
		// If we're on a board worth testing that has moves...
		if(b.numPegs() > bestNumPegs && b.moves.length > 0)
		{
			for(Board.Move m : b.moves)
			{
				Board next = b.executeMove(m);
				solve(next);
			}
		}
		else
		{
			// If we have no moves, see if this is the best we've done so far
			if(best == null || b.numPegs() > best.numPegs())
			{
				best = b;
				bestNumPegs = best.numPegs();
			}
		}
	}
}