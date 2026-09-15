<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item label="当前状态">
          <el-select v-model="searchInfo.status" clearable placeholder="请选择">
            <el-option label="待审核" value="待审核" />
            <el-option label="已付款" value="已付款" />
            <el-option label="未通过，已拒绝" value="未通过，已拒绝" />
          </el-select>
        </el-form-item>
        <el-form-item label="发房经纪人电话">
          <el-input v-model="searchInfo.publisherPhone" clearable placeholder="请输入发房经纪人电话" />
        </el-form-item>
        <el-form-item label="查看经纪人电话">
          <el-input v-model="searchInfo.viewerPhone" clearable placeholder="请输入查看经纪人电话" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="handleSearch">搜索</el-button>
          <el-button icon="refresh" @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="success" :disabled="!multipleSelection.length" @click="handleBatchAction('paid')">
          批量已付款
        </el-button>
        <el-button type="danger" :disabled="!multipleSelection.length" @click="handleBatchAction('reject')">
          批量未通过
        </el-button>
      </div>
      <el-table :data="tableData" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="48" />
        <el-table-column prop="id" label="编号" min-width="80" />
        <el-table-column prop="xiaoqu" label="小区" min-width="200" />
        <el-table-column prop="publisherName" label="发房中介" min-width="150" />
        <el-table-column prop="publisherPhone" label="发房中介电话" min-width="150" />
        <el-table-column label="录入时间" min-width="180">
          <template #default="scope">
            {{ formatTime(scope.row.resourceAt) }}
          </template>
        </el-table-column>
        <el-table-column prop="viewerName" label="查看中介" min-width="150" />
        <el-table-column prop="viewerPhone" label="查看中介电话" min-width="150" />
        <el-table-column label="查看时间" min-width="180">
          <template #default="scope">
            {{ formatTime(scope.row.viewAt) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="当前状态" min-width="150">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)" size="small">
              {{ scope.row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="最后操作时间" min-width="180">
          <template #default="scope">
            {{ formatTime(scope.row.lastOperateAt) }}
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
import { getHouseAuditList, houseAuditAction, houseAuditBatchAction } from '@/api/houseAudit'

defineOptions({
  name: 'HouseAuditList'
})

const defaultSearchInfo = () => ({
  status: '',
  publisherPhone: '',
  viewerPhone: ''
})

const searchInfo = ref(defaultSearchInfo())
const tableData = ref([])
const multipleSelection = ref([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const getStatusType = (status) => {
  if (status === '待审核') return 'info'
  if (status === '已通过' || status === '已通过，已付款') return 'success'
  if (status === '未通过，已拒绝') return 'danger'
  return 'info'
}

const getTableData = async () => {
  const res = await getHouseAuditList({
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
  const res = await houseAuditAction({ id: row.ID, action })
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
  const results = await Promise.all(
    multipleSelection.value.map((item) => {
      return houseAuditBatchAction({
        id: item.id,
        action
      })
    })
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

getTableData()
</script>

<style lang="scss" scoped></style>
