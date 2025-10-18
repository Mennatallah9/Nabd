package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	fmt.Println("🧪 Nabd Backend Test Suite")
	fmt.Println("==========================")
	fmt.Println()

	testPackages := []string{
		"./tests/utils",
		"./tests/models", 
		"./tests/controllers",
		"./tests/services",
	}

	var totalTests, passedTests int
	startTime := time.Now()

	for _, pkg := range testPackages {
		fmt.Printf("Running %s...\n", pkg)
		
		cmd := exec.Command("go", "test", "-v", pkg)
		output, err := cmd.CombinedOutput()
		
		if err != nil {
			fmt.Printf("❌ %s: FAILED\n", pkg)
			fmt.Printf("Error: %v\n", err)
			fmt.Printf("Output: %s\n", string(output))
		} else {
			// Count tests
			lines := strings.Split(string(output), "\n")
			pkgTests := 0
			pkgPassed := 0
			
			for _, line := range lines {
				if strings.Contains(line, "=== RUN") {
					pkgTests++
				}
				if strings.Contains(line, "--- PASS:") {
					pkgPassed++
				}
			}
			
			totalTests += pkgTests
			passedTests += pkgPassed
			
			fmt.Printf("✅ %s: %d/%d tests passed\n", pkg, pkgPassed, pkgTests)
		}
		fmt.Println()
	}

	duration := time.Since(startTime)
	
	fmt.Println("📊 Test Summary")
	fmt.Println("===============")
	fmt.Printf("Total Tests: %d\n", totalTests)
	fmt.Printf("Passed: %d\n", passedTests)
	fmt.Printf("Failed: %d\n", totalTests-passedTests)
	fmt.Printf("Duration: %v\n", duration.Round(time.Millisecond))
	
	if passedTests == totalTests {
		fmt.Println("🎉 All tests passed!")
		os.Exit(0)
	} else {
		fmt.Println("❌ Some tests failed!")
		os.Exit(1)
	}
}