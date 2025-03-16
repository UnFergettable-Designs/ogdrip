# Deployment Guide

This guide covers deploying OG Drip using Coolify and nixpacks.

## Prerequisites

- A Coolify instance (self-hosted or cloud)
- Git repository with your OG Drip fork
- Domain name (recommended)

## Environment Variables

### Frontend (.env)

```bash
PUBLIC_BACKEND_URL=https://api.your-domain.com
BACKEND_URL=https://api.your-domain.com
```

### Backend (.env)

```bash
PORT=8888
ADMIN_TOKEN=your-secure-admin-token
CHROME_PATH=/usr/bin/chromium
DATABASE_PATH=./data/ogdrip.db
OUTPUT_DIR=./outputs
```

## Deployment Steps

### 1. Prepare Your Repository

1. Fork or clone the OG Drip repository
2. Update environment variables as needed
3. Push your changes to your repository

### 2. Set Up Coolify

1. Log into your Coolify dashboard
2. Create a new project or select an existing one
3. Click "New Service"
4. Choose "Source: GitHub"
5. Select your OG Drip repository
6. Choose "Build Pack: Nixpacks"

### 3. Configure Build Settings

The repository includes a `nixpacks.toml` file that configures the build process. No additional
configuration is needed.

### 4. Configure Environment Variables

In Coolify's service settings, add the required environment variables:

```bash
# Frontend
PUBLIC_BACKEND_URL=https://your-domain.com
BACKEND_URL=https://your-domain.com

# Backend
PORT=8888
ADMIN_TOKEN=your-secure-admin-token
CHROME_PATH=/usr/bin/chromium
DATABASE_PATH=./data/ogdrip.db
OUTPUT_DIR=./outputs
```

### 5. Configure Domain and SSL

1. In your service settings, add your domain
2. Coolify will automatically handle SSL certificate generation
3. Configure your DNS records to point to your Coolify instance

### 6. Deploy

1. Click "Deploy" in Coolify
2. Monitor the build and deployment process
3. Once complete, your service will be available at your configured domain

## Monitoring and Maintenance

### Logs

Access logs through the Coolify dashboard:

1. Go to your service
2. Click on "Logs"
3. View real-time logs

### Updates

To update your deployment:

1. Push changes to your repository
2. Coolify will automatically detect changes
3. A new deployment will start automatically

### Backups

Configure backups in Coolify for:

- Database files
- Generated images
- Configuration

## Troubleshooting

### Common Issues

1. **Build Failures**

   - Check build logs in Coolify
   - Verify nixpacks.toml configuration
   - Ensure all dependencies are properly specified

2. **Runtime Errors**

   - Check application logs
   - Verify environment variables
   - Check Chrome/Chromium installation

3. **Performance Issues**
   - Monitor resource usage in Coolify
   - Consider scaling resources if needed
   - Check database and file system usage

## Security Considerations

1. **Environment Variables**

   - Use strong, unique ADMIN_TOKEN
   - Keep sensitive variables secure
   - Regularly rotate credentials

2. **Access Control**

   - Use HTTPS only
   - Configure proper CORS settings
   - Implement rate limiting if needed

3. **File System**
   - Monitor disk usage
   - Implement cleanup routines
   - Secure output directories

## Support

For additional help:

- Check the [GitHub Issues](https://github.com/yourusername/ogdrip/issues)
- Consult the [Coolify Documentation](https://docs.coolify.io)
- Join the community discussions
