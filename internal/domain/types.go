package domain

import evdev "github.com/gvalkov/golang-evdev"

type InputChan chan []evdev.InputEvent
