<template>
  <div>
    <warning-bar
      title="真房通知：请保持房源信息无异常可出租，平台可能会检测更正；多次呗举报会限制显示"
    />
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
        <el-form-item label="出租类型">
          <el-select v-model="searchInfo.rent_type" placeholder="出租类型">
            <el-option
              v-for="item in rentTypeOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
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
          <el-button type="primary" icon="search" @click="onSearch"> 查询 </el-button>
          <el-button icon="refresh" @click="onReset"> 重置 </el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="addHouse"> 新增房源 </el-button>
      </div>
      <el-table :data="tableData" row-key="ID" v-loading="false">
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
                >{{ scope.row.status }}</el-tag
              >
              <el-tag
                type="info"
                v-if="scope.row.approval_status === '未通过'"
                size="large"
                effect="dark"
                >未通过</el-tag
              >
            </div>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="维护情况"
          min-width="180"
          prop="updated_last_at"
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
        <el-table-column label="操作" :min-width="appStore.operateMinWith" fixed="right">
          <template #default="scope">
            <el-button type="primary" @click="handleEditHouse(scope.row)">编辑</el-button>
            <el-button type="danger" @click="handleDeleteHouse(scope.row)"
              >删除</el-button
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
      v-if="dialogFlag === 'add'"
      v-model="addHouseDialog"
      :size="appStore.drawerSize"
      :show-close="false"
      :close-on-press-escape="false"
      :close-on-click-modal="false"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">房源</span>
          <div>
            <el-button @click="closeAddHouseDialog">取 消</el-button>
            <el-button type="primary" @click="enterAddHouseDialog(houseFormRef)"
              >提交审核</el-button
            >
          </div>
        </div>
      </template>

      <el-form ref="houseFormRef" :rules="rules" :model="form" label-width="80px">
        <el-form-item label="出租类型" prop="rent_type">
          <el-radio-group v-model="form.rent_type" @change="handleRentTypeChange">
            <el-radio
              v-for="item in rentTypeOptions"
              :key="item.value"
              :label="item.label"
            />
          </el-radio-group>
        </el-form-item>
        <el-form-item label="小区名称" prop="xiaoqu_id">
          <el-select
            v-model="form.xiaoqu_id"
            filterable
            remote
            reserve-keyword
            placeholder="小区名称"
            :remote-method="remoteSearchXiaoqu"
            :loading="searchXiaoquLoading"
            style="width: 240px"
            @change="selectXiaoQu"
          >
            <el-option
              v-for="item in xiaoquOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="户室号" prop="door_no">
          <el-cascader :props="props" v-if="form.xiaoqu" @change="handleDoorNoChange" />
          <br />
          <el-text class="mx-1" type="info">户室信息将不在用户端展示具体信息</el-text>
        </el-form-item>
        <el-form-item label="房间号" prop="room_code" v-if="form.rent_type !== '整租'">
          <el-select
            v-model="form.room_code"
            class="m-2"
            placeholder="请选择房间号"
            style="width: 240px"
          >
            <el-option
              v-for="item in ['1', '2', '3', '4', '5', '6', '7', '8']"
              :key="item"
              :label="item"
              :value="item"
            />
          </el-select>
          <el-text class="mx-1" type="info"
            >注：自进门右手起，逆时针数，不区分空间功能，第一间为1号，房间有门即算。</el-text
          >
        </el-form-item>
        <el-form-item label="户型" prop="house_type" v-if="form.rent_type !== '合租'">
          <el-radio-group v-model="form.house_type" @change="handleHouseTypeChange">
            <el-radio
              v-for="item in houseTypeOptions[form.rent_type]"
              :key="item.value"
              :label="item.label"
            />
          </el-radio-group>
        </el-form-item>
        <el-form-item label="户型" prop="house_type" v-else>
          <el-radio-group v-model="form.house_type" @change="handleHouseTypeChange">
            <el-radio
              v-for="item in houseTypeOptions[form.rent_type]"
              :key="item.value"
              :label="item.label"
            />
          </el-radio-group>
        </el-form-item>
        <!-- <el-form-item label="面积" prop="area">
          <el-input type="number" v-model="form.area">
            <template #append>平方米</template>
          </el-input>
        </el-form-item> -->
        <el-divider />
        <el-form-item label="月租金" prop="price">
          <el-input type="number" v-model="form.price">
            <template #append>元/月</template>
          </el-input>
          <el-text class="mx-1" type="danger">年租月付的价格</el-text>
        </el-form-item>
        <el-form-item label="亮点" prop="feature">
          <el-checkbox-group v-model="form.feature">
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
            placeholder="例：分整租（次卧）物业费送宽带、佣金600年签"
          />
        </el-form-item>
        <el-divider />
        <warning-bar
          title="引起99%房源下架的图片规则：1、不得违反经纪公司logo发布规则（限1个，白色半透明且尺寸在25%以内） 2、不得盗图（含58、赶集、安居客等logo）3、不得有任何装饰、图文"
        />
        <el-form-item label="上传房源图片" prop="fileList">
          <el-upload
            v-model:file-list="form.fileList"
            list-type="picture-card"
            :on-preview="handlePictureCardPreview"
            :on-remove="handleRemove"
            :auto-upload="false"
            :multiple="true"
            :limit="12"
          >
            <el-icon><Plus /></el-icon>
          </el-upload>
          <el-dialog v-model="dialogVisible">
            <el-image
              w-full
              :src="dialogImageUrl"
              alt="Preview Image"
              style="width: 100%"
            />
          </el-dialog>
        </el-form-item>
        <el-form-item label="联系电话" prop="phone">
          <el-input v-model="form.phone" />
        </el-form-item>
        <el-form-item label="" prop="type">
          <el-checkbox-group v-model="form.type">
            <el-checkbox
              label="我已阅读并同意《房源发布规则》，承诺发布真实房源信息，本方可带看"
              name="type"
            />
          </el-checkbox-group>
        </el-form-item>
      </el-form>
    </el-drawer>

    <el-drawer
      v-if="dialogFlag === 'edit'"
      v-model="editHouseDialog"
      :size="appStore.drawerSize"
      :show-close="false"
      :close-on-press-escape="false"
      :close-on-click-modal="false"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">房源</span>
          <div>
            <el-button @click="closeEditHouseDialog">取 消</el-button>
            <el-button type="primary" @click="enterEditHouseDialog(houseFormRef)"
              >保存</el-button
            >
          </div>
        </div>
      </template>

      <el-form ref="houseFormRef" :rules="rules" :model="form" label-width="80px">
        <el-form-item label="出租类型" prop="rent_type">
          <el-radio-group
            v-model="form.rent_type"
            @change="handleRentTypeChange"
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
              v-for="item in houseTypeOptions[form.rent_type]"
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
          <el-input type="number" v-model="form.price">
            <template #append>元/月</template>
          </el-input>
          <el-text class="mx-1" type="danger">年租月付的价格</el-text>
        </el-form-item>
        <el-form-item label="亮点" prop="feature">
          <el-checkbox-group v-model="form.feature">
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
            placeholder="例：分整租（次卧）物业费送宽带、佣金600年签"
          />
        </el-form-item>
        <el-divider />
        <warning-bar
          title="引起99%房源下架的图片规则：1、不得违反经纪公司logo发布规则（限1个，白色半透明且尺寸在25%以内） 2、不得盗图（含58、赶集、安居客等logo）3、不得有任何装饰、图文"
        />
        <el-form-item label="上传房源图片" prop="fileList">
          <el-upload
            v-model:file-list="form.fileList"
            list-type="picture-card"
            :on-preview="handlePictureCardPreview"
            :on-remove="handleRemove"
            :auto-upload="false"
            :multiple="true"
            :limit="12"
          >
            <el-icon><Plus /></el-icon>
          </el-upload>
          <el-dialog v-model="dialogVisible">
            <el-image
              w-full
              :src="dialogImageUrl"
              alt="Preview Image"
              style="width: 100%"
            />
          </el-dialog>
        </el-form-item>
        <el-form-item label="联系电话" prop="phone">
          <el-input v-model="form.phone" />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref, watch } from "vue";
import { useAppStore } from "@/pinia";
import dayjs from "dayjs";

import { ElMessage, ElMessageBox } from "element-plus";
import { Picture as IconPicture } from "@element-plus/icons-vue";

import WarningBar from "@/components/warningBar/warningBar.vue";

import {
  getHouseListMy,
  changeHouseState,
  createHouse,
  deleteHouse,
  editHouse,
} from "@/api/house";
import { searchXiaoqu, getHouseOptions, uploadFile } from "@/api/center";
import { getBuildings, getUnits, getHouse } from "@/api/base";

import axios from "axios";

defineOptions({
  name: "HouseList",
});

const appStore = useAppStore();

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

// 查询列表相关
const searchInfo = ref({
  keyword: "",
  rent_type: "",
  approvalStatus: "",
});
const onSearch = () => {
  page.value = 1;
  getTableData();
};
const onReset = () => {
  searchInfo.value = {
    keyword: "",
    rent_type: "",
    approvalStatus: "",
  };
  getTableData();
};

//列表相关
const getTableData = async () => {
  const table = await getHouseListMy({
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
  handeleGetHouseOptions();
};

// 新增房源相关
const dialogFlag = ref("add");
const addHouseDialog = ref(false);
const houseFormRef = ref();

// 新增房源表单相关
const rules = ref({
  rent_type: [{ required: true, message: "请勾选出租类型", trigger: "change" }],
  xiaoqu_id: [{ required: true, message: "请选择小区", trigger: "change" }],
  door_no: [{ required: true, message: "请输入户室号", trigger: "change" }],
  room_code: [{ required: true, message: "请输入房间号", trigger: "change" }],
  house_type: [{ required: true, message: "请选择户型", trigger: "change" }],
  area: [{ required: true, message: "请输入面积", trigger: "change" }],
  price: [{ required: true, message: "请输入月租金", trigger: "change" }],
  phone: [
    { required: true, message: "请输入联系电话", trigger: "change" },
    {
      pattern: /^1([38][0-9]|4[014-9]|[59][0-35-9]|6[2567]|7[0-8])\d{8}$/,
      message: "请输入合法手机号",
      trigger: "change",
    },
  ],
  type: [{ required: true, message: "请勾选同意发布规则", trigger: "change" }],
});
const userForm = ref(null);
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
const selectXiaoQu = (val) => {
  const preForm = { ...form.value };
  houseFormRef.value.resetFields();
  form.value = {
    ...form.value,
    rent_type: preForm.rent_type,
    xiaoqu: "",
    xiaoqu_id: "",
    door_no: "",
  };
  setTimeout(() => {
    form.value.xiaoqu = xiaoquOptions.value.find((item) => item.value === val).label;
    form.value.xiaoqu_id = val;
  }, 0);
};

//处理楼栋号相关
let id = 0;
const props = {
  lazy: true,
  lazyLoad(node, resolve) {
    const { level } = node;

    switch (level) {
      case 0:
        getBuildings(form.value.xiaoqu_id).then((res) => {
          if (res.code === 0) {
            const nodes = res.data.map((item) => ({
              value: item.id,
              label: item.name + "号楼",
              leaf: level >= 2,
            }));
            // Invoke `resolve` callback to return the child nodes data and indicate the loading is finished.
            resolve(nodes);
          }
        });
        break;
      case 1:
        form.value.building_id = node.value;
        getUnits(node.value).then((res) => {
          if (res.code === 0) {
            const nodes = res.data.map((item) => ({
              value: item.id,
              label: item.name + "单元",
              leaf: level >= 2,
            }));
            // Invoke `resolve` callback to return the child nodes data and indicate the loading is finished.
            resolve(nodes);
          }
        });
        break;
      case 2:
        form.value.unit_id = node.value;
        getHouse(node.value).then((res) => {
          if (res.code === 0) {
            const nodes = res.data.map((item) => ({
              value: item.id,
              label: item.name + "室",
              leaf: level >= 2,
            }));
            // Invoke `resolve` callback to return the child nodes data and indicate the loading is finished.
            resolve(nodes);
          }
        });
        break;
      case 3:
        form.value.house_id = node.value;
        break;

      default:
        break;
    }
    console.log(form.value);
    form.value.door_no = node.value;
  },
};

//处理房源图片上传
const dialogImageUrl = ref("");
const dialogVisible = ref(false);

const handleRemove = (uploadFile, uploadFiles) => {
  console.log(uploadFile, uploadFiles);
};

const handlePictureCardPreview = (uploadFile) => {
  dialogImageUrl.value = uploadFile.url;
  dialogVisible.value = true;
};

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

    form.value.rent_type = rentTypeOptions.value[0].value;
  }
};

const handleRentTypeChange = (value) => {
  houseFormRef.value.resetFields();
  form.value = {
    ...form.value,
    rent_type: value,
    xiaoqu: "",
    door_no: "",
  };
};

const handleDoorNoChange = (value) => {
  const preForm = { ...form.value };

  houseFormRef.value.resetFields();
  form.value = {
    ...form.value,
    rent_type: preForm.rent_type,
    xiaoqu: preForm.xiaoqu,
    xiaoqu_id: preForm.xiaoqu_id,
    building_id: value[0],
    unit_id: value[1],
    house_id: value[2],
    door_no: value,
  };
};

const handleHouseTypeChange = (value) => {
  const preForm = { ...form.value };
  houseFormRef.value.resetFields();
  form.value = {
    ...form.value,
    rent_type: preForm.rent_type,
    xiaoqu: preForm.xiaoqu,
    xiaoqu_id: preForm.xiaoqu_id,
    door_no: preForm.door_no,
    building_id: preForm.building_id,
    unit_id: preForm.unit_id,
    house_type: value,
    house_id: preForm.house_id,
    room_code: preForm.room_code,
  };
};

//新增房源提交
const enterAddHouseDialog = async (formEl) => {
  console.log(formEl);
  if (!formEl) return;
  await formEl.validate(async (valid, fields) => {
    if (valid) {
      const values = { ...form.value };
      console.log(values);
      values.feature = values.feature.join(",");
      values.price = Number(values.price);

      if (values.fileList.length) {
        Promise.all(
          values.fileList.map((file) => {
            console.log(file);
            const form = new FormData();
            form.append("file", file.raw);
            return uploadFile(form);
          })
        ).then((res) => {
          if (res.find((item) => item.code !== 0)) {
            ElMessage.error("上传失败");
            return;
          }

          values.attachments.house = res.map((item) => ({
            url: item.data.url,
          }));

          delete values.type;
          delete values.fileList;
          delete values.door_no;
          createHouse(values).then((res) => {
            if (res.code === 0) {
              ElMessage.success("新增成功");
              addHouseDialog.value = false;
              getTableData();
              houseFormRef.value.resetFields();
            }
          });
        });

        return;
      }

      delete values.type;
      delete values.fileList;
      delete values.door_no;

      createHouse(values).then((res) => {
        if (res.code === 0) {
          ElMessage.success("新增成功");
          addHouseDialog.value = false;
          getTableData();
          houseFormRef.value.resetFields();
        }
      });
    } else {
      console.log("error submit!", form.value);
    }
  });
};
//打开新增房源弹框
const addHouse = () => {
  // handeleGetHouseOptions();
  dialogFlag.value = "add";
  addHouseDialog.value = true;
};
//关闭新增房源弹框
const closeAddHouseDialog = () => {
  houseFormRef.value.resetFields();
  addHouseDialog.value = false;
};

//编辑房源相关
const editHouseDialog = ref(false);
const handleEditHouse = (row) => {
  // handeleGetHouseOptions();
  dialogFlag.value = "edit";
  editHouseDialog.value = true;
  xiaoquOptions.value = [
    {
      value: row.xiaoqu_id,
      label: row.xiaoqu,
    },
  ];
  form.value = {
    ...row,
    feature: row.feature.split(","),
    fileList: row.attachments.house.map((item) => ({
      url: item.url,
    })),
  };
};
//关闭编辑房源弹框
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
  editHouseDialog.value = false;
};
//编辑房源提交
const enterEditHouseDialog = async (formEl) => {
  console.log(formEl);
  if (!formEl) return;
  await formEl.validate(async (valid, fields) => {
    if (valid) {
      const values = { ...form.value };
      console.log(values);
      values.feature = values.feature.join(",");
      values.price = Number(values.price);

      if (values.fileList.length) {
        Promise.all(
          values.fileList.map((file) => {
            console.log(file);
            const form = new FormData();
            form.append("file", file.raw);
            return uploadFile(form);
          })
        ).then((res) => {
          if (res.find((item) => item.code !== 0)) {
            ElMessage.error("上传失败");
            return;
          }
          values.attachments.house = res.map((item) => ({
            url: item.data.url,
          }));
          editHouse(values).then((res) => {
            if (res.code === 0) {
              ElMessage.success("编辑成功");
              editHouseDialog.value = false;
              getTableData();
            }
          });
        });
      }
    }
  });
};

const handleDeleteHouse = (row) => {
  ElMessageBox.confirm("确定要删除吗?", "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  }).then(() => {
    deleteHouse(row.ID).then((res) => {
      if (res.code === 0) {
        ElMessage.success("删除成功");
        getTableData();
      }
    });
  });
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

initPage();
</script>
<style lang="scss" scoped></style>
