// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Hello is a simple hello, world demonstration web server.
//
// It serves version information on /version and answers
// any other request like /name by saying "Hello, name!".
//
// See golang.org/x/example/outyet for a more sophisticated server.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	// Define command-line flags
	webpageURL := flag.String("url", "", "Webpage URL to capture")
	outputPath := flag.String("output", "screenshot.png", "Output file path for the screenshot")
	quality := flag.Int("quality", 90, "Screenshot quality (0-100)")

	// Parse command-line flags
	flag.Parse()

	// Validate input
	if *webpageURL == "" {
		fmt.Println("Please provide a webpage URL using the -url flag.")
		os.Exit(1)
	}

	// Create a Chrome context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Define a custom action to wait for network idle
	// waitNetworkIdle := chromedp.ActionFunc(func(ctx context.Context) error {
	// 	// JavaScript code to check if the network is idle
	// 	script := `
	// 		new Promise(resolve => {
	// 			if (window.performance && window.performance.getEntriesByType) {
	// 				const resources = window.performance.getEntriesByType('resource');
	// 				if (resources.length === 0) {
	// 					resolve(true);
	// 					return;
	// 				}
	// 				let lastResourceTimestamp = resources[resources.length - 1].responseEnd;
	// 				const interval = setInterval(() => {
	// 					const resources = window.performance.getEntriesByType('resource');
	// 					if (resources.length === 0 || resources[resources.length - 1].responseEnd === lastResourceTimestamp) {
	// 						clearInterval(interval);
	// 						resolve(true);
	// 					} else {
	// 						lastResourceTimestamp = resources[resources.length - 1].responseEnd;
	// 					}
	// 				}, 100); // Check every 100ms
	// 			} else {
	// 				resolve(true); // Fallback if performance API is not available
	// 			}
	// 		});
	// 	`
	// 	var result interface{}
	// 	err := chromedp.Run(ctx, chromedp.Evaluate(script, &result))
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// Check if result is a boolean and true
	// 	if boolValue, ok := result.(bool); ok && boolValue {
	// 		return nil // Network is idle
	// 	} else {
	// 		return fmt.Errorf("network not idle, result: %v", result)
	// 	}
	// })

	// Capture screenshot
	var buf []byte
	var htmlContent string
	if err := chromedp.Run(ctx,
		chromedp.Navigate(*webpageURL),
		chromedp.Sleep(1000*time.Millisecond),
		chromedp.OuterHTML("html", &htmlContent, chromedp.ByQuery),
		chromedp.FullScreenshot(&buf, int(*quality)),
	); err != nil {
		log.Fatal(err)
	}

	// Save screenshot to file
	if err := os.WriteFile(*outputPath, buf, 0644); err != nil {
		log.Fatal(err)
	}

	fmt.Println(htmlContent)
	fmt.Printf("Screenshot saved to %s\n", *outputPath)
}
