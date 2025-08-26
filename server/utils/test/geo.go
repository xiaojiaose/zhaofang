package test12

import (
	"math"
	"sort"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/quadtree"
)

const EarthRadius = 6371000.0

var GeoSearch *GeoService

type Community struct {
	ID       string
	Name     string
	Location orb.Point
}

func (c *Community) Point() orb.Point {
	return c.Location
}

type GeoService struct {
	communityTree *quadtree.Quadtree
}

func NewGeoService() *GeoService {
	bound := orb.Bound{
		Min: orb.Point{73.0, 18.0},  // 中国大致最西和最南
		Max: orb.Point{135.0, 54.0}, // 中国大致最东和最北
	}
	//bound := orb.Bound{orb.Point{-180, -90}, orb.Point{180, 90}}
	GeoSearch = &GeoService{
		communityTree: quadtree.New(bound),
	}
	return GeoSearch
}

func (g *GeoService) AddCommunity(community *Community) error {
	return g.communityTree.Add(community)
}

// 正确的实现方法
// FindNearbyCommunities 查找指定点附近的社区（新版本API）
func (g *GeoService) FindNearbyCommunities(point orb.Point, radius float64) ([]*Community, error) {
	// 将半径转换为经纬度偏移量
	latOffset := radius / 111000.0
	lngOffset := radius / (111000.0 * math.Cos(point[1]*math.Pi/180))

	searchBound := orb.Bound{
		Min: orb.Point{point[0] - lngOffset, point[1] - latOffset},
		Max: orb.Point{point[0] + lngOffset, point[1] + latOffset},
	}

	// 新版本 API: InBound 只需要一个参数，返回切片
	items := g.communityTree.InBound(nil, searchBound)

	var results []*Community

	for _, item := range items {
		candidate, ok := item.(*Community)
		if !ok {
			continue
		}

		distance := geo.Distance(candidate.Point(), point)
		if distance <= radius {
			results = append(results, candidate)
		}
	}

	// 按距离排序
	sort.Slice(results, func(i, j int) bool {
		return geo.Distance(results[i].Point(), point) < geo.Distance(results[j].Point(), point)
	})

	return results, nil
}

//func main() {
//	geoService := NewGeoService()
//	geoService.AddCommunity(&Community{
//		ID:       "1",
//		Name:     "社区1",
//		Location: orb.Point{116.3984, 39.9093},
//	})
//	// 测试
//	point := orb.Point{116.3974, 39.9093} // 北京
//	radius := 5000.0                      // 50公里
//
//	results, err := geoService.FindNearbyCommunities(point, radius)
//	if err != nil {
//		fmt.Println("Error:", err)
//		return
//	}
//
//	fmt.Printf("找到 %d 个社区在 %.0f 米范围内\n", len(results), radius)
//}
