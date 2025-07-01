# Security Overview

Security is a core consideration in OG Drip's design and implementation. This document outlines the
security measures, best practices, and considerations implemented throughout the system.

## Security Principles

OG Drip follows these fundamental security principles:

1. **Defense in Depth**: Multiple layers of security controls
2. **Least Privilege**: Minimal access rights for all components
3. **Fail Securely**: Secure defaults and graceful failure handling
4. **Security by Design**: Security considerations integrated from the start
5. **Regular Updates**: Keep dependencies and systems current

## Threat Model

### Assets

- **User Data**: URLs, generation history, metadata
- **Generated Images**: Potentially sensitive screenshots
- **Admin Access**: Administrative functionality and data
- **System Resources**: CPU, memory, disk space, network

### Threats

- **Unauthorized Access**: Accessing admin functionality without authentication
- **Data Exposure**: Leaking sensitive information through images or logs
- **Resource Abuse**: Using the service for malicious purposes or DoS attacks
- **Code Injection**: Exploiting input validation vulnerabilities
- **Data Tampering**: Modifying generated images or metadata

### Attack Vectors

- **API Abuse**: Excessive requests or malformed payloads
- **URL Manipulation**: Malicious URLs for SSRF or data exfiltration
- **Admin Interface**: Brute force or credential attacks
- **Browser Exploitation**: Exploiting ChromeDP or browser vulnerabilities

## Security Controls

### Authentication & Authorization

#### Admin Authentication

- **Bearer Token**: Secure token-based authentication for admin endpoints
- **Environment Variables**: Tokens stored securely outside codebase
- **Token Validation**: Proper validation and error handling

```go
// Example token validation
func validateAdminToken(token string) bool {
    expectedToken := os.Getenv("ADMIN_TOKEN")
    return subtle.ConstantTimeCompare([]byte(token), []byte(expectedToken)) == 1
}
```

#### API Access Control

- **Public Endpoints**: Limited functionality without authentication
- **Admin Endpoints**: Full functionality with proper authentication
- **Rate Limiting**: Prevent abuse of public endpoints

### Input Validation & Sanitization

#### URL Validation

- **Scheme Validation**: Only HTTP/HTTPS URLs allowed
- **Domain Filtering**: Block internal/private network addresses
- **Length Limits**: Prevent excessively long URLs
- **Malicious Pattern Detection**: Block known malicious patterns

```go
func validateURL(rawURL string) error {
    if len(rawURL) > 2048 {
        return errors.New("URL too long")
    }

    u, err := url.Parse(rawURL)
    if err != nil {
        return err
    }

    if u.Scheme != "http" && u.Scheme != "https" {
        return errors.New("invalid URL scheme")
    }

    // Block private networks
    if isPrivateNetwork(u.Host) {
        return errors.New("private network access denied")
    }

    return nil
}
```

#### Request Validation

- **JSON Schema Validation**: Validate request structure
- **Parameter Bounds**: Enforce reasonable limits on dimensions
- **Content-Type Validation**: Ensure proper content types

### Data Protection

#### Database Security

- **Prepared Statements**: Prevent SQL injection attacks
- **Connection Security**: Secure database connections
- **Data Encryption**: Sensitive data encrypted at rest (when applicable)
- **Access Controls**: Minimal database permissions

```go
// Example prepared statement
stmt, err := db.Prepare("INSERT INTO generations (url, image_path, created_at) VALUES (?, ?, ?)")
if err != nil {
    return err
}
defer stmt.Close()
```

#### File System Security

- **Isolated Storage**: Generated files stored in dedicated directory
- **Path Validation**: Prevent directory traversal attacks
- **File Permissions**: Restrictive file permissions
- **Cleanup Procedures**: Regular cleanup of old files

### Network Security

#### HTTPS Enforcement

- **TLS Configuration**: Strong TLS configuration in production
- **Certificate Management**: Proper SSL certificate handling
- **HSTS Headers**: HTTP Strict Transport Security

#### CORS Configuration

- **Origin Validation**: Proper CORS origin configuration
- **Method Restrictions**: Only necessary HTTP methods allowed
- **Credential Handling**: Secure credential handling in CORS

```javascript
// Frontend CORS configuration
const corsOptions = {
  origin: process.env.PUBLIC_BACKEND_URL,
  credentials: true,
  optionsSuccessStatus: 200,
};
```

### Browser Security

#### ChromeDP Security

- **Sandboxing**: Run browser in sandboxed environment
- **Resource Limits**: Limit CPU and memory usage
- **Timeout Controls**: Prevent long-running operations
- **User Agent**: Use standard user agent strings

```go
// Example ChromeDP security configuration
opts := append(chromedp.DefaultExecAllocatorOptions[:],
    chromedp.Flag("no-sandbox", true),
    chromedp.Flag("disable-gpu", true),
    chromedp.Flag("disable-dev-shm-usage", true),
    chromedp.Flag("disable-extensions", true),
)
```

#### Content Security

- **Content Filtering**: Basic content validation
- **Resource Limits**: Limit page load time and resources
- **Error Handling**: Secure error handling for browser operations

### Rate Limiting & DoS Protection

#### API Rate Limiting

- **IP-based Limits**: Rate limits per IP address
- **Token-based Limits**: Higher limits for authenticated users
- **Sliding Window**: Sophisticated rate limiting algorithm

#### Resource Protection

- **Request Timeouts**: Prevent long-running requests
- **Concurrent Limits**: Limit concurrent browser instances
- **Memory Limits**: Monitor and limit memory usage

### Logging & Monitoring

#### Security Logging

- **Authentication Events**: Log all authentication attempts
- **Access Logs**: Log all API access with details
- **Error Logging**: Log security-relevant errors
- **Audit Trail**: Maintain audit trail for admin actions

#### Monitoring

- **Anomaly Detection**: Monitor for unusual patterns
- **Alert System**: Alerts for security events
- **Log Analysis**: Regular analysis of security logs

### Dependency Security

#### Vulnerability Management

- **Dependency Scanning**: Regular security scans of dependencies
- **Update Process**: Regular updates of dependencies
- **Vulnerability Tracking**: Track and remediate known vulnerabilities

#### Supply Chain Security

- **Dependency Verification**: Verify integrity of dependencies
- **Minimal Dependencies**: Use minimal necessary dependencies
- **Trusted Sources**: Use only trusted package repositories

## Security Configuration

### Environment Variables

```env
# Security-related environment variables
ADMIN_TOKEN=your_secure_random_token_here
CORS_ORIGINS=https://yourdomain.com
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=3600
BROWSER_TIMEOUT=30
```

### Production Security Checklist

- [ ] Strong admin token configured
- [ ] HTTPS enabled with valid certificate
- [ ] CORS properly configured
- [ ] Rate limiting enabled
- [ ] Security headers configured
- [ ] Dependency vulnerabilities addressed
- [ ] Logging and monitoring configured
- [ ] Regular security updates scheduled
- [ ] Backup and recovery procedures tested

## Incident Response

### Security Incident Procedure

1. **Detection**: Identify potential security incident
2. **Assessment**: Evaluate severity and impact
3. **Containment**: Isolate affected systems
4. **Investigation**: Analyze logs and evidence
5. **Remediation**: Fix vulnerabilities and restore service
6. **Recovery**: Return to normal operations
7. **Lessons Learned**: Document and improve procedures

### Contact Information

- **Security Team**: security@yourdomain.com
- **Emergency Contact**: +1-XXX-XXX-XXXX
- **Incident Reporting**: Use GitHub Security Advisories for vulnerabilities

## Security Testing

### Regular Security Testing

- **Vulnerability Scanning**: Automated vulnerability scans
- **Penetration Testing**: Regular penetration testing
- **Code Review**: Security-focused code reviews
- **Dependency Audits**: Regular dependency security audits

### Security Testing Tools

- **Static Analysis**: ESLint security rules, Go security analyzers
- **Dependency Scanning**: npm audit, Go vulnerability database
- **Container Scanning**: Docker image vulnerability scanning
- **Dynamic Testing**: API security testing, browser automation testing

## Compliance & Standards

### Security Standards

- **OWASP Top 10**: Address common web application vulnerabilities
- **CWE/SANS Top 25**: Address most dangerous software errors
- **Security Headers**: Implement recommended security headers

### Privacy Considerations

- **Data Minimization**: Collect only necessary data
- **Data Retention**: Implement data retention policies
- **User Privacy**: Respect user privacy in generated images
- **Logging Privacy**: Avoid logging sensitive information

## Further Reading

- [Authentication & Authorization](auth.md) - Detailed authentication implementation
- [Best Practices](best-practices.md) - Security best practices for developers
- [Threat Model](threat-model.md) - Detailed threat analysis
- [Incident Response](incident-response.md) - Detailed incident response procedures

---

_For security issues, please report privately via GitHub Security Advisories or email
security@yourdomain.com_
