import java.util.ArrayList;
import java.util.Collection;
import java.util.List;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.Future;

/*
 * Nearly everything in this program is done with bytes
 * and plain arrays to save memory. This program is very
 * memory-intensive and we need to save as much space as possible.
 */
public class Game
{
	public static byte size, numHoles;
	public static ExecutorService threadPool;
	
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
			
			// Create a thread pool, where each thread solves a given board
			threadPool = Executors.newFixedThreadPool(numHoles);
		}
		catch(NumberFormatException e)
		{
			printUsageAndExit();
		}
		
		// List of board solvers for each thread
		Collection<BoardSolver> solvers = new ArrayList<BoardSolver>(numHoles);
		
		// Test each starting position
		for(byte i = 0; i < numHoles; i++)
			solvers.add(new BoardSolver(i));
		
		try
		{
		    // A future corresponding to each thread/board solver
			List<Future<Board>> results = threadPool.invokeAll(solvers);
			
			// Stop accepting threads, and force all threads to start
			threadPool.shutdown();
			Board best = null;
			
			/*
			 * Get the result for each future. We're iterating over all of them
			 * starting at index 0 rather than waiting for each one as they
			 * return. This is because we need to wait for the last one before
			 * we can continue anyway.
			 */
			for(Future<Board> f : results)
			{
				Board result = f.get();
				
				/*
				 * If the board we get is better than our current
				 * best, make our current best the one we just got
				 */
				if(best == null || result.numPegs() > best.numPegs())
					best = result;
			}
			
			// Print out the solution
			System.out.println((best.history[0].end + 1) + ", " + best.history.length);
			
			for(Board.Move m : best.history)
				System.out.println((m.start + 1) + ", " + (m.end + 1));
		}
		catch(InterruptedException | ExecutionException e)
		{
			// Do nothing
		}
	}
	
	private static void printUsageAndExit()
	{
		System.err.println("Usage: java Game -p [board size]");
		System.exit(-1);
	}
}
