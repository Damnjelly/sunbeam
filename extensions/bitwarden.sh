#!/bin/sh

# check if jq is installed
if ! [ -x "$(command -v jq)" ]; then
    echo "jq is not installed. Please install it." >&2
    exit 1
fi

if [ $# -eq 0 ]; then
    jq -n '{
        title: "Bitwarden Vault",
        description: "Search your Bitwarden passwords",
        preferences: [
            {
                name: "session",
                title: "Bitwarden Session",
                type: "string",
								optional: true
            },
						{
							name: "sessionPath",
							title: "Path to Bitwarden Session",
							type: "string",
							optional: true
						}
        ],
        commands: [
            {
                name: "list-passwords",
                title: "List Passwords",
                mode: "filter"
            }
        ]
    }'
    exit 0
fi

# check if bkt is installed
if ! [ -x "$(command -v bkt)" ]; then
    echo "bkt is not installed. Please install it." >&2
    exit 1
fi

BW_SESSION=$(echo "$1" | jq -r '.preferences.session')
if [ "$BW_SESSION" = "null" ]; then
		BW_SESSION=$(echo "$1" | jq -r '.preferences.sessionPath')
		if [ "$BW_SESSION" = "null" ]; then
				echo "Session token not set. Please set it in the sunbeam config file." >&2
				exit 1
		fi
		BW_SESSION=$(cat "$BW_SESSION")
fi

COMMAND=$(echo "$1" | jq -r '.command')
if [ "$COMMAND" = "list-passwords" ]; then
    bkt --ttl=1d -- bw --nointeraction list items --session "$BW_SESSION" | jq 'map({
        title: .name,
        subtitle: (.login.username // ""),
        actions: [
            {
                title: "Copy Password",
                type: "copy",
                text: (.login.password // ""),
                exit: true
            },
            {
                title: "Copy Username",
                key: "l",
                type: "copy",
                text: (.login.username // ""),
                exit: true
            }
        ]
    }) | { items: .}'
fi
