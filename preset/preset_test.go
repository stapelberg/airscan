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

package preset

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Same contents (aside from whitespace differences) as Apple’s
// AirScanScanner/41 (Image Preview app on Mac OS X 10.11):
const airscanScannerGolden = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<scan:ScanSettings xmlns:scan="http://schemas.hp.com/imaging/escl/2011/05/03" xmlns:pwg="http://www.pwg.org/schemas/2010/12/sm">
  <pwg:Version>2.0</pwg:Version>
  <pwg:ScanRegions pwg:MustHonor="true">
    <pwg:ScanRegion>
      <pwg:ContentRegionUnits>escl:ThreeHundredthsOfInches</pwg:ContentRegionUnits>
      <pwg:Width>2480</pwg:Width>
      <pwg:Height>3508</pwg:Height>
      <pwg:XOffset>0</pwg:XOffset>
      <pwg:YOffset>0</pwg:YOffset>
    </pwg:ScanRegion>
  </pwg:ScanRegions>
  <pwg:DocumentFormat>image/jpeg</pwg:DocumentFormat>
  <pwg:InputSource>Feeder</pwg:InputSource>
  <scan:ColorMode>Grayscale8</scan:ColorMode>
  <scan:XResolution>300</scan:XResolution>
  <scan:YResolution>300</scan:YResolution>
  <scan:Duplex>true</scan:Duplex>
</scan:ScanSettings>`

func TestScanSettings(t *testing.T) {
	settings := GrayscaleA4ADF()
	got, err := settings.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(airscanScannerGolden, got); diff != "" {
		t.Fatalf("unexpected ScanSettings request: diff (-want +got):\n%s", diff)
	}
}
