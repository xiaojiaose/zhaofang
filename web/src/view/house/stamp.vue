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
              style="width: 250px; height: 150px"
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
              <div style="text-align: left; cursor: pointer;" @click="handleCellClick(scope.row)">
                <el-text size="large" tag="b">{{ scope.row.xiaoqu }}</el-text>
                <br />
                <div>{{ scope.row.door_no }}</div>
              </div>
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
        <el-table-column align="left" label="审核状态" min-width="150">
          <template #default="scope">
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
              >{{ scope.row.status }}</el-tag
            >
            <el-tag
              type="info"
              v-if="scope.row.approval_status === '未通过'"
              size="large"
              effect="dark"
              >未通过</el-tag
            >
          </template>
        </el-table-column>
        <el-table-column align="left" label="出租状态" min-width="150">
          <template #default="scope">
            <el-tag
              type="danger"
              v-if="scope.row.status === '待出租'"
              size="large"
              effect="dark"
              >待出租</el-tag
            >
            <el-tag
              type="success"
              v-if="scope.row.status === '已出租'"
              size="large"
              effect="dark"
              >已出租</el-tag
            >
            <el-tag
              type="info"
              v-if="scope.row.status === '已下架'"
              size="large"
              effect="dark"
              >已下架</el-tag
            >
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

   <el-drawer
      v-model="houseDetailDialog"
      :size="appStore.drawerSize"
      :show-close="false"
      :close-on-press-escape="false"
      :close-on-click-modal="false"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">房源</span>
          <div>
            <el-button @click="closeEditHouseDialog" type="primary">关 闭</el-button>
          </div>
        </div>
      </template>

      <el-form ref="houseFormRef" :rules="rules" :model="form" label-width="80px">
        <el-form-item label="出租类型" prop="rent_type">
          <el-radio-group
            v-model="form.rent_type"
            disabled
          >
            <el-radio
              v-for="item in rentTypeOptions"
              :key="item.value"
              :label="item.label"
            />
          </el-radio-group>
        </el-form-item>
        <el-form-item label="小区名称" prop="xiaoqu_id">
          <el-input v-model="form.xiaoqu" disabled />
        </el-form-item>
        <el-form-item label="户室号" prop="door_no">
          <el-input v-model="form.door_no" disabled />
          <br />
          <el-text class="mx-1" type="info">户室信息将不在用户端展示具体信息</el-text>
        </el-form-item>
        <!-- <el-form-item label="房间号" prop="house_id" v-if="form.rent_type === '合租'">
          <el-select
            v-model="form.house_id"
            class="m-2"
            placeholder="请选择房间号"
            style="width: 240px"
          >
            <el-option
              v-for="item in [1,2,3,4,5,6,7,8,9,10]"
              :key="item"
              :label="item"
              :value="item"
            />
          </el-select>
          <el-text class="mx-1" type="info">注：自进门右手起，逆时针数，不区分空间功能，第一间为1号，房间有门即算。</el-text>
        </el-form-item> -->
        <el-form-item label="户型" prop="house_type" v-if="form.rent_type !== '合租'">
          <el-radio-group v-model="form.house_type" disabled>
            <el-radio
              v-for="item in houseTypeOptions[form.rent_type]"
              :key="item.value"
              :label="item.label"
            />
          </el-radio-group>
        </el-form-item>
        <el-form-item label="户型" prop="house_type" v-else>
          <el-radio-group v-model="form.house_type" disabled>
            <el-radio
              v-for="item in ['主卧', '次卧', '案间']"
              :key="item"
              :label="item"
            />
          </el-radio-group>
        </el-form-item>
        <!-- <el-form-item label="面积" prop="area">
          <el-input type="number" v-model="form.area" disabled>
            <template #append>平方米</template>
          </el-input>
        </el-form-item> -->
        <el-divider />
        <el-form-item label="月租金" prop="price">
          <el-input type="number" v-model="form.price" disabled>
            <template #append>元/月</template>
          </el-input>
          <el-text class="mx-1" type="danger">年租月付的价格</el-text>
        </el-form-item>
        <el-form-item label="亮点" prop="feature">
          <el-checkbox-group v-model="form.feature" disabled>
            <el-checkbox
              v-for="item in featureOptions[form.rent_type]"
              :key="item.value"
              :label="item.label"
              name="feature"
            />
          </el-checkbox-group>
        </el-form-item>
        <el-divider />
        <warning-bar
          title="不能出现任意联系方式（包括但不限于QQ、微信、电话、网址、MSN、邮箱等）；请勿添加其他小区广告，请勿输入与出租房源无关内容或非法信息。"
        />
        <el-form-item label="备注" prop="remarks">
          <el-input
            v-model="form.remarks"
            type="textarea"
            :rows="4"
            disabled
          />
        </el-form-item>
        <el-divider />
        <warning-bar
          title="引起99%房源下架的图片规则：1、不得违反经纪公司logo发布规则（限1个，白色半透明且尺寸在25%以内） 2、不得盗图（含58、赶集、安居客等logo）3、不得有任何装饰、图文"
        />
        <el-form-item label="房源图片" prop="fileList">
          <el-image
            v-for="item in form.fileList"
            :key="item.url"
            :src="item.url"
            fit="contain"
            :preview-src-list="form.fileList.map(item => item.url)"
            style="width: 150px; height: 150px; margin-right: 10px"
          />
        </el-form-item>
        <el-form-item label="联系电话" prop="phone">
          <el-input v-model="form.phone" disabled />
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

import { searchXiaoqu, getHouseOptions } from "@/api/center";

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

const initPage = async () => {
  getTableData();
};

initPage();

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



//点击某个单元格时
const houseFormRef = ref(null);
const handleCellClick = (row) => {
  handeleGetHouseOptions();
  houseDetailDialog.value = true;
  console.log(row)
  for (const key in row) {
    if (key === "feature") {
      form.value[key] = row[key]?.split(",")
    } else if (key === "attachments") {
      form.value['fileList'] = row[key].house?.map((item) => {
        return {
          url: item.url,
        }
      })
    } else {
      form.value[key] = row[key];
    }
    
  }
};
const houseDetailDialog = ref(false);
const form = ref({
  rent_type: "",
  xiaoqu: "",
  xiaoqu_id: "",
  door_no: "",
  house_type: "",
  price: "",
  feature: [],
  remarks: "",
  fileList: [],
  attachments: {
    house: [],
  },
  phone: "",
  type: [],
  house_id: "",
});
//房源选项回显
const rentTypeOptions = ref([]);
const houseTypeOptions = ref({});
const featureOptions = ref({});
const handeleGetHouseOptions = async () => {
  const res = await getHouseOptions();
  if (res.code === 0) {
    const data = res.data.houseType;
    rentTypeOptions.value = data
      .map((d) => d.name)
      .map((value) => {
        houseTypeOptions.value[value] = data
          .find((d) => d.name === value)
          .houseType.map((h) => ({
            value: h,
            label: h,
          }));

        featureOptions.value[value] = data
          .find((d) => d.name === value)
          .feature.map((f) => {
            return {
              value: f,
              label: f,
            };
          });

        return {
          value: value,
          label: value,
        };
      });
  }
};
//关闭详情房源弹框
const closeEditHouseDialog = () => {
  form.value = {
    rent_type: "",
    xiaoqu: "",
    xiaoqu_id: "",
    door_no: "",
    house_type: "",
    price: "",
    feature: [],
    remarks: "",
    fileList: [],
    attachments: {
      house: [],
    },
    phone: "",
    type: [],
    house_id: "",
  };

  houseFormRef.value.resetFields();
  houseDetailDialog.value = false;
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
