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
      <el-form :inline="true" :model="form">
        <el-form-item label="小区名称">
          <el-input v-model="form.buildingName" clearable placeholder="请输入小区名称" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleGrant">新增</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <el-table :data="tableData" border v-loading="loading">
        <el-table-column type="selection" width="55" />
        <el-table-column align="center" label="编号" prop="ID" width="80" />
        <el-table-column align="left" label="小区名称" prop="name" min-width="200" />
        <el-table-column align="left" label="所属区域" prop="area" min-width="120" />
        <el-table-column align="left" label="所属商圈" prop="districts" min-width="150" />
        <el-table-column align="left" label="小区坐标" prop="latitude" min-width="100" />
        <el-table-column align="center" label="房间数据" prop="roomData" width="200" />
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
          prop="lastOperationTime"
          width="200"
        />
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

    <el-dialog
      :title="dialogTitle"
      v-model="dialogVisible"
      width="700px"
      :close-on-click-modal="false"
      @close="closeDialog"
    >
      <el-form ref="form" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="小区名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入小区名称" />
        </el-form-item>
        <el-form-item label="城市" prop="city">
          <el-input v-model="form.city" placeholder="请输入城市" />
        </el-form-item>
        <el-form-item label="区域" prop="area">
          <el-input v-model="form.area" placeholder="请输入区域" />
        </el-form-item>
        <el-form-item label="位置" prop="position">
          <el-input v-model="form.position" placeholder="请输入位置" />
        </el-form-item>
        <el-form-item label="行政区" prop="districts">
          <el-input v-model="form.districts" placeholder="请输入行政区" />
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
          <el-switch v-model="form.able" active-text="启用" inactive-text="禁用" />
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="closeDialog">取 消</el-button>
        <el-button type="primary" @click="submitForm">确 定</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref } from "vue";
import { ElMessage } from "element-plus";
import {
  getDictBuildingList,
  getDictBuildingInfo,
  editXiaoqu as editXiaoquApi,
  addDictBuilding,
} from "@/api/dictBuilding";

defineOptions({
  name: "DictBuilding",
});

const defaultForm = () => ({
  ID: undefined,
  buildingName: "",
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
    // const _res = await getDictBuildingList({
    //   page: page.value,
    //   pageSize: pageSize.value,
    //   ...searchInfo.value,
    // });

    const res = {
      code: 0,
      data: {
        list: [
          {
            ID: 1,
            buildingName: "小区1",
            city: "城市1",
            area: "区域1",
            position: "位置1",
            districts: "行政区1",
            address: "地址1",
            latitude: "纬度1",
            longitude: "经度1",
            able: true,
            createTime: "2023-01-01 00:00:00",
            updateTime: "2023-01-01 00:00:00",
          },
        ],
        page: 1,
        pageSize: 10,
        total: 100,
      },
      msg: "success",
    };
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
  addDictBuilding(form.value).then((res) => {
    if (res.code === 0) {
      ElMessage.success("添加成功");
      form.value = defaultForm();
      getTableData();
    }
  });
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
  dialogTitle.value = "编辑小区";
  dialogVisible.value = true;
  // getDictBuildingInfo(row.ID).then((res) => {
  //   if (res.code === 0) {
  //     form.value = { ...res.data };
  //     dialogVisible.value = true;
  //   }
  // });
};

const closeDialog = () => {
  dialogVisible.value = false;
  form.value = defaultForm();
};

const submitForm = () => {
  editXiaoqu(form.value).then((res) => {
    if (res.code === 0) {
      ElMessage.success(dialogTitle.value === "新增小区" ? "新增成功" : "编辑成功");
      closeDialog();
      getTableData();
    }
  });
};

getTableData();
</script>

<style lang="scss" scoped></style>
