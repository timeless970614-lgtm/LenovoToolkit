//go:build windows

package backend

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// DispdiagCaptureState tracks the state of a dispdiag capture
type DispdiagCaptureState struct {
	IsCapturing  bool   `json:"isCapturing"`
	StartTime    string `json:"startTime"`
	OutputPath   string `json:"outputPath"`
	DurationSecs int    `json:"durationSecs"`
	Status       string `json:"status"`
	Error        string `json:"error"`
}

// DispdiagResult holds parsed results from a dispdiag .dat file
type DispdiagResult struct {
	OutputPath     string            `json:"outputPath"`
	OutputSize     string            `json:"outputSize"`
	CapturedAt     string            `json:"capturedAt"`
	DurationSecs   int               `json:"durationSecs"`
	FileName       string            `json:"fileName"`
	FileContent    string            `json:"fileContent"`
	Summary        DispdiagSummary   `json:"summary"`
	EDIDBlocks     int               `json:"edidBlocks"`
	Displays       []DispdiagDisplay `json:"displays"`
	LinkTraining   []DispdiagLinkStatus `json:"linkTraining"`
	BrightnessInfo []string          `json:"brightnessInfo"`
	Errors         []string          `json:"errors"`
	Warnings       []string          `json:"warnings"`
}

// DispdiagSummary holds the summary info
type DispdiagSummary struct {
	DriverVersion string `json:"driverVersion"`
	BuildVersion  string `json:"buildVersion"`
	DatVersion    string `json:"datVersion"`
	CaptureDate   string `json:"captureDate"`
	PassFail      string `json:"passFail"`
	KeyMessages   []string `json:"keyMessages"`
}

// DispdiagDisplay holds display info
type DispdiagDisplay struct {
	Name      string `json:"name"`
	Location  string `json:"location"`
	EDID      string `json:"edid"`
	Connected bool   `json:"connected"`
}

// DispdiagLinkStatus holds link training status
type DispdiagLinkStatus struct {
	Display    string `json:"display"`
	Status     string `json:"status"`
	Src        string `json:"src"`
	Dest       string `json:"dest"`
	Resolution string `json:"resolution"`
	RefreshRate string `json:"refreshRate"`
}

// RunDispdiag runs dispdiag.exe with given flags and returns the result
func RunDispdiag(outDir string, delaySecs int, dumpMode bool) DispdiagResult {
	if outDir == "" {
		outDir = filepath.Join(os.Getenv("TEMP"), "LenovoToolkit_Dispdiag")
	}

	// Create output directory
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return DispdiagResult{Errors: []string{fmt.Sprintf("Failed to create output dir: %v", err)}}
	}

	timestamp := time.Now().Format("20060102_150405")
	outFile := filepath.Join(outDir, fmt.Sprintf("dispdiag_%s.dat", timestamp))

	result := DispdiagResult{
		OutputPath:   outFile,
		DurationSecs: delaySecs,
	}

	// Build command
	args := []string{}
	if dumpMode {
		args = append(args, "-d")
	}
	if delaySecs > 0 {
		args = append(args, "-delay", fmt.Sprintf("%d", delaySecs))
	}
	args = append(args, "-out", outFile)

	// Run dispdiag (it writes to current dir if not using -out properly)
	// But dispdiag requires -out to be the LAST parameter
	// Also run from the output dir so that temp files land there
	cmd := exec.Command(filepath.Join(os.Getenv("SystemRoot"), "System32", "dispdiag.exe"), args...)
	cmd.Dir = outDir

	start := time.Now()
	output, err := cmd.CombinedOutput()
	elapsed := time.Since(start)

	_ = output // dispdiag stdout contains the output filename path
	result.CapturedAt = start.Format("2006-01-02 15:04:05")
	result.DurationSecs = int(elapsed.Seconds())

	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("dispdiag execution error: %v", err))
		// Check if file was still created
	}

	// Find the actual output file (dispdiag returns the filename on stdout)
	stdoutStr := strings.TrimSpace(string(output))
	if stdoutStr != "" {
		// stdout is usually the path to the generated file
		if strings.HasSuffix(strings.ToLower(stdoutStr), ".dat") {
			result.OutputPath = stdoutStr
		}
	}

	// If the exact outFile doesn't exist, try to find the latest .dat in the dir
	if _, statErr := os.Stat(result.OutputPath); os.IsNotExist(statErr) {
		result.OutputPath = ""
		matches, _ := filepath.Glob(filepath.Join(outDir, "*.dat"))
		latestTime := time.Time{}
		for _, m := range matches {
			if fi, err := os.Stat(m); err == nil {
				if fi.ModTime().After(latestTime) {
					latestTime = fi.ModTime()
					result.OutputPath = m
				}
			}
		}
	}

	result.FileName = filepath.Base(result.OutputPath)

	// Read file content and parse
	parseDispdiagOutput(&result)

	return result
}

// parseDispdiagOutput reads and parses the dispdiag .dat file
func parseDispdiagOutput(r *DispdiagResult) {
	if r.OutputPath == "" {
		r.Errors = append(r.Errors, "No output file found")
		return
	}

	// Get file size
	if fi, err := os.Stat(r.OutputPath); err == nil {
		sizeMB := float64(fi.Size()) / (1024 * 1024)
		r.OutputSize = fmt.Sprintf("%.2f MB", sizeMB)
	}

	// Read file content as hex dump (dispdiag .dat is binary)
	data, err := os.ReadFile(r.OutputPath)
	if err != nil {
		r.Errors = append(r.Errors, fmt.Sprintf("Failed to read output file: %v", err))
		return
	}

	// Generate hex dump (first 4KB) since .dat is binary
	r.FileContent = hexDump(data, 4096)

	// Extract printable text regions for parsing
	textContent := extractPrintable(data)
	parseDispdiagContent(r, textContent)
}

func parseDispdiagContent(r *DispdiagResult, content string) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	lineNum := 0
	var currentSection string

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Detect sections
		switch {
		case strings.Contains(strings.ToLower(trimmed), "header"):
			currentSection = "header"
		case strings.Contains(strings.ToLower(trimmed), "edid"):
			currentSection = "edid"
		case strings.Contains(strings.ToLower(trimmed), "summary"):
			currentSection = "summary"
		case strings.Contains(strings.ToLower(trimmed), "link training"):
			currentSection = "link_training"
		case strings.Contains(strings.ToLower(trimmed), "brightness"):
			currentSection = "brightness"
		case strings.Contains(strings.ToLower(trimmed), "mode status"):
			currentSection = "mode_status"
		}

		// Extract driver version
		if strings.Contains(trimmed, "Driver version") || strings.Contains(trimmed, "DriverVersion") {
			parts := strings.SplitN(trimmed, ":", 2)
			if len(parts) == 2 {
				r.Summary.DriverVersion = strings.TrimSpace(parts[1])
			}
		}

		// Extract build version
		if strings.Contains(trimmed, "Build version") || strings.Contains(trimmed, "BuildVersion") {
			parts := strings.SplitN(trimmed, ":", 2)
			if len(parts) == 2 {
				r.Summary.BuildVersion = strings.TrimSpace(parts[1])
			}
		}

		// Extract dat version
		if strings.Contains(trimmed, "Dat version") || strings.Contains(trimmed, "DatVersion") {
			parts := strings.SplitN(trimmed, ":", 2)
			if len(parts) == 2 {
				r.Summary.DatVersion = strings.TrimSpace(parts[1])
			}
		}

		// Extract key failures
		if strings.Contains(strings.ToLower(trimmed), "fail") || strings.Contains(strings.ToLower(trimmed), "error") {
			r.Errors = append(r.Errors, trimmed)
		}
		if strings.Contains(strings.ToLower(trimmed), "warning") {
			r.Warnings = append(r.Warnings, trimmed)
		}

		// EDID blocks
		if strings.Contains(trimmed, "EDID Block") || strings.Contains(trimmed, "EDID block") {
			r.EDIDBlocks++
		}

		// Link training pass/fail
		if strings.Contains(strings.ToLower(trimmed), "pass") {
			r.Summary.PassFail = "PASS"
		}
		if strings.Contains(strings.ToLower(trimmed), "fail") && currentSection == "link_training" {
			r.Summary.PassFail = "FAIL"
			link := DispdiagLinkStatus{Status: trimmed}
			r.LinkTraining = append(r.LinkTraining, link)
		}

		// Brightness info
		if currentSection == "brightness" && trimmed != "" {
			r.BrightnessInfo = append(r.BrightnessInfo, trimmed)
		}

		// Key messages / summary lines
		if currentSection == "summary" && trimmed != "" && len(r.Summary.KeyMessages) < 20 {
			r.Summary.KeyMessages = append(r.Summary.KeyMessages, trimmed)
		}
	}

	_ = scanner.Err()
}

// ExportDispdiagResult saves the dispdiag result as JSON for further analysis
func ExportDispdiagResult(result DispdiagResult, outputPath string) string {
	if outputPath == "" {
		outputPath = filepath.Join(os.Getenv("TEMP"), "LenovoToolkit_Dispdiag",
			fmt.Sprintf("dispdiag_analysis_%s.json", time.Now().Format("20060102_150405")))
	}

	// Ensure dir exists
	os.MkdirAll(filepath.Dir(outputPath), 0755)

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Sprintf("Export failed: %v", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Sprintf("Export failed: %v", err)
	}

	return outputPath
}

// GetDispdiagOutputDir returns the default output directory for dispdiag
func GetDispdiagOutputDir() string {
	return filepath.Join(os.Getenv("TEMP"), "LenovoToolkit_Dispdiag")
}

// OpenDispdiagLog finds and opens the latest dispdiag log file
func OpenDispdiagLog() string {
	outDir := GetDispdiagOutputDir()

	// Find the latest .dat file
	var latestPath string
	var latestTime time.Time

	// Check .dat files
	datMatches, _ := filepath.Glob(filepath.Join(outDir, "*.dat"))
	for _, m := range datMatches {
		if fi, err := os.Stat(m); err == nil {
			if fi.ModTime().After(latestTime) {
				latestTime = fi.ModTime()
				latestPath = m
			}
		}
	}

	// Also check .txt files (from -tee output)
	txtMatches, _ := filepath.Glob(filepath.Join(outDir, "dispdiag.txt"))
	for _, m := range txtMatches {
		if fi, err := os.Stat(m); err == nil {
			if fi.ModTime().After(latestTime) {
				latestTime = fi.ModTime()
				latestPath = m
			}
		}
	}

	// Also check .json files (from export)
	jsonMatches, _ := filepath.Glob(filepath.Join(outDir, "*.json"))
	for _, m := range jsonMatches {
		if fi, err := os.Stat(m); err == nil {
			if fi.ModTime().After(latestTime) {
				latestTime = fi.ModTime()
				latestPath = m
			}
		}
	}

	if latestPath == "" {
		return "No dispdiag log file found. Run Dispdiag first."
	}

	// Open Explorer and select the latest file
	cmd := exec.Command("explorer", "/select,"+latestPath)
	if err := cmd.Start(); err != nil {
		return fmt.Sprintf("Failed to open log: %v", err)
	}

	return latestPath
}

// hexDump produces a classic hex dump of binary data (limited to maxBytes)
func hexDump(data []byte, maxBytes int) string {
	if len(data) > maxBytes {
		data = data[:maxBytes]
	}
	var sb strings.Builder
	sb.WriteString("         0  1  2  3  4  5  6  7  8  9  A  B  C  D  E  F  |0123456789ABCDEF|\n")
	sb.WriteString("--------- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- |----------------|\n")
	for i := 0; i < len(data); i += 16 {
		sb.WriteString(fmt.Sprintf("%08X ", i))
		line := data[i:min(i+16, len(data))]
		// Hex part
		hexParts := make([]string, 16)
		ascii := make([]byte, 16)
		for j := 0; j < 16; j++ {
			if j < len(line) {
				hexParts[j] = fmt.Sprintf("%02X", line[j])
				if line[j] >= 32 && line[j] < 127 {
					ascii[j] = line[j]
				} else {
					ascii[j] = '.'
				}
			} else {
				hexParts[j] = "  "
				ascii[j] = ' '
			}
		}
		sb.WriteString(strings.Join(hexParts, " "))
		sb.WriteString(" |")
		sb.Write(ascii)
		sb.WriteString("|\n")
	}
	if len(data) >= maxBytes {
		sb.WriteString(fmt.Sprintf("\n... (hex dump truncated at %d bytes, %d total)\n", maxBytes, maxBytes))
	}
	return sb.String()
}

// extractPrintable extracts only printable ASCII (32-126) and common whitespace from binary data
func extractPrintable(data []byte) string {
	var sb strings.Builder
	for _, b := range data {
		if (b >= 32 && b < 127) || b == '\n' || b == '\r' || b == '\t' {
			sb.WriteByte(b)
		}
	}
	return sb.String()
}
