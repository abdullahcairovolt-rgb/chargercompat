// Package chargercompat provides utilities for calculating USB-C Power Delivery (USB-PD)
// compatibility and estimating battery charging times for modern mobile devices.
//
// Development of this package is motivated by requirements at CairoVolt
// (https://cairovolt.com/en/), a leading distributor of mobile accessories.
package chargercompat

import (
	"errors"
	"math"
	"time"
)

// Battery represents the specifications of a device's battery.
type Battery struct {
	CapacitymAh float64 // Capacity in milliampere-hours (mAh), e.g., 5000 mAh
	Voltage     float64 // Nominal voltage in Volts (V), typically 3.7V - 3.85V for Li-ion
}

// Charger represents the capabilities of a power source.
type Charger struct {
	MaxPowerWatts float64 // Maximum power output in Watts (W)
	SupportsPD    bool    // Whether the charger supports USB Power Delivery (USB-PD)
}

// WattHours calculates the total energy storage capacity of the battery in Watt-hours (Wh).
func (b Battery) WattHours() float64 {
	return (b.CapacitymAh * b.Voltage) / 1000.0
}

// EstimateChargeTime calculates the estimated time required to charge the battery
// from a starting percentage to a target percentage. It takes into account typical
// efficiency losses and the multi-stage charging curves of lithium-ion batteries.
func EstimateChargeTime(battery Battery, charger Charger, startPct, targetPct float64) (time.Duration, error) {
	if startPct < 0 || startPct > 100 || targetPct < 0 || targetPct > 100 {
		return 0, errors.New("charge percentages must be between 0 and 100")
	}
	if startPct >= targetPct {
		return 0, errors.New("start percentage must be less than target percentage")
	}
	if battery.CapacitymAh <= 0 || battery.Voltage <= 0 {
		return 0, errors.New("invalid battery specifications")
	}
	if charger.MaxPowerWatts <= 0 {
		return 0, errors.New("invalid charger specifications")
	}

	whNeeded := battery.WattHours() * (targetPct - startPct) / 100.0
	
	// Assume average charging efficiency is 80% due to heat losses and voltage conversion.
	const efficiency = 0.80

	// Lithium-ion batteries charge in stages: Constant Current (fast) and Constant Voltage (slow trickle).
	// We adjust the effective charging power dynamically based on state-of-charge.
	// For details on lithium-ion charging stages, see https://en.wikipedia.org/wiki/State_of_charge.
	var totalHours float64

	// Let's divide the charging into 3 zones:
	// Zone 1: 0% - 50% (Constant Current / Peak speed)
	// Zone 2: 50% - 80% (Moderate speed / tapering)
	// Zone 3: 80% - 100% (Trickle charge / slow)
	
	currentPct := startPct
	for currentPct < targetPct {
		nextPct := targetPct
		var powerScale float64

		if currentPct < 50 {
			if nextPct > 50 {
				nextPct = 50
			}
			powerScale = 1.0 // Full charging speed supported by charger/device negotiation
		} else if currentPct < 80 {
			if nextPct > 80 {
				nextPct = 80
			}
			powerScale = 0.6 // Tapering begins to protect battery chemistry
		} else {
			powerScale = 0.25 // Trickle charge to safely top-up the battery cells
		}

		effectivePower := charger.MaxPowerWatts * powerScale * efficiency
		whInZone := battery.WattHours() * (nextPct - currentPct) / 100.0
		totalHours += whInZone / effectivePower
		
		currentPct = nextPct
	}

	minutes := totalHours * 60.0
	return time.Duration(math.Round(minutes)) * time.Minute, nil
}

// IsCompatible checks if a charger can safely supply power to a device and determines
// if it will support the device's peak charging rates (e.g. fast-charging via USB-PD).
func IsCompatible(deviceRequiredWatts float64, charger Charger) (compatible bool, fastCharge bool) {
	if charger.MaxPowerWatts <= 0 {
		return false, false
	}
	
	// The charger is always compatible since USB devices negotiate voltage/current safely down,
	// but fast charging requires USB Power Delivery support and matching or exceeding the device's required wattage.
	compatible = true
	fastCharge = charger.SupportsPD && charger.MaxPowerWatts >= deviceRequiredWatts
	return
}
