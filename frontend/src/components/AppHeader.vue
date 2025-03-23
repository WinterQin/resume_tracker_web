<template>
  <el-header class="app-header">
    <div class="header-left">
      <router-link to="/dashboard" class="logo-container">
        <img src="@/assets/vue.svg" alt="Logo" class="logo-image"/>
        <span class="logo-text">求职跟踪系统</span>
      </router-link>
    </div>

    <!-- <div class="header-center">
      <el-menu
        mode="horizontal"
        :router="true"
        class="nav-menu"
      >
        <el-menu-item index="/">首页</el-menu-item>
        <el-menu-item index="/profile">个人页面</el-menu-item>
      </el-menu>
    </div> -->

    <div class="header-right">
      <el-dropdown @command="handleCommand">
        <span class="user-dropdown">
          Hello, {{ username }}
          <el-icon><CaretBottom /></el-icon>
        </span>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="profile">个人页面</el-dropdown-item>
            <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </el-header>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { CaretBottom } from '@element-plus/icons-vue'

const router = useRouter()
const username = ref('')

// 从 localStorage 获取用户信息
const getUserInfo = () => {
  const userInfo = localStorage.getItem('userInfo')
  if (userInfo) {
    const user = JSON.parse(userInfo)
    username.value = user.username // 或者根据实际的用户信息结构使用相应的字段
  }
}

// 组件挂载时获取用户信息
onMounted(() => {
  getUserInfo()
})

const handleCommand = (command) => {
  switch (command) {
    case 'profile':
      router.push('/profile')
      break
    case 'logout':
      // 登出时清除本地存储的用户信息
      localStorage.removeItem('token')
      localStorage.removeItem('userInfo')
      router.push('/login')
      break
  }
}
</script>

<style scoped>
.app-header {
  height: 60px;
  background-color: #fff;
  border-bottom: 1px solid #dcdfe6;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.header-left {
  display: flex;
  align-items: center;
}

.logo-container {
  display: flex;
  align-items: center;
  text-decoration: none;
  color: #303133;
}

.logo-image {
  height: 32px;
  width: 32px;
  margin-right: 10px;
}

.logo-text {
  font-size: 20px;
  font-weight: bold;
}

.header-center {
  flex: 1;
  display: flex;
  justify-content: center;
}

.nav-menu {
  border-bottom: none;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-dropdown {
  display: flex;
  align-items: center;
  cursor: pointer;
  color: #606266;
}

.user-dropdown:hover {
  color: #409EFF;
}

.el-icon {
  margin-left: 4px;
}
</style> 