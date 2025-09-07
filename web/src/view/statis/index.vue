<template>
  <div>
    <div class="gva-search-box">
      日期区间
      <el-date-picker
        v-model="value"
        type="daterange"
        unlink-panels
        range-separator="-"
        start-placeholder="开始日期"
        end-placeholder="结束日期"
        @change="handleDateChange"
      />
      <br />
      <br />
      <br />
      <el-space wrap :size="70">
        <div class="boxShadow w-[200px] px-7 py-3">
          <div>
            <div class="flex justify-between items-center">
              <div class="fz font-bold">新增经纪人</div>
            </div>
            <div class="mt-2 fz">{{statisData.add_saler}}</div>
          </div>
        </div>
        <div class="boxShadow w-[200px] px-7 py-3">
          <div>
            <div class="flex justify-between items-center">
              <div class="fz font-bold">使用的经纪人</div>
            </div>
            <div class="mt-2 fz">{{statisData.use_saler}}</div>
          </div>
        </div>
      </el-space>
      <br />
      <br />
      <br />
      <el-space wrap :size="70">
        <div class="boxShadow w-[200px] px-7 py-3">
          <div>
            <div class="flex justify-between items-center">
              <div class="fz font-bold">新增帖子</div>
            </div>
            <div class="mt-2 fz">{{statisData.add}}</div>
          </div>
        </div>
        <div class="boxShadow w-[200px] px-7 py-3">
          <div>
            <div class="flex justify-between items-center">
              <div class="fz font-bold">浏览帖子数</div>
            </div>
            <div class="mt-2 fz">{{statisData.view}}</div>
          </div>
        </div>
        <div class="boxShadow w-[200px] px-7 py-3">
          <div>
            <div class="flex justify-between items-center">
              <div class="fz font-bold">关注帖子数</div>
            </div>
            <div class="mt-2 fz">{{statisData.follow}}</div>
          </div>
        </div>
        <div class="boxShadow w-[200px] px-7 py-3">
          <div>
            <div class="flex justify-between items-center">
              <div class="fz font-bold">帖子分享次数</div>
            </div>
            <div class="mt-2 fz">{{statisData.shared}}</div>
          </div>
        </div>
        <div class="boxShadow px-7 py-3">
          <div>
            <div class="flex justify-between items-center">
              <div class="fz font-bold">联系方式被点击数</div>
            </div>
            <div class="mt-2 fz">{{statisData.click}}</div>
          </div>
        </div>
      </el-space>
    </div>
  </div>
</template>

<script setup>
import { ref } from "vue";
import { getStatis } from "@/api/statis";

const value = ref(["2025-09-01T00:00:00.000Z", new Date()]);
const statisData = ref({
  add_saler: 0,
  use_saler: 0,
  add: 0,
  view: 0,
  follow: 0,
  shared: 0,
  click: 0,
});

const getStatisData = async () => {
  const res = await getStatis({
    start: value.value[0],
    end: value.value[1],
  });
  if (res.code === 0) {
    statisData.value = res.data;
  }

};

const handleDateChange = (val) => {
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
</style>
