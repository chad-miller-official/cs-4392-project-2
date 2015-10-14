import java.util.ArrayList;
import java.util.List;

public class Board implements Cloneable
{
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
	private static byte[][] neighbors;
	
	static
	{
		byte counter = 0;
		neighbors = new byte[Game.size][];
		
		for(int i = 0; i < Game.size; i++)
		{
			neighbors[i] = new byte[i + 1];
			
			for(int j = 0; j < neighbors[i].length; j++)
			{
				neighbors[i][j] = counter;
				counter++;
			}
		}
	}
	
	public boolean[] holes;	// true = does not have peg; false = has peg
	public Move[] moves, history;
	
	/*
	 * Board constructor takes in an array of bytes, with each element being the
	 * hole number of a hole that is empty.
	 */
	public Board(byte[] empty)
	{
		this(empty, new Move[]{ });
	}
	
	/*
	 * This is a private constructor used only by Board.executeMove().
	 * It takes in an array of the moves we've used to get to this board state.
	 */
	private Board(byte[] empty, Move[] history)
	{
		holes = new boolean[Game.numHoles];
		
		// Create the holes array, which is just a 1-D array representing the board state
		for(byte i : empty)
			holes[i] = true;
		
		// Create the array of possible moves, unless there are no empty holes
		this.history = new Move[history.length];
		System.arraycopy(history, 0, this.history, 0, history.length);
		moves = (empty.length > 0) ? calculateMoves() : new Move[]{ };
	}
	
	private Move[] calculateMoves()
	{
		List<Move> retval = new ArrayList<Move>();
		
		// For each hole...
		for(byte i = 0; i < holes.length; i++)
		{
			// If it isn't empty...
			if(holes[i])
			{
				/*
				 * Get all its two-away neighbors that meet the following conditions:
				 * 1. The neighbor directly adjacent to it is not empty, and
				 * 2. The two-away neighbor in the same trajectory is empty.
				 */
				for(byte j : getTwoAwayNeighbors(i))
					retval.add(new Move(j, i));
			}
		}
		
		return retval.toArray(new Move[retval.size()]);
	}
	
	// Get the number of pegs on the board.
	public byte numPegs()
	{
		byte sum = 0;
		
		for(boolean h : holes)
		{
			if(!h)
				sum++;
		}
		
		return sum;
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
	public byte[] getTwoAwayNeighbors(byte center)
	{
		byte[] retval = new byte[]{ 0, 0, 0, 0, 0, 0 };
		byte count = 0;
		
		calcLoop:
		{
			for(byte i = 0; i < Game.size; i++)
			{
				for(byte j = 0; j < neighbors[i].length; j++)
				{
					if(neighbors[i][j] == center)
					{
						/*
						 * We do a lot of try-catches. That's because we may
						 * go off the board at some point, and because I'm
						 * too lazy to manually calculate whether an index
						 * is actually out of bounds or not.
						 */
						try
						{
							if(!holes[neighbors[i - 1][j - 1]] && !holes[neighbors[i - 2][j - 2]])
							{
								retval[count] = neighbors[i - 2][j - 2];
								count++;
							}
						}
						catch(ArrayIndexOutOfBoundsException e)
						{
							// Don't handle
						}
						
						try
						{
							if(!holes[neighbors[i - 1][j]] && !holes[neighbors[i - 2][j]])
							{
								retval[count] = neighbors[i - 2][j];
								count++;
							}
						}
						catch(ArrayIndexOutOfBoundsException e)
						{
							// Don't handle
						}
						
						try
						{
							if(!holes[neighbors[i][j - 1]] && !holes[neighbors[i][j - 2]])
							{
								retval[count] = neighbors[i][j - 2];
								count++;
							}
						}
						catch(ArrayIndexOutOfBoundsException e)
						{
							// Don't handle
						}
						
						try
						{
							if(!holes[neighbors[i][j + 1]] && !holes[neighbors[i][j + 2]])
							{
								retval[count] = neighbors[i][j + 2];
								count++;
							}
						}
						catch(ArrayIndexOutOfBoundsException e)
						{
							// Don't handle
						}
						
						try
						{
							if(!holes[neighbors[i + 1][j]] && !holes[neighbors[i + 2][j]])
							{
								retval[count] = neighbors[i + 2][j];
								count++;
							}
						}
						catch(ArrayIndexOutOfBoundsException e)
						{
							// Don't handle
						}
						
						try
						{
							if(!holes[neighbors[i + 1][j + 1]] && !holes[neighbors[i + 2][j + 2]])
							{
								retval[count] = neighbors[i + 2][j + 2];
								count++;
							}
						}
						catch(ArrayIndexOutOfBoundsException e)
						{
							// Don't handle
						}
						
						/*
						 * Always break at the end, because we calculated all the moves
						 * from the center
						 */
						break calcLoop;
					}
				}
			}
		}
		
		byte[] realRetval = new byte[count];
		System.arraycopy(retval, 0, realRetval, 0, realRetval.length);
		return realRetval;
	}
	
	/*
	 * This method doesn't actually affect the current board.
	 * This method actually constructs a new board based on
	 * the current board and the move that was executed.
	 */
	public Board executeMove(Move m)
	{
		/*
		 * This is just a way of calculating the empty holes
		 * on the board, so we can pass in an array of empty
		 * holes to the board that will be returned.
		 */
		boolean[] preset = new boolean[Game.numHoles];
		System.arraycopy(holes, 0, preset, 0, preset.length);
		preset[m.start] = true;
		preset[m.middle] = true;
		preset[m.end] = false;
		
		byte i;
		byte len = 0;
		
		// Get the number of empty holes
		for(i = 0; i < preset.length; i++)
		{
			if(preset[i])
				len++;
		}
		
		// Create an array that contains all the empty hole numbers
		byte[] presetArr = new byte[len];
		len = 0;
		
		// Populate that array
		for(i = 0; i < preset.length; i++)
		{
			if(preset[i])
			{
				presetArr[len] = i;
				len++;
			}
		}
		
		// Now, take the current history, and add the move we just executed to it
		Move[] updatedHistory = new Move[history.length + 1];
		System.arraycopy(history, 0, updatedHistory, 0, history.length);
		updatedHistory[updatedHistory.length - 1] = m;
		
		// And construct a new board with both the new empty holes and new history
		return new Board(presetArr, updatedHistory);
	}
	
	public class Move
	{
		public byte start, middle, end;
		
		public Move(byte start, byte end)
		{
			this.start = start;
			this.end = end;
			middle = getMiddle(this.start, this.end);
		}
		
		/*
		 * This method is similar to Board.getTwoAwayNeighbors(), but is a bit simpler.
		 * First, use the start as an "anchor" on the board.
		 * Iterate through every element in the neighbors array until we hit the start.
		 * Then, we test the six possible two-away neighbors that may be the end.
		 * When we find the right one, we return the one in between both of them.
		 */
		private byte getMiddle(byte start, byte end)
		{
			byte retval = 0;
			
			calcLoop:
			{
				for(byte i = 0; i < Game.size; i++)
				{
					for(byte j = 0; j < neighbors[i].length; j++)
					{
						if(neighbors[i][j] == start)
						{
							/*
							 * Like Board.getTwoAwayNeighbors(), we do a lot of
							 * try-catches, for the exact same reason.
							 */
							try
							{
								if(end == neighbors[i - 2][j - 2])
								{
									retval = neighbors[i - 1][j - 1];
									
									/*
									 * Like Board.getTwoAwayNeighbors(), we break
									 * for the same reason.
									 */
									break calcLoop;
								}
							}
							catch(ArrayIndexOutOfBoundsException e)
							{
								// Don't handle
							}
							
							try
							{
								if(end == neighbors[i - 2][j])
								{
									retval = neighbors[i - 1][j];
									break calcLoop;
								}
							}
							catch(ArrayIndexOutOfBoundsException e)
							{
								// Don't handle
							}
							
							try
							{								
								if(end == neighbors[i][j - 2])
								{
									retval = neighbors[i][j - 1];
									break calcLoop;
								}
							}
							catch(ArrayIndexOutOfBoundsException e)
							{
								// Don't handle
							}
								
							try
							{
								if(end == neighbors[i][j + 2])
								{
									retval = neighbors[i][j + 1];
									break calcLoop;
								}
							}
							catch(ArrayIndexOutOfBoundsException e)
							{
								// Don't handle
							}
							
							try
							{
								if(end == neighbors[i + 2][j])
								{
									retval = neighbors[i + 1][j];
									break calcLoop;
								}
							}
							catch(ArrayIndexOutOfBoundsException e)
							{
								// Don't handle
							}
								
							try
							{
								if(end == neighbors[i + 2][j + 2])
								{
									retval = neighbors[i + 1][j + 1];
									break calcLoop;
								}
							}
							catch(ArrayIndexOutOfBoundsException e)
							{
								// Don't handle
							}
						}
					}
				}
			}
						
			return retval;
		}
	}
}
