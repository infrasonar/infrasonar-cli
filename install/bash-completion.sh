#/usr/bin/env bash
_infrasonar_completions()
{
    local cur prev
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    if [[ "${COMP_WORDS[1]}" == "config" ]]; then
        if [[ "$COMP_CWORD" == "2" ]]; then
            local COMPLETES="list new update default delete"
            COMPREPLY=( $(compgen -W "$COMPLETES" -- ${COMP_WORDS[COMP_CWORD]}) )
            return 0
        fi

        if [[ "${COMP_WORDS[2]}" == "list" ]]; then
            if [[ "$cur" == --* ]]; then
                local COMPLETES="--more --help"
                COMPREPLY=( $(compgen -W "$COMPLETES" -- ${COMP_WORDS[COMP_CWORD]}) )
                return 0
            fi
        fi

        if [[ "${COMP_WORDS[2]}" == "new" ]]; then
            if [[ "$cur" == --* ]]; then
                local COMPLETES="--set-name --set-token --set-api --set-output --set-default --help"
                COMPREPLY=( $(compgen -W "$COMPLETES" -- ${COMP_WORDS[COMP_CWORD]}) )
                return 0
            fi
        fi

        if [[ "${COMP_WORDS[2]}" == "update" ]]; then
            if [[ "$cur" == --* ]]; then
                local COMPLETES="--config --set-token --set-api --set-output --set-default --help"
                COMPREPLY=( $(compgen -W "$COMPLETES" -- ${COMP_WORDS[COMP_CWORD]}) )
                return 0
            fi
        fi

        if [[ "${COMP_WORDS[2]}" == "default" ]]; then
            if [[ "$cur" == --* ]]; then
                local COMPLETES="--help"
                COMPREPLY=( $(compgen -W "$COMPLETES" -- ${COMP_WORDS[COMP_CWORD]}) )
                return 0
            fi
        fi

        if [[ "${COMP_WORDS[2]}" == "delete" ]]; then
            if [[ "$cur" == --* ]]; then
                local COMPLETES="--config --help"
                COMPREPLY=( $(compgen -W "$COMPLETES" -- ${COMP_WORDS[COMP_CWORD]}) )
                return 0
            fi
        fi
        return 0
    fi

    if [[ "${COMP_WORDS[1]}" == "get" ]]; then

        if [[ "$COMP_CWORD" == "2" ]]; then
            local COMPLETES="assets collectors me all-asset-kinds all-label-colors"
            COMPREPLY=( $(compgen -W "$COMPLETES" -- ${COMP_WORDS[COMP_CWORD]}) )
            return 0
        fi

        if [[ "$prev" == "-o" ]] || [[ "$prev" == "--output" ]]; then
            local COMPLETES="yaml json simple"
            COMPREPLY=( $(compgen -W "$COMPLETES" -- ${COMP_WORDS[COMP_CWORD]}) )
            return 0
        fi

        if [[ "$prev" == "-t" ]] || [[ "$prev" == "--target-filename" ]]; then
            compopt -o nospace

            if [[ "$cur" == "" ]]; then
                COMPREPLY=( $(compgen -d) )
            else
                COMPREPLY=( $(compgen -d -- "$cur") )
            fi

            # Add trailing slash to each completion
            for i in "${!COMPREPLY[@]}"; do
                COMPREPLY[$i]="${COMPREPLY[$i]}/"
            done
            return 0
        fi

        if [[ "$prev" == "-u" ]] || [[ "$prev" == "--use-config" ]]; then
            local COMPLETES=$(infrasonar config list 2>/dev/null)
            if [[ -z "$COMPLETES" ]]; then
                return 0
            fi

            COMPREPLY=( $(compgen -W "$COMPLETES" -- ${cur}) )
            return 0
        fi

        if [[ "${COMP_WORDS[2]}" == "assets" ]]; then
            if [[ "$prev" == "-c" ]] || [[ "$prev" == "--container" ]]; then
                return 0
            fi

            if [[ "$prev" == "-a" ]] || [[ "$prev" == "--asset" ]]; then
                return 0
            fi

            if [[ "$prev" == "-p" ]] || [[ "$prev" == "--properties" ]]; then
                return 0
            fi

            if [[ "$prev" == "-f" ]] || [[ "$prev" == "--filter" ]]; then
                local COMPLETES="kind== kind!= collector== collector!= label== label!= zone== zone!="
                compopt -o nospace
                COMPREPLY=( $(compgen -W "$COMPLETES" -- ${cur}) )
                return 0
            fi

            if [[ "$cur" == --* ]]; then
                local COMPLETES="--container --asset --properties --filter --include-defaults --output --target-filename --use-config --help"
                COMPREPLY=( $(compgen -W "$COMPLETES" -- ${COMP_WORDS[COMP_CWORD]}) )
                return 0
            fi
            return 0
        fi

        if [[ "${COMP_WORDS[2]}" == "collectors" ]]; then
            if [[ "$prev" == "-c" ]] || [[ "$prev" == "--container" ]]; then
                return 0
            fi

            if [[ "$prev" == "-p" ]] || [[ "$prev" == "--properties" ]]; then
                return 0
            fi

            if [[ "$prev" == "-k" ]] || [[ "$prev" == "--collector" ]]; then
                return 0
            fi

            if [[ "$cur" == --* ]]; then
                local COMPLETES="--container --collector --properties --output --target-filename --use-config --help"
                COMPREPLY=( $(compgen -W "$COMPLETES" -- ${COMP_WORDS[COMP_CWORD]}) )
                return 0
            fi
            return 0
        fi

        if [[ "${COMP_WORDS[2]}" == "me" ]]; then
            if [[ "$prev" == "-c" ]] || [[ "$prev" == "--container" ]]; then
                return 0
            fi

            if [[ "$prev" == "-p" ]] || [[ "$prev" == "--properties" ]]; then
                return 0
            fi

            if [[ "$cur" == --* ]]; then
                local COMPLETES="--container --properties --output --target-filename --use-config --help"
                COMPREPLY=( $(compgen -W "$COMPLETES" -- ${COMP_WORDS[COMP_CWORD]}) )
                return 0
            fi
            return 0
        fi

        if [[ "${COMP_WORDS[2]}" == "all-asset-kinds" ]] || [[ "${COMP_WORDS[2]}" == "all-label-colors" ]]; then
            if [[ "$cur" == --* ]]; then
                local COMPLETES="--output --target-filename --use-config --help"
                COMPREPLY=( $(compgen -W "$COMPLETES" -- ${COMP_WORDS[COMP_CWORD]}) )
                return 0
            fi
            return 0
        fi

        return 0
    fi

    if [[ "${COMP_WORDS[1]}" == "apply" ]]; then

        if [[ "$prev" == "-f" ]] || [[ "$prev" == "--filename" ]]; then
            local FILEPATH COMPLETES
            FILEPATH="$(dirname "${cur}")";

            if [[ "$cur" == "" ]]; then
                FILEPATH="."
            fi

            COMPLETES=$(find "$FILEPATH" -maxdepth 2 -type f \( -iname \*.json -o -iname \*.yaml -o -iname \*.yml \) 2>/dev/null)
            if [[ -z "$COMPLETES" ]]; then
                return 0
            fi
            COMPREPLY=( $(compgen -W "$COMPLETES" -- ${cur}) )
            return 0
        fi

        if [[ "$prev" == "-u" ]] || [[ "$prev" == "--use-config" ]]; then
            local COMPLETES=$(infrasonar config list 2>/dev/null)
            if [[ -z "$OPTIONS" ]]; then
                return 0
            fi
            COMPREPLY=( $(compgen -W "$COMPLETES" -- ${cur}) )
            return 0
        fi

        if [[ "$cur" == --* ]]; then
            local COMPLETES="--filename --dry-run --purge --use-config --help"
            COMPREPLY=( $(compgen -W "$COMPLETES" -- ${COMP_WORDS[COMP_CWORD]}) )
            return 0
        fi
        return 0
    fi

    local COMPLETES="version install config get apply"
    COMPREPLY=( $(compgen -W "$COMPLETES" -- ${COMP_WORDS[COMP_CWORD]}) )
    return 0
}

complete -F _infrasonar_completions infrasonar