#!/bin/bash
set -e

echo "=== Preparing for Coolify deployment ==="
# Change to the repo root directory
cd "$(dirname "$0")"

# Create .gitignore.coolify-deploy to temporarily allow dist files
cat > frontend/.gitignore.coolify-deploy << EOL
# Keep the standard gitignore but allow dist for Coolify deployment
!dist/
EOL

echo "=== Building frontend locally with SSR ==="
# Install dependencies at the monorepo level
echo "Installing dependencies..."
pnpm install

# Build the frontend
echo "Building frontend..."
cd frontend
NODE_ENV=production pnpm build
cd ..

# Move the original .gitignore aside temporarily
mv frontend/.gitignore frontend/.gitignore.original
# Use the deployment .gitignore instead
mv frontend/.gitignore.coolify-deploy frontend/.gitignore

echo "=== Ready to commit and push for Coolify deployment ==="
echo ""
echo "To complete the deployment:"
echo ""
echo "1. Commit the changes with dist files:"
echo "   git add frontend/dist frontend/.gitignore"
echo "   git commit -m 'Prepare frontend for Coolify deployment'"
echo ""
echo "2. Push to your repository:"
echo "   git push"
echo ""
echo "3. After successful deployment, restore original .gitignore:"
echo "   mv frontend/.gitignore.original frontend/.gitignore"
echo "   git add frontend/.gitignore"
echo "   git commit -m 'Restore original .gitignore'"
echo ""
echo "4. Push the restored .gitignore:"
echo "   git push"
echo ""
echo "NOTE: This approach temporarily includes the 'dist' directory in your Git"
echo "repository to make deployment easier. The final step restores your"
echo "original configuration."
