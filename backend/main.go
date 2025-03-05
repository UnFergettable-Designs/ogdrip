package main

import (
	"flag"
	"fmt"
	"log"
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
	} else {
		// Run in CLI mode
		fmt.Println("This is a placeholder for CLI functionality.")
		fmt.Println("To start the API service, use the -service flag.")
	}
} 