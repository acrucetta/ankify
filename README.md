# Ankify

Ankify is a command line tool written in Golang that reads a PDF, TXT, or an URL and outputs Anki questions for the user. The project uses Cobra CLI to provide a user-friendly interface for creating flashcards from your favorite PDFs.

## Installation

To install Ankify, you need to have Golang installed on your system. Clone the repository to your local machine, navigate to the project directory, and run the following command:

`go build`

This will compile the code and create an executable file named `ankify`.

## Usage

```
Usage:
  Anki ankify [file path] [flags]

Aliases:
  ankify, a

Flags:
  -h, --help          help for ankify
  -p, --pages ints    Page numbers to parse, e.g., '1,2,3' (default is 1)
  -t, --type string   Type of file to parse, either 'txt', 'pdf', or 'url'
```

Ankify can be run from the command line with the following command:

`go run main.go ankify input.pdf`

## Example

### URL Example

Here is an example of how to use Ankify to create flashcards from Paul Graham's "Need to Read" article (http://www.paulgraham.com/read.html).

`go run main.go -t=url ankify http://www.paulgraham.com/read.html`

```
Q: What makes writing valuable?; 
A: Writing is not just a way to convey ideas, but also a way to have them. A good writer will almost always discover new things in the process of writing, and there is no substitute for this kind of discovery. ; 

Q: What kind of thinking can be done without writing?; 
A: If you don't need to go too deeply into a problem, you can solve it without writing. If the problem can be described formally, it can sometimes be solved in one's head. ; 

Q: What is the importance of being good at reading?; 
A: You can't think well without writing well, and you can't write well without reading well. People who want to have ideas should still read. ; 

Q: What does the author suggest is the best way to solve a complicated, ill-defined problem?; 
A: Writing about it will almost always help to solve a complicated, ill-defined problem. Writing is a valuable tool for discovering new ideas and developing them further.['s']+a+'/r='+r+'/recdata='+csell_page_rec_data.join(',')); }  // Begin Yahoo Store Generated Code

Q: What is an example of good writing?;
A: Reading books and stories written by experienced authors is an example of good writing.;;

Q: What does it mean to be “good at reading”?;
A: Being “good at reading” means being able to extract meaning from the words, rather than simply extracting words from the page.;;

```

### PDF Example

Here is an example of how to use Ankify to create flashcards from a PDF file named `test.pdf`, this pdf contains the seminal paper on MapReduce by Google.

`go run main.go -t=pdf ankify test.pdf`

The output will be a set of Anki questions that you can import into the Anki software to study from.

```
Q: What is MapReduce?; 
A: MapReduce is a programming model and an associated implementation for processing and generating large data sets. It involves a map function that processes a key/value pair to generate a set of intermediate key/value pairs, and a reduce function that merges all intermediate values associated with the same intermediate key.; 

Q: What are the major contributions of MapReduce?; 
A: The major contributions of MapReduce are a simple and powerful interface that enables automatic parallelization and distribution of large-scale computations, combined with an implementation of this interface that achieves high performance on large clusters of commodity PCs.; 


Q: What are the benets of using MapReduce?; 
A: The benets of using MapReduce are that it allows programmers without any experience with parallel and distributed systems to easily utilize the resources of a large distributed system, and it is highly scalable, allowing a typical MapReduce computation to process many terabytes of data on thousands of machines.;


Q: What are two functions of MapReduce?; 
A: The two functions of MapReduce are a map function that processes a key/value pair to generate a set of intermediate key/value pairs, and a reduce function that merges all intermediate values associated with the same intermediate key.; 


Q: What is the purpose of the run-time system in MapReduce?; 
A: The purpose of the run-time system in MapReduce is to take care of the details of partitioning the input data, scheduling the program’s execution across a set of machines, handling machine failures, and managing the required inter-machine communication.
```

## Note

It is important to provide a valid PDF file as input to Ankify. The tool does not currently handle invalid or corrupted PDF files.
