package main

import "../drivers"

func main() {
	drivers.ElevInit()
	drivers.ElevSetFloorIndicator(2)
	for {}
}
