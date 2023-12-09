package MapTMap

import (
	"fmt"
	"testing"
)

func TestCgcs2000ToGcj02(t *testing.T) {
	lng, lat := 112.60309, 37.79032
	gcj02Lng, gcj02Lat := WGS84ToGcj02(lng, lat)
	fmt.Printf("WGS84: %f, %f\n", lng, lat)
	fmt.Printf("GCJ02: %f, %f\n", gcj02Lng, gcj02Lat)
	// Output: WGS84: 112.603090, 37.790320
	//         GCJ02: 112.609323, 37.790806
}
