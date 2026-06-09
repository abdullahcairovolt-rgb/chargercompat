package chargercompat

import (
	"testing"
	"time"
)

func TestWattHours(t *testing.T) {
	battery := Battery{CapacitymAh: 5000, Voltage: 3.8}
	expected := 19.0
	result := battery.WattHours()
	if result != expected {
		t.Errorf("expected %f Wh, got %f Wh", expected, result)
	}
}

func TestEstimateChargeTime(t *testing.T) {
	battery := Battery{CapacitymAh: 5000, Voltage: 3.8} // 19 Wh
	charger := Charger{MaxPowerWatts: 20, SupportsPD: true}

	// Calculate estimate from 0% to 100%
	duration, err := EstimateChargeTime(battery, charger, 0, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Expected calculation breakdown:
	// Wh total = 19 Wh. Efficiency = 0.8
	// Zone 1 (0-50%): 9.5 Wh needed. Power = 20 * 1.0 * 0.8 = 16W. Hours = 9.5 / 16 = 0.59375 hrs (35.625 mins)
	// Zone 2 (50-80%): 5.7 Wh needed. Power = 20 * 0.6 * 0.8 = 9.6W. Hours = 5.7 / 9.6 = 0.59375 hrs (35.625 mins)
	// Zone 3 (80-100%): 3.8 Wh needed. Power = 20 * 0.25 * 0.8 = 4W. Hours = 3.8 / 4.0 = 0.95 hrs (57 mins)
	// Total minutes = 35.625 + 35.625 + 57 = 128.25 mins. Round to 128 mins.
	expected := 128 * time.Minute
	if duration != expected {
		t.Errorf("expected %s duration, got %s", expected, duration)
	}
}

func TestIsCompatible(t *testing.T) {
	charger := Charger{MaxPowerWatts: 15, SupportsPD: false}
	compatible, fastCharge := IsCompatible(20, charger)
	if !compatible {
		t.Errorf("expected compatible to be true")
	}
	if fastCharge {
		t.Errorf("expected fastCharge to be false since PD is not supported and wattage is low")
	}

	chargerPD := Charger{MaxPowerWatts: 30, SupportsPD: true}
	compatible, fastCharge = IsCompatible(20, chargerPD)
	if !compatible {
		t.Errorf("expected compatible to be true")
	}
	if !fastCharge {
		t.Errorf("expected fastCharge to be true")
	}
}
