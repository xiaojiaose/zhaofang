package main

import (
	"bufio"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/core"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/initialize"
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
func main() {

	Init()
	filename := "/opt/zhaofang/output.txt"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("无法打开文件:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0

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
	if strings.Contains(line, "encryptHouseName") {
		arr := strings.Split(strings.Replace(line, " ", "", -1), ",")
		houseName := strings.Split(arr[3], ":")[1]
		houseId := strings.Split(arr[2], ":")[1]
		unitId := strings.Split(arr[1], ":")[1]

		//fmt.Printf("lineNumber: %d, houseName: %s, houseId: %s, unitId: %s\n", lineNumber, houseName, houseId, unitId)
		err := global.GVA_DB.Exec("insert into dict_house (house_open_id, encrypt_house_name, unit_open_id, created_at) values (?, ?, ?, ?)", houseId, houseName, unitId, time.Now()).Error
		if err != nil {
			global.GVA_LOG.Error("插入失败", zap.Error(err))
		}
	} else if strings.Contains(line, "encryptUnitName") {
		arr := strings.Split(strings.Replace(line, " ", "", -1), ",")
		unitName := strings.Split(arr[3], ":")[1]
		unitId := strings.Split(arr[2], ":")[1]
		buildingId := strings.Split(arr[1], ":")[1]
		err := global.GVA_DB.Exec("insert into dict_unit (building_open_id, unit_open_id, encrypt_unit_name, created_at) values (?, ?, ?, ?)", buildingId, unitId, unitName, time.Now()).Error
		if err != nil {
			global.GVA_LOG.Error("插入失败", zap.Error(err))
		}
	} else if strings.Contains(line, "encryptBuildingName") {
		arr := strings.Split(strings.Replace(line, " ", "", -1), ",")
		buildingName := strings.Split(arr[2], ":")[1]
		buildingId := strings.Split(arr[1], ":")[1]
		xiaoquId := strings.Split(arr[0], ":")[1]
		err := global.GVA_DB.Exec("insert into dict_building (community_id, building_open_id, encrypt_building_name, created_at) values (?, ?, ?, ?)", xiaoquId, buildingId, buildingName, time.Now()).Error
		if err != nil {
			global.GVA_LOG.Error("插入失败", zap.Error(err))
		}
	}
}
