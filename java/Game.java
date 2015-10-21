/*
 * Nearly everything in this program is done with bytes
 * and plain arrays to save memory. This program is very
 * memory-intensive and we need to save as much space as possible.
 */
public class Game
{
	public static byte size, numHoles;
	public static volatile Board best = null;
	public static volatile byte bestNumPegs = 0;
	public static Thread[] threads;
	
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
			threads = new Thread[numHoles];
		}
		catch(NumberFormatException e)
		{
			printUsageAndExit();
		}
		
		// Test each starting position
		for(byte i = 0; i < numHoles; i++)
		{
			threads[i] = new Thread(new BoardSolver(i));
			threads[i].start();
		}
		
		for(byte i = 0; i < numHoles; i++)
		{
			try
			{
				threads[i].join();
			}
			catch(InterruptedException e)
			{
				// Do nothing
			}
		}
		
		// Print out the final solution
		System.out.println((best.history[0].end + 1) + ", " + best.history.length);
		
		for(Board.Move m : best.history)
			System.out.println((m.start + 1) + ", " + (m.end + 1));
	}
	
	private static void printUsageAndExit()
	{
		System.err.println("Usage: java Game -p [board size]");
		System.exit(-1);
	}
	
	private static class BoardSolver implements Runnable
	{
		private Board b;
		
		public BoardSolver(byte start)
		{
			b = new Board(new byte[]{ start });
		}
		
		@Override
		public void run()
		{
			solve(b);
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
}
