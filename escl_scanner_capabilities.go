package airscan

type adf struct {
	AdfOptions          adfOptions           `xml:"AdfOptions"`
	AdfSimplexInputCaps *adfSimplexInputCaps `xml:"AdfSimplexInputCaps"`
	AdfDuplexInputCaps  *adfDuplexInputCaps  `xml:"AdfDuplexInputCaps"`
	FeederCapacity      int                  `xml:"FeederCapacity"`
	Justification       justification        `xml:"Justification"`
}

type adfOptions struct {
	AdfOption string `xml:"AdfOption"`
}

type adfSimplexInputCaps struct {
	MaxHeight             int              `xml:"MaxHeight"`
	MaxOpticalXResolution int              `xml:"MaxOpticalXResolution"`
	MaxOpticalYResolution int              `xml:"MaxOpticalYResolution"`
	MaxPhysicalHeight     int              `xml:"MaxPhysicalHeight"`
	MaxPhysicalWidth      int              `xml:"MaxPhysicalWidth"`
	MaxScanRegions        bool             `xml:"MaxScanRegions"`
	MaxWidth              int              `xml:"MaxWidth"`
	MinHeight             int              `xml:"MinHeight"`
	MinWidth              int              `xml:"MinWidth"`
	SettingProfiles       settingProfiles  `xml:"SettingProfiles"`
	SupportedIntents      supportedIntents `xml:"SupportedIntents"`
}

type adfDuplexInputCaps struct {
	FeedDirections   feedDirections   `xml:"FeedDirections"`
	MaxHeight        int              `xml:"MaxHeight"`
	MaxWidth         int              `xml:"MaxWidth"`
	MinHeight        int              `xml:"MinHeight"`
	MinWidth         int              `xml:"MinWidth"`
	SettingProfiles  settingProfiles  `xml:"SettingProfiles"`
	SupportedIntents supportedIntents `xml:"SupportedIntents"`
}

type feedDirections struct {
	FeedDirection []string `xml:"FeedDirection"`
}

type ccdChannels struct {
	CcdChannel []string `xml:"CcdChannel"`
}

type certifications struct {
	Name    string  `xml:"Name"`
	Version float64 `xml:"Version"`
}

type colorModes struct {
	ColorMode []string `xml:"ColorMode"`
}

type colorSpaces struct {
	ColorSpace string `xml:"ColorSpace"`
}

type contentTypes struct {
	ContentType []string `xml:"ContentType"`
}

type discreteResolution struct {
	XResolution int `xml:"XResolution"`
	YResolution int `xml:"YResolution"`
}

type discreteResolutions struct {
	DiscreteResolution discreteResolution `xml:"DiscreteResolution"`
}

type documentFormats struct {
	DocumentFormat    []string `xml:"DocumentFormat"`
	DocumentFormatExt []string `xml:"DocumentFormatExt"`
}

type justification struct {
	XImagePosition string `xml:"XImagePosition"`
	YImagePosition string `xml:"YImagePosition"`
}

type platen struct {
	PlatenInputCaps platenInputCaps `xml:"PlatenInputCaps"`
}

type platenInputCaps struct {
	MaxHeight             int              `xml:"MaxHeight"`
	MaxOpticalXResolution int              `xml:"MaxOpticalXResolution"`
	MaxOpticalYResolution int              `xml:"MaxOpticalYResolution"`
	MaxPhysicalHeight     int              `xml:"MaxPhysicalHeight"`
	MaxPhysicalWidth      int              `xml:"MaxPhysicalWidth"`
	MaxScanRegions        bool             `xml:"MaxScanRegions"`
	MaxWidth              int              `xml:"MaxWidth"`
	MinHeight             int              `xml:"MinHeight"`
	MinWidth              int              `xml:"MinWidth"`
	SettingProfiles       settingProfiles  `xml:"SettingProfiles"`
	SupportedIntents      supportedIntents `xml:"SupportedIntents"`
}

type scannerCapabilities struct {
	Adf            *adf           `xml:"Adf"`
	AdminURI       string         `xml:"AdminURI"`
	Certifications certifications `xml:"Certifications"`
	IconURI        string         `xml:"IconURI"`
	MakeAndModel   string         `xml:"MakeAndModel"`
	Manufacturer   string         `xml:"Manufacturer"`
	Platen         *platen        `xml:"Platen"`
	SerialNumber   string         `xml:"SerialNumber"`
	SharpenSupport sharpenSupport `xml:"SharpenSupport"`
	UUID           string         `xml:"UUID"`
	Version        float64        `xml:"Version"`
}

type settingProfile struct {
	CcdChannels          ccdChannels          `xml:"CcdChannels"`
	ColorModes           colorModes           `xml:"ColorModes"`
	ColorSpaces          colorSpaces          `xml:"ColorSpaces"`
	ContentTypes         contentTypes         `xml:"ContentTypes"`
	DocumentFormats      documentFormats      `xml:"DocumentFormats"`
	SupportedResolutions supportedResolutions `xml:"SupportedResolutions"`
}

type settingProfiles struct {
	SettingProfile settingProfile `xml:"SettingProfile"`
}

type sharpenSupport struct {
	Max    int  `xml:"Max"`
	Min    bool `xml:"Min"`
	Normal int  `xml:"Normal"`
	Step   bool `xml:"Step"`
}

type supportedIntents struct {
	Intent []string `xml:"Intent"`
}

type supportedResolutions struct {
	DiscreteResolutions discreteResolutions `xml:"DiscreteResolutions"`
}
