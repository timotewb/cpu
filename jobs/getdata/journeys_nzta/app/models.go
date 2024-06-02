package app

// ----------------------------------------------------------------------------------------
// chargers
// ----------------------------------------------------------------------------------------
type ChargersModel struct {
	Type     string `json:"type"`
	Features []struct {
		Type       string `json:"type"`
		Properties struct {
			ClassName            string `json:"ClassName"`
			LastEdited           string `json:"LastEdited"`
			Created              string `json:"Created"`
			SiteID               string `json:"SiteId"`
			Name                 string `json:"Name"`
			Operator             string `json:"Operator"`
			Address              string `json:"Address"`
			Is24Hours            int    `json:"Is24Hours"`
			CarParkCount         int    `json:"CarParkCount"`
			HasCarparkCost       int    `json:"HasCarparkCost"`
			MaxTimeLimit         string `json:"MaxTimeLimit"`
			HasTouristAttraction int    `json:"HasTouristAttraction"`
			AccessLocations      []struct {
				Lat float64 `json:"lat"`
				Lon float64 `json:"lon"`
			} `json:"AccessLocations"`
			ProviderDeleted int    `json:"ProviderDeleted"`
			HideFromFeed    int    `json:"HideFromFeed"`
			RegionID        int    `json:"RegionID"`
			ID              int    `json:"ID"`
			RecordClassName string `json:"RecordClassName"`
			ExternalID      int    `json:"ExternalId"`
			Uniq            string `json:"uniq"`
			Type            string `json:"type"`
			ID0             int    `json:"id"`
			LastUpdated     int    `json:"lastUpdated"`
			Region          string `json:"Region"`
			Connectors      []struct {
				Current           string `json:"Current"`
				KwRated           int    `json:"KwRated"`
				ConnectorType     string `json:"ConnectorType"`
				OperationStatus   string `json:"OperationStatus"`
				NextPlannedOutage string `json:"NextPlannedOutage"`
			} `json:"connectors"`
			OwnerName    string `json:"OwnerName"`
			ChargingCost string `json:"ChargingCost"`
			FeatureType  string `json:"featureType"`
			Regions      []struct {
				ID string `json:"id"`
			} `json:"regions"`
		} `json:"properties"`
		Geometry struct {
			Type        string    `json:"type"`
			Coordinates []float64 `json:"coordinates"`
		} `json:"geometry"`
	} `json:"features"`
}

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
