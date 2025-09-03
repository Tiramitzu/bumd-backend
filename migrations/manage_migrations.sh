#!/bin/bash

# Migration Management Script
# This script helps manage database migrations

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
MIGRATIONS_PATH="."
DATABASE_URL="postgres://postgres:postgres@192.168.191.10:5432/bumd?sslmode=disable"

# Function to show usage
show_usage() {
    echo -e "${BLUE}Migration Management Script${NC}"
    echo ""
    echo "Usage: $0 [COMMAND] [OPTIONS]"
    echo ""
    echo "Commands:"
    echo "  up [N]           Run N migrations up (default: all)"
    echo "  down [N]          Run N migrations down (default: 1)"
    echo "  goto V            Migrate to specific version V"
    echo "  force V           Set migration version to V without running migrations"
    echo "  version           Show current migration version"
    echo "  status            Show migration status"
    echo "  create NAME       Create new migration files"
    echo "  list              List all migration files"
    echo ""
    echo "Options:"
    echo "  -d, --database    Database connection string"
    echo "  -h, --help        Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 up                    # Run all pending migrations"
    echo "  $0 up 5                  # Run 5 migrations up"
    echo "  $0 down 3                # Rollback 3 migrations"
    echo "  $0 goto 10               # Migrate to version 10"
    echo "  $0 force 0               # Reset migration state to 0"
    echo "  $0 create add_new_table  # Create new migration files"
}

# Function to check if migrate command exists
check_migrate() {
    if ! command -v migrate &> /dev/null; then
        echo -e "${RED}Error: migrate command not found${NC}"
        echo "Please install golang-migrate:"
        echo "  go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
        exit 1
    fi
}

# Function to validate database URL
validate_database() {
    if [ -z "$DATABASE_URL" ]; then
        echo -e "${RED}Error: Database URL not provided${NC}"
        echo "Use -d or --database option to specify database connection string"
        exit 1
    fi
}

# Function to run migrations up
run_up() {
    local count=${1:-""}
    if [ -n "$count" ]; then
        echo -e "${BLUE}Running $count migrations up...${NC}"
        migrate -path "$MIGRATIONS_PATH" -database "$DATABASE_URL" step "$count"
    else
        echo -e "${BLUE}Running all pending migrations...${NC}"
        migrate -path "$MIGRATIONS_PATH" -database "$DATABASE_URL" up
    fi
}

# Function to run migrations down
run_down() {
    local count=${1:-1}
    echo -e "${BLUE}Rolling back $count migrations...${NC}"
    migrate -path "$MIGRATIONS_PATH" -database "$DATABASE_URL" down "$count"
}

# Function to goto specific version
goto_version() {
    local version=$1
    echo -e "${BLUE}Migrating to version $version...${NC}"
    migrate -path "$MIGRATIONS_PATH" -database "$DATABASE_URL" goto "$version"
}

# Function to force version
force_version() {
    local version=$1
    echo -e "${YELLOW}Forcing migration version to $version...${NC}"
    migrate -path "$MIGRATIONS_PATH" -database "$DATABASE_URL" force "$version"
}

# Function to show version
show_version() {
    echo -e "${BLUE}Current migration version:${NC}"
    migrate -path "$MIGRATIONS_PATH" -database "$DATABASE_URL" version
}

# Function to show status
show_status() {
    echo -e "${BLUE}Migration status:${NC}"
    migrate -path "$MIGRATIONS_PATH" -database "$DATABASE_URL" version
}

# Function to create new migration
create_migration() {
    local name=$1
    local timestamp=$(date +%Y%m%d%H%M%S)
    local version=$(printf "%06d" $(ls up/*.up.sql | wc -l | awk '{print $1 + 1}'))
    
    echo -e "${BLUE}Creating new migration: ${version}_${name}${NC}"
    
    # Create up migration
    cat > "up/${version}_${name}.up.sql" << EOF
-- Up migration for ${name}
-- Created: $(date)

-- Add your table creation or modification SQL here


EOF

    # Create down migration
    cat > "down/${version}_${name}.down.sql" << EOF
-- Down migration for ${name}
-- Created: $(date)

-- Add your rollback SQL here


EOF

    # Create symbolic links
    ln -sf "up/${version}_${name}.up.sql" "${version}_${name}.up.sql"
    ln -sf "down/${version}_${name}.down.sql" "${version}_${name}.down.sql"
    
    echo -e "${GREEN}Migration files created:${NC}"
    echo "  up/${version}_${name}.up.sql"
    echo "  down/${version}_${name}.down.sql"
    echo "  ${version}_${name}.up.sql (symlink)"
    echo "  ${version}_${name}.down.sql (symlink)"
}

# Function to list migrations
list_migrations() {
    echo -e "${BLUE}Available migrations:${NC}"
    echo ""
    echo -e "${GREEN}Up migrations:${NC}"
    ls -la up/*.up.sql | sort
    echo ""
    echo -e "${GREEN}Down migrations:${NC}"
    ls -la down/*.down.sql | sort
    echo ""
    echo -e "${GREEN}Symlinks in root:${NC}"
    ls -la *.sql | sort
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -d|--database)
            DATABASE_URL="$2"
            shift 2
            ;;
        -h|--help)
            show_usage
            exit 0
            ;;
        up|down|goto|force|version|status|create|list)
            COMMAND="$1"
            shift
            ;;
        *)
            ARGS+=("$1")
            shift
            ;;
    esac
done

# Check if migrate command exists
check_migrate

# Execute command
case $COMMAND in
    up)
        validate_database
        run_up "${ARGS[0]}"
        ;;
    down)
        validate_database
        run_down "${ARGS[0]}"
        ;;
    goto)
        validate_database
        if [ -z "${ARGS[0]}" ]; then
            echo -e "${RED}Error: Version number required${NC}"
            exit 1
        fi
        goto_version "${ARGS[0]}"
        ;;
    force)
        validate_database
        if [ -z "${ARGS[0]}" ]; then
            echo -e "${RED}Error: Version number required${NC}"
            exit 1
        fi
        force_version "${ARGS[0]}"
        ;;
    version)
        validate_database
        show_version
        ;;
    status)
        validate_database
        show_status
        ;;
    create)
        if [ -z "${ARGS[0]}" ]; then
            echo -e "${RED}Error: Migration name required${NC}"
            exit 1
        fi
        create_migration "${ARGS[0]}"
        ;;
    list)
        list_migrations
        ;;
    *)
        echo -e "${RED}Error: Unknown command '$COMMAND'${NC}"
        show_usage
        exit 1
        ;;
esac

echo -e "${GREEN}Migration operation completed successfully!${NC}"
