package models

import "time"

type JobsListing struct {
	Data []struct {
		Advertiser struct {
			ID          string `json:"id"`
			Description string `json:"description"`
		} `json:"advertiser"`
		Area               string `json:"area,omitempty"`
		AreaID             int    `json:"areaId,omitempty"`
		AreaWhereValue     string `json:"areaWhereValue,omitempty"`
		AutomaticInclusion bool   `json:"automaticInclusion"`
		Branding           struct {
			ID     string `json:"id"`
			Assets struct {
				Logo struct {
					Strategies struct {
						JdpLogo  string `json:"jdpLogo"`
						SerpLogo string `json:"serpLogo"`
					} `json:"strategies"`
				} `json:"logo"`
			} `json:"assets"`
		} `json:"branding,omitempty"`
		BulletPoints   []string `json:"bulletPoints"`
		Classification struct {
			ID          string `json:"id"`
			Description string `json:"description"`
		} `json:"classification"`
		CompanyName                    string `json:"companyName,omitempty"`
		CompanyProfileStructuredDataID int    `json:"companyProfileStructuredDataId"`
		DisplayStyle                   struct {
			Search string `json:"search"`
		} `json:"displayStyle"`
		DisplayType        string `json:"displayType"`
		ListingDateDisplay string `json:"listingDateDisplay"`
		Location           string `json:"location"`
		LocationID         int    `json:"locationId"`
		LocationWhereValue string `json:"locationWhereValue"`
		ID                 int    `json:"id"`
		IsPremium          bool   `json:"isPremium"`
		IsStandOut         bool   `json:"isStandOut"`
		JobLocation        struct {
			Label        string `json:"label"`
			CountryCode  string `json:"countryCode"`
			SeoHierarchy []struct {
				ContextualName string `json:"contextualName"`
			} `json:"seoHierarchy"`
		} `json:"jobLocation"`
		ListingDate time.Time `json:"listingDate"`
		Logo        struct {
			ID          string `json:"id"`
			Description any    `json:"description"`
		} `json:"logo"`
		RoleID      string `json:"roleId,omitempty"`
		Salary      string `json:"salary"`
		SolMetadata struct {
			SearchRequestToken string `json:"searchRequestToken"`
			Token              string `json:"token"`
			JobID              string `json:"jobId"`
			Section            string `json:"section"`
			SectionRank        int    `json:"sectionRank"`
			JobAdType          string `json:"jobAdType"`
			Tags               struct {
				MordorFlights string `json:"mordor__flights"`
				MordorS       string `json:"mordor__s"`
			} `json:"tags"`
		} `json:"solMetadata"`
		SubClassification struct {
			ID          string `json:"id"`
			Description string `json:"description"`
		} `json:"subClassification"`
		Suburb           string `json:"suburb,omitempty"`
		SuburbID         int    `json:"suburbId,omitempty"`
		SuburbWhereValue string `json:"suburbWhereValue,omitempty"`
		Teaser           string `json:"teaser"`
		Title            string `json:"title"`
		Tracking         string `json:"tracking"`
		WorkType         string `json:"workType"`
		WorkArrangements struct {
			Data []struct {
				ID    string `json:"id"`
				Label struct {
					Text string `json:"text"`
				} `json:"label"`
			} `json:"data"`
		} `json:"workArrangements,omitempty"`
		IsPrivateAdvertiser bool `json:"isPrivateAdvertiser"`
		Tags                []struct {
			Type  string `json:"type"`
			Label string `json:"label"`
		} `json:"tags,omitempty"`
	} `json:"data"`
	Title                string `json:"title"`
	TotalCount           int    `json:"totalCount"`
	TotalPages           int    `json:"totalPages"`
	PaginationParameters struct {
		SeekSelectAllPages bool `json:"seekSelectAllPages"`
		HadPremiumListings bool `json:"hadPremiumListings"`
	} `json:"paginationParameters"`
	Info struct {
		TimeTaken  int    `json:"timeTaken"`
		Source     string `json:"source"`
		Experiment string `json:"experiment"`
	} `json:"info"`
	UserQueryID string `json:"userQueryId"`
	SortMode    []struct {
		IsActive bool   `json:"isActive"`
		Name     string `json:"name"`
		Value    string `json:"value"`
	} `json:"sortMode"`
	SolMetadata struct {
		RequestToken     string   `json:"requestToken"`
		Token            string   `json:"token"`
		SortMode         string   `json:"sortMode"`
		Locations        []string `json:"locations"`
		LocationDistance int      `json:"locationDistance"`
		Categories       []string `json:"categories"`
		PageSize         int      `json:"pageSize"`
		PageNumber       int      `json:"pageNumber"`
		TotalJobCount    int      `json:"totalJobCount"`
		Tags             struct {
			MordorSearchMarket     string `json:"mordor:searchMarket"`
			MordorResultCountRst   string `json:"mordor:result_count_rst"`
			MordorResultCountVec   string `json:"mordor:result_count_vec"`
			MordorRt               string `json:"mordor:rt"`
			MordorVs               string `json:"mordor_vs"`
			MordorCountVec         string `json:"mordor:count_vec"`
			MordorFlights          string `json:"mordor__flights"`
			MordorRboPerfMngP95K20 string `json:"mordor__rbo_perfMng_p95_k20"`
			MordorCountRst         string `json:"mordor:count_rst"`
			MordorCountIr          string `json:"mordor:count_ir"`
			MordorResultCountIr    string `json:"mordor:result_count_ir"`
			ChaliceSearchAPISolID  string `json:"chalice-search-api:solId"`
		} `json:"tags"`
	} `json:"solMetadata"`
	Location struct {
		AreaDescription         string `json:"areaDescription"`
		AreaID                  int    `json:"areaId"`
		Description             string `json:"description"`
		LocationDescription     string `json:"locationDescription"`
		LocationID              int    `json:"locationId"`
		Matched                 bool   `json:"matched"`
		StateDescription        string `json:"stateDescription"`
		SuburbParentDescription string `json:"suburbParentDescription"`
		Type                    string `json:"type"`
		WhereID                 int    `json:"whereId"`
		Descriptions            struct {
			En struct {
				ContextualName string `json:"contextualName"`
			} `json:"en"`
			ID struct {
				ContextualName string `json:"contextualName"`
			} `json:"id"`
			Th struct {
				ContextualName string `json:"contextualName"`
			} `json:"th"`
		} `json:"descriptions"`
		ShortLocationName string `json:"shortLocationName"`
	} `json:"location"`
	Facets struct {
	} `json:"facets"`
	JoraCrossLink struct {
		CanCrossLink bool `json:"canCrossLink"`
	} `json:"joraCrossLink"`
	SearchParams struct {
		Sitekey               string `json:"sitekey"`
		Sourcesystem          string `json:"sourcesystem"`
		Userqueryid           string `json:"userqueryid"`
		Userid                string `json:"userid"`
		Usersessionid         string `json:"usersessionid"`
		Eventcapturesessionid string `json:"eventcapturesessionid"`
		Where                 string `json:"where"`
		Page                  string `json:"page"`
		Seekselectallpages    string `json:"seekselectallpages"`
		Classification        string `json:"classification"`
		Pagesize              string `json:"pagesize"`
		Include               string `json:"include"`
		Locale                string `json:"locale"`
		Solid                 string `json:"solid"`
	} `json:"searchParams"`
}