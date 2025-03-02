# Open Graph Image and Meta Tag Generator Design Document (Updated)

## 1. Overview

This project aims to create a tool that generates a PNG snapshot of a website along with the corresponding HTML meta tags necessary for the most common social media platforms. The tool emphasizes dynamic image generation, file storage and caching, usage tracking, robust security measures, comprehensive testing, and integrated error monitoring with Sentry to ensure overall project integrity during deployment.

## 2. Scope and Objectives

### Scope
- **Primary Functionality:**  
  - Generate a PNG image (screenshot) of a given website URL.
  - Generate default meta tags (Open Graph, Twitter Cards, etc.) for social media sharing.
  - Store generated images and securely provide download capabilities.
  - Package the generated image and HTML meta tag file into a zip file for download.
  - Track usage metrics for image generation and meta tag generation features.
  - Use Cloudflare Worker edge caching to serve the images quickly, even if the application is hosted on Coolify.
  - Automatically clean up generated files after a set duration (e.g., one hour) to conserve storage and ensure security.

### Extended Objectives
- **Image Storage & Downloading:**  
  - Generate, store, and securely serve screenshots.
  - Provide an endpoint to download the generated zip file containing the PNG image and HTML file.
  - Utilize Cloudflare Workers for edge caching, reducing origin load.
  - Automatically remove stored files (i.e., screenshots and zip packages) after a prescribed lifetime (such as one hour) using background cleanup jobs or lifecycle policies in storage services.
- **Usage Tracking:**  
  - Monitor and log API requests, image generation successes/failures, and download counts.
  - Record key metrics such as frequency of requests, response times, and feature usage (with appropriate anonymization).
- **Enhanced Security:**  
  - Emphasize input validation and sanitization for URLs.
  - Implement secure storage access for the generated images with signed URLs or tokens.
  - Maintain audit logs for critical actions (e.g., image generation, file downloads).
- **Comprehensive Testing:**  
  - Thoroughly test all components including image generation, meta tag creation, file packaging, API endpoints, caching, and security features.
  - Use automated unit, integration, end-to-end, performance, and security tests to ensure quality.
- **Integrated Error Monitoring with Sentry:**  
  - Integrate Sentry to capture and track runtime errors and performance issues.
  - Use Sentry’s dashboards to monitor errors, track release performance, and quickly identify bugs in production.
  - Incorporate Sentry error logging in API endpoints, background jobs (e.g., screenshot generation), and Cloudflare Worker interactions.

## 3. Input and Output Specifications

### Input
- **Primary Input:**  
  - Website URL to capture.
- **Optional Inputs (for future or advanced versions):**  
  - Specific social media platform selection.
  - Custom parameters for title, description or image overrides.
  
### Output
- **PNG Image:**  
  - A dynamically generated screenshot of the website.
  
- **HTML Meta Tags:**  
  - An HTML file containing meta tags that comply with default specifications for platforms such as Facebook and Twitter.
  
- **Downloadable Zip File:**  
  - A zip file containing both the PNG image and the HTML meta tag file.

- **Usage Data (Internal):**  
  - Logs and metrics for monitoring API usage and performance.
  - Sentry reports for error and performance monitoring.

## 4. Architecture and Components

### Frontend
- **User Interface:**
  - A web portal where users input a URL.
  - Preview of the generated PNG image with a download link.
  - Display of the corresponding meta tag snippet.
  - Potential dashboard elements for users/admins to view usage statistics (authenticated access).

### Backend
- **API Layer:**
  - RESTful endpoints to:
    - Submit URLs for screenshot and meta tag generation.
    - Retrieve secure download links for the generated zip file (PNG image + HTML meta tag file).
    - Access usage tracking metrics (with proper authentication).

- **Screenshot Generation Service:**
  - Uses a headless browser (e.g., Puppeteer) to capture a website screenshot.
  - Returns a PNG image for storage and distribution.
  - Integrate Sentry error logging to capture any runtime issues during screenshot generation.

- **Image Storage and Delivery:**
  - **Primary Storage Options:**
    - Local storage on the hosting server (Coolify) or using a cloud object storage solution.
    - Optionally, use Cloudflare R2 as an S3-compatible object storage if needed.
  - **Delivery and Caching:**
    - Serve images directly from the origin server through a secure endpoint.
    - Implement an edge caching layer using Cloudflare Workers.
    - The Cloudflare Worker will handle downloads, fetching images from the origin (Coolify) when needed, and caching them with an appropriate TTL (e.g., one hour).
  - **File Cleanup Strategy:**
    - Automatically delete generated screenshots and zip files after a prescribed lifetime (e.g., one hour) using background cleanup processes or lifecycle policies if using cloud storage.

- **Meta Tag Generator Service:**
  - Constructs meta tag HTML snippets based on website data and user inputs.
  - Logs errors with Sentry if meta tag generation fails.

- **File Packaging Service:**
  - Combines the generated PNG image and the HTML file with meta tags into a single zip file.
  - Provides a secure download endpoint for the zip file.
  - Integrate Sentry monitoring to log file packaging errors.
  - Ensure that packaged files have a limited lifetime before they are automatically cleaned up.

- **Usage Tracking System:**
  - Logs each API request.
  - Tracks image generation requests, download counts, and meta tag generation status.
  - Provides an analytics endpoint with aggregated data (admin-secured).

- **Error Monitoring with Sentry:**
  - Integrate Sentry in all backend components to capture errors, exceptions, and performance issues.
  - Monitor API endpoints, background jobs, and Cloudflare Worker interactions.
  - Use Sentry release tracking to associate errors with new code deployments.

## 5. Design Considerations

### Performance
- Use caching strategies to minimize redundant processing:
  - Cache generated images and meta tag outputs.
  - Use Cloudflare Workers to cache image requests at the edge for a defined duration (e.g., one hour).
- Integrate a content delivery network (CDN) for faster image delivery.

### Customizability
- Future support for custom meta tag templates and specific platform options.
- Ability to override default meta tag values with user-provided inputs.

### Security
- **Input Validation:**  
  - Sanitize URL inputs to avoid injection attacks.
  - Implement rate limiting and CAPTCHA where required.
- **Secure Storage & Access:**  
  - Use signed URLs or tokens for image download endpoints.
  - Ensure assets stored in object storage are not publicly accessible by default.
- **Audit Logging:**  
  - Log transactional events (e.g., image generation, downloads).
- **Edge Security:**  
  - The Cloudflare Worker will enforce cache control and secure endpoints to prevent unauthorized access.
- **Sentry Integration:**  
  - Use Sentry’s monitoring dashboard to track and respond to security-related exceptions or unexpected behaviors.

### Error Handling
- Provide robust error messages for invalid URL inputs or failed image generation.
- Implement fallback mechanisms to serve default images and meta tags if processing fails.
- Capture all exceptions with Sentry for real-time alerting and debugging.

## 6. API Design

### Endpoints
- **POST /api/v1/generate**
  - Accepts:
    - URL (required)
    - Optional platform selection and customization parameters.
  - Returns:
    - A secure URL for the generated screenshot (with a signed token if needed).
    - An HTML meta tag snippet.
    - A download URL pointing to a zip file containing both the PNG image and the HTML file.
    - A request tracking ID for analytics.
  - Integrate error reporting via Sentry for any generation issues encountered.

- **GET /api/v1/download/{fileId}**
  - Validates access using signed URLs or tokens.
  - Serves the zip file that contains:
    - The generated PNG image.
    - The HTML file with meta tags.
  - Uses Cloudflare Worker edge caching to serve the file efficiently.
  - Logs any download errors via Sentry.

- **GET /api/v1/usage**
  - (Admin/Authenticated only) Returns aggregated usage metrics.
  - Monitor any issues with Sentry integration on this secured endpoint.

### Edge Caching with Cloudflare Workers
- **Worker Role:**
  - Intercepts download requests for the zip file.
  - Checks the Cloudflare edge cache for a cached version of the requested file.
  - If not cached, the Worker fetches the file from the origin (Coolify), caches it with the defined TTL, and serves it.
  - Enforces secure access by checking for valid tokens or signed URLs and sets appropriate Cache-Control headers.
  - Integrates Sentry logging for any caching or retrieval errors.

## 7. Deployment & Environment

### Infrastructure
- **Hosting on Coolify:**
  - The main application and API are hosted on Coolify.
  - Local or cloud-based storage can be used for the generated images and zip files.
- **Cloudflare Integration:**
  - Use Cloudflare Workers for edge caching and secure delivery of files.
  - Optionally, leverage R2 for scalable object storage if required.
- **Containerization & Orchestration:**
  - Docker for isolated deployments and automated scaling.
- **CDN Integration:**
  - Cloudflare provides CDN features to accelerate file delivery.
- **Sentry Integration:**
  - Deploy Sentry DSN configuration as part of the environment.
  - Ensure Sentry is properly configured for each service and worker instance.

### CI/CD and Monitoring
- Implement automated testing for functionality, performance, and security.
- Deploy continuous integration pipelines to automate releases.
- Use telemetry and logging systems and Sentry to monitor API performance, error rates, and usage analytics.
- Integrate Sentry error reporting in the CI/CD pipeline to capture issues during deployment.
- Schedule or integrate cleanup jobs to automatically remove generated files (screenshots and zip packages) after their defined lifetime (e.g., one hour).

## 8. Testing and Quality Assurance

### Testing Strategy
- **Unit Tests:**  
  - Test individual components including:
    - URL input validation.
    - Screenshot generation service (e.g., mocking Puppeteer).
    - Meta tag generation—ensuring the correct HTML is rendered.
    - File packaging service to verify correct creation of zip files and subsequent cleanup.
- **Integration Tests:**  
  - Test the interaction between components:
    - End-to-end flow from API call to zip file download.
    - Edge caching via Cloudflare Workers.
    - Secure communication between API endpoints and storage services.
- **End-to-End Tests:**  
  - Simulate real-world use cases:
    - Full generation process from entering a URL to downloading the zip file.
    - Verify that generated meta tags render correctly on social media preview tools.
- **Performance and Load Tests:**  
  - Assess system performance under heavy load:
    - Evaluate caching efficiency.
    - Measure response times for API endpoints and file downloads.
- **Security Tests:**  
  - Perform penetration testing and vulnerability scans:
    - Test for injection flaws in the URL inputs.
    - Ensure that signed URLs or tokens are appropriately validated.
    - Verify rate limiting and other access control mechanisms.
- **Sentry Monitoring in Tests:**  
  - Include tests that verify errors are properly captured and reported to Sentry.
  - Simulate error conditions to confirm that Sentry alerts and logs are generated as expected.

### Automated Test Frameworks
- Use popular testing frameworks such as Jest, Mocha/Chai, or similar based on the development stack.
- Continuous Integration (CI) should run the full test suite on each commit or pull request.
- Integrate error monitoring tests and ensure Sentry DSN is correctly configured in test environments.

## 9. Future Enhancements

- **Enhanced User Customization:**  
  - Allow users to edit meta tag parameters and image templates.
- **Detailed Usage Analytics:**  
  - Extend the analytics dashboard with deeper insights and reporting.
- **Social Media Integration:**  
  - Explore direct integration with social media platforms for automated content posting.
- **Extended Caching Strategies:**  
  - Refine edge cache configurations and explore alternative caching layers as usage scales.
- **Expanded Sentry Monitoring:**  
  - Utilize Sentry performance monitoring features for deeper insights into request/response times and overall system health.
- **Dynamic Cleanup Policies:**  
  - Implement more dynamic or configurable file cleanup policies based on usage metrics and storage constraints.

## 10. Summary

This design document specifies an architecture that leverages both traditional server-side image generation and modern edge caching techniques using Cloudflare Workers. The solution offers fast, secure, and scalable image and meta tag generation, complete with packaging into a downloadable zip file. Comprehensive testing—from unit tests to end-to-end, performance, and security tests—as well as integrated error monitoring using Sentry, ensures that the project remains robust, reliable, and maintainable throughout its lifecycle. Additionally, a file cleanup mechanism automatically removes generated files after a predefined period (e.g., one hour) to optimize resource usage and enhance security.
