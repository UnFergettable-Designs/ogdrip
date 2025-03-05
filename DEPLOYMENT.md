# Deployment Guide for OG-Drip.com on Coolify

This guide outlines the steps to deploy the Open Graph Generator service on Coolify.

## Prerequisites

- A Coolify instance
- Domain names registered and pointed to your server (og-drip.com and www.og-drip.com)
- Basic knowledge of Docker and Nginx
- Go v1.23 or later (if deploying without Docker)

## Deployment Steps

### 1. Prepare your Environment Files

Before deploying, update your environment files with the correct values:

- `frontend/.env.production`: Update the Sentry DSN if you're using Sentry
- `backend/.env.production`: Update the ADMIN_TOKEN with a secure value and Sentry DSN if applicable

### 2. Deploy to Coolify

1. Log in to your Coolify dashboard
2. Add a new service
3. Select "Docker Compose" as the deployment method
4. Connect to your Git repository
5. Choose the `docker-compose.production.yml` file for deployment
6. Configure your domains:
   - Primary domain: www.og-drip.com
   - Additional domain: og-drip.com

### 3. Configure SSL/TLS

1. Enable "Auto SSL" in Coolify
2. Provide your email address for Let's Encrypt notifications
3. Wait for certificate generation

### 4. Set Up Persistent Volumes

1. In the Coolify dashboard, go to your service settings
2. Add persistent volumes:
   - Path: `/app/outputs`
   - Path: `/app/data` (for the database)

### 5. Configure Custom Nginx (Optional)

If you need to use the custom Nginx configuration:

1. Go to service settings > Advanced
2. Select "Custom Nginx Configuration"
3. Upload or paste the contents of `nginx.conf`

### 6. Monitor the Deployment

1. Check the deployment logs for any errors
2. Verify your site is accessible at https://www.og-drip.com
3. Test the API with a basic request to https://www.og-drip.com/api/health

## Troubleshooting

### CORS Issues

If you encounter CORS issues:

1. Verify the `ENABLE_CORS` environment variable is set to `true`
2. Check that the `BASE_URL` is set correctly to `https://www.og-drip.com`

### Certificate Issues

If you have SSL/TLS certificate issues:

1. Ensure your DNS records are properly configured
2. Check that both domains are registered in Coolify
3. Verify that ports 80 and 443 are accessible

### Volume Permissions

If you have issues with file permissions:

1. SSH into your Coolify server
2. Check the permissions on the volume directories
3. Run: `chmod -R 755 /path/to/volumes/outputs`

## Maintenance

### Backups

1. Set up a regular backup schedule for your database and generated files
2. The critical paths to back up are:
   - `/app/data` - Database files
   - `/app/outputs` - Generated images and HTML files

### Updates

When updating your application:

1. Make your changes to the codebase
2. Push to your repository
3. Redeploy through the Coolify dashboard

## Monitoring

Monitor your application health:

1. Set up regular checks to `/api/health` endpoint
2. Consider setting up alerts if the health check fails
3. Monitor disk space on the volumes to ensure you don't run out of space
