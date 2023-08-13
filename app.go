package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"


	"github.com/weswest/431wk8svelte/backend/vale-2.28.1/internal/core"
	"github.com/weswest/431wk8svelte/backend/vale-2.28.1/internal/lint"
	)
	
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's talk time!", name)
}

func (a *App) TestVale(text string) (string, error) {
	fmt.Println("TestVale function called with input:", text)
	return "TestVale Call", nil
}

func (a *App) CheckWithValePartial(inputText, ruleIs string) (string, error) {
	tmpFile, err := os.CreateTemp("", "vale-*.txt")
	if err != nil {
		log.Fatalf("Failed to create temp file: %s", err)
	}
	defer os.Remove(tmpFile.Name()) // clean up

	// Write the content to the temporary file
	_, err = tmpFile.Write([]byte(inputText))
	if err != nil {
		log.Fatalf("Failed to write to temp file: %s", err)
	}
	tmpFile.Close()

	// Base command arguments
	cmdArgs := []string{tmpFile.Name()}

	// Determine the appropriate command arguments based on ruleIs
	switch ruleIs {
	case "singular":
		cmdArgs = append(cmdArgs, "--config=styles/DataIsSingular/.vale.ini")
	case "plural":
		cmdArgs = append(cmdArgs, "--config=styles/DataIsPlural/.vale.ini")
	case "all":
		cmdArgs = append(cmdArgs, "--config=.vale.ini")
	default:
		return "", fmt.Errorf("Invalid ruleIs value: %s", ruleIs)
	}

	// Create the command.
	cmd := exec.Command("vale", cmdArgs...)

	// Run the command and capture the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitStatus := exitError.ExitCode()
			if exitStatus == 1 {
				// Vale found issues, log them without stopping the program.
				log.Println("Vale found issues.  See below.")
			} else {
				// Handle other non-zero exit statuses.
				log.Println("Command failed with non-content error:", exitStatus)
				return string(exitStatus), nil
			}
		} else {
			// Handle other types of errors.
			log.Println("Failed to run command:", err)
			thing := fmt.Sprintf("%s", err)
			return thing, nil
		}
	}

	if len(output) == 0 {
		return "Neener", nil
	}

	return "Cmd Run and Output Created", nil
}

func (a *App) CheckWithVale(inputText, ruleIs string) (string, error) {

	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "vale-*.txt")
	if err != nil {
		log.Fatalf("Failed to create temp file: %s", err)
	}
	defer os.Remove(tmpFile.Name()) // clean up

	// Write the content to the temporary file
	_, err = tmpFile.Write([]byte(inputText))
	if err != nil {
		log.Fatalf("Failed to write to temp file: %s", err)
	}
	tmpFile.Close()

	// Base command arguments
	cmdArgs := []string{tmpFile.Name()}

	// Determine the appropriate command arguments based on ruleIs
	switch ruleIs {
	case "singular":
		cmdArgs = append(cmdArgs, "--config=styles/DataIsSingular/.vale.ini")
	case "plural":
		cmdArgs = append(cmdArgs, "--config=styles/DataIsPlural/.vale.ini")
	case "all":
		cmdArgs = append(cmdArgs, "--config=.vale.ini")
	default:
		return "", fmt.Errorf("Invalid ruleIs value: %s", ruleIs)
	}

	// Create the command.
	cmd := exec.Command("vale", cmdArgs...)

	// Run the command and capture the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitStatus := exitError.ExitCode()
			if exitStatus == 1 {
				// Vale found issues, log them without stopping the program.
				log.Println("Vale found issues.  See below.")
			} else {
				// Handle other non-zero exit statuses.
				log.Fatalf("Command failed with non-content error:", exitStatus)
			}
		} else {
			// Handle other types of errors.
			log.Fatalf("Failed to run command:", err)
		}
	}

	segments := extractSegments(output)
	processedSegments := make([]string, len(segments))
	for i, segment := range segments {
		processedSegments[i] = processSegment(segment)
	}
	resultsForPrinting := strings.Join(processedSegments, "\n")
	if len(resultsForPrinting) == 0 {
		resultsForPrinting = "You are perfect just the way you are.  No errors detected."
	}
	return resultsForPrinting, nil
}

// Helper function to structure the vale output correctly

func stripANSI(input string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	return re.ReplaceAllString(input, "")
}

func trimText(text string) string {
	trimmedStr := strings.TrimSpace(text)

	// Remove extra spaces between words
	words := strings.Fields(trimmedStr)
	cleanStr := strings.Join(words, " ")
	return cleanStr
}

func processSegment(segment string) string {
	testLines := strings.Split(segment, "\n")

	//var area, level, rule, description string
	var area, level, description string

	for i, line := range testLines {
		trimmedTokens := strings.Split(trimText(line), " ")
		if i == 0 {
			area = trimmedTokens[0]
			level = stripANSI(trimmedTokens[1])
			//			rule = trimmedTokens[len(trimmedTokens)-1]
			description = strings.Join(trimmedTokens[2:len(trimmedTokens)-1], " ")
		} else {
			description += " " + strings.Join(trimmedTokens, " ")
		}
	}
	//	return fmt.Sprintf("(%s) %s:%s. %s", area, level, rule, description)
	return fmt.Sprintf("(%s) %s: %s", area, level, description)
}

func extractSegments(output []byte) []string {
	// Convert byte to string
	outputStr := string(output)

	// Split the string into lines
	lines := strings.Split(outputStr, "\n")

	// Check if there are at least 5 lines to remove
	if len(lines) >= 5 {
		// Remove the first two lines and the last three lines
		lines = lines[2 : len(lines)-3]
	}

	// Reassemble into text
	text := strings.Join(lines, "\n")

	// Use a regex pattern to match the X:XX pattern
	pattern := `(\d+:\d+)`
	re := regexp.MustCompile(pattern)

	// Find all matches
	matches := re.FindAllStringIndex(text, -1)

	// Split the text based on the matches
	var segments []string
	start := 0
	for _, match := range matches {
		end := match[0]
		segments = append(segments, text[start:end])
		start = end
	}
	segments = append(segments, text[start:])
	if len(segments) > 0 {
		segments = segments[1:]
	}

	return segments
}
