package main

import (
	"bufio"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/core"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/initialize"
	"github.com/flipped-aurora/gin-vue-admin/server/model/house"
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
	filename := "/opt/zhaofang/update.csv"
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

	unitId := arr[0]
	houseName := arr[4]
	houseId := ""
	var count int64

	result := global.GVA_DB.Model(&house.DictHouse{}).Where("unit_open_id = ? and encrypt_house_name = ?", unitId, houseName).Count(&count)
	if result.Error != nil {
		fmt.Printf("查询错误: %v\n", result.Error)
		return
	}
	if count == 0 {
		err := global.GVA_DB.Exec("insert into dict_house (house_open_id, encrypt_house_name, unit_open_id, created_at) values (?, ?, ?, ?)", houseId, houseName, unitId, time.Now()).Error
		if err != nil {
			global.GVA_LOG.Error("插入失败", zap.Error(err))
		}
	} else {
		fmt.Printf("unitId: %s, houseName: %s 找到 %d 条 存在的记录\n", unitId, houseName, count)
	}
}
