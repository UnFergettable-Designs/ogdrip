#!/bin/bash
set -e

# This script helps authenticate nerdctl with GitHub Container Registry
# It bypasses credential helpers by directly setting up a basic authentication

if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <github-username> <github-token>"
    echo ""
    echo "You need to provide your GitHub username and a personal access token"
    echo "The token needs 'read:packages' scope at minimum"
    echo ""
    echo "To create a token: https://github.com/settings/tokens/new"
    exit 1
fi

USERNAME="$1"
TOKEN="$2"

# Create nerdctl config directory
NERDCTL_DIR="$HOME/.local/share/nerdctl"
mkdir -p "$NERDCTL_DIR"

# Set up auth config
echo "Creating auth configuration for ghcr.io..."
ENCODED_AUTH=$(echo -n "$USERNAME:$TOKEN" | base64)

cat > "$NERDCTL_DIR/config.json" << EOF
{
  "auths": {
    "ghcr.io": {
      "auth": "$ENCODED_AUTH"
    }
  }
}
EOF

echo "Authentication configured for ghcr.io"
echo "You can now use docker-build/run-container.sh to run the container"
