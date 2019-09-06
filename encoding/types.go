//go:generate easyjson -all
package main

import (
	"image/color"
)

type MapOfStr map[string]string
type MapOfSlice map[string][]byte
type SliceOfStr []struct {
	Name string
	Data string
}
type SliceOfSlice []struct {
	Name string
	Data []byte
}

type EMapOfStr struct {
	Map MapOfStr
}
type EMapOfSlice struct {
	Map MapOfSlice
}
type ESliceOfStr struct {
	Slice SliceOfStr
}
type ESliceOfSlice struct {
	Slice SliceOfSlice
}

type Galaxy struct {
	//for systems - range of "system borders"
	SpawnDistance float64

	Points map[string]*GalaxyPoint

	//recalculated on Decode
	Ordered []*GalaxyPoint `json:"-"`
	maxLvl  int

	//used by update
	fixedTimeRest float64
}

type GalaxyPoint struct {
	//Id setted on load from file
	ID        string `json:"id,omitempty"`
	ParentID  string `json:"pid,omitempty"`
	IsVirtual bool   `json:"iv,omitempty"`

	//found on recalc
	//phys order level
	Level int `json:"lv,omitempty"`
	//graph order level, ignore
	GLevel int `json:"gl,omitempty"`

	Pos V2

	Orbit    float64 `json:"orb,omitempty"`
	Period   float64 `json:"per,omitempty"`
	AngPhase float64 `json:"ang,omitempty"`

	Type     string  `json:"t,omitempty"`
	SpriteAN string  `json:"sp,omitempty"`
	Size     float64 `json:"s,omitempty"`

	Mass   float64 `json:"m,omitempty"`
	GDepth float64 `json:"gd,omitempty"`

	//for warp points
	WarpSpawnDistance float64 `json:"wsd,omitempty"`
	WarpYellowOutDist float64 `json:"wyo,omitempty"`
	WarpGreenOutDist  float64 `json:"wgo,omitempty"`
	WarpGreenInDist   float64 `json:"wgi,omitempty"`
	WarpRedOutDist    float64 `json:"wro,omitempty"`
	//for warp points
	InnerColor color.RGBA `json:"wic,omitempty"`
	OuterColor color.RGBA `json:"woc,omitempty"`
	GreenColor color.RGBA `json:"wgc,omitempty"`

	ScanData string `json:"sd,omitempty"`

	Minerals   []int       `json:"mi,omitempty"`
	Emissions  []Emission  `json:"emm,omitempty"`
	Signatures []Signature `json:"sig,omitempty"`
	Color      color.RGBA  `json:"clr,omitempty"`

	//updated on Decode or add|del building
	//map[ownerName][]fullkey
	Mines      map[string][]string `json:"mns,omitempty"`
	FishHouses map[string]string   `json:"fhs,omitempty"`

	//for warp points
	//map[fullKey]message
	Beacons    map[string]string `json:"bcs,omitempty"`
	BlackBoxes map[string]string `json:"bbs,omitempty"`
}

type Emission struct {
	Type      string
	MainRange float64 `json:",omitempty"`
	MainValue float64 `json:",omitempty"`
	FarRange  float64 `json:",omitempty"`
	FarValue  float64 `json:",omitempty"`
}

type Signature struct {
	TypeName  string `json:"t"`
	SigString string `json:"s",omitempty`
	//deviation of this instance
	//supposed to be Len<=1
	Dev V2 `json:"d"`
}

type V2 struct {
	X float64 `json:",omitempty"`
	Y float64 `json:",omitempty"`
}

type PilotData struct {
	Ship        RBData  `json:"sh,omitempty"`
	SessionTime float64 `json:"ss,omitempty"`
	FlightTime  float64 `json:"ft,omitempty"`
	//for cosmo
	ThrustVector V2 `json:"tv,omitempty"`
	//for warp
	Distortion float64 `json:"wd,omitempty"`
	DistTurn   float64 `json:"dt,omitempty"`
	Dir        float64 `json:"dr,omitempty"`
	//warp position for return from zero system
	WarpPos V2 `json:"wp,omitempty"`

	//to Engi

	HeatProduction float64 `json:"hp,omitempty"`

	//do not reload same Msg, cz of ship.Pos extrapolate and SessionTime+=dt
	MsgID int `json:"id"`
}

type NaviData struct {
	//drop items
	BeaconCount int `json:"bc,omitempty"`
	//[]corpName, i.e. ["gd","gd","pre"]
	//[]planetName, i.e. ["CV8-85","RD4-42-13"]
	Mines   []string `json:"mn,omitempty"`
	Landing []string `json:"ld,omitempty"`

	//cosmo
	IsScanning    bool   `json:"is,omitempty"`
	IsDrop        bool   `json:"st,omitempty"`
	ScanObjectID  string `json:"so,omitempty"`
	IsOrbiting    bool   `json:"io,omitempty"`
	OrbitObjectID string `json:"oo,omitempty"`
	ActiveMarker  bool   `json:"ma,omitempty"`
	MarkerPos     V2     `json:"mp,omitempty"`

	//warp
	SonarDir   float64 `json:"sd,omitempty"`
	SonarRange float64 `json:"sr,omitempty"`
	SonarWide  float64 `json:"sw,omitempty"`

	CanLandhome bool
}

type RBData struct {
	Pos    V2
	Ang    float64
	Vel    V2
	AngVel float64
}

type EngiCounters struct {
	Fuel       float64 `json:"f,omitempty"`
	HoleSize   float64 `json:"h,omitempty"`
	Pressure   float64 `json:"p,omitempty"`
	Air        float64 `json:"a,omitempty"`
	Calories   float64 `json:"t,omitempty"`
	CO2        float64 `json:"co2,omitempty"`
	FlightTime float64 `json:"ft,omitempty"`
	Hitted     float64 `json:"ht,omitempty"`
}

type Boost struct {
	SysN     int
	LeftTime float64
	Power    float64
}

type EngiData struct {
	//[0.0 - 1.0]
	//0 for fully OKEY, 1 - for totally DEGRADED
	BSPDegrade BSPDegrade         `json:"deg,omitempty"`
	AZ         [8]float64         `json:"az,omitempty"`
	InV        [8]uint16          `json:"inv,omitempty"`
	Emissions  map[string]float64 `json:"emm,omitempty"`

	//Counters
	Counters EngiCounters `json:"c,omitempty"`
	Boosts   []Boost      `json:"b,omitempty"`
}

type BSPDegrade BSPParams

type BSPParams struct {
	//0...100
	March_engine struct {
		Thrust_max   float64 `json:"thrust"`
		Thrust_acc   float64 `json:"accel"`
		Thrust_slow  float64 `json:"slowdown"`
		Reverse_max  float64 `json:"thrust_rev"`
		Reverse_acc  float64 `json:"accel_rev"`
		Reverse_slow float64 `json:"slowdown_rev"`
		Heat_prod    float64 `json:"heat_prod"`
		AZ           float64 `json:"az_level"`
		Volume       float64 `json:"volume"`
	} `json:"march_engine"`

	Warp_engine struct {
		Distort_max            float64 `json:"distort"`
		Distort_acc            float64 `json:"distort_acc"`
		Distort_slow           float64 `json:"distort_slow"`
		Consumption            float64 `json:"consumption"`
		Warp_enter_consumption float64 `json:"warp_enter_consumption"`
		Turn_speed             float64 `json:"turn_speed"`
		Turn_consumption       float64 `json:"turn_consumption"`
		AZ                     float64 `json:"az_level"`
		Volume                 float64 `json:"volume"`
	} `json:"warp_engine"`

	Shunter struct {
		Turn_max    float64 `json:"turn"`
		Turn_acc    float64 `json:"turn_acc"`
		Turn_slow   float64 `json:"turn_slow"`
		Strafe_max  float64 `json:"strafe"`
		Strafe_acc  float64 `json:"strafe_acc"`
		Strafe_slow float64 `json:"strafe_slow"`
		Heat_prod   float64 `json:"heat_prod"`
		AZ          float64 `json:"az_level"`
		Volume      float64 `json:"volume"`
	} `json:"shunter"`

	Radar struct {
		Range_Max    float64 `json:"range_max"`
		Angle_Min    float64 `json:"angle_min"`
		Angle_Max    float64 `json:"angle_max"`
		Angle_Change float64 `json:"angle_change"`
		Range_Change float64 `json:"range_change"`
		Rotate_Speed float64 `json:"rotate_speed"`
		AZ           float64 `json:"az_level"`
		Volume       float64 `json:"volume"`
	} `json:"radar"`

	Scanner struct {
		DropRange float64 `json:"drop_range"`
		DropSpeed float64 `json:"drop_speed"`
		ScanRange float64 `json:"scan_range"`
		ScanSpeed float64 `json:"scan_speed"`
		AZ        float64 `json:"az_level"`
		Volume    float64 `json:"volume"`
	} `json:"scaner"`

	Fuel_tank struct {
		Fuel_volume     float64 `json:"fuel_volume"`
		Fuel_Protection float64 `json:"fuel_protection"`
		Radiation_def   float64 `json:"radiation_def"`
		AZ              float64 `json:"az_level"`
		Volume          float64 `json:"volume"`
	} `json:"fuel_tank"`

	Lss struct {
		Thermal_def       float64 `json:"thermal_def"`
		Co2_level         float64 `json:"co2_level"`
		Air_volume        float64 `json:"air_volume"`
		Air_prepare_speed float64 `json:"air_speed"`
		Lightness         float64 `json:"lightness"`
		AZ                float64 `json:"az_level"`
		Volume            float64 `json:"volume"`
	} `json:"lss"`

	Shields struct {
		Radiation_def   float64 `json:"radiation_def"`
		Disinfect_level float64 `json:"disinfect_level"`
		Mechanical_def  float64 `json:"mechanical_def"`
		Heat_reflection float64 `json:"heat_reflection"`
		Heat_capacity   float64 `json:"heat_capacity"`
		Heat_sink       float64 `json:"heat_sink"`
		AZ              float64 `json:"az_level"`
		Volume          float64 `json:"volume"`
	} `json:"shields"`
}
