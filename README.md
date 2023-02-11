# Ankify

Ankify is a command line tool written in Golang that reads a PDF document and outputs Anki questions for the user. The project uses Cobra CLI to provide a user-friendly interface for creating flashcards from your favorite PDFs.

## Installation

To install Ankify, you need to have Golang installed on your system. Clone the repository to your local machine, navigate to the project directory, and run the following command:

`go build`

This will compile the code and create an executable file named `ankify`.

## Usage

Ankify can be run from the command line with the following two main commands:

`go run main.go ankify input.pdf`

This command creates Anki questions from the PDF file `input.pdf`.

`go run main.go parse input.pdf`

This command parses the PDF file `input.pdf` and outputs the text to the terminal.

## Example

Here is an example of how to use Ankify to create flashcards from a PDF file named `test.pdf`, this pdf contains the seminal paper on MapReduce by Google.

`go run main.go ankify test.pdf`

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
