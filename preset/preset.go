// Copyright 2020 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package preset contains settings that have been verified to match,
// character-by-character, the requests that Apple’s AirScanScanner library
// produces.
//
// Given that AirScan is not standardized in an open standard, this is a good
// way of reaching a high level of compatibility with devices in the wild.
package preset

import "github.com/stapelberg/airscan"

// GrayscaleA4ADF scans an A4 document at 300 dpi from the Automated Document
// Feeder (ADF) in grayscale. Each call will return a struct that is safe to
// modify.
func GrayscaleA4ADF() *airscan.ScanSettings {
	return &airscan.ScanSettings{
		XmlnsScan: "http://schemas.hp.com/imaging/escl/2011/05/03",
		XmlnsPWG:  "http://www.pwg.org/schemas/2010/12/sm",
		Version:   "2.0",
		ScanRegions: airscan.ScanRegions{
			MustHonor: true,
			Regions: []*airscan.ScanRegion{
				// A4 at 300 dpi is 2480 x 3508 as per
				// https://www.papersizes.org/a-sizes-in-pixels.htm
				{
					ContentRegionUnits: "escl:ThreeHundredthsOfInches",
					Width:              2480,
					Height:             3508,
					XOffset:            0,
					YOffset:            0,
				},
			},
		},
		DocumentFormat: "image/jpeg",
		InputSource:    "Feeder",
		ColorMode:      "Grayscale8",
		XResolution:    300,
		YResolution:    300,
		Duplex:         true,
	}
}
