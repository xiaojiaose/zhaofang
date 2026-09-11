<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item label="小区名称">
          <el-input v-model="searchInfo.keyword" placeholder="小区名称" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSearch"> 查询 </el-button>
          <el-button icon="refresh" @click="onReset"> 重置 </el-button>
        </el-form-item>
      </el-form>
      <el-divider />
      <!-- <el-form :inline="true" :model="form">
        <el-form-item label="小区名称">
          <el-input v-model="form.name" clearable placeholder="请输入小区名称" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleGrant">新增</el-button>
        </el-form-item>
      </el-form> -->
      <el-button type="primary" @click="handleGrant">新增楼盘</el-button>
    </div>
    <div class="gva-table-box">
      <el-table :data="tableData" border v-loading="loading">
        <el-table-column type="selection" width="55" />
        <el-table-column align="center" label="编号" prop="ID" width="80" />
        <el-table-column align="left" label="小区名称" prop="name" min-width="200" />
        <el-table-column align="left" label="所属区域" prop="area" min-width="120" />
        <el-table-column align="left" label="所属商圈" prop="districts" min-width="150" />
        <el-table-column align="left" label="小区坐标" prop="latitude" min-width="200" >
          <template #default="scope">
            {{ scope.row.latitude }}, {{ scope.row.longitude }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="房间数据" prop="roomData" width="200">
          <template #default="scope">
            <el-button type="primary" size="small" link @click="editRoomData(scope.row)">
              管理房间数据
            </el-button>
          </template>
        </el-table-column>
        <el-table-column label="操作" align="center" width="200">
          <template #default="scope">
            <el-button type="primary" size="small" @click="editDictBuilding(scope.row)"
              >编辑</el-button
            >
          </template>
        </el-table-column>
        <el-table-column
          align="center"
          label="最后操作时间"
          prop="UpdatedAt"
          width="200"
        >
          <template #default="scope">
            {{ formatDate(scope.row.UpdatedAt) }}
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <el-drawer
      :title="dialogTitle"
      v-model="dialogVisible"
      :size="appStore.drawerSize"
      :close-on-click-modal="false"
      @close="closeDialog"
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="小区名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入小区名称" />
        </el-form-item>
        <!-- <el-form-item label="城市" prop="city">
          <el-input v-model="form.city" placeholder="请输入城市" />
        </el-form-item> -->
        <el-form-item label="区域" prop="area">
          <el-select v-model="form.area" placeholder="请选择区域">
            <el-option v-for="item in areaOptions" :key="item.id" :label="item.name" :value="item.name" />
          </el-select>
        </el-form-item>
        <el-form-item label="商圈" prop="districts">
          <el-select v-model="form.districts" placeholder="请选择商圈">
            <el-option v-for="item in (areaOptions.find(item => item.name === form.area)?.Districts || [])" :key="item.id" :label="item.name" :value="item.name" />
          </el-select>
        </el-form-item>
        <el-form-item label="位置" prop="position">
          <el-input v-model="form.position" placeholder="请输入位置" />
        </el-form-item>
        <el-form-item label="地址" prop="address">
          <el-input
            v-model="form.address"
            type="textarea"
            :rows="2"
            placeholder="请输入地址"
          />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="纬度" prop="latitude">
              <el-input v-model="form.latitude" placeholder="请输入纬度" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="经度" prop="longitude">
              <el-input v-model="form.longitude" placeholder="请输入经度" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="是否启用" prop="able">
          <el-switch v-model="form.able" />
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="closeDialog">取 消</el-button>
        <el-button type="primary" @click="submitForm">确 定</el-button>
      </span>
    </el-drawer>

    <el-drawer
      v-model="roomDataDialogVisible"
      :title="`管理房间数据 - ${currentXiaoqu?.name || ''}`"
      :size="appStore.drawerSize"
      :close-on-click-modal="false"
      @close="closeRoomDataDialog"
    >
      <div class="room-data-container">
        <div class="room-data-header">
          <el-button type="primary" @click="addBuilding">+ 添加楼栋</el-button>
        </div>

        <el-collapse accordion expand-icon-position="left">
          <el-collapse-item
            v-for="(building, buildingIndex) in dictTree.buildings"
            :key="building.id"
            :name="buildingIndex"
          >
            <template #title>
              <div class="collapse-title">
                <span class="building-name">{{ building.name }}号楼</span>
                <el-button
                  type="primary"
                  size="small"
                  link
                  @click.stop="editBuilding(building)"
                  class="ml-2"
                >
                  编辑
                </el-button>
                <el-button
                  type="primary"
                  size="small"
                  link
                  @click.stop="addUnit(building)"
                  class="ml-2"
                >
                  + 添加单元
                </el-button>
                <el-popconfirm
                  title="确定删除该楼栋及其下所有单元和房间吗？"
                  confirm-button-text="确定"
                  cancel-button-text="取消"
                  @confirm="deleteBuilding(building)"
                >
                  <template #reference>
                    <el-button
                      type="danger"
                      size="small"
                      link
                      @click.stop
                    >
                      删除楼栋
                    </el-button>
                  </template>
                </el-popconfirm>
              </div>
            </template>

            <div class="unit-container">
              <el-collapse accordion expand-icon-position="left">
                <el-collapse-item
                  v-for="(unit, unitIndex) in building.units"
                  :key="unit.id"
                  :name="unitIndex"
                >
                  <template #title>
                    <div class="collapse-title">
                      <span class="unit-name">{{ unit.name }}单元</span>
                      <el-button
                        type="primary"
                        size="small"
                        link
                        @click.stop="editUnit(building, unit)"
                        class="ml-2"
                      >
                        编辑
                      </el-button>
                      <el-button
                        type="primary"
                        size="small"
                        link
                        @click.stop="addHouse(building, unit)"
                        class="ml-2"
                      >
                        + 添加房间
                      </el-button>
                      <el-popconfirm
                        title="确定删除该单元及其下所有房间吗？"
                        confirm-button-text="确定"
                        cancel-button-text="取消"
                        @confirm="deleteUnit(unit)"
                      >
                        <template #reference>
                          <el-button
                            type="danger"
                            size="small"
                            link
                            @click.stop
                          >
                            删除单元
                          </el-button>
                        </template>
                      </el-popconfirm>
                    </div>
                  </template>

                  <div class="house-list">
                    <div
                      v-for="house in unit.houses"
                      :key="house.id"
                      class="house-item"
                    >
                      <span class="house-name">{{ house.name }}室</span>
                      <el-button
                        type="primary"
                        size="small"
                        link
                        @click="editHouse(house)"
                      >
                        编辑
                      </el-button>
                      <el-popconfirm
                        title="确定删除该房间吗？"
                        confirm-button-text="确定"
                        cancel-button-text="取消"
                        @confirm="deleteHouse(house)"
                      >
                        <template #reference>
                          <el-button
                            type="danger"
                            size="small"
                            link
                          >
                            删除
                          </el-button>
                        </template>
                      </el-popconfirm>
                    </div>
                    <el-empty v-if="!unit.houses?.length" description="暂无房间" :image-size="60" />
                  </div>
                </el-collapse-item>
              </el-collapse>
              <el-empty v-if="!building.units?.length" description="暂无单元，请点击上方按钮添加" :image-size="60" />
            </div>
          </el-collapse-item>
        </el-collapse>

        <el-empty v-if="!dictTree.buildings?.length" description="暂无楼栋数据，请点击上方按钮添加" :image-size="100" />
      </div>
    </el-drawer>

    <el-dialog
      v-model="editDialogVisible"
      :title="editDialogTitle"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form ref="editFormRef" :model="editForm" :rules="editRules" label-width="100px">
        <el-form-item :label="editFormLabel" prop="name">
          <el-input
            v-model="editForm.name"
            type="textarea"
            :rows="4"
            :placeholder="getBatchPlaceholder()"
          />
          <div class="batch-tip">
            <el-text type="info" size="small">
              💡 支持批量添加，每行一个或用逗号分隔，例如：{{ getBatchExample() }}
            </el-text>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="submitEdit">确 定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive } from "vue";
import { ElMessage } from "element-plus";
import { useAppStore } from "@/pinia";
import {formatDate} from '@/utils/format'
import {
  getxiaoquList,
  getDictBuildingInfo,
  addOrEditDictBuilding,
  getXiaoquDictTree,
  upsertXiaoquDict,
  deleteXiaoquDict,
} from "@/api/dictBuilding";
import { getAreaOptions } from "@/api/center";

defineOptions({
  name: "DictBuilding",
});

const appStore = useAppStore();

const defaultForm = () => ({
  ID: 0,
  name: "",
  city: undefined,
  area: undefined,
  position: undefined,
  districts: undefined,
  address: undefined,
  latitude: undefined,
  longitude: undefined,
  able: undefined,
});

const form = ref(defaultForm());
const formRef = ref(null);
const searchInfo = ref({
  keyword: "",
});
const tableData = ref([]);
const page = ref(1);
const pageSize = ref(10);
const total = ref(0);
const loading = ref(false);
const dialogVisible = ref(false);
const dialogTitle = ref("新增小区");

const rules = {
  name: [{ required: true, message: "请输入小区名称", trigger: "blur" }],
};

const getTableData = async () => {
  loading.value = true;
  try {
    const res = await getxiaoquList({
      page: page.value,
      pageSize: pageSize.value,
      ...searchInfo.value,
      cityId: "1",
    });

    if (res.code === 0) {
      tableData.value = res.data.list;
      total.value = res.data.total;
    }
  } finally {
    loading.value = false;
  }
};

const onSearch = () => {
  page.value = 1;
  getTableData();
};

const onReset = () => {
  searchInfo.value.keyword = "";
  page.value = 1;
  getTableData();
};

const handleGrant = () => {
  form.value = defaultForm();
  dialogTitle.value = "新增小区";
  dialogVisible.value = true;
};

const handleCurrentChange = (val) => {
  page.value = val;
  getTableData();
};

const handleSizeChange = (val) => {
  pageSize.value = val;
  getTableData();
};

const editDictBuilding = (row) => {
  form.value = { ...row };
  console.log(form.value)
  dialogTitle.value = "编辑小区";
  dialogVisible.value = true;
};

const closeDialog = () => {
  dialogVisible.value = false;
  form.value = defaultForm();
};

const submitForm = () => {
  addOrEditDictBuilding({...form.value, ID: form.value.ID || 0, cityId: "1"}).then((res) => {
    if (res.code === 0) {
      ElMessage.success(dialogTitle.value === "新增小区" ? "新增成功" : "编辑成功");
      closeDialog();
      getTableData();
    }
  });
};

const areaOptions = ref([]);
const getArea = async () => {
  getAreaOptions().then((res) => {
    if (res.code === 0) {
      Object.entries(res.data || {}).forEach(([areaName, areaItem]) => {
        areaOptions.value.push({ ...areaItem })
      });
      console.log(areaOptions.value)
    }
  });
}

getTableData();
getArea();

const roomDataDialogVisible = ref(false);
const currentXiaoqu = ref(null);
const dictTree = reactive({
  buildings: []
});
const activeBuildings = ref([]);
const activeUnits = reactive({});

const editDialogVisible = ref(false);
const editDialogTitle = ref("");
const editFormLabel = ref("");
const editFormRef = ref(null);
const editForm = reactive({
  id: undefined,
  name: "",
  type: "", // 'building' | 'unit' | 'house'
  parentId: undefined,
  parentOpenId: undefined,
  buildingOpenId: undefined
});
const editRules = {
  name: [{ required: true, message: "请输入名称", trigger: "blur" }]
};

const loadDictTree = async (row) => {
  try {
    const res = await getXiaoquDictTree(row.ID);
    if (res.code === 0) {
      currentXiaoqu.value = row;
      roomDataDialogVisible.value = true;
      dictTree.buildings = res.data.buildings || [];
    }
  } catch (error) {
    console.error("加载字典树失败:", error);
    ElMessage.error("加载房间数据失败");
  }
};

const editRoomData = (row) => {
  loadDictTree(row);
};

const closeRoomDataDialog = () => {
  roomDataDialogVisible.value = false;
  currentXiaoqu.value = null;
  dictTree.buildings = [];
  activeBuildings.value = [];
  Object.keys(activeUnits).forEach(key => delete activeUnits[key]);
};

const resetEditForm = () => {
  editForm.id = undefined;
  editForm.name = "";
  editForm.type = "";
  editForm.parentId = undefined;
  editForm.parentOpenId = undefined;
  editForm.buildingOpenId = undefined;
};

const addBuilding = () => {
  resetEditForm();
  editForm.type = "building";
  editDialogTitle.value = "添加楼栋";
  editFormLabel.value = "楼栋名称";
  editDialogVisible.value = true;
};

const addUnit = (building) => {
  resetEditForm();
  editForm.type = "unit";
  editForm.parentId = building.id;
  editForm.parentOpenId = building.buildingOpenId;
  editDialogTitle.value = "添加单元";
  editFormLabel.value = "单元名称";
  editDialogVisible.value = true;
};

const addHouse = (building, unit) => {
  resetEditForm();
  editForm.type = "house";
  editForm.parentId = unit.id;
  editForm.parentOpenId = unit.unitOpenId;
  editDialogTitle.value = "添加房间";
  editFormLabel.value = "房间号";
  editDialogVisible.value = true;
};

const editBuilding = (building) => {
  resetEditForm();
  editForm.id = building.id;
  editForm.name = building.name;
  editForm.type = "building";
  editForm.parentOpenId = building.buildingOpenId;
  editDialogTitle.value = "编辑楼栋";
  editFormLabel.value = "楼栋名称";
  editDialogVisible.value = true;
};

const editUnit = (building, unit) => {
  resetEditForm();
  editForm.id = unit.id;
  editForm.name = unit.name;
  editForm.type = "unit";
  editForm.parentId = building.id;
  editForm.parentOpenId = unit.unitOpenId;
  editForm.buildingOpenId = building.buildingOpenId;
  editDialogTitle.value = "编辑单元";
  editFormLabel.value = "单元名称";
  editDialogVisible.value = true;
};

const getBatchPlaceholder = () => {
  if (editForm.type === "building") {
    return "请输入楼栋名称，支持批量添加\n例如：1,2,3 或每行一个";
  } else if (editForm.type === "unit") {
    return "请输入单元名称，支持批量添加\n例如：A,B,C 或每行一个";
  } else {
    return "请输入房间号，支持批量添加\n例如：101,102,103 或每行一个";
  }
};

const getBatchExample = () => {
  if (editForm.type === "building") {
    return "1, 2, 3";
  } else if (editForm.type === "unit") {
    return "A单元, B单元, C单元";
  } else {
    return "101, 102, 103";
  }
};

const getBatchTypeName = () => {
  if (editForm.type === "building") {
    return "楼栋";
  } else if (editForm.type === "unit") {
    return "单元";
  } else {
    return "房间";
  }
};

const editHouse = (house) => {
  resetEditForm();
  editForm.id = house.id;
  editForm.name = house.name;
  editForm.type = "house";
  editForm.parentOpenId = house.houseOpenId;
  editDialogTitle.value = "编辑房间";
  editFormLabel.value = "房间号";
  editDialogVisible.value = true;
};

const submitEdit = async () => {
  if (!editFormRef.value) return;

  try {
    await editFormRef.value.validate();

    const data = {
      xiaoquId: currentXiaoqu.value.ID
    };

    const parseNames = (input) => {
      return input
        .split(/[,，\n]/)
        .map(name => name.trim())
        .filter(name => name.length > 0);
    };

    const names = parseNames(editForm.name);

    if (editForm.id) {
      if (editForm.type === "building") {
        data.buildingOps = [{ id: editForm.id, name: names[0] }];
      } else if (editForm.type === "unit") {
        data.unitOps = [{ id: editForm.id, name: names[0], buildingOpenId: editForm.buildingOpenId }];
      } else if (editForm.type === "house") {
        data.houseOps = [{ id: editForm.id, name: names[0] }];
      }
    } else {
      if (editForm.type === "building") {
        data.buildingOps = names.map(name => ({ name }));
      } else if (editForm.type === "unit") {
        data.unitOps = names.map(name => ({
          buildingOpenId: editForm.parentOpenId,
          name
        }));
      } else if (editForm.type === "house") {
        data.houseOps = names.map(name => ({
          unitOpenId: editForm.parentOpenId,
          name
        }));
      }
    }

    const res = await upsertXiaoquDict(data);

    if (res.code === 0) {
      ElMessage.success(editForm.id ? "修改成功" : `成功添加 ${names.length} 个${getBatchTypeName()}`);
      editDialogVisible.value = false;

      if (res.data?.buildings) {
        dictTree.buildings = res.data.buildings;
      } else {
        await loadDictTree(currentXiaoqu.value.ID);
      }
    }
  } catch (error) {
    if (error !== false) {
      console.error("提交编辑失败:", error);
      ElMessage.error("操作失败");
    }
  }
};

const deleteBuilding = async (building) => {
  try {
    const res = await deleteXiaoquDict({
      xiaoquId: currentXiaoqu.value.ID,
      buildingIds: [building.id]
    });

    if (res.code === 0) {
      ElMessage.success("删除成功");

      if (res.data?.buildings) {
        dictTree.buildings = res.data.buildings;
      } else {
        await loadDictTree(currentXiaoqu.value.ID);
      }
    } else {
      ElMessage.error(res.msg || "删除失败");
    }
  } catch (error) {
    console.error("删除楼栋失败:", error);
    ElMessage.error("删除失败");
  }
};

const deleteUnit = async (unit) => {
  try {
    const res = await deleteXiaoquDict({
      xiaoquId: currentXiaoqu.value.ID,
      unitIds: [unit.id]
    });

    if (res.code === 0) {
      ElMessage.success("删除成功");

      if (res.data?.buildings) {
        dictTree.buildings = res.data.buildings;
      } else {
        await loadDictTree(currentXiaoqu.value.ID);
      }
    } else {
      ElMessage.error(res.msg || "删除失败");
    }
  } catch (error) {
    console.error("删除单元失败:", error);
    ElMessage.error("删除失败");
  }
};

const deleteHouse = async (house) => {
  try {
    const res = await deleteXiaoquDict({
      xiaoquId: currentXiaoqu.value.ID,
      houseIds: [house.id]
    });

    if (res.code === 0) {
      ElMessage.success("删除成功");

      if (res.data?.buildings) {
        dictTree.buildings = res.data.buildings;
      } else {
        await loadDictTree(currentXiaoqu.value.ID);
      }
    } else {
      ElMessage.error(res.msg || "删除失败");
    }
  } catch (error) {
    console.error("删除房间失败:", error);
    ElMessage.error("删除失败");
  }
};
</script>

<style lang="scss" scoped>
.room-data-container {
  overflow-y: auto;

  .room-data-header {
    margin-bottom: 20px;
    display: flex;
    justify-content: flex-start;
  }

  .collapse-title {
    display: flex;
    align-items: center;
    gap: 12px;
    width: 100%;
    padding-right: 20px;

    .building-name,
    .unit-name {
      color: #303133;
    }

    .unit-count,
    .house-count {
      color: #909399;
      background-color: #f4f4f5;
      padding: 2px 8px;
      border-radius: 4px;
    }
  }

  .unit-container {
    padding-left: 20px;
  }

  .house-list {
    padding: 10px 0;
    display: flex;
    flex-wrap: wrap;
    gap: 10px;

    .house-item {
      display: flex;
      align-items: center;
      gap: 8px;
      padding: 8px 12px;
      background-color: #f5f7fa;
      border-radius: 6px;
      border: 1px solid #e4e7ed;

      .house-name {
        font-size: 14px;
        color: #606266;
      }

      &:hover {
        background-color: #ecf5ff;
        border-color: #b3d8ff;
      }
    }
  }
}

.ml-2 {
  margin-left: 8px;
}

.batch-tip {
  margin-top: 4px;
  line-height: 1.4;
}
</style>
