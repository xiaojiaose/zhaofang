<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item label="手机号">
          <el-input v-model="searchInfo.phone" placeholder="手机号" />
        </el-form-item>
        <el-form-item label="绑定状态">
          <el-select v-model="searchInfo.bind" placeholder="请选择绑定状态">
            <el-option
              v-for="item in bindStatusOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="handleSearch"> 查询 </el-button>
          <el-button icon="refresh" @click="onReset"> 重置 </el-button>
        </el-form-item>
      </el-form>
    </div>
    <el-divider />
    <div class="gva-search-box">
      <el-form ref="searchForm" :inline="true" :model="addSalerInfo">
        <el-form-item label="手机号">
          <el-input v-model="addSalerInfo.phone" placeholder="手机号" />
        </el-form-item>
        <el-form-item label="称呼">
          <el-input v-model="addSalerInfo.userName" placeholder="称呼" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="handleAddSaler">
            新增
          </el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <el-table :data="tableData" row-key="ID">
        <el-table-column align="left" label="编号" min-width="100" prop="ID" />
        <el-table-column align="left" label="绑定手机号" min-width="150" prop="phone" />
        <el-table-column align="left" label="操作" min-width="150">
          <template #default="scope">
            <el-button
              type="success"
              v-if="scope.row.enable === 1"
              @click="handleEnable(scope.row, 2)"
              >启用中</el-button
            >
            <el-button
              type="info"
              v-if="scope.row.enable === 2"
              @click="handleEnable(scope.row, 1)"
              >停用中</el-button
            >
          </template>
        </el-table-column>

        <el-table-column align="left" label="创建时间" min-width="150" prop="CreatedAt">
          <template #default="scope">
            {{ dayjs(scope.row.CreatedAt).format("YYYY-MM-DD HH:mm") }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="绑定时间" min-width="180" prop="CreatedAt">
          <template #default="scope">
            {{ dayjs(scope.row.CreatedAt).format("YYYY-MM-DD HH:mm") }}
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="绑定的微信号"
          min-width="180"
          prop="wxNo"
        >
          <template #default="scope"> </template>
        </el-table-column>
        <el-table-column
          label="称呼"
          :min-width="appStore.operateMinWith"
          fixed="right"
          prop="userName"
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
  </div>
</template>

<script setup>
import { ref } from "vue";
import { useAppStore } from "@/pinia";
import dayjs from "dayjs";
import { ElMessage, ElMessageBox } from "element-plus";

import { getSaler, addSaler, setUserInfo } from "@/api/user";

defineOptions({
  name: "Saler",
});

const appStore = useAppStore();

const bindStatusOptions = [
  {
    value: 2,
    label: "未绑定",
  },
  {
    value: 1,
    label: "已绑定",
  },
];

//搜索
const searchInfo = ref({
  phone: "",
});
const handleSearch = async () => {
  getSalerList();
};
const onReset = () => {
  searchInfo.value = {
    phone: "",
  };
  getSalerList();
};

//添加经纪人
const addSalerInfo = ref({
  phone: "",
  userName: "",
});
const handleAddSaler = async () => {
  const res = await addSaler(addSalerInfo.value);
  if (res.code === 0) {
    ElMessage.success("新增成功");
    getSalerList();
  }
};

const page = ref(1);
const total = ref(0);
const pageSize = ref(10);
const tableData = ref([]);

//获取经纪人列表
const getSalerList = async () => {
  const res = await getSaler({
    page: page.value,
    pageSize: pageSize.value,
    phone: searchInfo.value.phone,
    bind: searchInfo.value.bind,
  });

  if (res.code === 0) {
    tableData.value = res.data.list;
    total.value = res.data.total;
  }
};

// 启用停用
const handleEnable = async (row, enable) => {
  ElMessageBox.confirm("确定要" + (enable === 1 ? "启用" : "停用") + "该经纪人吗?", "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  }).then(async () => {
    const res = await setUserInfo({
    ID: row.ID,
    enable: enable,
  });
  if (res.code === 0) {
    ElMessage.success("操作成功");
    getSalerList();
  }
  });
};

const initPage = async () => {
  getSalerList();
};

initPage();


// 分页
const handleSizeChange = (val) => {
  pageSize.value = val;
  getSalerList();
};

const handleCurrentChange = (val) => {
  page.value = val;
  getSalerList();
};

</script>
<style lang="scss" scoped></style>
