#!/bin/ash -efu

if [ "$GODEP_IMPORT_PATH" = "github.com/Sirupsen/logrus" ]; then
	exit $EXIT_CHOOSE_OLD
fi

[ -n "$GODEP_NEW_REV" ] &&
	exit $EXIT_CHOOSE_NEW ||
	exit $EXIT_CHOOSE_OLD
