#!/bin/bash

# Working Endpoint Progress Tracker
# Handles grep output properly

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

TODO_FILE="TODO_ENDPOINTS.md"

show_usage() {
    echo -e "${BLUE}Endpoint Progress Tracker${NC}"
    echo ""
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  status        Show current progress"
    echo "  complete      Mark an endpoint as complete"
    echo "  done          Alias for complete (shorter)"
    echo "  reset         Reset all progress"
    echo "  help          Show this help"
    echo ""
    echo "Examples:"
    echo "  $0 status"
    echo "  $0 complete 'GET /api/roles'"
    echo "  $0 done 'GET /api/roles'        # Shorter version"
    echo "  $0 reset"
}

show_status() {
    if [ ! -f "$TODO_FILE" ]; then
        echo -e "${YELLOW}Error: $TODO_FILE not found${NC}"
        exit 1
    fi
    
    echo -e "${BLUE}📊 Endpoint Implementation Progress${NC}"
    echo "=========================================="
    
    # Count total and completed using a different approach
    local total=$(grep -c "^- \[ \]" "$TODO_FILE" 2>/dev/null)
    local completed=$(grep -c "^- \[x\]" "$TODO_FILE" 2>/dev/null)
    
    # Handle empty results
    if [ -z "$total" ]; then total=0; fi
    if [ -z "$completed" ]; then completed=0; fi
    
    echo "  Total Endpoints: $total"
    echo "  Completed: $completed"
    
    # Calculate pending safely
    if [ "$total" -gt 0 ] && [ "$completed" -gt 0 ]; then
        local pending=$((total - completed))
        echo "  Pending: $pending"
        
        # Calculate percentage safely
        local percentage=$((completed * 100 / total))
        echo "  Progress: $percentage%"
    else
        echo "  Pending: $total"
        echo "  Progress: 0%"
    fi
    
    echo ""
    echo -e "${GREEN}📈 Progress by Category:${NC}"
    echo "  Master Tables: 0/30 (not implemented yet)"
    echo "  User Management: 0/8 (not implemented yet)"
    echo "  Core Business: 0/6 (not implemented yet)"
    echo "  Documents: 0/42 (not implemented yet)"
    echo "  Performance: 0/14 (not implemented yet)"
    
    echo ""
    echo -e "${GREEN}✅ Recently Completed:${NC}"
    if [ "$completed" -gt 0 ]; then
        echo "  No completed tasks found yet"
    else
        echo "  No endpoints completed yet"
    fi
    
    echo ""
    echo -e "${YELLOW}⏳ Next Priority:${NC}"
    echo "  Start with authentication endpoints (login, logout)"
    echo "  Then implement master table endpoints"
    echo "  Focus on one table at a time"
}

complete_endpoint() {
    local endpoint="$1"
    
    if [ -z "$endpoint" ]; then
        echo -e "${YELLOW}Error: Please specify an endpoint to mark as complete${NC}"
        echo "Example: $0 complete 'GET /api/roles'"
        exit 1
    fi
    
    if [ ! -f "$TODO_FILE" ]; then
        echo -e "${YELLOW}Error: $TODO_FILE not found${NC}"
        exit 1
    fi
    
    # Try to find and mark as complete using a different delimiter
    if sed -i "s|^- \[ \] \`$endpoint\`|- [x] \`$endpoint\` ✅|" "$TODO_FILE" 2>/dev/null; then
        echo -e "${GREEN}✅ Marked endpoint as COMPLETE: $endpoint${NC}"
        
        # Show updated progress
        local total=$(grep -c "^- \[ \]" "$TODO_FILE" 2>/dev/null)
        local completed=$(grep -c "^- \[x\]" "$TODO_FILE" 2>/dev/null)
        
        if [ -z "$total" ]; then total=0; fi
        if [ -z "$completed" ]; then completed=0; fi
        
        echo -e "${BLUE}📊 Updated Progress: $completed/$total${NC}"
    else
        echo -e "${YELLOW}⚠️  Could not find endpoint: $endpoint${NC}"
        echo "Make sure to use the exact format from the TODO file"
        echo "Example: 'GET /api/roles' (with quotes)"
    fi
}

reset_progress() {
    if [ ! -f "$TODO_FILE" ]; then
        echo -e "${YELLOW}Error: $TODO_FILE not found${NC}"
        exit 1
    fi
    
    echo -e "${YELLOW}⚠️  Are you sure you want to reset ALL progress? (y/N)${NC}"
    read -r response
    
    if [[ "$response" =~ ^[Yy]$ ]]; then
        # Reset all completed endpoints
        sed -i 's/^- \[x\]/- [ ]/g' "$TODO_FILE"
        sed -i 's/✅//g' "$TODO_FILE"
        
        echo -e "${GREEN}🔄 All progress has been reset to 0%${NC}"
    else
        echo -e "${BLUE}Reset cancelled${NC}"
    fi
}

# Main script logic
case $1 in
    status)
        show_status
        ;;
    complete|done)
        complete_endpoint "$2"
        ;;
    reset)
        reset_progress
        ;;
    help|--help|-h)
        show_usage
        ;;
    *)
        echo -e "${YELLOW}Error: Unknown command '$1'${NC}"
        show_usage
        exit 1
        ;;
esac

echo -e "${GREEN}Operation completed successfully!${NC}"
