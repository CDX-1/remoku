# Macros

Macros are a powerful tool for creating your own automations to control your TV

## Creating Macros

Macros use the JSON format and are made up of an array of commands.

For example, the following macro will open the home screen, wait 2 seconds, and then open the YouTube app:

```json
[
    {
        "type": "keypress",
        "value": "home"
    },
    {
        "type": "sleep",
        "duration": 2
    },
    {
        "type": "launch",
        "value": "<youtubeAppId>"
    }
]
```

## Using Macros

To use a macro, run the following command:

```bash
remoku macro <path_to_macro.json>
```

## Valid Inputs

The following inputs are valid for macros:

* Any of the keypress inputs listed in the [Sending Keypresses](README.md#sending-keypresses) section
* `sleep`   - wait for the specified duration
* `keypress`   - press a key on the remote
* `launch`  - launch an app