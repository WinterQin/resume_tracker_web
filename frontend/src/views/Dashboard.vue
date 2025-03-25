<template>
  <div class="app-container">
    <AppHeader />
    <div class="dashboard-container">
      <!-- 上半部分 -->
      <el-row :gutter="12">
        <!-- 左侧统计卡片 2x2 布局 -->
        <el-col :span="14">
          <el-row :gutter="12">
            <el-col :span="12">
              <el-card class="stat-card">
                <template #header>
                  <div class="card-header">
                    <span>总申请数</span>
                  </div>
                </template>
                <div class="stat-value">{{ statusCount.total }}</div>
              </el-card>
            </el-col>
            <el-col :span="12">
              <el-card class="stat-card">
                <template #header>
                  <div class="card-header">
                    <span>已录用</span>
                  </div>
                </template>
                <div class="stat-value">{{ statusCount.accepted }}</div>
              </el-card>
            </el-col>
            
          </el-row>

          <el-row :gutter="12" class="mt-12">
            <el-col :span="12">
              <el-card class="stat-card">
                <template #header>
                  <div class="card-header">
                    <span>面试/笔试中</span>
                  </div>
                </template>
                <div class="stat-value">{{ statusCount.interview + statusCount.written }}</div>
              </el-card>
            </el-col>

            <el-col :span="12">
              <el-card class="stat-card">
                <template #header>
                  <div class="card-header">
                    <span>已拒绝</span>
                  </div>
                </template>
                <div class="stat-value">{{ statusCount.rejected }}</div>
              </el-card>
            </el-col>
          </el-row>
        </el-col>

        <!-- 右侧饼图 -->
        <el-col :span="10">
          <el-card class="chart-card">
            <template #header>
              <div class="card-header">
                <span>状态分布</span>
              </div>
            </template>
            <div ref="chartRef" style="height: 300px"></div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 下半部分：最近申请列表 -->
      <el-row class="mt-12">
        <el-col :span="24">
          <el-card>
            <template #header>
              <div class="card-header">
                <span>最近申请</span>
                <el-button text @click="$router.push('/applications')">查看全部</el-button>
              </div>
            </template>
            <el-table :data="recentApplications" style="width: 100%" align="right">
              <el-table-column prop="company" label="公司" width="220" align="right" />
              <el-table-column prop="position" label="职位" width="220" align="right" />
              <el-table-column prop="apply_url" label="投递网址" align="center">
                <template #default="{ row }">
                  <el-link type="primary" :href="row.event_link" target="_blank" :underline="false"
                    v-if="row.event_link">
                    <el-icon>
                      <Link />
                    </el-icon>
                    {{ row.event_link }}
                  </el-link>
                </template>
              </el-table-column>
              <el-table-column prop="status" label="状态" width="120" align="left">
                <template #default="{ row }">
                  <el-tag :type="getStatusType(row.status)" :style="getStatusStyle(row.status)">
                    {{ getStatusText(row.status) }}
                  </el-tag>
                </template>
              </el-table-column>
              <!-- <el-table-column prop="apply_date" label="申请时间" width="120" align="left">
                <template #default="{ row }">
                  {{ formatDate(row.apply_date) }}
                </template>
              </el-table-column> -->
            </el-table>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Link } from '@element-plus/icons-vue'
import { applicationApi } from '../api/application'
import { formatDate } from '../utils/date'
import * as echarts from 'echarts'
import { ElMessage } from 'element-plus'
import AppHeader from '../components/AppHeader.vue'


const recentApplications = ref([])
const chartRef = ref(null)
let chart = null
const loading = ref(false)

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

// 添加自定义颜色样式
const getStatusStyle = (status) => {
  if (status === 'interview') {
    return { backgroundColor: '#A020F0', color: '#ffffff' }
  }
  return {}
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


const statusCount = reactive({
  total: 0,
  submitted:0,
  written: 0,
  interview: 0,
  accepted: 0,
  rejected: 0
})

// const statusCount = {
//     total:0,
//     submitted: 0,
//     written: 0,
//     interview: 0,
//     accepted: 0,
//     rejected: 0
//   }

const fetchData = async () => {
  try {
    loading.value = true
    const response = await applicationApi.getRecentApplications()

    const applications = response.applications || []
    recentApplications.value = applications

    // 获取统计数据仍然需要所有申请记录
    const allApplicationsResponse = await applicationApi.getStatistics()
    const StatisticsInfo = allApplicationsResponse.statistics

    console.log(StatisticsInfo)

    // 计算统计数据
    statusCount.total = StatisticsInfo.submitted+ StatisticsInfo.written+StatisticsInfo.interview+StatisticsInfo.accepted+StatisticsInfo.rejected
    statusCount.submitted = StatisticsInfo.submitted
    statusCount.written = StatisticsInfo.written
    statusCount.interview = StatisticsInfo.interview
    statusCount.accepted = StatisticsInfo.accepted
    statusCount.rejected = StatisticsInfo.rejected

    updateChart(StatisticsInfo)
  } catch (error) {
    console.error('获取数据失败:', error)
    ElMessage.error('获取数据失败，请刷新重试')
  } finally {
    loading.value = false
  }
}

const updateChart = (StatisticsInfo) => {
  if (!chartRef.value) return




  // 计算统计数据
  // statusCount.written = stats.written
  // statusCount.interview = stats.interview
  // statusCount.accepted = stats.accepted
  // statusCount.rejected = stats.rejected
  // statusCount.submitted = stats.submitted

  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      left: 'left'
    },
    series: [
      {
        type: 'pie',
        radius: '50%',
        data: [
          {
            value: statusCount.submitted,
            name: '已投递',
            itemStyle: { color: '#409EFF' }  // 蓝色
          },
          {
            value: statusCount.written,
            name: '笔试中',
            itemStyle: { color: '#E6A23C' }  // 橙色
          },
          {
            value: statusCount.interview,
            name: '面试中',
            itemStyle: { color: '#A020F0' }  // 紫色
          },
          {
            value: statusCount.accepted,
            name: '已录用',
            itemStyle: { color: '#67C23A' }  // 绿色
          },
          {
            value: statusCount.rejected,
            name: '已拒绝',
            itemStyle: { color: '#F56C6C' }  // 红色
          }
        ],
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)'
          }
        }
      }
    ]
  }

  if (!chart) {
    chart = echarts.init(chartRef.value)
  }
  chart.setOption(option)
}

onMounted(() => {
  fetchData()
  window.addEventListener('resize', () => {
    chart?.resize()
  })
})
</script>

<style scoped>
.app-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.dashboard-container {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
}

.upper-section {
  margin-bottom: 20px;
}

.stat-card {
  height: 180px;
  margin-bottom: 12px;
}

.chart-card {
  height: 382px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stat-value {
  font-size: 36px;
  font-weight: bold;
  color: #409EFF;
  text-align: center;
  margin-top: 20px;
}

.mt-12 {
  margin-top: 12px;
}

:deep(.el-card__body) {
  padding: 15px;
}

:deep(.el-link) {
  display: flex;
  align-items: center;
  gap: 5px;
}

:deep(.el-table) {
  margin: 0;
}
</style>