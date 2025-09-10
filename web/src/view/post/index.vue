<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
                  <el-form-item label="所属小区">
          <el-select
            v-model="searchInfo.xiaoquId"
            filterable
            remote
            reserve-keyword
            placeholder="小区名称"
            :remote-method="remoteSearchXiaoqu"
            :loading="searchXiaoquLoading"
            clearable
          >
            <el-option
              v-for="item in xiaoquOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="出租类型">
          <el-select v-model="searchInfo.rentType" placeholder="出租类型" clearable>
            <el-option
              v-for="item in rentTypeOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="户室号">
          <el-input v-model="searchInfo.doorNo" placeholder="户室号 例 1201" clearable />
        </el-form-item>
        <el-form-item label="经纪人手机号">
          <el-input v-model="searchInfo.phone" placeholder="经纪人手机号" clearable />
        </el-form-item>
        <el-form-item label="微信号">
          <el-input v-model="searchInfo.wxNo" placeholder="微信号" clearable />
        </el-form-item>
        <el-form-item label="审核状态">
          <el-select v-model="searchInfo.approvalStatus" placeholder="审核状态" clearable>
            <el-option
              v-for="item in approvalStatusOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="是否有图">
          <el-select v-model="searchInfo.hasPic" placeholder="是否有图" clearable>
            <el-option
              v-for="item in hasImgOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="最后编辑时间"> 
          <el-date-picker
            v-model="value"
            type="daterange"
            style="width: 340px"
            clearable
            unlink-panels
            range-separator="-"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            @change="handleDateChange"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="handleSearch"> 查询 </el-button>
          <el-button icon="refresh" @click="onReset"> 重置 </el-button>
        </el-form-item>
      </el-form>
      <el-divider />
      <!-- 查看 所在小区 出租类型 户型 户室号 价格 审核状态 是否有图 访问次数 分享次数 关注数 所属手机号 微信号 昵称 最后编辑时间 -->
       <div class="gva-table-box">
        <el-table :data="tableData" style="width: 100%">
            <!-- 操作 -->
            <el-table-column label="详情" width="80">
              <template #default="scope" >
                <el-button type="primary" link size="default" @click="handleDetail(scope.row)">查看</el-button>
              </template>
            </el-table-column>
          <el-table-column prop="xiaoqu" label="所在小区" />
          <el-table-column prop="rent_type" label="出租类型" />
          <el-table-column prop="door_no" width="200" label="户室号">
            <template #default="scope">
              {{scope.row.door_no}}{{ scope.row.room_code?'-'+scope.row.room_code+'号房':'' }}
            </template>
          </el-table-column>
          <el-table-column prop="price" label="价格" />
          <el-table-column prop="approval_status" label="审核状态" >
            <template #default="scope">
              {{!scope.row.approval_status?'待审核':scope.row.approval_status}}
            </template>
          </el-table-column>
          <el-table-column prop="hasPic" label="是否有图">
            <template #default="scope">
              {{scope.row.hasPic ? '有' : '无'}}
            </template>
          </el-table-column>
          <el-table-column prop="view" label="访问次数" />
          <!-- <el-table-column prop="shared" label="分享次数" /> -->
          <el-table-column prop="follow" label="关注数" />
          <el-table-column prop="phone" label="所属手机号" />
          <el-table-column prop="wxNo" label="微信号" />
          <el-table-column prop="wxNickName" label="昵称" />
          <el-table-column prop="updated_last_at" label="最后编辑时间" >
            <template #default="scope">
              {{dayjs(scope.row.updated_last_at).format('YYYY-MM-DD HH:mm:ss')}}
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
            <el-button @click="houseDetailDialog = false" type="primary">关 闭</el-button>
          </div>
        </div>
      </template>

      <el-form ref="houseFormRef" :model="form" label-width="80px">
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
        <el-form-item label="户型" prop="house_type">
          <el-radio-group v-model="form.house_type" disabled>
            <el-radio
              v-for="item in houseTypeOptions[form.rent_type]"
              :key="item.value"
              :label="item.label"
            />
          </el-radio-group>
        </el-form-item>
        <!-- <el-form-item label="户型" prop="house_type" v-if="form.rent_type !== '合租'">
          <el-radio-group v-model="form.house_type" disabled>
            <el-radio
              v-for="item in rentTypeOptions[form.rent_type]"
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
        </el-form-item> -->
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
import {ref, onMounted} from "vue"
import { getPostList } from "@/api/statis"
import dayjs from "dayjs";
import { searchXiaoqu, getHouseOptions } from "@/api/center";
import { useAppStore } from "@/pinia";
import CustomPic from "@/components/customPic/index.vue";
import WarningBar from "@/components/warningBar/warningBar.vue";

const appStore = useAppStore();
const houseFormRef = ref(null)
const searchInfo = ref({
  xiaoquId: undefined,
  rentType: undefined,
  doorNo: undefined,
  phone: undefined,
  wxNo: undefined,
  approvalStatus: undefined,
  hasPic: undefined,
  updatedAtLast: undefined,
  updatedAtStart: undefined,
})
const tableData = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
//通过 未通过 待审核
const approvalStatusOptions = ref([
  {
    value: '通过',
    label: '通过'
  },
  {
    value: '未通过',
    label: '未通过'
  },
  {
    value: '待审核',
    label: '待审核'
  }
])
const hasImgOptions = ref([
    {
        value: true,
        label: '有'
    },
    {
        value: false,
        label: '无'
    }
])
const value = ref([]);

onMounted(()=>{
  handeleGetHouseOptions()
  getTableData()
})

const handleDateChange = (val) => {
  if (val&&val.length === 2) {
    searchInfo.value.updatedAtStart = val[0]
    searchInfo.value.updatedAtLast = val[1]
  }else{
    searchInfo.value.updatedAtStart = null
    searchInfo.value.updatedAtLast = null
  }
}

const handleSearch = () => {
    page.value = 1
  getTableData()
}

const onReset = () => {
  searchInfo.value = {
    xiaoquId: undefined,
    rentType: undefined,
    doorNo: undefined,
    phone: undefined,
    wxNo: undefined,
    approvalStatus: undefined,
    hasPic: undefined,
    updatedAtLast: undefined,
    updatedAtStart: undefined,
  }
  value.value = []
  page.value = 1
  getTableData()
}
// 分页
const handleSizeChange = (val) => {
  pageSize.value = val;
  getTableData();
};

const handleCurrentChange = (val) => {
  page.value = val;
  getTableData();
};

const getTableData = async ()=>{
  try {
    let data = {
      ...searchInfo.value,
      page: page.value,
      pageSize: pageSize.value
    }
    const res = await getPostList(data)
    if (res.code === 0) {
        tableData.value = res.data.list
        total.value = res.data.total
    }
  } catch (error) {
    console.log(error)
  }
}
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
    console.log(rentTypeOptions.value,houseTypeOptions.value,featureOptions.value,'111');
    
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
const handleDetail = (detail) => {
  houseDetailDialog.value = true
  form.value = {...detail}
  form.value.feature = detail.feature.split(',')
  form.value.fileList = detail.attachments.house.map(item => ({
    url: item.url,
    name: item.name,
  }))
}

</script>
<style lang="scss" scoped></style>
