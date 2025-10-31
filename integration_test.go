package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// TestIntegration discovers and runs all integration test cases in testdata/
func TestIntegration(t *testing.T) {
	binaryPath := buildBinary(t)
	defer os.Remove(binaryPath)

	testCases := discoverTestCases(t)
	if len(testCases) == 0 {
		t.Skip("No test cases found in testdata/")
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runTestCase(t, binaryPath, tc)
		})
	}
}

// testCase represents a single integration test case
type testCase struct {
	name         string
	dir          string
	inputFile    string
	expectedFile string
	argsFile     string
}

// buildBinary builds the dockerfile-json binary and returns its path
func buildBinary(t *testing.T) string {
	t.Helper()

	binaryPath := filepath.Join(t.TempDir(), "dockerfile-json")
	cmd := exec.Command("go", "build", "-tags=dfrunsecurity", "-o", binaryPath, ".")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build binary: %v\nOutput: %s", err, output)
	}

	return binaryPath
}

// discoverTestCases scans testdata/ for test case directories.
// A scenario directory must contain a Containerfile and is always a test case.
// All subdirectories are also test cases, using the parent's Containerfile.
func discoverTestCases(t *testing.T) []testCase {
	t.Helper()

	testdataDir := "testdata"
	entries, err := os.ReadDir(testdataDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		t.Fatalf("Failed to read testdata directory: %v", err)
	}

	var cases []testCase
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		scenarioDir := filepath.Join(testdataDir, entry.Name())
		containerfile := filepath.Join(scenarioDir, "Containerfile")

		// Scenario directory must have a Containerfile
		if _, err := os.Stat(containerfile); os.IsNotExist(err) {
			continue
		}

		// Main scenario
		cases = append(cases, testCase{
			name:         entry.Name(),
			dir:          scenarioDir,
			inputFile:    containerfile,
			expectedFile: filepath.Join(scenarioDir, "expected.json"),
			argsFile:     filepath.Join(scenarioDir, "args.txt"),
		})

		// Sub-scenarios
		subEntries, err := os.ReadDir(scenarioDir)
		if err != nil {
			t.Fatalf("Failed to read scenario directory %s: %v", scenarioDir, err)
		}

		for _, sub := range subEntries {
			if !sub.IsDir() {
				continue
			}

			subDir := filepath.Join(scenarioDir, sub.Name())
			cases = append(cases, testCase{
				name:         entry.Name() + "/" + sub.Name(),
				dir:          subDir,
				inputFile:    containerfile,
				expectedFile: filepath.Join(subDir, "expected.json"),
				argsFile:     filepath.Join(subDir, "args.txt"),
			})
		}
	}

	return cases
}

// runTestCase executes a single test case
func runTestCase(t *testing.T, binaryPath string, tc testCase) {
	t.Helper()

	args := readArgs(t, tc.argsFile)
	args = append(args, tc.inputFile)

	// Execute the binary
	cmd := exec.Command(binaryPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to execute binary: %v\nOutput: %s", err, output)
	}

	// Update mode: write output to expected.json
	// Use UPDATE_TESTDATA=1 or UPDATE_TESTDATA=true to update golden files
	updateTestdata := os.Getenv("UPDATE_TESTDATA")
	if updateTestdata == "1" || strings.EqualFold(updateTestdata, "true") {
		if err := updateGoldenFile(tc.expectedFile, output); err != nil {
			t.Fatalf("Failed to update golden file: %v", err)
		}
		t.Logf("Updated %s", tc.expectedFile)
		return
	}

	// Comparison mode: compare with expected.json
	expected, err := os.ReadFile(tc.expectedFile)
	if err != nil {
		t.Fatalf("Failed to read expected file %s: %v\nRun with UPDATE_TESTDATA=1 to generate it", tc.expectedFile, err)
	}
	if err := compareJSON(t, expected, output); err != nil {
		t.Fatalf("Output mismatch:\n%v", err)
	}
}

// readArgs reads arguments from args.txt file (one per line)
func readArgs(t *testing.T, argsFile string) []string {
	t.Helper()

	data, err := os.ReadFile(argsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		t.Fatalf("Failed to read args file: %v", err)
	}

	lines := strings.Split(string(data), "\n")
	var args []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			args = append(args, line)
		}
	}

	return args
}

// updateGoldenFile writes the output to the golden file with pretty formatting.
// It formats the JSON without unmarshaling to preserve deterministic key ordering
// from the binary's struct field order.
func updateGoldenFile(path string, data []byte) error {
	var buf bytes.Buffer
	if err := json.Indent(&buf, data, "", "  "); err != nil {
		return err
	}
	return os.WriteFile(path, buf.Bytes(), 0644)
}

// compareJSON compares two JSON byte slices by pretty-printing them
// and comparing as strings. This provides more readable diffs.
func compareJSON(t *testing.T, expected, actual []byte) error {
	t.Helper()

	var expectedBuf bytes.Buffer
	if err := json.Indent(&expectedBuf, expected, "", "  "); err != nil {
		return fmt.Errorf("failed to format expected JSON: %w", err)
	}

	var actualBuf bytes.Buffer
	if err := json.Indent(&actualBuf, actual, "", "  "); err != nil {
		return fmt.Errorf("failed to format actual JSON: %w", err)
	}

	if diff := cmp.Diff(expectedBuf.String(), actualBuf.String()); diff != "" {
		return fmt.Errorf("JSON mismatch (-expected +actual):\n%s", diff)
	}

	return nil
}
