package main

import (
	"bufio"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/core"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/initialize"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"
	"time"
)

func Init() {
	global.GVA_VP = core.Viper() // 初始化Viper
	initialize.OtherInit()
	global.GVA_LOG = core.Zap() // 初始化zap日志库
	zap.ReplaceGlobals(global.GVA_LOG)
	global.GVA_DB = initialize.Gorm() // gorm连接数据库
	global.InitZincSearch(
		global.GVA_CONFIG.ZincSearch.Url,
		global.GVA_CONFIG.ZincSearch.Username,
		global.GVA_CONFIG.ZincSearch.Password,
		global.GVA_CONFIG.ZincSearch.ResourceIndex,
	)
}

// 增加小区 脚本，小区手动录入到数据库
func main() {

	Init()
	filename := "/Users/xiongda/GolandProjects/zhaofang/server/utilities/add_xq/xiaoqu.txt"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("无法打开文件:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 1

	for scanner.Scan() {
		line := scanner.Text()
		// 处理每一行
		processLine(line, lineNumber)
		// 如果是超大文件，可以定期输出进度
		if lineNumber%20000 == 0 {
			fmt.Printf("=========> 已处理 %d 行\n", lineNumber)
		}
		lineNumber++
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("无法打开文件:", err)
		return
	}
	fmt.Printf("文件读取完成，共 %d 行\n", lineNumber)
	return
}

func processLine(line string, lineNumber int) {
	// 20024,海纳光合谷(裕华),114.615679,38.014228,精英中学,开发区,学苑路
	arr := strings.Split(strings.Replace(line, " ", "", -1), ",")
	if len(arr) != 7 {
		fmt.Errorf("lineNumber: %d, line: %s", lineNumber, line)
		return
	}
	communityId, _ := strconv.Atoi(arr[0])
	communityName := arr[1]
	latitude := arr[2]
	longitude := arr[3]
	district := arr[4]
	area := arr[5]
	address := arr[6]

	var xq system.XiaoQu
	tx := global.GVA_DB.Model(&system.XiaoQu{}).Where("community_id = ?", communityId).First(&xq)
	if tx.Error != nil {
		fmt.Printf("++++++ %s, %d\n", tx.Error.Error(), communityId)
	} else {
		fmt.Printf("==== %d\n", communityId)
	}
	var dis system.District
	_ = global.GVA_DB.Model(&system.District{}).Where("name = ?", district).First(&dis).Error

	if xq.ID == 0 {
		xq = system.XiaoQu{
			CommunityId: communityId,
			Name:        communityName,
			Latitude:    latitude,
			Longitude:   longitude,
			Districts:   district,
			DistrictIds: fmt.Sprintf("|%d|", dis.ID),
			Area:        area,
			Address:     address,
		}
		err := global.GVA_DB.Model(&system.XiaoQu{}).Create(&xq).Error
		if err != nil {
			fmt.Printf("创建XiaoQu错误: %v\n", err)
			return
		}
	} else {
		if len(xq.Districts) > 0 {
			return
		}
		xq.Name = communityName
		xq.Latitude = latitude
		xq.Longitude = longitude
		xq.Area = area
		xq.Address = address
		xq.Districts = district
		xq.DistrictIds = fmt.Sprintf("|%d|", dis.ID)
		xq.CreatedAt = time.Now()
		xq.UpdatedAt = time.Now()
		e := global.GVA_DB.Save(&xq)
		if e.Error != nil {
			fmt.Printf("更新XiaoQu错误: %v\n", e.Error)
			return
		}
	}
}
