# wacomctl

A command-line utility for controlling Wacom tablets on Linux.

## Description

The `wacomctl` tool is a simple Go application that can map and control the Wacom stylus to different monitors or disable/enable it as needed.

## Prerequisites

- Linux system with X11
- Connected Wacom tablet
- `xsetwacom` and `xinput` packages installed
- Go 1.x (for building)

## Installation

1. Clone or download the source code
2. Build the executable:

   ```bash
   go build -o wacomctl wacomctl.go
   ```

3. Make the executable available in your PATH (optional):

   ```bash
   sudo mv wacomctl /usr/local/bin/
   ```

## Usage

```bash
wacomctl
```

Without arguments, the program opens an interactive selector with the detected monitors and additional options to map to all monitors, turn off, or turn on the stylus.

### Commands

- **`interactive`** - Opens the interactive selector to choose between detected monitors, both, turn off, and turn on
- **`map <monitor>`** - Maps the stylus to the provided monitor (for example, `HDMI-1`)
- **`both`** - Maps the stylus to all active monitors (full desktop)
- **`off`** - Turns the stylus off (disables the device)
- **`on`** - Turns the stylus on (enables the device)

### Examples

```bash
# Interactive selector
wacomctl

# Map to a specific monitor
wacomctl map HDMI-1

# Map to all monitors
wacomctl both

# Turn the stylus off
wacomctl off

# Turn the stylus on
wacomctl on
```

## How it works

The program uses the following system tools:

- **`xsetwacom`** - To detect and configure Wacom devices
- **`xrandr`** - To list connected monitors
- **`xinput`** - To enable/disable the device

### Automatic detection

- Automatically detects the stylus device ID using `xsetwacom --list devices`
- Identifies connected VGA and HDMI monitors using `xrandr --listmonitors`
- Maps the active stylus area to the specified monitor

## Limitations

- Works only on X11 systems (not Wayland)
- Requires monitors to be connected and active
- Detects only the first stylus device found
- Supports only VGA and HDMI monitors (not DisplayPort, DVI, etc.)

## Troubleshooting

### Stylus device not found

```text
Stylus device not found.
```

- Check that the Wacom tablet is connected
- Run `xsetwacom --list devices` to verify the device is recognized

### Monitor not found

```text
No monitor VGA/HDMI found.
```

- Check that the monitor is connected and active
- Run `xrandr --listmonitors` to see the available monitors

### Mapping error

```text
Failed to map stylus: ...
```

- Check that you have the appropriate permissions
- Make sure `xsetwacom` is installed

## License

This project is in the public domain. Use it freely.
