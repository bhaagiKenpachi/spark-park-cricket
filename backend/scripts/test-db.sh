#!/bin/bash

# Test Database Management Script
# Handles setup, auto-setup, and verification of testing_db schema

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
log() { echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"; }
success() { echo -e "${GREEN}✅${NC} $1"; }
warning() { echo -e "${YELLOW}⚠️${NC} $1"; }
error() { echo -e "${RED}❌${NC} $1"; }

usage() {
    echo "Test Database Management Script"
    echo ""
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  setup         Setup test database schema (manual instructions)"
    echo "  auto-setup    Attempt automatic setup via API"
    echo "  verify        Verify schema is properly configured"
    echo "  status        Check database connectivity and schema status"
    echo "  help          Show this help message"
    echo ""
    echo "Required Environment Variables:"
    echo "  SUPABASE_URL      Supabase project URL"
    echo "  SUPABASE_API_KEY  Supabase API key"
    echo ""
    echo "Examples:"
    echo "  $0 setup"
    echo "  $0 verify"
    echo "  $0 status"
}

# Check prerequisites
check_prerequisites() {
    if [ -z "$SUPABASE_URL" ]; then
        error "SUPABASE_URL environment variable is not set"
        exit 1
    fi

    if [ -z "$SUPABASE_API_KEY" ]; then
        error "SUPABASE_API_KEY environment variable is not set"
        exit 1
    fi

    SCHEMA_FILE="internal/database/migrations/complete_schema_testing_db.sql"
    if [ ! -f "$SCHEMA_FILE" ]; then
        error "Schema file not found: $SCHEMA_FILE"
        exit 1
    fi

    success "Prerequisites check passed"
}

# Test database connectivity
test_connectivity() {
    log "Testing database connectivity..."
    if curl -s -H "apikey: $SUPABASE_API_KEY" -H "Authorization: Bearer $SUPABASE_API_KEY" "$SUPABASE_URL/rest/v1/" > /dev/null; then
        success "Database connectivity test successful"
        return 0
    else
        error "Database connectivity test failed"
        echo "Please check your SUPABASE_URL and SUPABASE_API_KEY"
        return 1
    fi
}

# Check if testing_db schema exists and is accessible
check_schema_access() {
    log "Testing testing_db schema access..."
    if curl -s -H "apikey: $SUPABASE_API_KEY" -H "Authorization: Bearer $SUPABASE_API_KEY" \
        "$SUPABASE_URL/rest/v1/users?select=id&limit=1" \
        -H "Accept-Profile: testing_db" > /dev/null 2>&1; then
        return 0
    else
        return 1
    fi
}

# Verify schema tables
verify_schema_tables() {
    log "Testing key tables in testing_db schema..."
    
    TABLES=("users" "series" "matches" "innings" "overs" "balls")
    for table in "${TABLES[@]}"; do
        if curl -s -H "apikey: $SUPABASE_API_KEY" -H "Authorization: Bearer $SUPABASE_API_KEY" \
            "$SUPABASE_URL/rest/v1/$table?select=id&limit=1" \
            -H "Accept-Profile: testing_db" > /dev/null 2>&1; then
            success "Table '$table' is accessible"
        else
            error "Table '$table' is not accessible"
            return 1
        fi
    done
    return 0
}

# Show manual setup instructions
show_manual_setup() {
    echo ""
    echo "📋 Manual Setup Instructions:"
    echo "1. Go to Supabase Dashboard: $SUPABASE_URL/project/default/sql"
    echo "2. Execute the SQL from: $SCHEMA_FILE"
    echo "3. Ensure the 'testing_db' schema is created with all required tables"
    echo ""
    echo "📄 Schema file preview:"
    echo "=========================================="
    head -20 "$SCHEMA_FILE"
    echo "=========================================="
    echo ""
}

# Setup command (manual instructions)
setup_command() {
    log "Setting up test database schema..."
    
    check_prerequisites
    test_connectivity
    
    warning "Manual schema setup required"
    show_manual_setup
    
    success "Test database setup script completed"
    warning "Manual schema setup required - see instructions above"
}

# Auto-setup command (attempt via API)
auto_setup_command() {
    log "Attempting to auto-setup test database schema..."
    
    check_prerequisites
    test_connectivity
    
    SCHEMA_FILE="internal/database/migrations/complete_schema_testing_db.sql"
    SCHEMA_CONTENT=$(cat "$SCHEMA_FILE")
    
    log "Attempting to create testing_db schema..."
    
    # Try to execute the schema using Supabase's SQL execution endpoint
    RESPONSE=$(curl -s -w "%{http_code}" -X POST \
      -H "apikey: $SUPABASE_API_KEY" \
      -H "Authorization: Bearer $SUPABASE_API_KEY" \
      -H "Content-Type: application/json" \
      -d "{\"query\": \"$SCHEMA_CONTENT\"}" \
      "$SUPABASE_URL/rest/v1/rpc/exec_sql" 2>/dev/null || echo "000")
    
    HTTP_CODE="${RESPONSE: -3}"
    
    if [ "$HTTP_CODE" = "200" ]; then
        success "Schema setup completed successfully"
    else
        warning "Automated schema setup failed (HTTP: $HTTP_CODE)"
        echo "This is expected if SQL execution API is not enabled"
    fi
    
    # Verify the schema was created
    log "Verifying schema setup..."
    if check_schema_access; then
        success "testing_db schema is accessible"
        success "Test database schema setup completed successfully!"
    else
        error "testing_db schema is not accessible"
        show_manual_setup
        error "Integration tests will fail without proper schema setup!"
        exit 1
    fi
}

# Verify command
verify_command() {
    log "Verifying test database schema..."
    
    check_prerequisites
    test_connectivity
    
    if check_schema_access; then
        success "testing_db schema is accessible"
        
        if verify_schema_tables; then
            success "All test database schema verification completed successfully!"
            success "testing_db schema is properly configured and ready for integration tests"
        else
            error "Schema setup may be incomplete!"
            exit 1
        fi
    else
        error "testing_db schema is not accessible or doesn't exist"
        echo ""
        error "CRITICAL: The testing_db schema is not properly set up!"
        echo ""
        show_manual_setup
        error "Integration tests will fail without proper schema setup!"
        exit 1
    fi
}

# Status command
status_command() {
    log "Checking database status..."
    
    check_prerequisites
    
    if test_connectivity; then
        if check_schema_access; then
            success "testing_db schema is accessible"
            if verify_schema_tables; then
                success "All tables are accessible - schema is properly configured"
            else
                warning "Some tables are not accessible - schema may be incomplete"
            fi
        else
            warning "testing_db schema is not accessible - needs setup"
            echo ""
            echo "Run: $0 setup (for manual setup instructions)"
            echo "Run: $0 auto-setup (to attempt automatic setup)"
        fi
    else
        error "Database connectivity failed - check environment variables"
    fi
}

# Main execution
main() {
    case "${1:-help}" in
        setup)
            setup_command
            ;;
        auto-setup)
            auto_setup_command
            ;;
        verify)
            verify_command
            ;;
        status)
            status_command
            ;;
        help|--help|-h)
            usage
            ;;
        *)
            error "Unknown command: $1"
            usage
            exit 1
            ;;
    esac
}

# Run main function
main "$@"
