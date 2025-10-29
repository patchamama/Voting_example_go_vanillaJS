#!/bin/bash

# Voting System Database Setup Script
# Compatible with macOS and Linux
# Supports MySQL, PostgreSQL, and MongoDB

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default configuration
DB_USER="${DB_USER:-voting_user}"
DB_PASSWORD="${DB_PASSWORD:-voting_password}"
DB_HOST="${DB_HOST:-localhost}"
DB_NAME="${DB_NAME:-voting_system}"

# Print colored message
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Detect OS
detect_os() {
    if [[ "$OSTYPE" == "darwin"* ]]; then
        OS="macos"
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        OS="linux"
    else
        print_error "Unsupported operating system: $OSTYPE"
        exit 1
    fi
    print_info "Detected OS: $OS"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Install Homebrew on macOS
install_homebrew() {
    if ! command_exists brew; then
        print_info "Installing Homebrew..."
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
        print_success "Homebrew installed"
    else
        print_info "Homebrew already installed"
    fi
}

# Setup MySQL
setup_mysql() {
    print_info "Setting up MySQL..."
    
    # Install MySQL
    if ! command_exists mysql; then
        if [[ "$OS" == "macos" ]]; then
            install_homebrew
            print_info "Installing MySQL via Homebrew..."
            brew install mysql
            brew services start mysql
        else
            print_info "Installing MySQL..."
            sudo apt-get update
            sudo apt-get install -y mysql-server
            sudo systemctl start mysql
            sudo systemctl enable mysql
        fi
        print_success "MySQL installed"
    else
        print_info "MySQL already installed"
    fi
    
    # Wait for MySQL to be ready
    print_info "Waiting for MySQL to be ready..."
    sleep 3
    
    # Create database and user
    print_info "Creating database and user..."
    
    if [[ "$OS" == "macos" ]]; then
        # macOS - no root password by default
        mysql -u root <<EOF
CREATE DATABASE IF NOT EXISTS ${DB_NAME} CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER IF NOT EXISTS '${DB_USER}'@'localhost' IDENTIFIED BY '${DB_PASSWORD}';
GRANT ALL PRIVILEGES ON ${DB_NAME}.* TO '${DB_USER}'@'localhost';
FLUSH PRIVILEGES;
EOF
    else
        # Linux - try with sudo first
        sudo mysql -u root <<EOF
CREATE DATABASE IF NOT EXISTS ${DB_NAME} CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER IF NOT EXISTS '${DB_USER}'@'localhost' IDENTIFIED BY '${DB_PASSWORD}';
GRANT ALL PRIVILEGES ON ${DB_NAME}.* TO '${DB_USER}'@'localhost';
FLUSH PRIVILEGES;
EOF
    fi
    
    print_success "MySQL database '${DB_NAME}' created"
    print_success "MySQL user '${DB_USER}' created with password '${DB_PASSWORD}'"
    
    # Test connection
    if mysql -u "${DB_USER}" -p"${DB_PASSWORD}" -e "USE ${DB_NAME};" 2>/dev/null; then
        print_success "MySQL connection test successful"
    else
        print_warning "MySQL connection test failed - you may need to set the root password"
    fi
}

# Setup PostgreSQL
setup_postgresql() {
    print_info "Setting up PostgreSQL..."
    
    # Install PostgreSQL
    if ! command_exists psql; then
        if [[ "$OS" == "macos" ]]; then
            install_homebrew
            print_info "Installing PostgreSQL via Homebrew..."
            brew install postgresql@14
            brew services start postgresql@14
            # Add to PATH
            export PATH="/usr/local/opt/postgresql@14/bin:$PATH"
        else
            print_info "Installing PostgreSQL..."
            sudo apt-get update
            sudo apt-get install -y postgresql postgresql-contrib
            sudo systemctl start postgresql
            sudo systemctl enable postgresql
        fi
        print_success "PostgreSQL installed"
    else
        print_info "PostgreSQL already installed"
    fi
    
    # Wait for PostgreSQL to be ready
    print_info "Waiting for PostgreSQL to be ready..."
    sleep 3
    
    # Create database and user
    print_info "Creating database and user..."
    
    if [[ "$OS" == "macos" ]]; then
        # macOS - current user is superuser
        createdb ${DB_NAME} 2>/dev/null || true
        psql -d ${DB_NAME} <<EOF
CREATE USER ${DB_USER} WITH PASSWORD '${DB_PASSWORD}';
GRANT ALL PRIVILEGES ON DATABASE ${DB_NAME} TO ${DB_USER};
GRANT ALL ON SCHEMA public TO ${DB_USER};
EOF
    else
        # Linux - use postgres user
        sudo -u postgres psql <<EOF
CREATE DATABASE ${DB_NAME};
CREATE USER ${DB_USER} WITH PASSWORD '${DB_PASSWORD}';
GRANT ALL PRIVILEGES ON DATABASE ${DB_NAME} TO ${DB_USER};
\c ${DB_NAME}
GRANT ALL ON SCHEMA public TO ${DB_USER};
EOF
    fi
    
    print_success "PostgreSQL database '${DB_NAME}' created"
    print_success "PostgreSQL user '${DB_USER}' created with password '${DB_PASSWORD}'"
    
    # Test connection
    if PGPASSWORD="${DB_PASSWORD}" psql -U "${DB_USER}" -d "${DB_NAME}" -c "\q" 2>/dev/null; then
        print_success "PostgreSQL connection test successful"
    else
        print_warning "PostgreSQL connection test failed"
    fi
}

# Setup MongoDB
setup_mongodb() {
    print_info "Setting up MongoDB..."
    
    # Install MongoDB
    if ! command_exists mongod; then
        if [[ "$OS" == "macos" ]]; then
            install_homebrew
            print_info "Installing MongoDB via Homebrew..."
            brew tap mongodb/brew
            brew install mongodb-community
            brew services start mongodb-community
        else
            print_info "Installing MongoDB..."
            # Import MongoDB public GPG key
            wget -qO - https://www.mongodb.org/static/pgp/server-6.0.asc | sudo apt-key add -
            # Create list file
            echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu $(lsb_release -cs)/mongodb-org/6.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-6.0.list
            # Update and install
            sudo apt-get update
            sudo apt-get install -y mongodb-org
            sudo systemctl start mongod
            sudo systemctl enable mongod
        fi
        print_success "MongoDB installed"
    else
        print_info "MongoDB already installed"
    fi
    
    # Wait for MongoDB to be ready
    print_info "Waiting for MongoDB to be ready..."
    sleep 5
    
    print_success "MongoDB is running (no authentication required by default)"
    print_info "Database '${DB_NAME}' will be created automatically on first use"
    
    # Test connection
    if mongosh --eval "db.version()" >/dev/null 2>&1 || mongo --eval "db.version()" >/dev/null 2>&1; then
        print_success "MongoDB connection test successful"
    else
        print_warning "MongoDB connection test failed"
    fi
}

# Setup Go dependencies
setup_go_dependencies() {
    print_info "Setting up Go dependencies..."
    
    if ! command_exists go; then
        print_error "Go is not installed. Please install Go 1.16+ first."
        print_info "Visit: https://golang.org/doc/install"
        exit 1
    fi
    
    print_info "Go version: $(go version)"
    
    # Initialize Go module if not exists
    if [ ! -f "go.mod" ]; then
        print_info "Initializing Go module..."
        go mod init voting-system
    fi
    
    # Install dependencies
    print_info "Installing Go dependencies..."
    go get github.com/go-sql-driver/mysql
    go get github.com/lib/pq
    go get go.mongodb.org/mongo-driver/mongo
    go get golang.org/x/crypto/bcrypt
    
    print_success "Go dependencies installed"
}

# Create .env file
create_env_file() {
    local db_type=$1
    local db_port=$2
    
    print_info "Creating .env file..."
    
    cat > .env <<EOF
# Database Configuration
DB_TYPE=${db_type}
DB_USER=${DB_USER}
DB_PASSWORD=${DB_PASSWORD}
DB_HOST=${DB_HOST}
DB_PORT=${db_port}
DB_NAME=${DB_NAME}
EOF
    
    print_success ".env file created"
}

# Create run script
create_run_script() {
    local db_type=$1
    
    print_info "Creating run script..."
    
    cat > run.sh <<'EOF'
#!/bin/bash

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Run the application
go run main.go
EOF
    
    chmod +x run.sh
    print_success "run.sh script created"
}

# Main menu
show_menu() {
    echo ""
    echo "======================================"
    echo "  Voting System Database Setup"
    echo "======================================"
    echo ""
    echo "Select database to setup:"
    echo "  1) MySQL (default)"
    echo "  2) PostgreSQL"
    echo "  3) MongoDB"
    echo "  4) All databases"
    echo "  5) Exit"
    echo ""
}

# Main setup function
main() {
    detect_os
    
    show_menu
    read -p "Enter your choice [1-5]: " choice
    
    case $choice in
        1)
            setup_mysql
            setup_go_dependencies
            create_env_file "mysql" "3306"
            create_run_script "mysql"
            print_success "MySQL setup complete!"
            ;;
        2)
            setup_postgresql
            setup_go_dependencies
            create_env_file "postgresql" "5432"
            create_run_script "postgresql"
            print_success "PostgreSQL setup complete!"
            ;;
        3)
            setup_mongodb
            setup_go_dependencies
            create_env_file "mongodb" "27017"
            create_run_script "mongodb"
            print_success "MongoDB setup complete!"
            ;;
        4)
            setup_mysql
            setup_postgresql
            setup_mongodb
            setup_go_dependencies
            create_env_file "mysql" "3306"
            create_run_script "mysql"
            print_success "All databases setup complete!"
            print_info "Default database is set to MySQL"
            ;;
        5)
            print_info "Exiting..."
            exit 0
            ;;
        *)
            print_error "Invalid choice"
            exit 1
            ;;
    esac
    
    echo ""
    print_info "================================================"
    print_success "Setup completed successfully!"
    print_info "================================================"
    echo ""
    print_info "Configuration:"
    echo "  Database Type: $(grep DB_TYPE .env | cut -d'=' -f2)"
    echo "  Database Name: ${DB_NAME}"
    echo "  Username: ${DB_USER}"
    echo "  Password: ${DB_PASSWORD}"
    echo "  Host: ${DB_HOST}"
    echo ""
    print_info "To run the application:"
    echo "  ./run.sh"
    echo ""
    print_info "Or manually:"
    echo "  source .env"
    echo "  go run main.go"
    echo ""
    print_info "Swagger UI will be available at:"
    echo "  http://127.0.0.1:8000/swagger/"
    echo ""
}

# Run main function
main