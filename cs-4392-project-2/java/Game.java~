
/*
 * Nearly everything in this program is done with bytes
 * and plain arrays to save memory. This program is very
 * memory-intensive and we need to save as much space as possible.
 */
public class Game
{
	public static byte size, numHoles;
	public static Board best = null;
	public static byte bestNumPegs = 0;
	
	public static void main(String[] args)
	{
		// Make sure we get the correct arguments
		if(args.length != 2 || !args[0].equals("-s"))
			printUsageAndExit();
		
		try
		{
			// Get the board size, and use it to calculate the corresponding triangle number
			size = Byte.parseByte(args[1]);
			numHoles = (byte) ((size * (size + 1)) / 2);
		}
		catch(NumberFormatException e)
		{
			printUsageAndExit();
		}
		
		// Test each starting position
		for(byte i = 0; i < numHoles; i++)
		{
			Board b = new Board(new byte[]{ i });
			solveBoard(b);
		}
		
		// Print out the final solution
		System.out.println((best.history[0].end + 1) + ", " + best.history.length);
		
		for(Board.Move m : best.history)
			System.out.println((m.start + 1) + ", " + (m.end + 1));
	}
	
	/*
	 * This is a recursive algorithm. It is very memory-intensive.
	 * Basically, we start with a board, and test every possible subsequent board
	 * we can get to from this board, recursively.
	 */
	private static void solveBoard(Board b)
	{
		// If we're on a board worth testing that has moves...
		if(b.numPegs() > bestNumPegs && b.moves.length > 0)
		{
			for(Board.Move m : b.moves)
			{
				Board next = b.executeMove(m);
				solveBoard(next);
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
	
	private static void printUsageAndExit()
	{
		System.err.println("Usage: java Game -s [board size]");
		System.exit(-1);
	}
}
