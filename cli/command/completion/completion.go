/* Copied from https://raw.githubusercontent.com/kubernetes/kubernetes/master/pkg/kubectl/cmd/completion.go */
package completion

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

const defaultBoilerPlate = `
# Copyright 2017 The AMP Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
`

var (
	completion_long = `
		Output shell completion code for the specified shell (bash or zsh).
		The shell code must be evalutated to provide interactive
		completion of amp commands.  This can be done by sourcing it from
		the .bash_profile.

		Note: this requires the bash-completion framework, which is not installed
		by default on Mac.  This can be installed by using homebrew:

		    $ brew install bash-completion

		Once installed, bash_completion must be evaluated.  This can be done by adding the
		following line to the .bash_profile

		    $ source $(brew --prefix)/etc/bash_completion

		Note for zsh users: [1] zsh completions are only supported in versions of zsh >= 5.2`

	completion_example = `
		# Install bash completion on a Mac using homebrew

		    $ brew install bash-completion

		Once installed, bash_completion must be evaluated.  This can be done by adding the
		following line to the .bash_profile, and reloading it

		    source $(brew --prefix)/etc/bash_completion

		    $ source $HOME/.bash_profile

		# Write bash completion code to a file

		    $ amp completion bash > ~/.config/amp/completion.bash.inc

		Add this line in your $HOME/.bash_profile, and reload it

		    source "$HOME/.config/amp/completion.bash.inc"

		    $ source $HOME/.bash_profile

		# Load the amp completion code for zsh[1] into the current shell
		source <(amp completion zsh)`
)

var (
	completion_shells = map[string]func(out io.Writer, cmd *cobra.Command) error{
		"bash": runCompletionBash,
		"zsh":  runCompletionZsh,
	}
)

func NewCompletionCommand(c cli.Interface) *cobra.Command {
	shells := []string{}
	boilerPlate := ""
	for s := range completion_shells {
		shells = append(shells, s)
	}

	cmd := &cobra.Command{
		Use:     "completion SHELL",
		Short:   "Output shell completion code for the specified shell (bash or zsh)",
		Long:    completion_long,
		Example: completion_example,
		Run: func(cmd *cobra.Command, args []string) {
			err := RunCompletion(os.Stdout, boilerPlate, cmd, args)
			if err != nil {
				fmt.Println(err)
			}
		},
		ValidArgs: shells,
	}

	return cmd
}

func RunCompletion(out io.Writer, boilerPlate string, cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("shell not specified")
	}
	if len(args) > 1 {
		return errors.New("too many arguments. Expected only the shell type")
	}
	run, found := completion_shells[args[0]]
	if !found {
		return errors.New(fmt.Sprintf("unsupported shell type %q", args[0]))
	}

	if len(boilerPlate) == 0 {
		boilerPlate = defaultBoilerPlate
	}
	if _, err := out.Write([]byte(boilerPlate)); err != nil {
		return err
	}
	return run(out, cmd.Parent())
}

func runCompletionBash(out io.Writer, c *cobra.Command) error {
	return c.GenBashCompletion(out)
}

func runCompletionZsh(out io.Writer, c *cobra.Command) error {
	zsh_initialization := `
__amp_bash_source() {
	alias shopt=':'
	alias _expand=_bash_expand
	alias _complete=_bash_comp
	emulate -L sh
	setopt kshglob noshglob braceexpand

	source "$@"
}

__amp_type() {
	# -t is not supported by zsh
	if [ "$1" == "-t" ]; then
		shift

		# fake Bash 4 to disable "complete -o nospace". Instead
		# "compopt +-o nospace" is used in the code to toggle trailing
		# spaces. We don't support that, but leave trailing spaces on
		# all the time
		if [ "$1" = "__amp_compopt" ]; then
			echo builtin
			return 0
		fi
	fi
	type "$@"
}

__amp_compgen() {
	local completions w
	completions=( $(compgen "$@") ) || return $?

	# filter by given word as prefix
	while [[ "$1" = -* && "$1" != -- ]]; do
		shift
		shift
	done
	if [[ "$1" == -- ]]; then
		shift
	fi
	for w in "${completions[@]}"; do
		if [[ "${w}" = "$1"* ]]; then
			echo "${w}"
		fi
	done
}

__amp_compopt() {
	true # don't do anything. Not supported by bashcompinit in zsh
}

__amp_declare() {
	if [ "$1" == "-F" ]; then
		whence -w "$@"
	else
		builtin declare "$@"
	fi
}

__amp_ltrim_colon_completions()
{
	if [[ "$1" == *:* && "$COMP_WORDBREAKS" == *:* ]]; then
		# Remove colon-word prefix from COMPREPLY items
		local colon_word=${1%${1##*:}}
		local i=${#COMPREPLY[*]}
		while [[ $((--i)) -ge 0 ]]; do
			COMPREPLY[$i]=${COMPREPLY[$i]#"$colon_word"}
		done
	fi
}

__amp_get_comp_words_by_ref() {
	cur="${COMP_WORDS[COMP_CWORD]}"
	prev="${COMP_WORDS[${COMP_CWORD}-1]}"
	words=("${COMP_WORDS[@]}")
	cword=("${COMP_CWORD[@]}")
}

__amp_filedir() {
	local RET OLD_IFS w qw

	__debug "_filedir $@ cur=$cur"
	if [[ "$1" = \~* ]]; then
		# somehow does not work. Maybe, zsh does not call this at all
		eval echo "$1"
		return 0
	fi

	OLD_IFS="$IFS"
	IFS=$'\n'
	if [ "$1" = "-d" ]; then
		shift
		RET=( $(compgen -d) )
	else
		RET=( $(compgen -f) )
	fi
	IFS="$OLD_IFS"

	IFS="," __debug "RET=${RET[@]} len=${#RET[@]}"

	for w in ${RET[@]}; do
		if [[ ! "${w}" = "${cur}"* ]]; then
			continue
		fi
		if eval "[[ \"\${w}\" = *.$1 || -d \"\${w}\" ]]"; then
			qw="$(__amp_quote "${w}")"
			if [ -d "${w}" ]; then
				COMPREPLY+=("${qw}/")
			else
				COMPREPLY+=("${qw}")
			fi
		fi
	done
}

__amp_quote() {
    if [[ $1 == \'* || $1 == \"* ]]; then
        # Leave out first character
        printf %q "${1:1}"
    else
        printf %q "$1"
    fi
}

autoload -U +X bashcompinit && bashcompinit

# use word boundary patterns for BSD or GNU sed
LWORD='[[:<:]]'
RWORD='[[:>:]]'
if sed --help 2>&1 | grep -q GNU; then
	LWORD='\<'
	RWORD='\>'
fi

__amp_convert_bash_to_zsh() {
	sed \
	-e 's/declare -F/whence -w/' \
	-e 's/_get_comp_words_by_ref "\$@"/_get_comp_words_by_ref "\$*"/' \
	-e 's/local \([a-zA-Z0-9_]*\)=/local \1; \1=/' \
	-e 's/flags+=("\(--.*\)=")/flags+=("\1"); two_word_flags+=("\1")/' \
	-e 's/must_have_one_flag+=("\(--.*\)=")/must_have_one_flag+=("\1")/' \
	-e "s/${LWORD}_filedir${RWORD}/__amp_filedir/g" \
	-e "s/${LWORD}_get_comp_words_by_ref${RWORD}/__amp_get_comp_words_by_ref/g" \
	-e "s/${LWORD}__ltrim_colon_completions${RWORD}/__amp_ltrim_colon_completions/g" \
	-e "s/${LWORD}compgen${RWORD}/__amp_compgen/g" \
	-e "s/${LWORD}compopt${RWORD}/__amp_compopt/g" \
	-e "s/${LWORD}declare${RWORD}/__amp_declare/g" \
	-e "s/\\\$(type${RWORD}/\$(__amp_type/g" \
	<<'BASH_COMPLETION_EOF'
`
	out.Write([]byte(zsh_initialization))

	buf := new(bytes.Buffer)
	c.GenBashCompletion(buf)
	out.Write(buf.Bytes())

	zsh_tail := `
BASH_COMPLETION_EOF
}

__amp_bash_source <(__amp_convert_bash_to_zsh)
`
	out.Write([]byte(zsh_tail))
	return nil
}
