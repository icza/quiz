# quiz


Q: Given a list of words like https://github.com/NodePrime/quiz/blob/master/word.list find the longest compound-word in the list, which is also a concatenation of other sub-words that exist in the list. The program should allow the user to input different data. The finished solution shouldn't take more than one hour. Any programming language can be used, but Go is preferred.


Fork this repo, add your solution and documentation on how to compile and run your solution, and then issue a Pull Request. 

Obviously, we are looking for a fresh solution, not based on others' code.

Original task can be found at https://github.com/NodePrime/quiz

---

# Solution 

This is a simple command line application written in Go which takes a source input file
that contains a list of words (one word in each line), and prints the
longest compound-word in the list, which is also a concatenation of other sub-words that exist in the list.
 
The program allows the user to input different data. The input file can be
specified/overridden using the `-src` flag.
By default input file should be sorted alphabetically. If it is not sorted, the `-unsorted` flag must be provided.

An optional `-parallel` flag may be provided in which case the program will perform a parallel search
utilizing all CPU cores available on the machine.

Provide the `-h` flag to see program usage:

    go run findcompword.go -h

## Instructions

This is a simple command line tool. You can run it with

    go run findcompword.go

To perform tests, benchmarks:

    go test -bench .

## Notes

1. The task doesn't mention anything about how many words may combine the compound-word,
so this solution finds the longest word that may be assembled from any number of sub-words (>= 2).
2. The task also doesn't mention if words contain only the letters of the English alphabet,
so this solution properly handles unicode characters (does not assume English alphabet letters only).
3. In case there are multiple equally longest compound-words, the result may be any of these.
If `-parallel` flag is not provided, the first in alphabetical order is returned.
4. This solution reads all words into memory and stores them in a slice.
This is the only memory requirement.

## Speed Analyis

Tested on Intel Core i5-4440 3.1 GHz (4 cores), Go 1.5. To verify, run

    go test -bench .

Running the tool on the list found at https://github.com/NodePrime/quiz/blob/master/word.list
which contains 263533 words, the result (`"antidisestablishmentarianisms"`)
is found in **~0.7 ms**! Less than a 1 ms out of ~263k words, that is wicked fast!
Using the `-parallel` flag, it is slower (~5ms). This is partly due to the fact that the above list
contains the longest compound-word in the first 3% of the list.
When using a randomly generated list (10 million words, with length 5..44, having exactly one compound-word
relatively at the end of the list), result is found in ~47 sec; when using `-parallel` it is found in ~15 sec (with 4 CPU cores).

## Implementation details

The program first reads all words into a slice, then operates on this word slice onward.

**Simple solution** (default when not using the `-parallel` flag) is implemented in `findLongest()`.
It simply iterates over all words, and if a longer compoud-word is found
(longer than previously found), it is stored as the new candidate.
Shorter words than our current candidate are not tested for obvious reasons.

The "is-compound" test is implemented in the `compound()` function, and it uses a recursive definition:
A word is compound-word if:
- There is another word that is a _prefix_ of it
- And the rest after the prefix is also a _word_ or a _compound word_.

The **Parallel solution** is implemented in `findLongest()` function.
The most obvious parallel solution would simply cut the word list to multiple segments, and run the simple solution
on all segments, then "merge" the results.

But I chose a different approach:
I start multiple worker goroutines, and distribute the most expensive "is-compound" jobs to the workers.
Reasoning is that whenever a compound-word is found, ongoing words only needs to be tested if they are longer than this.
Separating the search would not take advantage of this optimization without distributing compound findings.
But then it is just easier to do distribute the compound checks.
