package test12

import (
	"fmt"
	"github.com/paulmach/orb"
)

func main() {
	geoService := NewGeoService()
	geoService.AddCommunity(&Community{
		ID:       "1",
		Name:     "社区1",
		Location: orb.Point{116.3984, 39.9093},
	})
	// 测试
	point := orb.Point{116.3974, 39.9093} // 北京
	radius := 5000.0                      // 50公里

	results, err := geoService.FindNearbyCommunities(point, radius)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("找到 %d 个社区在 %.0f 米范围内\n", len(results), radius)
}
