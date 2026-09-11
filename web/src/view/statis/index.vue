<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" class="form-inline">
        <el-form-item label="日期区间">
          <el-date-picker
            v-model="value"
            type="daterange"
            unlink-panels
            range-separator="-"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            @change="getStatisData"
            size="default"
          />
        </el-form-item>
        <el-form-item label="手机号统计">
          <el-input v-model="phone" placeholder="搜索房源手机号" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="getStatisData">查询</el-button>
          <el-button @click="handleResetPhone">清空手机号</el-button>
        </el-form-item>
      </el-form>

      <template v-if="!isPhone">
        <div class="section-title">全站汇总</div>
        <el-space wrap :size="32">
          <div
            class="boxShadow w-[200px] px-7 py-3"
            v-for="item in summaryCards"
            :key="item.label"
          >
            <div class="flex justify-between items-center">
              <div class="fz font-bold">{{ item.label }}</div>
            </div>
            <div class="mt-2 fz">{{ item.value }}</div>
          </div>
        </el-space>
      </template>

      <template v-if="isPhone">
        <div class="section-title">手机号维度统计</div>
        <el-space wrap :size="32">
          <div
            class="boxShadow w-[220px] px-7 py-3"
            v-for="item in phoneCards"
            :key="item.label"
          >
            <div class="flex justify-between items-center">
              <div class="fz font-bold">{{ item.label }}</div>
            </div>
            <div class="mt-2 fz">{{ item.value }}</div>
          </div>
        </el-space>
      </template>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from "vue";
import { getStatis } from "@/api/statis";

defineOptions({
  name: "Index",
});

const value = ref(["2025-09-01T00:00:00.000Z", new Date()]);
const phone = ref("");
const isPhone = ref(false);
const statisData = ref({
  add_saler: 0,
  use_saler: 0,
  add: 0,
  view: 0,
  follow: 0,
  shared: 0,
  click: 0,
  reward_apply: 0,
});
const phoneSummary = ref({
  add: 0,
  view: 0,
  shared: 0,
  click: 0,
  reward_apply: 0,
});

const summaryCards = computed(() => [
  { label: "新增经纪人", value: statisData.value.add_saler },
  { label: "使用的经纪人", value: statisData.value.use_saler },
  { label: "新增帖子", value: statisData.value.add },
  { label: "浏览帖子数", value: statisData.value.view },
  // { label: "关注帖子数", value: statisData.value.follow },
  { label: "帖子分享次数", value: statisData.value.shared },
  { label: "联系方式被点击数", value: statisData.value.click },
  { label: "申请出房有礼次数", value: statisData.value.reward_apply },
]);

const phoneCards = computed(() => [
  { label: "新增帖子数", value: phoneSummary.value.add },
  { label: "浏览帖子数", value: phoneSummary.value.view },
  { label: "帖子分享数", value: phoneSummary.value.shared },
  { label: "联系方式点击数", value: phoneSummary.value.click },
  { label: "申请出房有礼数", value: phoneSummary.value.reward_apply },
]);

const getStatisData = async () => {
  const res = await getStatis({
    start: value.value[0],
    end: value.value[1],
    phone: phone.value,
  });
  if (res.code === 0) {
    isPhone.value = phone.value !== "";
    statisData.value = res.data.summary || {};
    phoneSummary.value = res.data.phoneSummary || {
      add: 0,
      view: 0,
      shared: 0,
      click: 0,
      reward_apply: 0,
    };
  }
};

const handleResetPhone = () => {
  isPhone.value = false;
  phone.value = "";
  getStatisData();
};

getStatisData();
</script>

<style lang="scss" scoped>
.boxShadow {
  box-shadow: rgba(0, 0, 0, 0.02) 0px 1px 3px 0px, rgba(27, 31, 35, 0.15) 0px 0px 0px 1px;

  .fz {
    font-size: 20px;
  }
}

.section-title {
  margin: 28px 0 18px;
  font-size: 16px;
  font-weight: 600;
}
</style>
