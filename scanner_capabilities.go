package airscan

type Adf struct {
	AdfOptions          AdfOptions           `xml:"AdfOptions"`
	AdfSimplexInputCaps *AdfSimplexInputCaps `xml:"AdfSimplexInputCaps"`
	AdfDuplexInputCaps  *AdfDuplexInputCaps  `xml:"AdfDuplexInputCaps"`
	FeederCapacity      int                  `xml:"FeederCapacity"`
	Justification       Justification        `xml:"Justification"`
}

type AdfOptions struct {
	AdfOption string `xml:"AdfOption"`
}

type AdfSimplexInputCaps struct {
	MaxHeight             int              `xml:"MaxHeight"`
	MaxOpticalXResolution int              `xml:"MaxOpticalXResolution"`
	MaxOpticalYResolution int              `xml:"MaxOpticalYResolution"`
	MaxPhysicalHeight     int              `xml:"MaxPhysicalHeight"`
	MaxPhysicalWidth      int              `xml:"MaxPhysicalWidth"`
	MaxScanRegions        bool             `xml:"MaxScanRegions"`
	MaxWidth              int              `xml:"MaxWidth"`
	MinHeight             int              `xml:"MinHeight"`
	MinWidth              int              `xml:"MinWidth"`
	SettingProfiles       SettingProfiles  `xml:"SettingProfiles"`
	SupportedIntents      SupportedIntents `xml:"SupportedIntents"`
}

type AdfDuplexInputCaps struct {
	FeedDirections   FeedDirections   `xml:"FeedDirections"`
	MaxHeight        int              `xml:"MaxHeight"`
	MaxWidth         int              `xml:"MaxWidth"`
	MinHeight        int              `xml:"MinHeight"`
	MinWidth         int              `xml:"MinWidth"`
	SettingProfiles  SettingProfiles  `xml:"SettingProfiles"`
	SupportedIntents SupportedIntents `xml:"SupportedIntents"`
}

type FeedDirections struct {
	FeedDirection []string `xml:"FeedDirection"`
}

type CcdChannels struct {
	CcdChannel []string `xml:"CcdChannel"`
}

type Certifications struct {
	Name    string  `xml:"Name"`
	Version float64 `xml:"Version"`
}

type ColorModes struct {
	ColorMode []string `xml:"ColorMode"`
}

type ColorSpaces struct {
	ColorSpace string `xml:"ColorSpace"`
}

type ContentTypes struct {
	ContentType []string `xml:"ContentType"`
}

type DiscreteResolution struct {
	XResolution int `xml:"XResolution"`
	YResolution int `xml:"YResolution"`
}

type DiscreteResolutions struct {
	DiscreteResolution DiscreteResolution `xml:"DiscreteResolution"`
}

type DocumentFormats struct {
	DocumentFormat    []string `xml:"DocumentFormat"`
	DocumentFormatExt []string `xml:"DocumentFormatExt"`
}

type Justification struct {
	XImagePosition string `xml:"XImagePosition"`
	YImagePosition string `xml:"YImagePosition"`
}

type Platen struct {
	PlatenInputCaps PlatenInputCaps `xml:"PlatenInputCaps"`
}

type PlatenInputCaps struct {
	MaxHeight             int              `xml:"MaxHeight"`
	MaxOpticalXResolution int              `xml:"MaxOpticalXResolution"`
	MaxOpticalYResolution int              `xml:"MaxOpticalYResolution"`
	MaxPhysicalHeight     int              `xml:"MaxPhysicalHeight"`
	MaxPhysicalWidth      int              `xml:"MaxPhysicalWidth"`
	MaxScanRegions        bool             `xml:"MaxScanRegions"`
	MaxWidth              int              `xml:"MaxWidth"`
	MinHeight             int              `xml:"MinHeight"`
	MinWidth              int              `xml:"MinWidth"`
	SettingProfiles       SettingProfiles  `xml:"SettingProfiles"`
	SupportedIntents      SupportedIntents `xml:"SupportedIntents"`
}

type ScannerCapabilities struct {
	Adf            *Adf           `xml:"Adf"`
	AdminURI       string         `xml:"AdminURI"`
	Certifications Certifications `xml:"Certifications"`
	IconURI        string         `xml:"IconURI"`
	MakeAndModel   string         `xml:"MakeAndModel"`
	Manufacturer   string         `xml:"Manufacturer"`
	Platen         *Platen        `xml:"Platen"`
	SerialNumber   string         `xml:"SerialNumber"`
	SharpenSupport SharpenSupport `xml:"SharpenSupport"`
	UUID           string         `xml:"UUID"`
	Version        float64        `xml:"Version"`
}

type SettingProfile struct {
	CcdChannels          CcdChannels          `xml:"CcdChannels"`
	ColorModes           ColorModes           `xml:"ColorModes"`
	ColorSpaces          ColorSpaces          `xml:"ColorSpaces"`
	ContentTypes         ContentTypes         `xml:"ContentTypes"`
	DocumentFormats      DocumentFormats      `xml:"DocumentFormats"`
	SupportedResolutions SupportedResolutions `xml:"SupportedResolutions"`
}

type SettingProfiles struct {
	SettingProfile SettingProfile `xml:"SettingProfile"`
}

type SharpenSupport struct {
	Max    int  `xml:"Max"`
	Min    bool `xml:"Min"`
	Normal int  `xml:"Normal"`
	Step   bool `xml:"Step"`
}

type SupportedIntents struct {
	Intent []string `xml:"Intent"`
}

type SupportedResolutions struct {
	DiscreteResolutions DiscreteResolutions `xml:"DiscreteResolutions"`
}
