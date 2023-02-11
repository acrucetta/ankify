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

Here is an example of how to use Ankify to create flashcards from a PDF file named `sample.pdf`:

`go run main.go ankify sample.pdf`

The output will be a set of Anki questions that you can import into the Anki software to study from.

## Note

It is important to provide a valid PDF file as input to Ankify. The tool does not currently handle invalid or corrupted PDF files.
