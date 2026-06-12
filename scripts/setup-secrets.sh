#!/bin/bash

# Setup Docker Secrets
# Creates secret files for production deployment

set -e

BLUE='\033[36m'
GREEN='\033[32m'
YELLOW='\033[33m'
RESET='\033[0m'

SECRETS_DIR="./secrets"

echo -e "${BLUE}E-Commerce Platform - Docker Secrets Setup${RESET}"
echo "============================================="
echo ""

# Create secrets directory
mkdir -p "$SECRETS_DIR"
chmod 700 "$SECRETS_DIR"

echo -e "${YELLOW}This script will generate random secrets for production use.${RESET}"
echo -e "${YELLOW}Press Enter to continue, or Ctrl+C to cancel...${RESET}"
read

# Generate Postgres password
echo -e "${BLUE}Generating PostgreSQL password...${RESET}"
openssl rand -hex 32 > "$SECRETS_DIR/postgres_password.txt"
chmod 600 "$SECRETS_DIR/postgres_password.txt"

# Generate Redis password
echo -e "${BLUE}Generating Redis password...${RESET}"
openssl rand -hex 32 > "$SECRETS_DIR/redis_password.txt"
chmod 600 "$SECRETS_DIR/redis_password.txt"

# Generate JWT secret
echo -e "${BLUE}Generating JWT secret...${RESET}"
openssl rand -hex 64 > "$SECRETS_DIR/jwt_secret.txt"
chmod 600 "$SECRETS_DIR/jwt_secret.txt"

# Stripe keys (must be entered manually)
echo ""
echo -e "${YELLOW}Stripe API keys must be entered manually.${RESET}"
echo -e "${YELLOW}Get your keys from: https://dashboard.stripe.com/apikeys${RESET}"
echo ""
read -p "Enter Stripe Secret Key: " stripe_key
echo -n "$stripe_key" > "$SECRETS_DIR/stripe_secret_key.txt"
chmod 600 "$SECRETS_DIR/stripe_secret_key.txt"

read -p "Enter Stripe Webhook Secret: " webhook_secret
echo -n "$webhook_secret" > "$SECRETS_DIR/stripe_webhook_secret.txt"
chmod 600 "$SECRETS_DIR/stripe_webhook_secret.txt"

# SMTP password
echo ""
read -p "Enter SMTP Password: " smtp_password
echo -n "$smtp_password" > "$SECRETS_DIR/smtp_password.txt"
chmod 600 "$SECRETS_DIR/smtp_password.txt"

echo ""
echo -e "${GREEN}✓ Secrets created successfully!${RESET}"
echo ""
echo "Secret files:"
ls -lh "$SECRETS_DIR"
echo ""
echo -e "${YELLOW}IMPORTANT:${RESET}"
echo "- Keep these files secure and never commit to git"
echo "- Add secrets/ to .gitignore"
echo "- Backup secrets to secure storage"
echo "- Rotate secrets regularly"
