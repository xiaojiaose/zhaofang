<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item label="小区名称">
          <el-select
            v-model="searchInfo.xiaoquId"
            filterable
            remote
            reserve-keyword
            placeholder="小区名称"
            :remote-method="remoteSearchXiaoqu"
            :loading="searchXiaoquLoading"
            style="width: 240px"
          >
            <el-option
              v-for="item in xiaoquOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="户室号">
          <el-input v-model="searchInfo.keyword" placeholder="户室号" />
        </el-form-item>
        <el-form-item label="经纪人手机号">
          <el-input v-model="searchInfo.phone" placeholder="经纪人手机号" />
        </el-form-item>
        <el-form-item label="审核状态">
          <el-select v-model="searchInfo.approvalStatus" placeholder="请选择审核状态">
            <el-option
              v-for="item in approvalStatusOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit"> 查询 </el-button>
          <el-button icon="refresh" @click="onReset"> 重置 </el-button>
        </el-form-item>
      </el-form>
    </div>
    <el-divider />
    <div class="gva-table-box">
      <el-space wrap>
        <el-button type="primary" @click="handleBatchPass">批量通过</el-button>
        <el-button type="primary" @click="handleBatchDown">批量下架</el-button>
      </el-space>
    </div>
    <div class="gva-table-box">
      <el-table :data="tableData" row-key="ID" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="房源信息" min-width="300">
          <template #default="scope">
            <el-image
              :src="scope.row.attachments?.house[0]?.url"
              fit="cover"
              style="width: 300px; height: 150px"
            >
              <template #error>
                <div class="image-slot">
                  <el-icon><icon-picture /></el-icon>
                </div>
              </template>
            </el-image>
          </template>
        </el-table-column>
        <el-table-column align="left" label="" min-width="150">
          <template #default="scope">
            <el-text size="large" tag="b">{{ scope.row.xiaoqu }}</el-text>
            <div>{{ scope.row.door_no }}</div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="出租类型" min-width="150">
          <template #default="scope">
            {{ scope.row.rent_type }} {{ scope.row.house_type }}
          </template>
        </el-table-column>

        <el-table-column align="left" label="平台推广" min-width="150">
          <template #default="scope">
            <el-text type="danger">{{ scope.row.price }}元/月</el-text>
          </template>
        </el-table-column>
        <el-table-column align="left" label="状态" min-width="180">
          <template #default="scope">
            <el-tag
              type="info"
              v-if="scope.row.status === '已下架'"
              size="large"
              effect="dark"
              >已下架</el-tag
            >
            <div v-else>
              <el-tag
                type="danger"
                v-if="scope.row.approval_status === ''"
                size="large"
                effect="dark"
                >待审核</el-tag
              >
              <el-tag
                type="success"
                v-if="scope.row.approval_status === '通过'"
                size="large"
                effect="dark"
                >通过</el-tag
              >
              <!-- <el-tag type="danger" v-else size="large" effect="dark">{{
                scope.row.approval_status
              }}</el-tag>
              <el-tag type="success" v-if="scope.row.status === '已出租'" size="large"
                >已出租</el-tag
              > -->
            </div>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="维护情况"
          min-width="180"
          prop="updated_last_at"
          sortable
        >
          <template #default="scope">
            距上次编辑过去{{ dayjs().diff(dayjs(scope.row.updated_last_at), "day") }}天
          </template>
        </el-table-column>
        <el-table-column label="推广" :min-width="appStore.operateMinWith" fixed="right">
          <template #default="scope">
            <el-button
              type="primary"
              @click="changeState(scope.row, 1)"
              v-if="scope.row.status === '已下架'"
              >上架</el-button
            >
            <el-button type="danger" @click="changeState(scope.row, 2)" v-else
              >下架</el-button
            >
          </template>
        </el-table-column>
        <el-table-column
          label="发布时间"
          :min-width="appStore.operateMinWith"
          fixed="right"
          sortable
        >
          <template #default="scope">
            {{ dayjs(scope.row.created_at).format("YYYY-MM-DD HH:mm:ss") }}
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
    <!-- 重置密码对话框 -->
    <el-dialog
      v-model="resetPwdDialog"
      title="重置密码"
      width="500px"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
    >
      <el-form :model="resetPwdInfo" ref="resetPwdForm" label-width="100px">
        <el-form-item label="用户账号">
          <el-input v-model="resetPwdInfo.userName" disabled />
        </el-form-item>
        <el-form-item label="用户昵称">
          <el-input v-model="resetPwdInfo.nickName" disabled />
        </el-form-item>
        <el-form-item label="新密码">
          <div class="flex w-full">
            <el-input
              class="flex-1"
              v-model="resetPwdInfo.password"
              placeholder="请输入新密码"
              show-password
            />
            <el-button
              type="primary"
              @click="generateRandomPassword"
              style="margin-left: 10px"
            >
              生成随机密码
            </el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeResetPwdDialog">取 消</el-button>
          <el-button type="primary" @click="confirmResetPassword">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <el-drawer
      v-model="addUserDialog"
      :size="appStore.drawerSize"
      :show-close="false"
      :close-on-press-escape="false"
      :close-on-click-modal="false"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">用户</span>
          <div>
            <el-button @click="closeAddUserDialog">取 消</el-button>
            <el-button type="primary" @click="enterAddUserDialog">确 定</el-button>
          </div>
        </div>
      </template>

      <el-form ref="userForm" :rules="rules" :model="userInfo" label-width="80px">
        <el-form-item v-if="dialogFlag === 'add'" label="用户名" prop="userName">
          <el-input v-model="userInfo.userName" />
        </el-form-item>
        <el-form-item v-if="dialogFlag === 'add'" label="密码" prop="password">
          <el-input v-model="userInfo.password" />
        </el-form-item>
        <el-form-item label="昵称" prop="nickName">
          <el-input v-model="userInfo.nickName" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="userInfo.phone" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="userInfo.email" />
        </el-form-item>
        <el-form-item label="用户角色" prop="authorityId">
          <el-cascader
            v-model="userInfo.authorityIds"
            style="width: 100%"
            :options="authOptions"
            :show-all-levels="false"
            :props="{
              multiple: true,
              checkStrictly: true,
              label: 'authorityName',
              value: 'authorityId',
              disabled: 'disabled',
              emitPath: false,
            }"
            :clearable="false"
          />
        </el-form-item>
        <el-form-item label="启用" prop="disabled">
          <el-switch
            v-model="userInfo.enable"
            inline-prompt
            :active-value="1"
            :inactive-value="2"
          />
        </el-form-item>
        <el-form-item label="头像" label-width="80px">
          <SelectImage v-model="userInfo.headerImg" />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import { setUserAuthorities, register, deleteUser } from "@/api/user";

import { getHouseList, changeHouseState, batchPass } from "@/api/house";

import { getAuthorityList } from "@/api/authority";
import CustomPic from "@/components/customPic/index.vue";
import WarningBar from "@/components/warningBar/warningBar.vue";
import { setUserInfo, resetPassword } from "@/api/user.js";

import { nextTick, ref, watch } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { Picture as IconPicture } from "@element-plus/icons-vue";
import SelectImage from "@/components/selectImage/selectImage.vue";
import { useAppStore } from "@/pinia";
import dayjs from "dayjs";

import { searchXiaoqu } from "@/api/center";

defineOptions({
  name: "House",
});

const appStore = useAppStore();

const searchInfo = ref({
  approvalStatus: "",
  phone: "",
  keyword: "",
});

const approvalStatusOptions = [
  {
    value: "通过",
    label: "通过",
  },
  {
    value: "未通过",
    label: "未通过",
  },
  {
    value: "待审批",
    label: "待审批",
  },
];

const onSubmit = () => {
  page.value = 1;
  getTableData();
};

const onReset = () => {
  searchInfo.value = {
    approvalStatus: "",
    phone: "",
    keyword: "",
  };
  getTableData();
};

const page = ref(1);
const total = ref(0);
const pageSize = ref(10);
const tableData = ref([]);
// 分页
const handleSizeChange = (val) => {
  pageSize.value = val;
  getTableData();
};

const handleCurrentChange = (val) => {
  page.value = val;
  getTableData();
};

// 查询
const getTableData = async () => {
  const table = await getHouseList({
    page: page.value,
    pageSize: pageSize.value,
    ...searchInfo.value,
  });

  if (table.code === 0) {
    tableData.value = table.data.list;
    total.value = table.data.total;
    page.value = table.data.page;
    pageSize.value = table.data.pageSize;
  }
};

watch(
  () => tableData.value,
  () => {
    setAuthorityIds();
  }
);

const initPage = async () => {
  getTableData();
};

initPage();

// 重置密码对话框相关
const resetPwdDialog = ref(false);
const resetPwdForm = ref(null);
const resetPwdInfo = ref({
  ID: "",
  userName: "",
  nickName: "",
  password: "",
});

// 生成随机密码
const generateRandomPassword = () => {
  const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*";
  let password = "";
  for (let i = 0; i < 12; i++) {
    password += chars.charAt(Math.floor(Math.random() * chars.length));
  }
  resetPwdInfo.value.password = password;
  // 复制到剪贴板
  navigator.clipboard
    .writeText(password)
    .then(() => {
      ElMessage({
        type: "success",
        message: "密码已复制到剪贴板",
      });
    })
    .catch(() => {
      ElMessage({
        type: "error",
        message: "复制失败，请手动复制",
      });
    });
};

// 确认重置密码
const confirmResetPassword = async () => {
  if (!resetPwdInfo.value.password) {
    ElMessage({
      type: "warning",
      message: "请输入或生成密码",
    });
    return;
  }

  const res = await resetPassword({
    ID: resetPwdInfo.value.ID,
    password: resetPwdInfo.value.password,
  });

  if (res.code === 0) {
    ElMessage({
      type: "success",
      message: res.msg || "密码重置成功",
    });
    resetPwdDialog.value = false;
  } else {
    ElMessage({
      type: "error",
      message: res.msg || "密码重置失败",
    });
  }
};

// 关闭重置密码对话框
const closeResetPwdDialog = () => {
  resetPwdInfo.value.password = "";
  resetPwdDialog.value = false;
};
const setAuthorityIds = () => {
  tableData.value &&
    tableData.value.forEach((user) => {
      user.authorityIds =
        user.authorities &&
        user.authorities.map((i) => {
          return i.authorityId;
        });
    });
};

const deleteUserFunc = async (row) => {
  ElMessageBox.confirm("确定要删除吗?", "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  }).then(async () => {
    const res = await deleteUser({ id: row.ID });
    if (res.code === 0) {
      ElMessage.success("删除成功");
      await getTableData();
    }
  });
};

// 弹窗相关
const userInfo = ref({
  userName: "",
  password: "",
  nickName: "",
  headerImg: "",
  authorityId: "",
  authorityIds: [],
  enable: 1,
});

const rules = ref({
  userName: [
    { required: true, message: "请输入用户名", trigger: "blur" },
    { min: 5, message: "最低5位字符", trigger: "blur" },
  ],
  password: [
    { required: true, message: "请输入用户密码", trigger: "blur" },
    { min: 6, message: "最低6位字符", trigger: "blur" },
  ],
  nickName: [{ required: true, message: "请输入用户昵称", trigger: "blur" }],
  phone: [
    {
      pattern: /^1([38][0-9]|4[014-9]|[59][0-35-9]|6[2567]|7[0-8])\d{8}$/,
      message: "请输入合法手机号",
      trigger: "blur",
    },
  ],
  email: [
    {
      pattern: /^([0-9A-Za-z\-_.]+)@([0-9a-z]+\.[a-z]{2,3}(\.[a-z]{2})?)$/g,
      message: "请输入正确的邮箱",
      trigger: "blur",
    },
  ],
  authorityId: [{ required: true, message: "请选择用户角色", trigger: "blur" }],
});
const userForm = ref(null);
const enterAddUserDialog = async () => {
  userInfo.value.authorityId = userInfo.value.authorityIds[0];
  userForm.value.validate(async (valid) => {
    if (valid) {
      const req = {
        ...userInfo.value,
      };
      if (dialogFlag.value === "add") {
        const res = await register(req);
        if (res.code === 0) {
          ElMessage({ type: "success", message: "创建成功" });
          await getTableData();
          closeAddUserDialog();
        }
      }
      if (dialogFlag.value === "edit") {
        const res = await setUserInfo(req);
        if (res.code === 0) {
          ElMessage({ type: "success", message: "编辑成功" });
          await getTableData();
          closeAddUserDialog();
        }
      }
    }
  });
};

const addUserDialog = ref(false);
const closeAddUserDialog = () => {
  userForm.value.resetFields();
  userInfo.value.headerImg = "";
  userInfo.value.authorityIds = [];
  addUserDialog.value = false;
};

const dialogFlag = ref("add");

const addUser = () => {
  dialogFlag.value = "add";
  addUserDialog.value = true;
};

const tempAuth = {};
const changeAuthority = async (row, flag, removeAuth) => {
  if (flag) {
    if (!removeAuth) {
      tempAuth[row.ID] = [...row.authorityIds];
    }
    return;
  }
  await nextTick();
  const res = await setUserAuthorities({
    ID: row.ID,
    authorityIds: row.authorityIds,
  });
  if (res.code === 0) {
    ElMessage({ type: "success", message: "角色设置成功" });
  } else {
    if (!removeAuth) {
      row.authorityIds = [...tempAuth[row.ID]];
      delete tempAuth[row.ID];
    } else {
      row.authorityIds = [removeAuth, ...row.authorityIds];
    }
  }
};

const openEdit = (row) => {
  dialogFlag.value = "edit";
  userInfo.value = JSON.parse(JSON.stringify(row));
  addUserDialog.value = true;
};

const switchEnable = async (row) => {
  userInfo.value = JSON.parse(JSON.stringify(row));
  await nextTick();
  const req = {
    ...userInfo.value,
  };
  const res = await setUserInfo(req);
  if (res.code === 0) {
    ElMessage({
      type: "success",
      message: `${req.enable === 2 ? "禁用" : "启用"}成功`,
    });
    await getTableData();
    userInfo.value.headerImg = "";
    userInfo.value.authorityIds = [];
  }
};

//表格勾选
const selectedRows = ref([]);
const handleSelectionChange = (val) => {
  selectedRows.value = val;
};

//批量审核通过
const handleBatchPass = async () => {
  const res = await batchPass({
    ids: selectedRows.value.map((i) => i.ID),
    state: 1,
  });
  if (res.code === 0) {
    ElMessage.success("批量通过成功");
    await getTableData();
  }
};

const handleBatchDown = async () => {
  const res = await batchPass({
    ids: selectedRows.value.map((i) => i.ID),
    state: 2,
  });
  if (res.code === 0) {
    ElMessage.success("批量下架成功");
    await getTableData();
  }
};

const changeState = async (row, status) => {
  switch (status) {
    case 1: //上架
      ElMessageBox.confirm("确定要上架吗?", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      }).then(async () => {
        const res = await changeHouseState({ ids: [row.ID], state: 1 });
        if (res.code === 0) {
          ElMessage.success("上架成功");
          await getTableData();
        }
      });
      break;
    case 2: //下架
      ElMessageBox.confirm("确定要下架吗?", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      }).then(async () => {
        const res = await changeHouseState({ ids: [row.ID], state: 2 });
        if (res.code === 0) {
          ElMessage.success("下架成功");
          await getTableData();
        }
      });
      break;
  }
};

// 搜索小区
const searchXiaoquLoading = ref(false);
const xiaoquOptions = ref([]);
const remoteSearchXiaoqu = async (queryString, callback) => {
  if (queryString) {
    searchXiaoquLoading.value = true;
    const res = await searchXiaoqu(queryString);
    searchXiaoquLoading.value = false;
    console.log(res);
    xiaoquOptions.value = res.data.list.map((item) => {
      return {
        value: item.ID,
        label: item.name,
      };
    });
  } else {
    xiaoquOptions.value = [];
  }
};
</script>

<style lang="scss">
.header-img-box {
  @apply w-52 h-52 border border-solid border-gray-300 rounded-xl flex justify-center items-center cursor-pointer;
}
.gva-table-box .image-slot {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
  background: #e5e5e5;
  font-size: 30px;
}
</style>
