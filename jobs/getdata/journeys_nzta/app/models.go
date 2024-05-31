package app

// ----------------------------------------------------------------------------------------
// cameras
// ----------------------------------------------------------------------------------------
type CamerasModel struct {
	Type     string `json:"type"`
	Features []struct {
		Type       string `json:"type"`
		Properties struct {
			ClassName        string `json:"ClassName"`
			LastEdited       string `json:"LastEdited"`
			Created          string `json:"Created"`
			ExternalID       int    `json:"ExternalId"`
			Name             string `json:"Name"`
			Description      string `json:"Description"`
			Offline          int    `json:"Offline"`
			UnderMaintenance int    `json:"UnderMaintenance"`
			ImageURL         string `json:"ImageUrl"`
			ThumbURL         string `json:"ThumbUrl"`
			Latitude         string `json:"Latitude"`
			Longitude        string `json:"Longitude"`
			Direction        string `json:"Direction"`
			SortGroup        string `json:"SortGroup"`
			TasJourneyID     int    `json:"TasJourneyId"`
			RegionID         int    `json:"RegionID"`
			TasRegionID      int    `json:"TasRegionId"`
			ID               int    `json:"id"`
			Uniq             string `json:"uniq"`
			Type             string `json:"type"`
			LastUpdated      int    `json:"lastUpdated"`
		} `json:"properties"`
		Geometry struct {
			Type        string    `json:"type"`
			Coordinates []float64 `json:"coordinates"`
		} `json:"geometry"`
	} `json:"features"`
	LastUpdated int `json:"lastUpdated"`
}
