<template>
  <div class="applications-container">
    <AppHeader />
    <div class="header">
      <h2>申请列表</h2>
      <div class="header-right">
        <el-input
          v-model="searchQuery"
          placeholder="搜索公司名称"
          class="search-input"
          clearable
          @clear="handleSearch"
          @input="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-button type="primary" @click="showCreateDialog">新建申请</el-button>
      </div>
    </div>

    <el-table :data="filteredApplications" v-loading="loading" style="width: 100%" :empty-text="emptyText">
      <el-table-column prop="company" label="公司" min-width="120" />
      <el-table-column prop="position" label="职位" min-width="120" />
      <el-table-column prop="status" label="状态" width="100">
        <template #header>
          <el-popover
            placement="bottom"
            :width="200"
            trigger="click"
          >
            <template #reference>
              <span>状态 <el-icon><ArrowDown /></el-icon></span>
            </template>
            <el-checkbox-group v-model="selectedStatuses" @change="handleStatusFilter">
              <el-checkbox v-for="option in statusOptions" :key="option.value" :label="option.value">
                {{ option.label }}
              </el-checkbox>
            </el-checkbox-group>
          </el-popover>
        </template>
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)" :style="getStatusStyle(row.status)">
            {{ getStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="event_link" label="面试链接" min-width="150">
        <template #default="{ row }">
          <el-link v-if="row.event_link && row.event_link !== '无'" type="primary" :href="row.event_link" target="_blank"
            :underline="false">
            <el-icon>
              <Link />
            </el-icon>
            {{ row.event_link }}
          </el-link>
          <span v-else>{{ row.event_link }}</span>
        </template>
      </el-table-column>

      <el-table-column type="expand" fixed="right" width="50">
        <template #default="props">
          <el-descriptions class="expanded-info" :column="1" border>
            <el-descriptions-item label="备注">{{ props.row.Notes }}</el-descriptions-item>
          </el-descriptions>
        </template>
      </el-table-column>

      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button-group>
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
            <el-button size="small" type="success" @click="handleUpdateStatus(row)">更新状态</el-button>
          </el-button-group>
        </template>
      </el-table-column>
    </el-table>

    <!-- 添加分页组件 -->
    <div class="pagination-container">
      <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize" :page-sizes="[10, 15, 20]"
        :total="total" layout="total, sizes, prev, pager, next, jumper" @size-change="handleSizeChange"
        @current-change="handleCurrentChange" />
    </div>

    <!-- 创建/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑申请' : '新建申请'" width="600px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <!-- 必填字段 -->
        <el-form-item label="公司" prop="company">
          <el-input v-model="form.company" placeholder="请输入公司名称" />
        </el-form-item>
        <el-form-item label="职位" prop="position">
          <el-input v-model="form.position" placeholder="请输入职位名称" />
        </el-form-item>


        <!-- 可选字段 -->
        <el-divider>可选信息</el-divider>
        <el-form-item label="面试链接" prop="event_link">
          <el-input v-model="form.event_link" placeholder="请输入面试链接" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.notes" type="textarea" :rows="3" placeholder="请输入备注信息" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit" :loading="submitting">
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 更新状态对话框 -->
    <el-dialog v-model="statusDialogVisible" title="更新状态" width="400px">
      <el-form :model="statusForm" label-width="80px">
        <el-form-item label="状态">
          <el-select v-model="statusForm.status" placeholder="请选择状态">
            <el-option v-for="item in statusOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="statusDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitStatusUpdate" :loading="submitting">
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Link, Search, ArrowDown } from '@element-plus/icons-vue'
import { applicationApi } from '../api/application'
import { formatDate } from '../utils/date'
import AppHeader from '../components/AppHeader.vue'

const applications = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const statusDialogVisible = ref(false)
const submitting = ref(false)
const isEdit = ref(false)
const formRef = ref(null)
const currentApplication = ref(null)
const emptyText = ref('暂无数据')

// 分页相关的响应式变量
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

const form = reactive({
  company: '',
  position: '',
  event_link: '',
  notes: ''
})

const statusForm = reactive({
  status: ''
})

const rules = {
  company: [{ required: true, message: '请输入公司名称', trigger: 'blur' }],
  position: [{ required: true, message: '请输入职位名称', trigger: 'blur' }],
  event_link: [{ required: false, message: '请输入链接', trigger: 'blur' }]
}

const statusOptions = [
  { label: '已投递', value: 'submitted' },
  { label: '笔试中', value: 'written' },
  { label: '面试中', value: 'interview' },
  { label: '已录用', value: 'accepted' },
  { label: '已拒绝', value: 'rejected' }
]

const getStatusType = (status) => {
  const types = {
    submitted: 'primary',    // 蓝色
    written: 'warning',      // 橙色
    interview: '',           // 空字符串表示使用自定义颜色
    accepted: 'success',     // 绿色
    rejected: 'danger'       // 红色
  }
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = {
    submitted: '已投递',
    written: '笔试中',
    interview: '面试中',
    accepted: '已录用',
    rejected: '已拒绝'
  }
  return texts[status] || status
}

const getStatusStyle = (status) => {
  if (status === 'interview') {
    return { backgroundColor: '#A020F0', color: '#ffffff' }
  }
  return {}
}

// 搜索和筛选相关的响应式变量
const searchQuery = ref('')
const selectedStatuses = ref(['submitted', 'written', 'interview', 'accepted'])

// 直接使用后端数据，不再需要本地过滤
const filteredApplications = computed(() => {
  return applications.value
})

const fetchApplications = async () => {
  try {
    loading.value = true
    emptyText.value = '数据加载中...'

    const response = await applicationApi.getApplications({
      page: currentPage.value,
      pageSize: pageSize.value,
      search: searchQuery.value,
      statuses: selectedStatuses.value.join(',')
    })

    if (!response || !response.applications) {
      emptyText.value = '暂无申请记录'
      applications.value = []
      return
    }

    applications.value = response.applications.map(item => ({
      id: item.ID,
      company: item.company,
      position: item.position,
      status: item.status,
      event_link: item.event_link,
      Notes: item.notes,
    }))

    total.value = response.total
  } catch (error) {
    console.error('获取申请列表失败:', error)
    emptyText.value = '数据加载失败，请刷新重试'
    ElMessage.error('获取申请列表失败：' + (error.message || '未知错误'))
  } finally {
    loading.value = false
  }
}

const showCreateDialog = () => {
  isEdit.value = false
  form.company = ''
  form.position = ''
  form.event_link = ''

  form.notes = ''
  dialogVisible.value = true
}

const handleEdit = (row) => {
  isEdit.value = true
  currentApplication.value = row
  form.company = row.company
  form.position = row.position

  form.event_link = row.event_link

  form.notes = row.Notes
  dialogVisible.value = true
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除这条申请记录吗？', '提示', {
      type: 'warning'
    })
    await applicationApi.deleteApplication(row.id)
    ElMessage.success('删除成功')
    fetchApplications()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
      ElMessage.error(error.response?.data?.error || '删除失败')
    }
  }
}

const handleUpdateStatus = (row) => {
  currentApplication.value = row
  statusForm.status = row.status
  statusDialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    submitting.value = true

    const submitData = {
      company: form.company,
      position: form.position,
      event_link: form.event_link,
      notes: form.notes || '无'
    }

    if (isEdit.value) {
      await applicationApi.updateApplication(currentApplication.value.id, submitData)
      ElMessage.success('更新成功')
    } else {
      await applicationApi.createApplication(submitData)
      ElMessage.success('创建成功')
    }

    dialogVisible.value = false
    fetchApplications()
  } catch (error) {
    console.error('操作失败:', error)
    ElMessage.error(error.response?.data?.error || '操作失败')
  } finally {
    submitting.value = false
  }
}

const submitStatusUpdate = async () => {
  try {
    submitting.value = true
    await applicationApi.updateStatus({
      id: currentApplication.value.id,
      status: statusForm.status
    })
    ElMessage.success('状态更新成功')
    statusDialogVisible.value = false
    fetchApplications()
  } catch (error) {
    console.error('状态更新失败:', error)
    ElMessage.error(error.response?.data?.error || '状态更新失败')
  } finally {
    submitting.value = false
  }
}

// 处理页码改变
const handleCurrentChange = (val) => {
  currentPage.value = val
  fetchApplications()
}

// 处理每页条数改变
const handleSizeChange = (val) => {
  pageSize.value = val
  currentPage.value = 1 // 重置到第一页
  fetchApplications()
}

// 处理搜索
const handleSearch = () => {
  currentPage.value = 1 // 重置到第一页
  fetchApplications()
}

// 处理状态筛选
const handleStatusFilter = () => {
  currentPage.value = 1 // 重置到第一页
  fetchApplications()
}

onMounted(() => {
  fetchApplications()
})
</script>

<style scoped>
.applications-container {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.search-input {
  width: 200px;
}

.expanded-info {
  margin: 20px;
}

:deep(.el-descriptions__label) {
  width: 120px;
  font-weight: bold;
  background-color: #f5f7fa;
}

:deep(.el-descriptions__content) {
  padding: 12px;
}

:deep(.el-descriptions-item__content) {
  white-space: pre-wrap;
}

:deep(.el-link) {
  display: flex;
  align-items: center;
  gap: 5px;
}

:deep(.el-table__expand-icon) {
  transform: rotate(90deg);
}

:deep(.el-table__expand-icon .el-icon) {
  font-size: 12px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

:deep(.el-checkbox-group) {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

:deep(.el-table-column__header) {
  cursor: pointer;
}

:deep(.el-table-column__header .el-icon) {
  margin-left: 4px;
}
</style>