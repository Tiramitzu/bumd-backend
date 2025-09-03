#!/bin/bash

# Setup script for endpoint tracker autocomplete
# This script will add autocomplete to your bash session

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}🚀 Setting up Endpoint Tracker Autocomplete${NC}"
echo "================================================"

# Check if completion script exists
if [ ! -f "endpoint_tracker_completion.sh" ]; then
    echo -e "${YELLOW}Error: endpoint_tracker_completion.sh not found${NC}"
    exit 1
fi

# Source the completion script
echo -e "${BLUE}📝 Sourcing completion script...${NC}"
source endpoint_tracker_completion.sh

# Test if it's working
echo -e "${BLUE}🧪 Testing autocomplete...${NC}"
echo "Type: ./endpoint_tracker.sh <TAB> to see available commands"
echo "Type: ./endpoint_tracker.sh done <TAB> to see available endpoints"

# Add to .bashrc for permanent setup
echo ""
echo -e "${YELLOW}💡 For permanent setup, add this line to your ~/.bashrc:${NC}"
echo "source $(pwd)/endpoint_tracker_completion.sh"

# Check if already in .bashrc
if grep -q "endpoint_tracker_completion.sh" ~/.bashrc 2>/dev/null; then
    echo -e "${GREEN}✅ Already configured in ~/.bashrc${NC}"
else
    echo ""
    echo -e "${YELLOW}🔧 Would you like me to add it to ~/.bashrc now? (y/N)${NC}"
    read -r response
    
    if [[ "$response" =~ ^[Yy]$ ]]; then
        echo "" >> ~/.bashrc
        echo "# Endpoint Tracker Autocomplete" >> ~/.bashrc
        echo "source $(pwd)/endpoint_tracker_completion.sh" >> ~/.bashrc
        echo -e "${GREEN}✅ Added to ~/.bashrc${NC}"
        echo -e "${BLUE}📝 Restart your terminal or run: source ~/.bashrc${NC}"
    else
        echo -e "${BLUE}📝 You can add it manually later${NC}"
    fi
fi

echo ""
echo -e "${GREEN}🎉 Autocomplete setup complete!${NC}"
echo ""
echo -e "${BLUE}📋 Quick Reference:${NC}"
echo "  ./endpoint_tracker.sh <TAB>     - Show commands"
echo "  ./endpoint_tracker.sh done <TAB> - Show endpoints"
echo "  ./endpoint_tracker.sh done 'GET /api/roles' - Mark complete"
echo "  ./endpoint_tracker.sh status    - Check progress"
