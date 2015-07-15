package distance
// this should be in a seperate "units" package
import (
	"fmt"
)

type Distance float64
type ImperialDistance float64
type NauticalDistance float64
type USSurveyDistance float64

const (
	Meter      Distance = 1.0          // m
	Kilometer           = 1000 * Meter // km
	Decimeter           = Meter / 10   // dm
	Centimeter          = Meter / 100  // cm
	Millimeter          = Meter / 1000 // mm
	// EarthRadius               = 6378.1 * Kilometer
	EarthRadius = 6371.0 * Kilometer

	Yard   = 0.9144 * ImperialDistance(Meter) // yd
	Foot   = Yard / 3                         // ft
	Inch   = Foot / 12                        // in
	Pica   = Inch / 6                         // /p
	PPoint = Pica / 12                        // p
	Mile   = 1760 * Yard                      // mi

	NauticalMile = 1852 * NauticalDistance(Meter) // nmi
	Fathom       = 2 * NauticalDistance(Yard)     // ftm
	Cable        = 120 * Fathom                   // cb

	USInch = USSurveyDistance(Inch)
	USFoot = 1200.0 / 3937 * USSurveyDistance(Meter)
	USYard = USSurveyDistance(Yard)
	USMile = 8 /* fur */ * 10 /* ch */ * 66 * USFoot
)

func (d Distance) String() string {
	if d >= Kilometer {
		return fmt.Sprintf("%f", float64(d/Kilometer)) + "km"
	}
	if d <= Millimeter {
		return fmt.Sprintf("%f", float64(d/Millimeter)) + "mm"
	}
	if d <= Meter {
		return fmt.Sprintf("%f", float64(d/Centimeter)) + "cm"
	}
	return fmt.Sprintf("%f", float64(d/Meter)) + "m"
}

func (d NauticalDistance) String() string {
	return fmt.Sprintf("%f", float64(d/NauticalMile)) + "nmi"
}

func (d USSurveyDistance) String() string {
	if d >= USMile {
		return fmt.Sprintf("%f", float64(d/USMile)) + "mi"
	}
	if d >= USYard {
		return fmt.Sprintf("%f", float64(d/USYard)) + "yd"
	}
	if d > USFoot {
		return fmt.Sprintf("%f", float64(d/USFoot)) + "ft"
	}
	return fmt.Sprintf("%f", float64(d/USInch)) + "in"
}

func (d ImperialDistance) String() string {
	if d >= Mile {
		return fmt.Sprintf("%f", float64(d/Mile)) + "mi"
	}
	if d >= Yard {
		return fmt.Sprintf("%f", float64(d/Yard)) + "yd"
	}
	if d >= Foot {
		return fmt.Sprintf("%f", float64(d/Foot)) + "ft"
	}
	return fmt.Sprintf("%f", float64(d/Inch)) + "in"
}

// vim: ts=4 sw=4 noexpandtab nolist syn=go
