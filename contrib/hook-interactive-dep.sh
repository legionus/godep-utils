#!/bin/bash -efu

PROG="${0##*/}"
STATE="/tmp/godep-merge/state"

state() {
	mkdir -p -- "${STATE%/*}"
	echo "$1" > "$STATE"
}

usage() {
	printf '\nCommands: n/o/i = new/old/ignore path; N/O/I = new/old/ignore rest all paths; a = abort\n'
}

handle() {
	local cmd="${1-}"
	case "$cmd" in
		n|N|new)
			[ "$cmd" != 'N' ] || state n
			printf ' - %s\n' "use${GODEP_NEW_COMMENT:+ $GODEP_NEW_COMMENT} $GODEP_NEW_REV"
			exit $EXIT_CHOOSE_NEW
			;;
		o|O|old)
			[ "$cmd" != 'O' ] || state o
			printf ' - %s\n' "use${GODEP_OLD_COMMENT:+ $GODEP_OLD_COMMENT} $GODEP_OLD_REV"
			exit $EXIT_CHOOSE_OLD
			;;
		i|I|ignore)
			[ "$cmd" != 'I' ] || state i
			printf ' - %s\n' "ignore $GODEP_IMPORT_PATH"
			exit $EXIT_CHOOSE_IGN
			;;
		a|abort)
			printf ' - abort merge\n'
			exit 1
			;;
		*)
			usage
			;;
	esac
}

if [ -f "$STATE" ]; then
	read cmd < "$STATE"
	handle "$cmd"
fi

usage
while read -p "hook($GODEP_IMPORT_PATH)\$ " cmd; do
	handle "$cmd"
done
echo
