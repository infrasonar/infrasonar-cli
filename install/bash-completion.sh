#/usr/bin/env bash
_infrasonar_completions()
{
    local cur prev prevprev opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    prevprev="${COMP_WORDS[COMP_CWORD-2]}"

    if [[ "${COMP_WORDS[1]}" == "config" ]]; then
        if [[ "$COMP_CWORD" == "2" ]]; then
            # CMD: config
            local CONFIG_COMPLETES="list new update default delete"
            COMPREPLY=( $(compgen -W "$CONFIG_COMPLETES" -- ${COMP_WORDS[COMP_CWORD]}) )
            return 0
        fi

        if [[ "${COMP_WORDS[2]}" == "list" ]]; then
            if [[ "$cur" == --* ]]; then
                local LIST_CONFIG_COMPLETES="--more --help"
                COMPREPLY=( $(compgen -W "$LIST_CONFIG_COMPLETES" -- ${COMP_WORDS[COMP_CWORD]}) )
                return 0
            fi
        fi

        if [[ "${COMP_WORDS[2]}" == "new" ]]; then
            if [[ "$cur" == --* ]]; then
                local NEW_CONFIG_COMPLETES="--set-name --set-token --set-api --set-output --set-default --help"
                COMPREPLY=( $(compgen -W "$NEW_CONFIG_COMPLETES" -- ${COMP_WORDS[COMP_CWORD]}) )
                return 0
            fi
        fi

        if [[ "${COMP_WORDS[2]}" == "update" ]]; then
            if [[ "$cur" == --* ]]; then
                local UPD_CONFIG_COMPLETES="--config --set-token --set-api --set-output --set-default --help"
                COMPREPLY=( $(compgen -W "$UPD_CONFIG_COMPLETES" -- ${COMP_WORDS[COMP_CWORD]}) )
                return 0
            fi
        fi

        if [[ "${COMP_WORDS[2]}" == "default" ]]; then
            if [[ "$cur" == --* ]]; then
                local DEF_CONFIG_COMPLETES="--help"
                COMPREPLY=( $(compgen -W "$DEF_CONFIG_COMPLETES" -- ${COMP_WORDS[COMP_CWORD]}) )
                return 0
            fi
        fi

        if [[ "${COMP_WORDS[2]}" == "delete" ]]; then
            if [[ "$cur" == --* ]]; then
                local DEL_CONFIG_COMPLETES="--config --help"
                COMPREPLY=( $(compgen -W "$DEL_CONFIG_COMPLETES" -- ${COMP_WORDS[COMP_CWORD]}) )
                return 0
            fi
        fi

        return 0
    fi

    if [[ "${COMP_WORDS[1]}" == "get" ]]; then

        if [[ "$COMP_CWORD" == "2" ]]; then
            # CMD: get
            local GET_COMPLETES="assets collectors me all-asset-kinds all-label-colors"
            COMPREPLY=( $(compgen -W "$GET_COMPLETES" -- ${COMP_WORDS[COMP_CWORD]}) )
            return 0
        fi

        if [[ "$prev" == "-o" ]] || [[ "$prev" == "--output" ]]; then
            # CMD: get --output
            local OUTPUT_COMPLETES="yaml json simple"
            COMPREPLY=( $(compgen -W "$OUTPUT_COMPLETES" -- ${COMP_WORDS[COMP_CWORD]}) )
            return 0
        fi

        if [[ "$prev" == "-u" ]] || [[ "$prev" == "--use-config" ]]; then
            # CMD: get --config
            local OPTIONS=$(infrasonar config list 2>/dev/null)

            # Handle potential errors (e.g., empty output)
            if [[ -z "$OPTIONS" ]]; then
                return 0
            fi

            # Generate completions from the captured output
            COMPREPLY=( $(compgen -W "$OPTIONS" -- ${cur}) )
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

            if [[ "$cur" == --* ]]; then
                local GET_ASSETS_COMPLETES="--container --assets --properties --filter --include-defaults --output --target-filename --use-config --help"
                COMPREPLY=( $(compgen -W "$GET_ASSETS_COMPLETES" -- ${COMP_WORDS[COMP_CWORD]}) )
                return 0
            fi

            return 0
        fi

        return 0
    fi

    if [[ "${COMP_WORDS[1]}" == "apply" ]]; then

        if [[ "$prev" == "-f" ]] || [[ "$prev" == "--filename" ]]; then
            local FILEPATH="$(dirname "${cur}")";

            if [[ "$cur" == "" ]]; then
                FILEPATH="."
            fi

            local FILES=$(find "$FILEPATH" -maxdepth 2 -type f \( -iname \*.json -o -iname \*.yaml -o -iname \*.yml \) 2>/dev/null)
            local OPTIONS="$DIRS $FILES"

            # Handle potential errors (e.g., empty output)
            if [[ -z "$OPTIONS" ]]; then
                return 0
            fi

            # Generate completions from the captured output
            COMPREPLY=( $(compgen -W "$OPTIONS" -- ${cur}) )
            return 0
        fi

        if [[ "$prev" == "-u" ]] || [[ "$prev" == "--use-config" ]]; then
            # CMD: get --config
            local OPTIONS=$(infrasonar config list 2>/dev/null)

            # Handle potential errors (e.g., empty output)
            if [[ -z "$OPTIONS" ]]; then
                return 0
            fi

            # Generate completions from the captured output
            COMPREPLY=( $(compgen -W "$OPTIONS" -- ${cur}) )
            return 0
        fi

        if [[ "$cur" == --* ]]; then
            local GET_ASSETS_COMPLETES="--filename --dry-run --purge --use-config --help"
            COMPREPLY=( $(compgen -W "$GET_ASSETS_COMPLETES" -- ${COMP_WORDS[COMP_CWORD]}) )
            return 0
        fi

        return 0
    fi

    # Handle other cases (main command completions)
    local COMPLETES="version install config get apply"
    COMPREPLY=( $(compgen -W "$COMPLETES" -- ${COMP_WORDS[COMP_CWORD]}) )

    return 0
}

complete -F _infrasonar_completions infrasonar

# https://opensource.com/article/18/3/creating-bash-completion-script