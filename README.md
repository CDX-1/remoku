# Remoku

A CLI tool for controlling Roku devices over a local network.

[Here is a video demo.](https://www.youtube.com/watch?v=Pt1i_EhmWlQ&feature=youtu.be)

## Features

- Scan your network for Roku devices using the SSDP protocol
- Simulate remote controller keypresses (ex: down, volume up, home, etc)
- Show a list of all apps installed on the Roku device and launch any of them
- Create and execute macros (JSON files with a list of commands)
- Interactive mode with keyboard controls

## Requirements

* Go v1.25.6
> Older/newer versions may work but Remoku has only been tested with this version.

## Installation

1. You must have [Go](https://go.dev/dl/) installed on your computer first.
2. Run `go install github.com/CDX-1/remoku/cmd/remoku@latest`
3. Run `remoku` to use the tool

#### On Your Roku Device

1. Ensure that in `Settings > System > Advanced system settings` that 'Control by mobile apps' is set to 'Enabled'
2. Make sure that your computer and the Roku device are on the same network
3. You can check your Rokus IP address by going to `Settings > Network > About > IP Address`

## Usage

### Sending Keypresses

**Valid inputs include:**
* home
* back
* up
* down
* left
* right
* ok (or 'select')
* play
* pause
* vup
* vdown
* mute
* power (or 'off')
* on

```bash
remoku press up     # move selector up
remoku press home   # open home screen
remoku press vup    # increase volume
```

### Using Apps

```bash
remoku apps             # get a list of installed apps
remoku launch <appId>   # launch an app
```

### Using Macros

Want to create your own macros? See [MACROS.md](MACROS.md) for more information.

```bash
remoku macro macros/example.json   # execute commands in example.json
```

### Interactive Mode

```bash
remoku interactive  # launches interactive mode
```

### Specifying a device manually

```bash
remoku press home --ip <ip_address>  # opens home screen on the device at <ip_address>
```
