package install

const bashCompletion = `#/usr/bin/env bash
_infrasonar_completions()
{
    local cur="${COMP_WORDS[COMP_CWORD]}"
    local prev="${COMP_WORDS[COMP_CWORD-1]}"

    if [[ "$prev" == "get" ]]; then
        local ADD_COMPLETES="assets all-asset-kinds"
        COMPREPLY=( $(compgen -W "$ADD_COMPLETES" -- ${COMP_WORDS[COMP_CWORD]}) )
    else
        # Handle other cases (main command completions)
        local COMPLETES="version install config get"
        COMPREPLY=( $(compgen -W "$COMPLETES" -- ${COMP_WORDS[COMP_CWORD]}) )
    fi

    return 0
}

complete -F _infrasonar_completions infrasonar
`
