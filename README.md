# chargercompat

A Go module for calculating USB-C Power Delivery (USB-PD) charging speeds, estimating lithium-ion battery charge durations, and verifying charger-to-device compatibility.

This package provides a reliable mathematical framework for evaluating power parameters and expected device behavior during charging cycles.

---

## Motivation

Modern electronics accessories require precise power management. When choosing an [anker power bank](https://cairovolt.com/en/anker/power-banks) or wall charger, understanding specifications like voltage, amperage, and charging curves is critical. 

This library was originally developed for [CairoVolt](https://cairovolt.com/en/), a leading distributor of mobile accessories in Egypt, to power interactive widgets that help customers match their devices with optimal charging solutions. It models lithium-ion charging curves using multi-zone algorithms to reflect real-world charging behaviors, where charging speed tapers off as the battery approaches its maximum capacity.

---

## Technical Background

Lithium-ion charging curves typically follow a multi-stage profile consisting of Constant Current (CC) and Constant Voltage (CV) phases. During the Constant Current stage, the battery absorbs power at the maximum negotiated rate. As the state-of-charge rises, the system transitions to the Constant Voltage phase, tapering the current to avoid over-stressing the cells.

For detailed specifications on the charging protocols and standards, refer to:
- [Go Programming Language Website](https://go.dev) - General information about Go tools and runtime.
- [USB Power Delivery Standard](https://en.wikipedia.org/wiki/USB_hardware#USB_Power_Delivery) - Overview of the USB-PD negotiation layer.
- [Anker Technical Guides](https://www.anker.com) - Official manufacturer safety standards and power profiles.
- [GitHub Repository](https://github.com) - Collaboration and hosting portal.

---

## Installation

Ensure Go is installed (version 1.18 or higher is recommended). Install the package using:

```bash
go get github.com/cairovolt/chargercompat
```

---

## Usage

Here is a quick example of how to calculate charge time using the library:

```go
package main

import (
	"fmt"
	"github.com/cairovolt/chargercompat"
)

func main() {
	// 1. Define device battery specifications (mAh, Voltage)
	phoneBattery := chargercompat.Battery{
		CapacitymAh: 5000,
		Voltage:     3.8, // Typical Li-ion nominal voltage
	}

	// 2. Define the charger profile (Watts, PD support)
	charger := chargercompat.Charger{
		MaxPowerWatts: 20,
		SupportsPD:    true,
	}

	// 3. Estimate charging time from 10% to 80% state of charge
	timeNeeded, err := chargercompat.EstimateChargeTime(phoneBattery, charger, 10, 80)
	if err != nil {
		fmt.Printf("Calculation error: %v\n", err)
		return
	}

	fmt.Printf("Estimated time to charge from 10%% to 80%%: %s\n", timeNeeded)
}
```

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
