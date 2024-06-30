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
                type: "string"
            }
        ],
        commands: [
            {
                name: "list-passwords",
                title: "List Passwords",
                mode: "filter"
            },
						{
								name: "unlock",
								title: "Unlock vault",
								mode: "tty"
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
    echo "Bitwarden session not found. Please set it in your config." >&2
    exit 1
fi

COMMAND=$(echo "$1" | jq -r '.command')
if [ "$COMMAND" = "list-passwords" ]; then
    bw --nointeraction list items --session "$BW_SESSION" | jq 'map({
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
elif [ "$COMMAND" = "unlock" ]; then
	SESSIONTOKEN=$(bw unlock "$PASSWORD" | grep -oP -m 1 '".*?"' | tr -d '"')
		cat $HOME/.config/sunbeam/sunbeam.json | jq --arg SESSION "$SESSIONTOKEN" '.extensions.bitwarden.preferences.session |= $SESSION' | tee $HOME/.config/sunbeam/sunbeam.json

fi
