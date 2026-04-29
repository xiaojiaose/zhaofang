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
        <el-form-item label="上架和录入房源权限">
          <el-select v-model="searchInfo.isPublish" placeholder="请选择上架和录入房源权限">
            <el-option
              v-for="item in publishOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="找房超市">
          <el-select
            v-model="searchInfo.isFindHouseSupermarket"
            placeholder="请选择找房超市"
          >
            <el-option
              v-for="item in findStoreOptions"
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
        <el-form-item label="用户名">
          <el-input v-model="addSalerInfo.userName" placeholder="用户名" />
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
        <el-table-column align="left" label="绑定的微信号" min-width="180" prop="wxNo">
          <template #default="scope">
            <!-- <img :src="scope.row.headerImg" alt="" />
            <span>{{ scope.row.wxNo }}</span> -->
            <div class="flex flex-col space-y-2 py-2">
              <el-image :src="scope.row.headerImg" fit="cover" class="w-16 h-16 rounded-md shadow-sm">
                <template #error>
                  <div class="image-viewer-slot image-slot">
                    <el-icon><icon-picture /></el-icon>
                  </div>
                </template>
              </el-image>
              <div class="text-left" v-if="scope.row.wxNo">
                <div
                  class="px-2 py-0.5 mt-1 bg-gray-100 rounded text-[10px] font-mono text-gray-500"
                >
                  {{ scope.row.wxNo }}
                </div>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="微信昵称" min-width="150" prop="wxNickName" />
        <el-table-column
          label="用户名"
          :min-width="appStore.operateMinWith"
          prop="userName"
        />
        <el-table-column
          label="编辑用户"
          :min-width="appStore.operateMinWith"
          fixed="right"
        >
          <template #default="scope">
            <el-button type="primary" size="mini" @click="handleEdit(scope.row)"
              >编辑</el-button
            >
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
      v-model="editSalerDialog"
      :size="appStore.drawerSize"
      :show-close="false"
      :close-on-press-escape="false"
      :close-on-click-modal="false"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">经纪人</span>
          <div>
            <el-button @click="editSalerDialog = false">取 消</el-button>
            <el-button type="primary" @click="enterEditSalerDialog(editSalerFormRef)"
              >保存</el-button
            >
          </div>
        </div>
      </template>

      <el-form
        ref="editSalerFormRef"
        :rules="rules"
        :model="editSalerForm"
        label-width="150px"
      >
        <el-form-item label="绑定手机号" prop="phone">
          <el-input v-model="editSalerForm.phone" disabled />
        </el-form-item>
        <el-form-item label="用户名" prop="userName">
          <el-input v-model="editSalerForm.userName" />
        </el-form-item>
        <el-form-item label="微信号" prop="wxNo">
          <el-input v-model="editSalerForm.wxNo" />
        </el-form-item>
        <el-form-item label="微信昵称" prop="wxNickName">
          <el-input v-model="editSalerForm.wxNickName" />
        </el-form-item>
        <el-form-item label="头像" prop="headerImg">
          <SelectImage v-model="editSalerForm.headerImg" />
        </el-form-item>
        <el-form-item label="上架和录入房源权限" prop="isPublish">
          <el-switch v-model="editSalerForm.isPublish" />
        </el-form-item>
        <el-form-item label="上架额度" v-if="editSalerForm.isPublish === true">
          <el-input-number v-model="editSalerForm.publishQuotaTotal" :min="0" />
        </el-form-item>
        <el-form-item label="找房超市">
          <el-switch v-model="editSalerForm.isFindHouseSupermarket" />
        </el-form-item>
        <el-form-item label="联系次数">
          <el-input-number v-model="editSalerForm.contactViewQuotaTotal" :min="0" />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref } from "vue";
import { useAppStore } from "@/pinia";
import dayjs from "dayjs";
import { ElMessage, ElMessageBox } from "element-plus";
import { Picture as IconPicture } from '@element-plus/icons-vue'
import SelectImage from "@/components/selectImage/selectImage.vue";

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

const publishOptions = [
  {
    value: 1,
    label: "有",
  },
  {
    value: 2,
    label: "没有",
  },
];

const findStoreOptions = [
  {
    value: 1,
    label: "找房超市",
  },
  {
    value: 2,
    label: "否",
  },
];

//搜索
const searchInfo = ref({
  phone: ""
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
  const res = await addSaler({
    ...addSalerInfo.value
  });
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
    isFindHouseSupermarket: searchInfo.value.isFindHouseSupermarket,
    isPublish: searchInfo.value.isPublish,
  });

  if (res.code === 0) {
    tableData.value = res.data.list;
    total.value = res.data.total;
  }
};

// 启用停用
const handleEnable = async (row, enable) => {
  ElMessageBox.confirm(
    "确定要" + (enable === 1 ? "启用" : "停用") + "该经纪人吗?",
    "提示",
    {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning",
    }
  ).then(async () => {
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

// 编辑经纪人相关
const editSalerDialog = ref(false);
const editSalerForm = ref();
const editSalerFormRef = ref();

const handleEdit = (row) => {
  editSalerDialog.value = true;
  editSalerForm.value = { ...row, isPublish: row.publishQuotaTotal > 0 };
};

const enterEditSalerDialog = async (formEl) => {
  console.log(formEl);
  if (!formEl) return;
  await formEl.validate(async (valid, fields) => {
    if (valid) {
      const values = { ...editSalerForm.value };
      console.log(values);
      setUserInfo({
        ...values,
        ID: values.ID,
        wxNo: values.wxNo,
        phone: values.phone,
        userName: values.userName,
        enable: values.enable,
        isFindHouseSupermarket: values.isFindHouseSupermarket,
        contactViewQuotaTotal: values.contactViewQuotaTotal,
        publishQuotaTotal: values.isPublish ? values.publishQuotaTotal : 0,
        isPublish: values.isPublish
      }).then((res) => {
        if (res.code === 0) {
          ElMessage.success("编辑成功");
          editSalerDialog.value = false;
          getSalerList();
        }
      });
    }
  });
};

// 编辑经纪人表单相关
const rules = ref({
  phone: [{ required: true, message: "请输入手机号", trigger: "change" }],
  userName: [{ required: true, message: "请输入称呼", trigger: "change" }],
  wxNo: [{ required: true, message: "请输入微信号", trigger: "change" }],
  enable: [{ required: true, message: "请选择上架和录入房源权限", trigger: "change" }],
  isFindHouseSupermarket: [
    { required: true, message: "请选择找房超市", trigger: "change" },
  ],
});
</script>
<style lang="scss" scoped></style>
