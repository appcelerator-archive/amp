#!/bin/bash

shopt -s extglob
unset FAIL
pgrep="pcregrep -o -e "

#==========================================================
# time
#==========================================================

# elapsed prints the number of seconds elapsed since the script
# began execution. It uses the bash special variable $SECONDS,
# which can be explicitly set to a starting value before
# calling this function, if desired (ex: `SECONDS=0`).
# If passed the previous value, elapsed will print the
# difference between them.
elapsed() {
  [[ $# -eq 1 ]] && printf $((( $SECONDS - $1 ))) || printf $SECONDS
}

# elapsed_hms calls elapsed and prints the result formated as
# 00h:00m:00s
elapsed_hms() {
  local s=$(elapsed $1)
  printf "%02dh:%02dm:%02ds" $(($s/3600)) $(($s%3600/60)) $(($s%60))
}

#==========================================================
# string comparisons
#==========================================================

streq() {
  [[ "$1" == "$2" ]]
}

strneq() {
  [[ "$1" != "$2" ]]
}

strlt() {
  [[ "$1" < "$2" ]]
}

strle() {
  [[ "$1" < "$2" || "$1" == "$2" ]]
}

strgt() {
  [[ "$1" > "$2" ]]
}

strge() {
  [[ "$1" > "$2" || "$1" == "$2" ]]
}

#==========================================================
# numeric comparisons
#==========================================================

numeq() {
  (( "$1" == "$2" ))
}

numneq() {
  (( "$1" != "$2" ))
}

numlt() {
  (( "$1" < "$2" ))
}

numle() {
  (( "$1" <= "$2" ))
}

numgt() {
  (( "$1" > "$2" ))
}

numge() {
  (( "$1" >= "$2" ))
}

#==========================================================
# assert
#==========================================================

# private implementation
_assert_result() {
  local status=$?
  local caller="${FUNCNAME[1]}"
  local pre=""
  local msg=""

  if [[ $FAIL -ne 0 ]]; then
    if [[ $status -eq 0 ]]; then
      # it SHOULD have failed, so it IS an error!
      status=1
      pre="failed because an error was expected: "
    else
      # error was expected, so it's NOT an error!
      status=0
    fi
  fi
  if [[ $status -ne 0 ]]; then
    if [[ $# -eq 0 ]]; then
      msg="failed: exit status: $status"
    elif [[ $# -eq 2 ]]; then
      msg=$(printf "actual: '%s', expected: '%s'\n" "$1" "$2")
    else
      msg="${*: -1:1}"
    fi
    echo "[$caller] $pre${msg:-assertion failed}"
  fi

  unset FAIL
  return $status
}

# assert tests that the previous command or function didn't exit with an error
# a custom error message can be supplied as an argument
assert() {
  _assert_result "$1"
}

assert_streq() {
  streq "$1" "$2"
  _assert_result "$1" "$2" "${3:-false: $1 == $2}"
}

assert_strneq() {
  strneq "$1" "$2"
  _assert_result "$1" "$2" "false: $1 != $2"
}

assert_strlt() {
  strlt "$1" "$2"
  _assert_result "$1" "$2" "false: $1 < $2"
}

assert_strle() {
  strle "$1" "$2"
  _assert_result "$1" "$2" "false: $1 <= $2"
}

assert_strgt() {
  strgt "$1" "$2"
  _assert_result "$1" "$2" "false: $1 > $2"
}

assert_strge() {
  strge "$1" "$2"
  _assert_result "$1" "$2" "false: $1 >= $2"
}

assert_numeq() {
  numeq "$1" "$2"
  _assert_result "$1" "$2" "false: $1 == $2"
}

assert_numneq() {
  numneq "$1" "$2"
  _assert_result "$1" "$2" "false: $1 != $2"
}

assert_numlt() {
  numlt "$1" "$2"
  _assert_result "$1" "$2" "false: $1 < $2"
}

assert_numle() {
  numle "$1" "$2"
  _assert_result "$1" "$2" "false: $1 <= $2"
}

assert_numgt() {
  numgt "$1" "$2"
  _assert_result "$1" "$2" "false: $1 > $2"
}

assert_numge() {
  numge "$1" "$2"
  _assert_result "$1" "$2" "false: $1 >= $2"
}

assert_numin() {
  numlt $2 $1 && numlt $1 $3
  _assert_result "$1" "$2" "$3" "false: $2 < $1 < $3"
}

assert_numin#() {
  numle $2 $1 && numle $1 $3
  _assert_result "$1" "$2" "$3" "false: $2 <= $1 <= $3"
}

#==========================================================
# string
#==========================================================

# trim substring from front
# $1: string
# $2: pattern (default: [:space:])
trimf() {
  echo "${1#${2:-+([[:space:]])}}"
  #echo "${1/${2:-+([[:space:]])}/}"
}

# trim substring from front (greedy)
# $1: string
# $2: pattern (default: [:space:])
trimf#() {
  echo "${1##${2:-+([[:space:]])}}"
}

# trim substring from end
# $1: string
# $2: pattern (default: [:space:])
trimb() {
  echo "${1%${2:-+([[:space:]])}}"
}

# trim substring from end (greedy)
# $1: string
# $2: pattern (default: [:space:])
trimb#() {
  echo "${1%%${2:-+([[:space:]])}}"
}

# trim substring from front and back
# $1: string
# $2: pattern (default: [:space:])
trim() {
  s="${1#${2:-+([[:space:]])}}"
  echo "${s%${2:-+([[:space:]])}}"
}

# trim substring from front and back (greedy)
# $1: string
# $2: pattern (default: [:space:])
trim#() {
  s="${1##${2:-+([[:space:]])}}"
  echo "${s%%${2:-+([[:space:]])}}"
}


#==========================================================
# path and file helpers
#==========================================================

# get immediate child directories for specified path(s)
# $1: one or more paths (optional; default=$PWD)
dir_children() {
  echo $(find "${1:-$PWD}" -mindepth 1 -maxdepth 1 -type d | sort)
}
