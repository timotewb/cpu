package app

type ResponseWrapper struct {
	Cnt  int        `json:"cnt"`
	List []Response `json:"list",omitempty`
}

type Response struct {
	Coord      Coordinates      `json:"coord",omitempty`
	Weather    []WeatherDetails `json:"weather",omitempty`
	Base       string           `json:"base",omitempty`
	Main       MainDetails      `json:"main",omitempty`
	Visibility int32            `json:"visibility",omitempty`
	Wind       WindDetails      `json:"wind",omitempty`
	Clouds     CloudDetails     `json:"clouds",omitempty`
	Rain       RainDetails      `json:"rain",omitempty`
	Snow       SnowDetails      `json:"snow",omitempty`
	Dt         int32            `json:"dt",omitempty`
	Sys        SysDetails       `json:"sys",omitempty`
	Timezone   int32            `json:"timezone",omitempty`
	Id         int32            `json:"id",omitempty`
	Name       string           `json:"name",omitempty`
	Cod        int16            `json:"cod",omitempty`
}

type Coordinates struct {
	Lon float32 `json:"lon",omitempty`
	Lat float32 `json:"lat",omitempty`
}

type WeatherDetails struct {
	Id          int32  `json:"id",omitempty`
	Main        string `json:"main",omitempty`
	Description string `json:"description",omitempty`
	Icon        string `json:"icon",omitempty`
}

type MainDetails struct {
	Temp       float32 `json:"temp",omitempty`
	Feels_like float32 `json:"feels_like",omitempty`
	Pressure   int32   `json:"pressure",omitempty`
	Humidity   int32   `json:"humidity",omitempty`
	Temp_min   float32 `json:"temp_min",omitempty`
	Temp_max   float32 `json:"temp_max",omitempty`
	Sea_level  float32 `json:"sea_level",omitempty`
	Grnd_level float32 `json:"grnd_level",omitempty`
}

type WindDetails struct {
	Speed float32 `json:"speed",omitempty`
	Deg   float32 `json:"deg",omitempty`
}

type CloudDetails struct {
	All int16 `json:"all",omitempty`
}

type RainDetails struct {
	Rain1h int16 `json:"1h",omitempty`
	Rain3h int16 `json:"3h",omitempty`
}

type SnowDetails struct {
	Snow1h int16 `json:"1h",omitempty`
	Snow3h int16 `json:"3h",omitempty`
}

type SysDetails struct {
	Type    int16  `json:"type",omitempty`
	Id      int32  `json:"id",omitempty`
	Message string `json:"message",omitempty`
	Country string `json:"country",omitempty`
	Sunrise int32  `json:"sunrise",omitempty`
	Sunset  int32  `json:"sunset",omitempty`
}
