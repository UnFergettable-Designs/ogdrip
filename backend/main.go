package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// Check if we're running in API service mode
	serviceMode := flag.Bool("service", false, "Run in API service mode")
	
	// Parse command-line flags
	flag.Parse()
	
	// Determine which mode to run in
	if *serviceMode {
		// Run in API service mode
		fmt.Println("Starting Open Graph API service...")
		ServiceMain()
	} else if flag.NArg() == 0 && len(os.Args) > 1 {
		// If no positional arguments and we have some flags, run the generator
		fmt.Println("Running Open Graph generator...")
		ServerMain()
	} else {
		// Display usage information
		fmt.Println("Open Graph Generator/API")
		fmt.Println("=======================")
		fmt.Println("Usage:")
		fmt.Println("  1. Run as generator: og-generator -url=https://example.com -output=output.png")
		fmt.Println("  2. Run as API service: og-generator -service")
		fmt.Println("")
		fmt.Println("For generator options, run: og-generator -help")
		os.Exit(1)
	}
} 