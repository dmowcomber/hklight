package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/eliukblau/pixterm/ansimage"
	colorful "github.com/lucasb-eyer/go-colorful"
	rpio "github.com/stianeikeland/go-rpio"
)

const (
	accessoryPin = "32191123"
	lightGPIOPin = 17
)

var (
	hue float64
	sat float64
)

func turnLightOn() {
	if os.Getenv("GARCH") == "arm" {
		pin := rpio.Pin(lightGPIOPin)
		pin.Output() // Output mode
		pin.High()   // Set pin High
	}

	log.Println("Turn Light On")
}

func turnLightOff() {
	if os.Getenv("GARCH") == "arm" {
		pin := rpio.Pin(lightGPIOPin)
		pin.Output() // Output mode
		pin.Low()    // Set pin Low
	}

	log.Println("Turn Light Off")
}

func drawColor() {
	cf := colorful.Hsv(hue, sat/100, 1.0)
	fmt.Printf("RGB values: %v, %v, %v\n", cf.R, cf.G, cf.B)

	m := image.NewRGBA(image.Rect(0, 0, 100, 100))
	c := color.RGBA{uint8(cf.R*256) - 1, uint8(cf.G*256) - 1, uint8(cf.B*256) - 1, 255}
	draw.Draw(m, m.Bounds(), &image.Uniform{c}, image.ZP, draw.Src)

}

func main() {
	if os.Getenv("GARCH") == "arm" {
		// initialize rpio
		err := rpio.Open()
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	ansimage.ClearTerminal()
	drawColor()

	info := accessory.Info{
		Name:         "Light",
		Manufacturer: "dmowcomber",
	}

	acc := accessory.NewLightbulb(info)

	acc.Lightbulb.On.OnValueRemoteUpdate(func(on bool) {
		if on == true {
			turnLightOn()
		} else {
			turnLightOff()
		}
	})
	acc.Lightbulb.Brightness.OnValueRemoteUpdate(func(i int) {
		fmt.Println(i)
	})
	acc.Lightbulb.Saturation.OnValueRemoteUpdate(func(s float64) {
		sat = s
		drawColor()
	})
	acc.Lightbulb.Hue.OnValueRemoteUpdate(func(h float64) {
		hue = h
		drawColor()
	})

	t, err := hc.NewIPTransport(hc.Config{Pin: accessoryPin}, acc.Accessory)
	if err != nil {
		log.Fatal(err)
	}

	hc.OnTermination(func() {
		t.Stop()
	})
	fmt.Printf("accessory pin: %s\n", accessoryPin)
	t.Start()
}
