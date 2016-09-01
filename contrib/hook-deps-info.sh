#!/bin/ash -efu

printf 'import-path: %s ' "$GODEP_IMPORT_PATH"
if [ -d "vendor/$GODEP_IMPORT_PATH" ]; then
	printf '(%s)\n' "vendor/$GODEP_IMPORT_PATH"
elif [ -d "$GOPATH/src/$GODEP_IMPORT_PATH" ]; then
	printf '(%s)\n' "$GOPATH/src/$GODEP_IMPORT_PATH"
else
	printf '(not found)\n'
fi
