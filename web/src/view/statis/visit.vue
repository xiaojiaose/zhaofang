<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item label="经纪人手机号">
          <el-input v-model="searchInfo.phone" placeholder="经纪人手机号" />
        </el-form-item>
        <el-form-item label="微信号">
          <el-input v-model="searchInfo.wxNo" placeholder="微信号" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSearch"> 查询 </el-button>
          <el-button icon="refresh" @click="onReset"> 重置 </el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <el-table :data="tableData" row-key="ID">
        <el-table-column align="left" label="访问手机号" min-width="150" prop="phone" />
        <el-table-column align="left" label="访问微信号" min-width="150" prop="wxNo" />
        <el-table-column align="left" label="最近访问时间" min-width="150" prop="lastVisitTime">
          <template #default="scope">
            {{ dayjs(scope.row.lastVisitTime).format("YYYY-MM-DD HH:mm") }}
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
  </div>
</template>

<script setup>
import { ref } from "vue";
import dayjs from "dayjs";

import { getVisitList } from "@/api/statis";


const searchInfo = ref({
  phone: "",
  wxNo: "",
});
const tableData = ref([]);
const total = ref(0);
const page = ref(1);
const pageSize = ref(10);


const onSearch = () => {
  console.log(searchInfo.value);
};

const onReset = () => {
  searchInfo.value = {
    phone: "",
    wxNo: "",
  };
};

const getTableData = async () => {
  const res = await getVisitList({
    page: page.value,
    pageSize: pageSize.value,
    ...searchInfo.value,
  });
  if (res.code === 0) {
    tableData.value = res.data.list;
    total.value = res.data.total;
  }
};

const handleCurrentChange = (val) => {
  page.value = val;
  getTableData();
}

const handleSizeChange = (val) => {
  pageSize.value = val;
  getTableData();
}


getTableData();


</script>
<style lang="scss" scoped></style>
