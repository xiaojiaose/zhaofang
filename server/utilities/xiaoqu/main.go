package main

import (
	"bufio"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/core"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/initialize"
	"github.com/flipped-aurora/gin-vue-admin/server/model/house"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"os"
	"strings"
	"time"
)

func Init() {
	global.GVA_VP = core.Viper() // 初始化Viper
	initialize.OtherInit()
	global.GVA_LOG = core.Zap() // 初始化zap日志库
	zap.ReplaceGlobals(global.GVA_LOG)
	global.GVA_DB = initialize.Gorm() // gorm连接数据库
	global.InitZincSearch(global.GVA_CONFIG.ZincSearch.Url, global.GVA_CONFIG.ZincSearch.Username, global.GVA_CONFIG.ZincSearch.Password)
}

// 增加小区 脚本，小区手动录入到数据库
func main() {

	Init()
	filename := "/opt/zhaofang/918.txt"
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
	arr := strings.Split(strings.Replace(line, " ", "", -1), ",")
	if len(arr) != 5 {
		fmt.Errorf("lineNumber: %d, line: %s", lineNumber, line)
		return
	}
	xqName := arr[1]
	buildingName := arr[2]
	unitName := arr[3]
	houseName := arr[4]

	var xq system.XiaoQu
	tx := global.GVA_DB.Model(&system.XiaoQu{}).Where("name = ?", xqName).First(&xq)
	if tx.Error != nil {
		fmt.Printf("查询XiaoQu错误: %v\n", tx.Error)
		return
	}

	xiaoquId := xq.CommunityId
	//err := global.GVA_DB.Exec("insert into dict_building (community_id, building_open_id, encrypt_building_name) values (?, ?, ?)",
	//	xiaoquId, uuid.New().String(), buildingName).Error

	var building house.DictBuilding
	err := global.GVA_DB.Model(&house.DictBuilding{}).Where("community_id = ? and encrypt_building_name = ?", xiaoquId, buildingName).Find(&building).Error
	if err != nil {
		fmt.Printf("find building 错误: %v\n", tx.Error)
		return
	}
	if building.ID == 0 {
		building = house.DictBuilding{
			CommunityID:         int(xiaoquId),
			BuildingOpenID:      uuid.New().String(),
			EncryptBuildingName: buildingName,
			CreatedAt:           time.Now(),
		}
		global.GVA_DB.Model(&house.DictBuilding{}).Create(&building)
	}

	var unit house.DictUnit
	err = global.GVA_DB.Model(&house.DictUnit{}).Where("building_open_id = ? and encrypt_unit_name = ?", building.BuildingOpenID, unitName).Find(&unit).Error
	if err != nil {
		fmt.Printf("find unit 错误: %v\n", tx.Error)
		return
	}
	if unit.ID == 0 {
		unit = house.DictUnit{
			BuildingOpenID:  building.BuildingOpenID,
			UnitOpenID:      uuid.New().String(),
			EncryptUnitName: unitName,
			CreatedAt:       time.Now(),
		}
		global.GVA_DB.Model(&house.DictUnit{}).Create(&unit)
	}

	err = global.GVA_DB.Exec("insert into dict_house (house_open_id, encrypt_house_name, unit_open_id, created_at) values (?, ?, ?, ?)", uuid.New().String(), houseName, unit.UnitOpenID, time.Now()).Error
	if err != nil {
		global.GVA_LOG.Error("插入失败", zap.Error(err))
	}
}
