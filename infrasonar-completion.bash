#/usr/bin/env bash
_dothis_completions()
{
  COMPREPLY+=("now")
  COMPREPLY+=("tomorrow")
  COMPREPLY+=("never")
}

complete -F _infrasonar_completions dothis

# https://opensource.com/article/18/3/creating-bash-completion-script