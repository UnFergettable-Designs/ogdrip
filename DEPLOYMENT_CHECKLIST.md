# Deployment Checklist

This checklist ensures your OG Drip deployment is production-ready.

## Pre-Deployment Setup

### 1. Environment Files
Create these files before deploying:

**backend/.env.production**
```env
ADMIN_TOKEN=your_secure_admin_token_here
PORT=8888
HOST=0.0.0.0
DATABASE_PATH=./data/ogdrip.db
BROWSER_TIMEOUT=30
CORS_ORIGINS=https://yourdomain.com
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=3600
CHROME_PATH=/nix/store/*-chromium-*/bin/chromium
DISPLAY=:99
```

**frontend/.env.production**
```env
PUBLIC_BACKEND_URL=https://yourdomain.com
BACKEND_URL=https://yourdomain.com
NODE_ENV=production
```

### 2. Security Configuration
- [ ] Generate secure ADMIN_TOKEN (use: `openssl rand -hex 32`)
- [ ] Set proper CORS_ORIGINS to your domain
- [ ] Configure rate limiting
- [ ] Enable HTTPS/SSL

### 3. Domain Setup
- [ ] Configure DNS to point to your server
- [ ] Set up SSL certificates
- [ ] Test domain accessibility

## Coolify Deployment

### Quick Setup
- [ ] Create new service in Coolify
- [ ] Connect Git repository
- [ ] Select "Nixpacks" as build pack
- [ ] Set environment variables
- [ ] Add domain and enable Auto SSL
- [ ] Deploy

### Environment Variables in Coolify
Set these in the Coolify dashboard:
```
ADMIN_TOKEN=your_secure_admin_token_here
PUBLIC_BACKEND_URL=https://your-domain.com
BACKEND_URL=https://your-domain.com
CORS_ORIGINS=https://your-domain.com
```

## Docker Deployment

### Prerequisites
- [ ] Docker and Docker Compose installed
- [ ] Environment files created
- [ ] Domain configured

### Commands
```bash
# Build and start
docker-compose -f docker-compose.yml up -d

# Check logs
docker-compose logs -f

# Health check
curl https://yourdomain.com/api/health
```

## Post-Deployment Verification

### Health Checks
- [ ] Backend API responding: `/api/health`
- [ ] Frontend loading properly
- [ ] Image generation working
- [ ] Database connections working

### Security Tests
- [ ] HTTPS redirect working
- [ ] CORS headers correct
- [ ] Rate limiting active
- [ ] Admin endpoints protected

### Performance Tests
- [ ] Page load times acceptable
- [ ] Image generation under 30 seconds
- [ ] Multiple concurrent requests handled
- [ ] Memory usage stable

### Monitoring Setup
- [ ] Health monitoring configured
- [ ] Log aggregation working
- [ ] Backup schedule configured
- [ ] Alert notifications set up

## Troubleshooting

### Common Issues
1. **Build fails**: Check Go/Node versions in nixpacks.toml
2. **Chrome not starting**: Verify Chromium installation and DISPLAY variable
3. **CORS errors**: Check CORS_ORIGINS matches your domain
4. **Health checks fail**: Verify backend is running on port 8888

### Debug Commands
```bash
# Check service status
./healthcheck.sh

# View application logs
docker logs container_name

# Test API manually
curl -X POST https://yourdomain.com/api/generate \
  -H "Content-Type: application/json" \
  -d '{"title":"Test","description":"Test description"}'
```

## Production Maintenance

### Regular Tasks
- [ ] Monitor resource usage
- [ ] Check error logs weekly
- [ ] Update dependencies monthly
- [ ] Test backup/restore quarterly

### Security Updates
- [ ] Keep base images updated
- [ ] Monitor security advisories
- [ ] Rotate admin tokens annually
- [ ] Review access logs

---

✅ **Ready for Production**: All items checked and verified
⚠️ **Needs Attention**: Some items need configuration
❌ **Not Ready**: Critical items missing
