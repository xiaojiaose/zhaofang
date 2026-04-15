<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="form">
        <el-form-item label="经纪人电话">
          <el-input v-model="form.userPhone" clearable />
        </el-form-item>
        <el-form-item label="增加次数">
          <el-input-number v-model="form.amount" :min="1" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleGrant">提交</el-button>
        </el-form-item>
      </el-form>
      <el-divider />
      <el-form :inline="true" :model="searchInfo">
        <el-form-item label="经纪人电话">
          <el-input v-model="searchInfo.userPhone" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <el-table :data="tableData">
        <el-table-column prop="userPhone" label="经纪人电话" min-width="150" />
        <el-table-column label="操作类型" min-width="100">
          <template #default="scope">
            {{ scope.row.action === 'grant' ? '新增次数' : '消耗次数' }}
          </template>
        </el-table-column>
        <el-table-column prop="changeAmount" label="变更次数" min-width="100" />
        <el-table-column prop="usedAmount" label="使用次数" min-width="100" />
        <el-table-column prop="remainingAmount" label="剩余次数" min-width="100" />
        <el-table-column prop="relatedResourceId" label="关联房源ID" min-width="110" />
        <el-table-column prop="remark" label="备注" min-width="220" show-overflow-tooltip />
        <el-table-column label="创建时间" min-width="180">
          <template #default="scope">
            {{ formatTime(scope.row.CreatedAt) }}
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
import { ElMessage } from 'element-plus'
import { getContactQuotaList, grantContactQuota } from '@/api/contactQuota'

defineOptions({
  name: 'ContactQuotaList'
})

const defaultForm = () => ({
  userPhone: '',
  amount: 1,
  remark: ''
})

const form = ref(defaultForm())
const searchInfo = ref({
  userPhone: ''
})
const tableData = ref([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const handleGrant = async () => {
  const res = await grantContactQuota(form.value)
  if (res.code === 0) {
    ElMessage.success('添加成功')
    form.value = defaultForm()
    getTableData()
  }
}

const getTableData = async () => {
  const res = await getContactQuotaList({
    page: page.value,
    pageSize: pageSize.value,
    ...searchInfo.value
  })
  if (res.code === 0) {
    tableData.value = res.data.list
    total.value = res.data.total
  }
}

const handleSearch = () => {
  page.value = 1
  getTableData()
}

const handleReset = () => {
  searchInfo.value.userPhone = ''
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
