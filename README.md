# Overview of Project

This is an assignment for Week 8 in the Northwestern Masters in Data Science [MSDS-431 Data Engineering with Go](https://msdsgo.netlify.app/data-engineering-with-go/) course.

The purpose of this assignment was to build a desktop application using wails, svelte, and vale that can take in text and run the vale linter on it to identify cases where "data" is incorrectly being referenced as singular or plural.

## TL;DR: Project Isn't Finished

Unfortunately, this assignment is currently incomplete.  The program is fully built, and runs just fine in the wails dev environment, but I haven't yet figured out how to get the program that is built from wails build to run successfully.

The root of the issue is that vale is a command-line tool that is distributed as an executable.  Its use is straightforward: navigate to a folder with a file you want to lint and type:

```bash
vale <filename>
```

And vale will run the rules you've specified and produce some attractive output in your terminal.

The problem comes in, though, with running the linter as part of a go program, and in particular as part of a wails-built program.  See... the dev and build environments have different permissions.

I wrote go code that mimicked a command-line input, so the go code would execute the vale command, process the outputs, and return it to the wails/svelte frontend for display.  This, however, was premised on having global access to the "vale" command.

The final built app has a different set of permissions and *doesn't* have global access to the vale command.  So my work of getting go to execute the command line instruction, ingest the result, and then deliver the result to svelte no longer worked.

## Things Tried

I was 15 hours in to this assignment getting the different pieces to work before I learned about this wails dev/wails build difference.  I spent another 12 hours on the following alternatives trying to get this to work, to no avail.  So as I shared in the video, I just simply ran out of time and am submitting what I have.

None of the following options worked for me:

1. **Access vale as a go package and access methods directly**
    * Goal: The code was all written in go so why not just access it without going through the command line?
    * Problem: All of the important code is located in an internal folder, and programmatically go prohibits importing of "internal" code
2. **Copy the entirety of the vale program into my project folder and access methods directly.**
    * Goal: Skip the rigamarole of trying to run vale as a command-line program and do it all through go functions
    * Problem: I ran into program-breaking issues with the package imports.  I even went through and updated every *.go file in the entire vale project and couldn't get it to work
    * Problem secondary: I felt a little dirty copying their entire program
3. **Package the vale program as part of the wails executable.**
    * Goal: Put the mac and windows versions in a "binaries" folder, have them be available to the program after compiling, and then figure out the user's system and then identify the right vale program to call.
    * Problem: I was unable to figure out how to adjust the wails settings to include the executables as part of the program.
4. **Package the vale program as text within the go program and then dynamically create binaries as needed.**
    * Goal: Get around the issues with item #3 by packaging the stuff directly within the go program.  Note this actually worked, and ballooned the size of the application since it had four separate vale programs embedded in it
    * Problem: I couldn't figure out how to get the binaries - read in as bytes - to be processed as a program at the time it was needed
5. **Hard-code the vale path so the executable knew exactly where vale was.**
    * Goal: Rather than having a relative path, just presume that the user has vale already installed and can point the computer to vale specifically
    * Problem Main: I tried to create a hard-linked path on my personal machine, and the program was still unable to run vale in build (but could still run it in dev)
    * Problem Secondary: I tried to write a function that would run "which vale" or "where vale", as appropriate, but the function which worked in the dev environment failed to run correctly in the built application
6. **Explore whether I was doing the wrong thing writing the text to a temp file for read-in thru vale.**
    * Goal: Change the structure of the vale command to not need to read in a file.
    * Problem: The read-in wasn't the part that was broken: the finding vale part was broken.  I had low hopes for this because all of the erroring really pointed to vale itself not being found, but I figured I'd try this anyway


# Project Structure

This project is comprised of three different pieces:

1. The wails/svelte frontend
2. Bespoke vale rules to identify whether "data" is treated as plural or singular
3. The golang-based backend.

## 1. Wails/Svelte Frontend

The frontend is structured so there is a text input box, an "Evaluate" button, and then the result of the vale linter printed below.  There is also a radio toggle that allows one to identify whether "data" should be treated as singular or plural, as well as a debugging option to print results from both sets of rules.

I enriched this view with the following:
1. A file upload button so that one could choose .md or .txt files.  These got read in as text and then dropped in the input box
2. A handful of samples that, when clicked, would also populate the input box.  Very helpful for testing.

## 2. Bespoke Vale rules

Vale creates its rules via yaml files which are then managed via .vale.ini codes.  I was unable to come up with a programmatic / generalized way of having vale look at a piece of text and identify whether "data" was being used as a singular or as a plural form of the word, so I tackled this by creating a set of hard-written rules.

One set of rules was created as a "DataIsSingular" bundle and another as "DataIsPlural".  These were manually kept in alignment/opposition.

The svelte frontend passed the selection of singular/plural to the go backend.  The selection then created some forking paths for whether the Plural or Singular versions of the code was called.

## 3. Golang Backend

The go code was written to receive the text from the svelte frontend, then run "vale \<text>" as a command line (it's a little more complicated than that, especially with the branching vale rulesets, but this is the gist).

The output of the vale code then required substantial processing to be passed back to the svelte frontend.  Shrug.

As identified above, this approach worked fine when using the "wails dev" dev environment, which had the permissions to run "vale \<text>".  However, this approach didn't work in the wails build environment.


# FYI - assignment details motivating this work

### Management Problem

A technology startup sees assisted writing as a market opportunity. Current products such as wordtune and grammarly are not sufficiently extensible to serve the needs of professionals who have their own writing styles. 

As a marketing manager and potential user puts it:

*I can get a Grammarly to check my grammar and spelling. I can use ChatGPT to ghost-write if I like, although that's a scary proposition. What I can't do is get a program to write like me. Sure, I want to use correct grammar, but I don't want to sound like everyone else. I want to sound like me.*

To attract potential investors, the startup needs to build a product prototype. Members of the startup want show investors that it is possible to develop an extensible writing assistant and that their company can do the development. 

Having seen Joseph Kato's Vale Links to an external site., a syntax-aware linter for prose, data engineers at the startup believe they can use Go to build a product prototype. It should be a desktop application that demonstrates the ability to be tailored to the styles of individual writers.

### Assignment Requirements 

Take on the role of the company's data engineers in building a product prototype. It should be a desktop application that can be distributed as a single executable load module for Windows 10/11, Apple Mac, or Linux computers. It should be capable of reading a plain text file of at least 200 words and checking for spelling and grammatical errors. But, more importantly, it should be customizable, with the ability to adapt to the writing style of each individual user/writer. 

The application should be developed in using the Wails desktop development platform with a Svelte frontend. Use Vale as the prose analytics engine on the Go backend.

To demonstrate the ability to conform to the styles of individual writers, the application should ask users whether they want to treat "data" as a singular or plural noun. (Note that proper English has "data" as a plural noun. The singular of "data" is "datum."). When presented with a document in which the word "data" is used incorrectly (as defined by user preference), the sentence containing the word should be flagged as erroneous.

### Grading Guidelines (100 Total Points)

* Coding rules, organization, and aesthetics (20 points). Effective use of Go modules and idiomatic Go. Code should be readable, easy to understand. Variable and function names should be meaningful, specific rather than abstract. They should not be too long or too short. Avoid useless temporary variables and intermediate results. Code blocks and line breaks should be clean and consistent. Break large code blocks into smaller blocks that accomplish one task at a time. Utilize readable and easy-to-follow control flow (if/else blocks and for loops). Distribute the not rather than the switch (and/or) in complex Boolean expressions. Programs should be self-documenting, with comments explaining the logic behind the code (McConnell 2004, 777–817).
* Testing and software metrics (20 points). Employ unit tests of critical components, generating synthetic test data when appropriate. Generate program logs and profiles when appropriate. Monitor memory and processing requirements of code components and the entire program. If noted in the requirements definition, conduct a Monte Carlo performance benchmark.
* Design and development (20 points). Employ a clean, efficient, and easy-to-understand design that meets all aspects of the requirements definition and serves the use case. When possible, develop general-purpose code modules that can be reused in other programming projects.
* Documentation (20 points). Effective use of Git/GitHub, including a README.md Markdown file for each repository, noting the roles of programs and data and explaining how to test and use the application.
* Application (20 points). Delivery of an executable load module or application (.exe file for Windows or .app file for MacOS). The application should run to completion without issues. If user input is required, the application should check for valid/usable input and should provide appropriate explanation to users who provide incorrect input. The application should employ clean design for the user experience and user interface (UX/UI).

### Assignment Deliverables

* Text string showing the link (URL) to the GitHub repository for the assignment
* README.md Markdown text file documentation for the assignment
* Zip compressed file with executable load module for the program/application (.exe for Windows or .app for MacOS)