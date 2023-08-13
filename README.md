# Project Overview

The purpose of this assignment was to build a desktop application using wails, svelte, and vale that can take in text and run the vale linter on it to identify cases where "data" is incorrectly being referenced as singular or plural.



## Things Tried

1. **Access vale as a go package and access methods directly**
    * Goal: The code was all written in go so why not just access it without going through the command line?
    * Problem: All of the important code is located in an internal folder, and programmatically go prohibits importing of "internal" code
2. **Copy the entirety of the vale program into my project folder and access methods directly.**
    * Goal: Skip the rigamarole of trying to run vale as a command-line program and do it all through go functions
    * Problem: I'm not smart enough to figure out all of the fiddly details that were necessary to get the program to execute.  The vale people are much smarter than me
    * Problem secondary: I felt a little dirty copying their entire program
3. **Package the vale program as part of the wails executable.**
    * Goal: Put the mac and windows versions in a "binaries" folder, have them be available to the program after compiling, and then figure out the user's system and then identify the right vale program to call.
    * Problem: I was unable to figure out how to adjust the wails settings to include the executables as part of the program.
4. **Package the vale program as text within the go program and then dynamically create binaries as needed.**
    * Goal: Get around the issues with item #1 by packaging the stuff directly within the go program.  Note this actually worked, and ballooned the size of the application since it had four separate vale programs embedded in it
    * Problem: I couldn't figure out how to get the binaries - read in as bytes - to be processed as a program at the time it was needed
5. **Hard-code the vale path so the executable knew exactly where vale was.**
    * Goal: Rather than having a relative path, just presume that the user has vale already installed and can point the computer to vale specifically
    * Problem Main: I tried to create a hard-linked path on my personal machine, and for reasons I ran out of time to debug, the built application still couldn't execute the vale program
    * Problem Secondary: I tried to write a function that would run "which vale" or "where vale", as appropriate, but the function which worked in the dev environment failed to run correctly in the built application
