#!/bin/bash

# FlowPay VPS Installation Script
# Ubuntu 22.04 LTS
# Usage: curl -sSL https://raw.githubusercontent.com/your-repo/flowpay/main/scripts/vps-install.sh | bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_step() {
    echo -e "\n${BLUE}==>${NC} $1\n"
}

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    log_error "Please run as root (use sudo)"
    exit 1
fi

log_step "FlowPay VPS Installation"
log_info "This script will install FlowPay on Ubuntu 22.04"
log_warn "Press Ctrl+C to cancel, or Enter to continue..."
read

# Update system
log_step "Updating system packages..."
apt update && apt upgrade -y

# Install dependencies
log_step "Installing dependencies..."
apt install -y curl git wget unzip ca-certificates gnupg lsb-release

# Install Docker
log_step "Installing Docker..."
if ! command -v docker &> /dev/null; then
    curl -fsSL https://get.docker.com -o get-docker.sh
    sh get-docker.sh
    rm get-docker.sh
    systemctl enable docker
    systemctl start docker
    log_info "Docker installed successfully"
else
    log_info "Docker already installed"
fi

# Install Docker Compose
log_step "Installing Docker Compose..."
if ! command -v docker-compose &> /dev/null; then
    apt install -y docker-compose
    log_info "Docker Compose installed successfully"
else
    log_info "Docker Compose already installed"
fi

# Install Nginx
log_step "Installing Nginx..."
if ! command -v nginx &> /dev/null; then
    apt install -y nginx
    systemctl enable nginx
    systemctl start nginx
    log_info "Nginx installed successfully"
else
    log_info "Nginx already installed"
fi

# Install Certbot
log_step "Installing Certbot for SSL..."
if ! command -v certbot &> /dev/null; then
    apt install -y certbot python3-certbot-nginx
    log_info "Certbot installed successfully"
else
    log_info "Certbot already installed"
fi

# Create application directory
log_step "Creating application directory..."
mkdir -p /opt/flowpay
cd /opt/flowpay

# Get domain name
log_step "Configuration"
read -p "Enter your domain name (e.g., flowpay.example.com): " DOMAIN
read -p "Enter your email for SSL certificate: " EMAIL

# Clone repository (or you can upload manually)
log_warn "Please clone your repository or upload files to /opt/flowpay"
log_info "Example: git clone https://github.com/your-username/flowpay.git ."

# Create .env file
log_step "Creating .env file..."
cat > .env << EOF
# Database Configuration
DB_HOST=postgres
DB_PORT=5432
DB_USER=flowpay
DB_PASSWORD=$(openssl rand -base64 32 | tr -d "=+/" | cut -c1-32)
DB_NAME=flowpay_db
DB_SSLMODE=disable

# Server Configuration
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
ALLOWED_ORIGINS=https://${DOMAIN}

# JWT Configuration
JWT_SECRET=$(openssl rand -base64 32)
JWT_EXPIRY_HOURS=168

# Telegram Bot Configuration
TELEGRAM_BOT_TOKEN=your-telegram-bot-token-here

# Application Configuration
APP_ENV=production
APP_URL=https://${DOMAIN}
APP_NAME=FlowPay

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
EOF

log_info ".env file created with secure passwords"
log_warn "⚠️  IMPORTANT: Save these credentials!"
cat .env

# Configure Nginx
log_step "Configuring Nginx..."
cat > /etc/nginx/sites-available/flowpay << 'NGINX_EOF'
upstream flowpay_backend {
    server 127.0.0.1:8080;
}

server {
    listen 80;
    server_name DOMAIN_PLACEHOLDER;

    location /.well-known/acme-challenge/ {
        root /var/www/html;
    }

    location / {
        return 301 https://$server_name$request_uri;
    }
}

server {
    listen 443 ssl http2;
    server_name DOMAIN_PLACEHOLDER;

    ssl_certificate /etc/letsencrypt/live/DOMAIN_PLACEHOLDER/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/DOMAIN_PLACEHOLDER/privkey.pem;

    add_header Strict-Transport-Security "max-age=31536000" always;

    location / {
        proxy_pass http://flowpay_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location /health {
        proxy_pass http://flowpay_backend;
        access_log off;
    }
}
NGINX_EOF

sed -i "s/DOMAIN_PLACEHOLDER/${DOMAIN}/g" /etc/nginx/sites-available/flowpay
ln -sf /etc/nginx/sites-available/flowpay /etc/nginx/sites-enabled/
rm -f /etc/nginx/sites-enabled/default
nginx -t && systemctl reload nginx

# Get SSL certificate
log_step "Getting SSL certificate..."
certbot --nginx -d ${DOMAIN} --non-interactive --agree-tos -m ${EMAIL}

# Setup auto-renewal
systemctl enable certbot.timer

# Configure firewall
log_step "Configuring firewall..."
ufw allow 22/tcp
ufw allow 80/tcp
ufw allow 443/tcp
ufw --force enable

# Create backup script
log_step "Setting up backups..."
mkdir -p /opt/backups/flowpay

cat > /opt/scripts/backup-flowpay.sh << 'BACKUP_EOF'
#!/bin/bash
BACKUP_DIR="/opt/backups/flowpay"
DATE=$(date +%Y%m%d_%H%M%S)
docker exec flowpay-postgres pg_dump -U flowpay flowpay_db > $BACKUP_DIR/backup_$DATE.sql
gzip $BACKUP_DIR/backup_$DATE.sql
find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +7 -delete
BACKUP_EOF

chmod +x /opt/scripts/backup-flowpay.sh

# Add to crontab
(crontab -l 2>/dev/null; echo "0 2 * * * /opt/scripts/backup-flowpay.sh") | crontab -

log_step "Installation complete! 🎉"
log_info "Next steps:"
log_info "1. Upload your code to /opt/flowpay"
log_info "2. Edit .env file: nano /opt/flowpay/.env"
log_info "3. Add TELEGRAM_BOT_TOKEN to .env"
log_info "4. Start application: cd /opt/flowpay && docker-compose -f docker-compose.production.yml up -d"
log_info "5. Check logs: docker-compose logs -f"
log_info "6. Open https://${DOMAIN} in browser"

log_warn "\n⚠️  Remember to:"
log_warn "- Save .env credentials"
log_warn "- Configure Telegram bot"
log_warn "- Test the application"

exit 0
