<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item label="发布人确认">
          <el-select v-model="searchInfo.publisherConfirmStatus" clearable>
            <el-option label="待确认" value="待确认" />
            <el-option label="已拒绝" value="已拒绝" />
            <el-option label="已确认" value="已确认" />
          </el-select>
        </el-form-item>
        <el-form-item label="审核状态">
          <el-select v-model="searchInfo.auditStatus" clearable>
            <el-option label="未进入审核" value="未进入审核" />
            <el-option label="待审核" value="待审核" />
            <el-option label="审核中" value="审核中" />
            <el-option label="审通过待发放" value="审通过待发放" />
            <el-option label="已发放" value="已发放" />
            <el-option label="未通过" value="未通过" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" :disabled="!multipleSelection.length" @click="handleBatchAction('processing')">
          批量设为审核中
        </el-button>
        <el-button type="success" :disabled="!multipleSelection.length" @click="handleBatchAction('approve')">
          批量通过
        </el-button>
        <el-button type="warning" :disabled="!multipleSelection.length" @click="handleBatchAction('paid')">
          批量已发放
        </el-button>
        <el-button type="danger" :disabled="!multipleSelection.length" @click="handleBatchAction('reject')">
          批量未通过
        </el-button>
      </div>
      <el-table :data="tableData" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="48" />
        <el-table-column prop="xiaoqu" label="小区" min-width="160" />
        <el-table-column prop="doorNo" label="户室号" min-width="120" />
        <el-table-column prop="applyUserPhone" label="申请人手机号" min-width="140" />
        <el-table-column prop="applyUserWxNo" label="申请人微信号" min-width="140" />
        <el-table-column prop="publisherUserPhone" label="发布人手机号" min-width="140" />
        <el-table-column prop="publisherUserWxNo" label="发布人微信号" min-width="140" />
        <el-table-column prop="publisherConfirmStatus" label="发布人确认状态" min-width="140" />
        <el-table-column prop="auditStatus" label="审核状态" min-width="140" />
        <el-table-column label="申请时间" min-width="180">
          <template #default="scope">
            {{ formatTime(scope.row.CreatedAt) }}
          </template>
        </el-table-column>
        <el-table-column label="最后操作时间" min-width="180">
          <template #default="scope">
            {{ formatUnixMilli(scope.row.lastOperatedAtUnixMilli) }}
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="220" show-overflow-tooltip />
        <el-table-column label="操作" min-width="260" fixed="right">
          <template #default="scope">
            <el-button link type="primary" @click="handleAction(scope.row, 'processing')">审核中</el-button>
            <el-button link type="success" @click="handleAction(scope.row, 'approve')">通过</el-button>
            <el-button link type="warning" @click="handleAction(scope.row, 'paid')">已发放</el-button>
            <el-button link type="danger" @click="handleAction(scope.row, 'reject')">未通过</el-button>
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
import { ref } from 'vue'
import dayjs from 'dayjs'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getRewardList, rewardAction, rewardBatchAction } from '@/api/reward'

defineOptions({
  name: 'RewardList'
})

const defaultSearchInfo = () => ({
  publisherConfirmStatus: '已确认',
  auditStatus: ''
})

const searchInfo = ref(defaultSearchInfo())
const tableData = ref([])
const multipleSelection = ref([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const getTableData = async () => {
  const res = await getRewardList({
    page: page.value,
    pageSize: pageSize.value,
    ...searchInfo.value
  })
  if (res.code === 0) {
    tableData.value = res.data.list
    total.value = res.data.total
  }
}

const handleAction = async (row, action) => {
  const res = await rewardAction({ id: row.ID, action })
  if (res.code === 0) {
    ElMessage.success('操作成功')
    getTableData()
  }
}

const handleBatchAction = async (action) => {
  if (!multipleSelection.value.length) {
    ElMessage.warning('请先选择记录')
    return
  }
  await ElMessageBox.confirm('确认批量执行该操作吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
  const results = await rewardBatchAction(
    multipleSelection.value.map((item) => item.ID),
    action
  )
  const failed = results.filter((item) => item.code !== 0)
  if (failed.length) {
    ElMessage.warning(`已完成 ${results.length - failed.length} 条，失败 ${failed.length} 条`)
  } else {
    ElMessage.success('批量操作成功')
  }
  getTableData()
}

const handleSelectionChange = (rows) => {
  multipleSelection.value = rows
}

const handleSearch = () => {
  page.value = 1
  getTableData()
}

const handleReset = () => {
  searchInfo.value = defaultSearchInfo()
  page.value = 1
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const formatTime = (value) => {
  if (!value) {
    return '-'
  }
  return dayjs(value).format('YYYY-MM-DD HH:mm:ss')
}

const formatUnixMilli = (value) => {
  if (!value) {
    return '-'
  }
  return dayjs(Number(value)).format('YYYY-MM-DD HH:mm:ss')
}

getTableData()
</script>
