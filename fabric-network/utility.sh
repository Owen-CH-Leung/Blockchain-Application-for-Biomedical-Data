#!/bin/bash

COL_RED=`tput setaf 1`
COL_GREEN=`tput setaf 2`
COL_YELLOW=`tput setaf 3`
COL_RESET=`tput sgr0`

function print() {
  echo "$1"
}

# Info returns text in green
function Info() {
  print "${COL_GREEN}${1}${COL_RESET}"
}

# Alert returns text in yello
function Alert() {
  print "${COL_YELLOW}${1}${COL_RESET}"
}
# Error print text in red
function Error() {
  print "${COL_RED}${1}${COL_RESET}"
}

function errorCatch() {
  if [ $1 -ne 0 ]; then
    Error "$2"
  fi
}

export -f Info
export -f Error
export -f Alert
export -f errorCatch