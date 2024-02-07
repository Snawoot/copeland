# copeland

Copeland's preferential voting method implementation. Outputs Copeland's score and ranking for each candidate.

## Usage

### Create ballot files

Just create plain text files with list of preferences for each voter. One candidate name per line.

Example:

a1.txt

```
Bob
Charlie
Alice
Donald
```

a2.txt

```
Donald
Alice
Charlie
Bob
```

a3.txt

```
Alice
Charlie
Bob
Donald
```

### Run the program

```
copeland a1.txt a2.txt a3.txt
```

Example output:

```
Registered names:
	ALICE
	BOB
	CHARLIE
	DONALD

Scores:
	Rank 1:
		3	ALICE
	Rank 2:
		2	CHARLIE
	Rank 3:
		1	BOB
	Rank 4:
		0	DONALD
```

Note: candidate which scored against every other candidate (having score *N-1* with default scoring settings, where *N* is number of candidates), also happens to be a Condorcet winner.

## Synopsis

```
$ copeland -h
Usage:

copeland [OPTION]... [BALLOT FILE]...

  Process files with preference lists and output Copeland's ranking.

Options:
  -normalize-case
    	normalize case (default true)
  -score-loss float
    	score for tie against opponent
  -score-tie float
    	score for tie against opponent (default 0.5)
  -score-win float
    	score for win against opponent (default 1)
  -version
    	show program version and exit
```
