# Streamdeckd Modules

This repository contains a collection of plugins I have made for my local streamdeckd setup, they are published as is with no intention for support, as examples to aid creating plugins.

## Overview

The modules in this repository provide various functionalities for the Stream Deck:

- **Toggle**: Create toggle buttons that can run different commands based on state
- **Lights**: Control home automation systems (like Home Assistant)
- **CCTV**: Display camera feeds on Stream Deck buttons
- **NoOp**: A simple "no operation" placeholder button
- **Volume**: Control PulseAudio volume levels for system audio devices
- **PlayerCtlVolume**: Control volume for media players via MPRIS

## Modules

### Toggle

The Toggle module provides a button that can toggle between two states (up/down). It runs a check command to determine its current state and displays different icons accordingly. When pressed, it executes either an "up command" or "down command" depending on the current state.

**Configuration Fields:**
- Up Icon: Image to display when in "up" state
- Down Icon: Image to display when in "down" state
- Check Command: Shell command to determine the current state
- Up Command: Command to execute when toggling to "up" state
- Down Command: Command to execute when toggling to "down" state

### Lights

The Lights module allows for controlling lights in home-assistant. It sends HTTP requests to control smart lights or other entities. This plugin exposes no icon handler, but could be combined with a toggle icon handler to display a different image based on the state of the light

**Configuration Fields:**
- Domain: The domain of the entity (e.g., "light", "switch")
- Service: The service to call (e.g., "toggle", "turn_on")
- Entity ID: The ID of the entity to control
- API Key: Authentication token for the home automation system
- Base URL: The base URL of the home automation system

### CCTV

The CCTV module fetches images from a URL (likely a security camera feed) and displays them on a Stream Deck button. It continuously updates the image at regular intervals.

**Configuration Fields:**
- URL: The URL of the camera feed to display

### NoOp

The NoOp module is a simple "no operation" module that creates a blank button. It doesn't perform any action when pressed and doesn't have any configurable fields.

### Volume

The Volume module provides integration with PulseAudio to control audio devices. It can display volume levels and mute status on Stream Deck LCD screens, and allows controlling volume via buttons, or the knobs & touch screen on the StreamDeck+.

**Configuration Fields:**
- Device Type: Type of audio device to control (sink, source, sink_input, source_output)
- Input Name: Name of the specific audio input/output to control
- Props: Properties to identify the audio device
- Unmuted Icon: Image to display when not muted
- Muted Icon: Image to display when muted

### PlayerCtlVolume

The PlayerCtlVolume module provides integration with media players via MPRIS. It allows controlling the volume of media players and displaying their volume level on Stream Deck LCD screens.

**Configuration Fields:**
- Player Name: Name of the media player to control (optional, controls active player if not specified)
- Icon: Image to display on the button