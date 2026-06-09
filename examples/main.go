package main

import (
	"fmt"
	"github.com/cairovolt/chargercompat"
)

func main() {
	// Specifications of a typical smartphone battery (e.g., 5000 mAh capacity at 3.8V nominal voltage)
	phoneBattery := chargercompat.Battery{
		CapacitymAh: 5000,
		Voltage:     3.8,
	}

	// Specifications of a standard 20W USB-C charger supporting Power Delivery
	charger := chargercompat.Charger{
		MaxPowerWatts: 20,
		SupportsPD:    true,
	}

	fmt.Printf("Phone Battery Energy: %.2f Wh\n", phoneBattery.WattHours())

	// Estimate charge time from 20% to 80% (healthy daily charging band)
	timeNeeded, err := chargercompat.EstimateChargeTime(phoneBattery, charger, 20, 80)
	if err != nil {
		fmt.Printf("Error calculating charging time: %v\n", err)
		return
	}

	fmt.Printf("Estimated time to charge from 20%% to 80%%: %s\n", timeNeeded)

	// Check compatibility for a 25W device requirement
	compatible, fastCharge := chargercompat.IsCompatible(25, charger)
	fmt.Printf("Charger Compatibility status: Compatible: %t, Fast Charging: %t\n", compatible, fastCharge)
}
