package api

import (
	"go-kvm/internal/domain"

	evdev "github.com/gvalkov/golang-evdev"
	log "github.com/sirupsen/logrus"
)

type Input struct {
	Devices   []*evdev.InputDevice
	InputChan domain.InputChan
}

func NewInput(inputChan domain.InputChan) *Input {
	return &Input{
		InputChan: inputChan,
	}
}

func (s *Input) Start() {
	s.findDevices()

	if len(s.Devices) == 0 {
		log.Warnf("no input devices found")
		return
	}

	for _, dev := range s.Devices {
		go s.readDevice(dev)
	}
}

func (s *Input) readDevice(dev *evdev.InputDevice) {
	for {
		inputEvents, err := dev.Read()
		if err != nil {
			log.Errorf("readDevice error reading: %v", err)
			continue
		}
		select {
		case s.InputChan <- inputEvents:
		default:
			//log.Warnf("inputchan blocked")
		}
	}
}

func (s *Input) findDevices() {
	devices, _ := evdev.ListInputDevices()

	for _, dev := range devices {

	capabilitiesLoop:
		for cType, eCodes := range dev.Capabilities {

			switch cType.Name {

			case "EV_KEY":
				for i := range eCodes {
					//log.Printf("EV_KEY %d %s\n", eCodes[i].Code, eCodes[i].Name)
					if "KEY_LEFTCTRL" == eCodes[i].Name {
						s.Devices = append(s.Devices, dev)
						break capabilitiesLoop
					}
				}

			case "EV_REL":
				for i := range eCodes {
					//log.Printf("EV_REL %d %s\n", eCodes[i].Code, eCodes[i].Name)
					if "REL_X" == eCodes[i].Name {
						s.Devices = append(s.Devices, dev)
						break capabilitiesLoop
					}
				}
			}
		}
	}
}
