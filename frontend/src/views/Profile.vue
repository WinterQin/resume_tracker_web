<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { userApi } from '../api/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import AppHeader from '../components/AppHeader.vue'

const router = useRouter()
const userStore = useUserStore()
const formRef = ref(null)
const passwordFormRef = ref(null)
const isEditing = ref(false)
const updating = ref(false)

// 表单数据
const form = reactive({
  username: '',
  email: '',
  age: 0,
  gender: '',
  phone: ''
})

// 密码表单数据
const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 表单验证规则
const rules = {
  age: [
    { type: 'number', message: '年龄必须为数字' }
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
  ]
}

// 密码表单验证规则
const passwordRules = {
  oldPassword: [
    { required: true, message: '请输入旧密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能小于6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== passwordForm.newPassword) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

// 获取用户信息
const fetchUserProfile = async () => {
  try {
    const data = await userApi.getProfile()
    Object.assign(form, data)
  } catch (error) {
    ElMessage.error('获取用户信息失败')
  }
}

// 编辑个人信息
const handleEdit = () => {
  isEditing.value = true
}

// 取消编辑
const cancelEdit = () => {
  isEditing.value = false
  fetchUserProfile() // 重新获取数据，放弃修改
}

// 保存个人信息
const handleSave = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    await userApi.updateProfile({
      age: form.age,
      gender: form.gender,
      phone: form.phone
    })
    ElMessage.success('保存成功')
    isEditing.value = false
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '保存失败')
  }
}

// 更新密码
const handleUpdatePassword = async () => {
  if (!passwordFormRef.value) return

  try {
    await passwordFormRef.value.validate()
    updating.value = true
    await userApi.updatePassword({
      old_password: passwordForm.oldPassword,
      new_password: passwordForm.newPassword
    })
    ElMessage.success('密码更新成功')
    // 清空表单
    passwordForm.oldPassword = ''
    passwordForm.newPassword = ''
    passwordForm.confirmPassword = ''
    passwordFormRef.value.resetFields()
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '密码更新失败')
  } finally {
    updating.value = false
  }
}

// 删除账号
const handleDeleteAccount = async () => {
  try {
    await ElMessageBox.confirm(
      '此操作将永久删除您的账号，所有相关数据将被清除且无法恢复，是否继续？',
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    await userApi.deleteAccount()
    ElMessage.success('账号已删除')
    userStore.logout()
    router.push('/login')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.error || '删除账号失败')
    }
  }
}

onMounted(() => {
  fetchUserProfile()
})
</script>

<template>
      <AppHeader />
  <div class="profile-container">
    <el-card class="profile-card">
      <template #header>
        <div class="card-header">
          <h2>个人资料</h2>
          <el-button type="primary" @click="handleEdit" v-if="!isEditing">编辑</el-button>
          <div v-else>
            <el-button type="success" @click="handleSave">保存</el-button>
            <el-button @click="cancelEdit">取消</el-button>
          </div>
        </div>
      </template>

      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px" :disabled="!isEditing">
        <el-form-item label="用户名">
          <el-input v-model="form.username" disabled></el-input>
        </el-form-item>

        <el-form-item label="邮箱">
          <el-input v-model="form.email" disabled></el-input>
        </el-form-item>

        <el-form-item label="年龄" prop="age">
          <el-input-number v-model="form.age" :min="0" :max="150"></el-input-number>
        </el-form-item>

        <el-form-item label="性别" prop="gender">
          <el-select v-model="form.gender" placeholder="请选择性别">
            <el-option label="男" value="男"></el-option>
            <el-option label="女" value="女"></el-option>
            <el-option label="其他" value="其他"></el-option>
          </el-select>
        </el-form-item>

        <el-form-item label="手机号码" prop="phone">
          <el-input v-model="form.phone" placeholder="请输入手机号码"></el-input>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="profile-card">
      <template #header>
        <div class="card-header">
          <h2>修改密码</h2>
        </div>
      </template>

      <el-form :model="passwordForm" :rules="passwordRules" ref="passwordFormRef" label-width="100px">
        <el-form-item label="旧密码" prop="oldPassword">
          <el-input v-model="passwordForm.oldPassword" type="password" show-password></el-input>
        </el-form-item>

        <el-form-item label="新密码" prop="newPassword">
          <el-input v-model="passwordForm.newPassword" type="password" show-password></el-input>
        </el-form-item>

        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input v-model="passwordForm.confirmPassword" type="password" show-password></el-input>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleUpdatePassword" :loading="updating">更新密码</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="profile-card danger-zone">
      <template #header>
        <div class="card-header">
          <h2>危险区域</h2>
        </div>
      </template>

      <div class="danger-content">
        <div class="danger-item">
          <div class="danger-text">
            <h3>删除账号</h3>
            <p>此操作将永久删除您的账号，所有相关数据将被清除且无法恢复。</p>
          </div>
          <el-button type="danger" @click="handleDeleteAccount">删除账号</el-button>
        </div>
      </div>
    </el-card>
  </div>
</template>

<style scoped>
.profile-container {
  max-width: 800px;
  margin: 20px auto;
  padding: 0 20px;
}

.profile-card {
  margin-bottom: 20px;
}

.profile-card:last-child {
  margin-bottom: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header h2 {
  margin: 0;
  font-size: 18px;
  color: #303133;
}

.danger-zone {
  border: 1px solid #f56c6c;
}

.danger-zone :deep(.el-card__header) {
  background-color: #fef0f0;
}

.danger-content {
  padding: 20px 0;
}

.danger-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.danger-text h3 {
  margin: 0 0 8px 0;
  color: #f56c6c;
}

.danger-text p {
  margin: 0;
  color: #606266;
  font-size: 14px;
}
</style>
