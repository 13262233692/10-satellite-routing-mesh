package model

import (
	"math"
	"sync"
)

const MaxSatellites = 5000

const MaxLaserLinksPerSatellite = 4

const MaxLinkDistanceMeters = 5000000

type SatelliteID uint32

type Vec3 struct {
	X float64
	Y float64
	Z float64
}

type Ephemeris struct {
	ID        SatelliteID
	Position  Vec3
	Velocity  Vec3
	Timestamp int64
}

type Satellite struct {
	ID           SatelliteID
	Name         string
	OrbitalPlane uint16
	IndexInPlane uint16
	CurrentEphem Ephemeris
	Lasers       uint8
}

var ephemerisPool = sync.Pool{
	New: func() interface{} {
		return &Ephemeris{}
	},
}

func GetEphemeris() *Ephemeris {
	return ephemerisPool.Get().(*Ephemeris)
}

func PutEphemeris(e *Ephemeris) {
	ephemerisPool.Put(e)
}

func (v Vec3) Distance(other Vec3) float64 {
	dx := v.X - other.X
	dy := v.Y - other.Y
	dz := v.Z - other.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func (v Vec3) DistanceSquared(other Vec3) float64 {
	dx := v.X - other.X
	dy := v.Y - other.Y
	dz := v.Z - other.Z
	return dx*dx + dy*dy + dz*dz
}
