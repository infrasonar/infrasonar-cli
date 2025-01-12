#/usr/bin/env bash
_infrasonar_completions()
{
    local cur="${COMP_WORDS[COMP_CWORD]}"
    local prev="${COMP_WORDS[COMP_CWORD-1]}"
    local prevprev="${COMP_WORDS[COMP_CWORD-2]}"

    if [[ "${COMP_WORDS[1]}" == "config" ]]; then
        return 0
    fi

    if [[ "${COMP_WORDS[1]}" == "get" ]]; then

        if [[ "$COMP_CWORD" == "2" ]]; then
            # CMD: get
            local GET_COMPLETES="assets all-asset-kinds"
            COMPREPLY=( $(compgen -W "$GET_COMPLETES" -- ${COMP_WORDS[COMP_CWORD]}) )
            return 0
        fi

        if [[ "$prev" == "-o" ]] || [[ "$prev" == "--output" ]]; then
            # CMD: get --output
            local OUTPUT_COMPLETES="yaml json simple"
            COMPREPLY=( $(compgen -W "$OUTPUT_COMPLETES" -- ${COMP_WORDS[COMP_CWORD]}) )
            return 0
        fi

        if [[ "$prev" == "--config" ]]; then
            # CMD: get --config
            local options=$(infrasonar config list 2>/dev/null)

            # Handle potential errors (e.g., empty output)
            if [[ -z "$options" ]]; then
                return 0
            fi

            # Generate completions from the captured output
            COMPREPLY=( $(compgen -W "$options" -- ${cur}) )
            return 0
        fi

        if [[ "${COMP_WORDS[2]}" == "assets" ]]; then
            if [[ "$prev" == "-c" ]] || [[ "$prev" == "--container" ]]; then
                # CMD: get assets --container
                return 0
            fi

            if [[ "$prev" == "-f" ]] || [[ "$prev" == "--filter" ]]; then
                # CMD: get assets --filter
                local FILTER_ARGS="kind== kind!= collector== collector!= label== label!= zone== zone!="
                COMPREPLY=( $(compgen -o nospace -W "$FILTER_ARGS" -- ${cur}) )
                return 0
            fi

            if [[ "$prev" == "-p" ]] || [[ "$prev" == "--properties" ]]; then
                # CMD: get assets --properties
                return 0
            fi

            # local c_seen=0
            # # Check if the argument has already been seen
            # for ((i=1; i<=$COMP_CWORD; i++)); do
            #     if [[ "${COMP_WORDS[$i]}" == "--container" ]]; then
            #         c_seen=1
            #     fi
            # done

            # CMD: get assets
            # if [[ $arg_seen -eq 0 ]]; then
            #     local GET_ASSETS_COMPLETES="--container --properties --filter --output --config"
            #     COMPREPLY=( $(compgen -W "$GET_ASSETS_COMPLETES" -- ${COMP_WORDS[COMP_CWORD]}) )
            # fi
            local GET_ASSETS_COMPLETES="--container --properties --filter --output --config"
            COMPREPLY=( $(compgen -W "$GET_ASSETS_COMPLETES" -- ${COMP_WORDS[COMP_CWORD]}) )
            return 0
        fi

        return 0
    fi

    # Handle other cases (main command completions)
    local COMPLETES="version install config get"
    COMPREPLY=( $(compgen -W "$COMPLETES" -- ${COMP_WORDS[COMP_CWORD]}) )

    return 0
}

complete -F _infrasonar_completions infrasonar

# https://opensource.com/article/18/3/creating-bash-completion-script