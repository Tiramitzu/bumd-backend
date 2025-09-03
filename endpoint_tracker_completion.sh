# Bash completion script for endpoint_tracker.sh
# Source this file or add to your .bashrc

_endpoint_tracker_completion() {
    local cur prev opts cmds
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    
    # Available commands
    cmds="status complete done reset help"
    
    # If this is the first argument, complete with commands
    if [[ $COMP_CWORD -eq 1 ]]; then
        COMPREPLY=( $(compgen -W "${cmds}" -- "${cur}") )
        return 0
    fi
    
    # If the first argument is complete/done, suggest endpoints
    if [[ "${prev}" == "complete" || "${prev}" == "done" ]]; then
        # Read endpoints from TODO file
        if [ -f "TODO_ENDPOINTS.md" ]; then
            local endpoints=$(grep -o '`[^`]*`' TODO_ENDPOINTS.md | sed 's/`//g' | sort -u)
            COMPREPLY=( $(compgen -W "${endpoints}" -- "${cur}") )
        fi
        return 0
    fi
    
    return 0
}

# Register the completion function
complete -F _endpoint_tracker_completion endpoint_tracker.sh
complete -F _endpoint_tracker_completion ./endpoint_tracker.sh
