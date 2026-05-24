# Sinkor Vivatron TCP Controller + HomeAssistant

Simple Go utility for controlling compatible RGBW LED controllers over TCP.

The application connects directly to the device on port `5577` and sends raw binary commands for:

- turning the light on
- turning the light off
- setting RGBW color values

---

# Requirements

- Go 1.18+ (recommended)
- Network access to the LED controller

---

# Compilation

## Linux AMD64

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o surplifeha vivatron.go
```

## Raspberry Pi / ARM64

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o surplifeha vivatron.go
```

## Windows

```bash
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o surplifeha.exe vivatron.go
```

---

# Build Options Explained

| Variable | Description |
|---|---|
| `CGO_ENABLED=0` | Builds a static binary without C dependencies |
| `GOOS=linux` | Target operating system |
| `GOARCH=amd64` | Target CPU architecture |
| `-o surplifeha` | Output binary name |

---

# Usage

```bash
./surplifeha <ip> <action> [arguments]
```

---

# Actions

## Turn ON

```bash
./surplifeha 192.168.1.50 on
```

## Turn OFF

```bash
./surplifeha 192.168.1.50 off
```

## Set RGBW Color

```bash
./surplifeha 192.168.1.50 color 255 0 0 0
```

Example above sets:

- Red = 255
- Green = 0
- Blue = 0
- White = 0

Result: pure red color.

---

# RGBW Value Range

Each color channel accepts values:

```text
0-255
```

Where:

- `0` = channel disabled
- `255` = maximum brightness

---

# Example Colors

## White

```bash
./surplifeha 192.168.1.50 color 0 0 0 255
```

## Blue

```bash
./surplifeha 192.168.1.50 color 0 0 255 0
```

## Warm White Mix

```bash
./surplifeha 192.168.1.50 color 255 180 100 80
```

---

# Protocol Details

- TCP Port: `5577`
- Binary packet format
- Includes packet checksum
- 3-second TCP timeout

---

# Notes

The utility is intended for compatible WiFi RGBW LED controllers using the TCP protocol on port `5577`.

No authentication or encryption is used.

# Home Assistant

## Overview

This project allows you to control a Surplife RGBW light from Home Assistant.

Features:
- control of Red, Green, Blue, White channels (RGBW)
- power on/off switch
- slider control using input numbers
- automatic updates when values change
- integration with a local control script

---

## Requirements

- Home Assistant installed (HA OS or Docker)
- access to `/config` folder
- executable script:
  /config/scripts/vivatron
- Surplife device reachable in local network (example IP: 192.168.1.217)

---

## Installation

### 1. Shell Command (RGBW control)

Add to configuration.yaml:

shell_command:
  surplife_set_rgbw: >
    /config/scripts/surplifeha 192.168.1.217 color
    {{ states('input_number.surplife_r') | int }}
    {{ states('input_number.surplife_g') | int }}
    {{ states('input_number.surplife_b') | int }}
    {{ states('input_number.surplife_w') | int }}

---

### 2. Power Switch

switch:
  - platform: command_line
    switches:
      surplife_power:
        friendly_name: "Surplife Power"
        unique_id: surplife_power
        command_on: "/config/scripts/surplifeha 192.168.1.217 on"
        command_off: "/config/scripts/surplifeha 192.168.1.217 off"

---

### 3. RGBW Sliders

input_number:
  surplife_r:
    name: "Red"
    min: 0
    max: 255
    step: 1

  surplife_g:
    name: "Green"
    min: 0
    max: 255
    step: 1

  surplife_b:
    name: "Blue"
    min: 0
    max: 255
    step: 1

  surplife_w:
    name: "White"
    min: 0
    max: 255
    step: 1

---

### 4. Automation

Add to automations.yaml or via UI:

alias: Update Surplife RGBW
trigger:
  - platform: state
    entity_id:
      - input_number.surplife_r
      - input_number.surplife_g
      - input_number.surplife_b
      - input_number.surplife_w

action:
  - service: shell_command.surplife_set_rgbw

---

### 5. Dashboard (Lovelace)

type: entities
title: Surplife Light
entities:
  - entity: switch.surplife_power
    name: Power

  - entity: input_number.surplife_r
  - entity: input_number.surplife_g
  - entity: input_number.surplife_b
  - entity: input_number.surplife_w

---

## Setup Steps

1. Copy configuration into configuration.yaml
2. Restart Home Assistant
3. Add automation
4. Add dashboard card
5. Ensure script is executable:

chmod +x /config/scripts/surplifeha

---

## How it works

- Changing sliders updates RGBW values
- Automation triggers shell command
- Script sends values to Surplife device
- Power switch sends on/off commands directly

---

## Troubleshooting

- Check Home Assistant logs if automation does not trigger
- Verify device IP address (192.168.1.217)
- Test script manually from terminal
- Ensure script has execute permissions




