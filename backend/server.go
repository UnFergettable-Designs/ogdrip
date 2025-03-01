package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

// OpenGraphData represents the data needed for Open Graph meta tags
type OpenGraphData struct {
	Title       string
	Description string
	ImageURL    string
	PageURL     string
	Type        string
	SiteName    string
	ImageWidth  int
	ImageHeight int
	TwitterCard string
	LocalImage  bool // If true, ImageURL is a local path
}

// Default values
const (
	defaultImageWidth  = 1200
	defaultImageHeight = 630
	defaultType        = "website"
	defaultTwitterCard = "summary_large_image"
)

// generateMetaTags creates HTML with Open Graph meta tags
func generateMetaTags(data OpenGraphData) string {
	metaTemplate := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    
    <!-- Open Graph / Facebook -->
    <meta property="og:type" content="{{.Type}}">
    <meta property="og:url" content="{{.PageURL}}">
    <meta property="og:title" content="{{.Title}}">
    <meta property="og:description" content="{{.Description}}">
    <meta property="og:image" content="{{.ImageURL}}">
    {{if .SiteName}}<meta property="og:site_name" content="{{.SiteName}}">{{end}}
    <meta property="og:image:width" content="{{.ImageWidth}}">
    <meta property="og:image:height" content="{{.ImageHeight}}">
    
    <!-- Twitter -->
    <meta name="twitter:card" content="{{.TwitterCard}}">
    <meta name="twitter:url" content="{{.PageURL}}">
    <meta name="twitter:title" content="{{.Title}}">
    <meta name="twitter:description" content="{{.Description}}">
    <meta name="twitter:image" content="{{.ImageURL}}">
    
    <!-- LinkedIn -->
    <meta name="linkedin:title" content="{{.Title}}">
    <meta name="linkedin:description" content="{{.Description}}">
    <meta name="linkedin:image" content="{{.ImageURL}}">
    
    <!-- Additional helpful meta tags -->
    <meta name="description" content="{{.Description}}">
</head>
<body>
    <h1>Open Graph Preview for: {{.Title}}</h1>
    <p>This page contains the Open Graph meta tags for your content.</p>
    
    <div style="margin: 20px 0;">
        <h2>Preview:</h2>
        <div style="border: 1px solid #ccc; border-radius: 8px; overflow: hidden; max-width: 600px;">
            <img src="{{.ImageURL}}" style="width: 100%; height: auto;" alt="Open Graph preview image">
            <div style="padding: 15px;">
                <h3 style="margin: 0 0 10px; font-size: 18px;">{{.Title}}</h3>
                <p style="margin: 0; color: #666; font-size: 14px;">{{.Description}}</p>
                <p style="margin: 5px 0 0; color: #999; font-size: 12px;">{{.PageURL}}</p>
            </div>
        </div>
    </div>
    
    <div style="margin: 30px 0;">
        <h2>HTML Code:</h2>
        <pre style="background: #f4f4f4; padding: 15px; border-radius: 5px; overflow: auto;"><code>
&lt;!-- Open Graph / Facebook --&gt;
&lt;meta property="og:type" content="{{.Type}}"&gt;
&lt;meta property="og:url" content="{{.PageURL}}"&gt;
&lt;meta property="og:title" content="{{.Title}}"&gt;
&lt;meta property="og:description" content="{{.Description}}"&gt;
&lt;meta property="og:image" content="{{.ImageURL}}"&gt;
{{if .SiteName}}&lt;meta property="og:site_name" content="{{.SiteName}}"&gt;{{end}}
&lt;meta property="og:image:width" content="{{.ImageWidth}}"&gt;
&lt;meta property="og:image:height" content="{{.ImageHeight}}"&gt;

&lt;!-- Twitter --&gt;
&lt;meta name="twitter:card" content="{{.TwitterCard}}"&gt;
&lt;meta name="twitter:url" content="{{.PageURL}}"&gt;
&lt;meta name="twitter:title" content="{{.Title}}"&gt;
&lt;meta name="twitter:description" content="{{.Description}}"&gt;
&lt;meta name="twitter:image" content="{{.ImageURL}}"&gt;

&lt;!-- LinkedIn --&gt;
&lt;meta name="linkedin:title" content="{{.Title}}"&gt;
&lt;meta name="linkedin:description" content="{{.Description}}"&gt;
&lt;meta name="linkedin:image" content="{{.ImageURL}}"&gt;
        </code></pre>
    </div>
</body>
</html>`

	tmpl, err := template.New("metatags").Parse(metaTemplate)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Fatalf("Error executing template: %v", err)
	}

	return buf.String()
}

// startLocalServer starts a local HTTP server to serve the HTML and image
func startLocalServer(htmlContent string, imagePath string, port string) string {
	serverURL := fmt.Sprintf("http://localhost:%s", port)
	
	// Create a file server handler
	fs := http.FileServer(http.Dir("."))
	
	// Define handlers
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(htmlContent))
			return
		}
		// Serve other requests through the file server
		fs.ServeHTTP(w, r)
	})
	
	// Start the server in a goroutine
	go func() {
		fmt.Printf("Starting local server at %s\n", serverURL)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()
	
	// Wait a moment for the server to start
	time.Sleep(500 * time.Millisecond)
	
	return serverURL
}

// ServerMain is the entry point for the OG generator functionality
func ServerMain() {
	// First check if API service mode is active
	isAPIService := false
	// Check command line args for service flag
	for _, arg := range os.Args {
		if arg == "-api-service" || arg == "-api-service=true" {
			isAPIService = true
			break
		}
	}
	// Also check environment variable
	if os.Getenv("OG_API_SERVICE") == "true" {
		isAPIService = true
	}

	// Create a new FlagSet and define all the flags
	fs := flag.NewFlagSet("og-generator", flag.ContinueOnError)
	fs.SetOutput(ioutil.Discard) // Suppress error output
	
	// Define all flags
	webpageURL := fs.String("url", "", "Webpage URL to capture")
	outputPath := fs.String("output", "outputs/og_image.png", "Output file path for the screenshot")
	outputHTML := fs.String("html", "outputs/og_meta.html", "Output file for HTML with meta tags")
	quality := fs.Int("quality", 90, "Screenshot quality (0-100)")
	waitTime := fs.Int("wait", 8000, "Wait time in milliseconds before taking screenshot")
	selector := fs.String("selector", "body", "CSS selector to wait for before capturing")
	debug := fs.Bool("debug", false, "Enable debug mode with additional logging")
	verbose := fs.Bool("verbose", false, "Enable verbose logging")
	title := fs.String("title", "", "Title for Open Graph meta tags")
	description := fs.String("description", "", "Description for Open Graph meta tags")
	ogType := fs.String("type", defaultType, "Type for Open Graph meta tags")
	siteName := fs.String("site", "", "Site name for Open Graph meta tags")
	targetURL := fs.String("target-url", "", "Target URL for the content (where it will be hosted)")
	imgWidth := fs.Int("width", defaultImageWidth, "Width of the Open Graph image")
	imgHeight := fs.Int("height", defaultImageHeight, "Height of the Open Graph image")
	twitterCard := fs.String("twitter-card", defaultTwitterCard, "Twitter card type")
	preview := fs.Bool("preview", false, "Start a local server to preview the Open Graph implementation")
	port := fs.String("port", "8080", "Port for the preview server")
	isApiService := fs.Bool("api-service", isAPIService, "Set to true when running as part of the API service")
	
	// Parse the command line arguments
	_ = fs.Parse(os.Args[1:])
	
	if *verbose {
		log.Printf("Command-line arguments: %v", os.Args)
		log.Printf("Output paths: image=%s, html=%s", *outputPath, *outputHTML)
		log.Printf("Key parameters: url=%s, title=%s, api-service=%v", *webpageURL, *title, *isApiService)
	}

	// Check for required inputs - but don't exit, just log the error
	if *webpageURL == "" && *title == "" {
		log.Printf("Warning: Neither a webpage URL (-url) nor a title (-title) was provided for Open Graph content.")
		// Use defaults instead of exiting
		*title = "Generated Open Graph Content"
	}

	// Create absolute path for output files
	absOutputPath, err := filepath.Abs(*outputPath)
	if err != nil {
		log.Printf("Error creating absolute path: %v", err)
		return
	}
	
	absHTMLPath, err := filepath.Abs(*outputHTML)
	if err != nil {
		log.Printf("Error creating absolute path: %v", err)
		return
	}

	// Ensure parent directory exists for each output file
	outputDir := filepath.Dir(absOutputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Printf("Error creating output directory for image: %v", err)
		return
	}

	htmlDir := filepath.Dir(absHTMLPath)
	if htmlDir != outputDir {
		if err := os.MkdirAll(htmlDir, 0755); err != nil {
			log.Printf("Error creating output directory for HTML: %v", err)
			return
		}
	}

	// Log the paths we're using
	log.Printf("Using output paths: Image=%s, HTML=%s", absOutputPath, absHTMLPath)

	// Generate image if a URL is provided
	if *webpageURL != "" {
		// Fix common URL issues
		fixedURL := *webpageURL
		// Handle double protocol issue (like https://www://example.com)
		if strings.Contains(fixedURL, "://") && strings.Count(fixedURL, "://") > 1 {
			parts := strings.SplitN(fixedURL, "://", 2)
			if len(parts) >= 2 {
				secondPart := parts[1]
				if strings.HasPrefix(secondPart, "www://") {
					secondPart = strings.Replace(secondPart, "www://", "www.", 1)
				}
				fixedURL = parts[0] + "://" + secondPart
			}
		}

		// Add http:// prefix if no protocol is present
		if !strings.Contains(fixedURL, "://") {
			fixedURL = "http://" + fixedURL
			if *verbose {
				log.Printf("Added http:// prefix to URL: %s", fixedURL)
			}
		}

		// Verify URL is valid
		_, err := url.Parse(fixedURL)
		if err != nil {
			log.Printf("Invalid URL: %s\nError: %v\n", fixedURL, err)
			// Return early instead of exiting
			return
		}

		if *verbose {
			fmt.Printf("Using URL: %s\n", fixedURL)
		}

		// Create a Chrome context with options for better rendering
		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("disable-web-security", true),
			chromedp.Flag("disable-background-networking", false),
			chromedp.Flag("enable-features", "NetworkService,NetworkServiceInProcess"),
			chromedp.Flag("disable-background-timer-throttling", true),
			chromedp.Flag("disable-backgrounding-occluded-windows", true),
			chromedp.Flag("disable-breakpad", true),
			chromedp.Flag("disable-client-side-phishing-detection", true),
			chromedp.Flag("disable-default-apps", true),
			chromedp.Flag("disable-dev-shm-usage", true),
			chromedp.Flag("disable-extensions", true),
			chromedp.Flag("disable-features", "site-per-process,TranslateUI,BlinkGenPropertyTrees"),
			chromedp.Flag("disable-hang-monitor", true),
			chromedp.Flag("disable-ipc-flooding-protection", true),
			chromedp.Flag("disable-popup-blocking", true),
			chromedp.Flag("disable-prompt-on-repost", true),
			chromedp.Flag("disable-renderer-backgrounding", true),
			chromedp.Flag("disable-sync", true),
			chromedp.Flag("force-color-profile", "srgb"),
			chromedp.Flag("metrics-recording-only", true),
			chromedp.Flag("safebrowsing-disable-auto-update", true),
			chromedp.Flag("enable-automation", true),
			chromedp.Flag("password-store", "basic"),
			chromedp.Flag("use-mock-keychain", true),
			// Additional rendering optimization flags
			chromedp.Flag("disable-accelerated-2d-canvas", false),
			chromedp.Flag("enable-gpu-rasterization", true),
			chromedp.Flag("disable-gpu-vsync", true),
			// Enhanced SSL error handling
			chromedp.Flag("ignore-certificate-errors", true),
			chromedp.Flag("allow-insecure-localhost", true),
			chromedp.Flag("allow-running-insecure-content", true),
			chromedp.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36"),
		)

		allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
		defer cancel()

		// Set up logging options
		logOpts := []chromedp.ContextOption{chromedp.WithLogf(log.Printf)}
		if *verbose {
			// Add more verbose logging if requested
			logOpts = append(logOpts, chromedp.WithDebugf(log.Printf))
		}

		// Create a Chrome context
		ctx, cancel := chromedp.NewContext(allocCtx, logOpts...)
		defer cancel()

		// Set a generous timeout
		ctx, cancel = context.WithTimeout(ctx, 45*time.Second)
		defer cancel()

		// Capture screenshot
		var buf []byte
		var htmlContent string
		
		fmt.Printf("Navigating to %s and waiting for content to load...\n", fixedURL)
		
		if err := chromedp.Run(ctx,
			// Navigate to the URL
			chromedp.Navigate(fixedURL),
			
			// Wait for the specified selector to be visible
			chromedp.WaitVisible(*selector, chromedp.ByQuery),
			
			// Wait for document to be ready
			chromedp.ActionFunc(func(ctx context.Context) error {
				err := chromedp.Evaluate(`
					new Promise((resolve) => {
						if (document.readyState === 'complete') {
							resolve();
						} else {
							window.addEventListener('load', resolve);
						}
					})
				`, nil).Do(ctx)
				return err
			}),
			
			// Ensure text elements are visible by forcing display properties
			chromedp.ActionFunc(func(ctx context.Context) error {
				script := `
					new Promise((resolve) => {
						// Force all elements to be visible
						const textElements = document.querySelectorAll('p, h1, h2, h3, h4, h5, h6, span, div, a, button, input, textarea, label');
						
						textElements.forEach(el => {
							// Check if element or its ancestors might have text content
							if (el.textContent && el.textContent.trim() !== '') {
								// Check computed style
								const style = window.getComputedStyle(el);
								if (style.display === 'none') {
									console.log('Forcing display for element:', el);
									el.style.setProperty('display', 'block', 'important');
								}
								if (style.visibility === 'hidden') {
									console.log('Forcing visibility for element:', el);
									el.style.setProperty('visibility', 'visible', 'important');
								}
								if (parseFloat(style.opacity) === 0) {
									console.log('Forcing opacity for element:', el);
									el.style.setProperty('opacity', '1', 'important');
								}
							}
						});
						
						// Force all potential text-containing elements to render
						document.querySelectorAll('[style*="display:none"], [style*="display: none"]').forEach(el => {
							if (el.textContent && el.textContent.trim() !== '') {
								el.style.setProperty('display', 'block', 'important');
							}
						});
						
						// Allow a bit of time for changes to take effect
						setTimeout(resolve, 500);
					})
				`
				return chromedp.Evaluate(script, nil).Do(ctx)
			}),
			
			// Simulate user scrolling to trigger lazy-loaded content
			chromedp.ActionFunc(func(ctx context.Context) error {
				script := `
					new Promise((resolve) => {
						// Simulate scrolling to trigger any lazy-loaded content
						const scrollHeight = Math.max(
							document.body.scrollHeight, document.documentElement.scrollHeight,
							document.body.offsetHeight, document.documentElement.offsetHeight,
							document.body.clientHeight, document.documentElement.clientHeight
						);
						
						// Scroll in increments to trigger events
						const increment = Math.max(window.innerHeight / 2, 200);
						let currentScroll = 0;
						
						const scrollInterval = setInterval(() => {
							window.scrollTo(0, currentScroll);
							currentScroll += increment;
							
							if (currentScroll >= scrollHeight) {
								clearInterval(scrollInterval);
								// Scroll back to top
								window.scrollTo(0, 0);
								resolve();
							}
						}, 100);
					})
				`
				return chromedp.Evaluate(script, nil).Do(ctx)
			}),
			
			// Simulate hovering on elements to trigger any hover effects
			chromedp.ActionFunc(func(ctx context.Context) error {
				script := `
					new Promise((resolve) => {
						// Find all interactive elements
						const elements = document.querySelectorAll('a, button, [role="button"], [tabindex]');
						elements.forEach(el => {
							// Dispatch mouseenter and mouseover events
							el.dispatchEvent(new MouseEvent('mouseenter', {
								view: window,
								bubbles: true,
								cancelable: true
							}));
						});
						resolve();
					})
				`
				return chromedp.Evaluate(script, nil).Do(ctx)
			}),
			
			// Additional wait time to ensure all content is fully loaded
			chromedp.Sleep(time.Duration(*waitTime) * time.Millisecond),
			
			// Get the HTML content for debugging
			chromedp.OuterHTML("html", &htmlContent, chromedp.ByQuery),
			
			// Take the screenshot
			chromedp.FullScreenshot(&buf, int(*quality)),
		); err != nil {
			if strings.Contains(err.Error(), "ERR_SSL_PROTOCOL_ERROR") || strings.Contains(err.Error(), "ERR_CERT") {
				log.Fatalf("SSL Certificate Error accessing %s: %v\nTry checking if the domain name is correct or if the site has valid SSL.", fixedURL, err)
			} else {
				log.Fatalf("Error executing chromedp: %v", err)
			}
		}

		// Save screenshot to file
		if err := ioutil.WriteFile(absOutputPath, buf, 0644); err != nil {
			log.Fatal(err)
		}

		// In debug mode, save the HTML content for inspection
		if *debug {
			debugFile := absOutputPath + ".html"
			if err := ioutil.WriteFile(debugFile, []byte(htmlContent), 0644); err != nil {
				log.Printf("Warning: Failed to save debug HTML: %v", err)
			} else {
				fmt.Printf("Debug HTML saved to %s\n", debugFile)
			}
		}

		fmt.Printf("Screenshot saved to %s\n", absOutputPath)
		
		// Try to extract title and description from the page if not provided
		if *title == "" {
			var extractedTitle string
			if err := chromedp.Run(ctx, chromedp.EvaluateAsDevTools(`document.querySelector('title').innerText`, &extractedTitle)); err == nil && extractedTitle != "" {
				*title = extractedTitle
				fmt.Printf("Extracted title from page: %s\n", *title)
			}
		}
		
		if *description == "" {
			var extractedDesc string
			if err := chromedp.Run(ctx, chromedp.EvaluateAsDevTools(`document.querySelector('meta[name="description"]')?.content || document.querySelector('meta[property="og:description"]')?.content || ""`, &extractedDesc)); err == nil && extractedDesc != "" {
				*description = extractedDesc
				fmt.Printf("Extracted description from page: %s\n", *description)
			}
		}
		
		if *siteName == "" {
			var extractedSiteName string
			if err := chromedp.Run(ctx, chromedp.EvaluateAsDevTools(`document.querySelector('meta[property="og:site_name"]')?.content || document.domain || ""`, &extractedSiteName)); err == nil && extractedSiteName != "" {
				*siteName = extractedSiteName
				fmt.Printf("Extracted site name from page: %s\n", *siteName)
			}
		}
	}

	// Use default values if fields are still empty
	if *title == "" {
		*title = "Open Graph Generated Content"
	}
	
	if *description == "" {
		*description = "Content shared with Open Graph meta tags"
	}
	
	// Determine the target URL
	pageURL := *targetURL
	if pageURL == "" {
		if *webpageURL != "" {
			pageURL = *webpageURL
		} else {
			// Use a placeholder if no URL is provided
			pageURL = "https://example.com/"
		}
	}
	
	// Determine image URL - check that the image file exists first
	imageURL := ""
	_, err = os.Stat(absOutputPath)
	imageExists := !os.IsNotExist(err)
	
	if imageExists {
		if *verbose {
			log.Printf("Image file exists at: %s", absOutputPath)
		}
	} else {
		if *verbose {
			log.Printf("Image file does not exist at: %s", absOutputPath)
		}
		fmt.Println("Warning: No image was generated. Using a placeholder in the meta tags.")
	}
	
	if *preview {
		// For preview, use relative path to the image
		imageURL = "/" + filepath.Base(absOutputPath)
	} else if *isApiService {
		// When running as part of the API service, just use the filename
		// The service will handle constructing the full URL
		imageURL = "/" + filepath.Base(absOutputPath)
	} else {
		// For production use the full URL
		imageURL = pageURL
		if !strings.HasSuffix(imageURL, "/") {
			imageURL += "/"
		}
		imageURL += filepath.Base(absOutputPath)
	}
	
	// If we don't have an image, use a placeholder
	if !imageExists && *webpageURL == "" {
		fmt.Println("Warning: No image was generated. Using a placeholder in the meta tags.")
		// Set a placeholder URL for the image
		imageURL = "https://via.placeholder.com/1200x630?text=" + url.QueryEscape(*title)
	}
	
	// Create the Open Graph data
	ogData := OpenGraphData{
		Title:       *title,
		Description: *description,
		ImageURL:    imageURL,
		PageURL:     pageURL,
		Type:        *ogType,
		SiteName:    *siteName,
		ImageWidth:  *imgWidth,
		ImageHeight: *imgHeight,
		TwitterCard: *twitterCard,
		LocalImage:  *preview,
	}
	
	htmlOutput := generateMetaTags(ogData)
	
	// Save HTML to file, using the absolute path to ensure it's saved to the correct location
	if err := ioutil.WriteFile(absHTMLPath, []byte(htmlOutput), 0644); err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("HTML with Open Graph meta tags saved to %s\n", absHTMLPath)
	
	if *verbose {
		log.Printf("Files generated - Image: %s, HTML: %s", absOutputPath, absHTMLPath)
	}
	
	// If preview mode is enabled, start a local server
	if *preview {
		serverURL := startLocalServer(htmlOutput, absOutputPath, *port)
		fmt.Printf("Preview available at: %s\n", serverURL)
		fmt.Printf("Press Ctrl+C to stop the server.\n")
		
		// Keep the program running until it's terminated
		select {}
	}
	
	if *verbose && *isApiService {
		log.Printf("ServerMain completed successfully. Returning to API service.")
	}
}
