#!/bin/bash
set -e

# This script fixes the nerdctl credential helper issue
# It creates a dummy credential helper script that nerdctl can call

echo "Creating dummy credential helper for nerdctl..."

# Find nerdctl config directory
NERDCTL_DIR="$HOME/.local/share/nerdctl"
mkdir -p "$NERDCTL_DIR"

# Create the config file without any credential store
cat > "$NERDCTL_DIR/config.json" << EOF
{
  "auths": {},
  "credsStore": ""
}
EOF

# Create a directory for the dummy credential helper
mkdir -p "$HOME/bin"

# Create a dummy credential helper script
cat > "$HOME/bin/docker-credential-undefined" << 'EOF'
#!/bin/bash
# This is a dummy credential helper that always returns an empty credential
echo "{}"
EOF

# Make it executable
chmod +x "$HOME/bin/docker-credential-undefined"

# Add it to the PATH if it's not already there
if [[ ":$PATH:" != *":$HOME/bin:"* ]]; then
    echo "export PATH=\$HOME/bin:\$PATH" >> "$HOME/.bashrc"
    echo "export PATH=\$HOME/bin:\$PATH" >> "$HOME/.zshrc"
    export PATH="$HOME/bin:$PATH"
fi

echo "Dummy credential helper created at $HOME/bin/docker-credential-undefined"
echo "Your PATH has been updated to include $HOME/bin"
echo "You may need to restart your terminal or run: export PATH=$HOME/bin:$PATH"
echo "Now you can run nerdctl commands without credential helper errors"
