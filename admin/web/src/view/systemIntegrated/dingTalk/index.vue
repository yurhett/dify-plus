<template>
  <div class="system">
    <el-form
      ref="form"
      :model="config"
      label-width="240px"
    >
      <div class="page-header mb-6">
        <h2 class="text-xl font-bold">
          钉钉应用集成配置
        </h2>
        <p class="text-gray-500 mt-2">
          配置钉钉扫码登录相关参数
        </p>
      </div>

      <el-tabs class="dingtalk-tabs">
        <div class="card">
          <div class="card-header flex items-center justify-between">
            <span class="text-lg font-medium">启用状态</span>
            <div class="flex items-center">
              <el-switch
                v-model="config.status"
                active-text="已启用"
                :disabled="!isConfigValid"
                @change="handleStatusChange"
              />
            </div>
          </div>

          <el-divider />

          <div class="card-section">
            <div class="section-title">
              钉钉回调域名配置
            </div>
            <div class="text-gray-600 mb-3">
              <p>回调域名：此信息将在创建钉钉扫码登录应用时使用，可至<span class="text-blue-500 cursor-pointer" @click="goToSecuritySettings">开发配置-安全设置</span>进行修改</p>
            </div>

            <div class="flex items-center">
              <el-input v-model="host" disabled readonly class="flex-1" />
              <el-button type="primary" class="ml-2" icon="copy-document" @click="copyHost">
                复制
              </el-button>
            </div>
          </div>

          <el-divider />

          <div class="card-section">
            <div class="section-title">
              应用信息配置
            </div>
            <div class="mb-4">
              <el-button v-if="!openEdit" type="primary" class="config-btn" icon="setting" @click="openConfig">
                配置链接应用信息
              </el-button>
            </div>
            <div class="bg-gray-50 dark:bg-slate-800 p-5 border dark:border-slate-700 rounded-lg">
              <div class="flex items-center mb-4">
                <span class="info-label">CorpID:</span>
                <el-input v-if="openEdit" v-model="config.corp_id" class="info-value flex-1" />
                <span v-else class="info-value">{{ config.corp_id || '未配置' }}</span>
              </div>
              <div class="flex items-center mb-4">
                <span class="info-label">AppID:</span>
                <el-input v-if="openEdit" v-model="config.app_id" class="info-value flex-1" />
                <span v-else class="info-value">{{ config.app_id }}</span>
              </div>
              <div class="flex items-center mb-4">
                <span class="info-label">AgentID:</span>
                <el-input v-if="openEdit" v-model="config.agent_id" class="info-value flex-1" />
                <span v-else class="info-value">{{ config.agent_id || '未配置' }}</span>
              </div>
              <div class="flex items-center mb-4">
                <span class="info-label">AppKey:</span>
                <el-input v-if="openEdit" v-model="config.app_key" class="info-value flex-1" />
                <span v-else class="info-value">{{ config.app_key || '未配置' }}</span>
              </div>
              <div class="flex items-center mb-4">
                <span class="info-label">AppSecret:</span>
                <el-input v-if="openEdit" v-model="config.app_secret" class="info-value flex-1" />
                <span v-else class="info-value">{{ config.app_secret || '未配置' }}</span>
              </div>
              <div class="float-right">
                <el-button type="primary" plain icon="connection" @click="testConnection">
                  测试连接
                </el-button>
                <el-button v-if="openEdit" type="primary" icon="goods-filled" @click="update">
                  保存
                </el-button>
              </div>
              <div class="clear-both" />
            </div>
          </div>

          <el-divider />

          <div class="card-section">
            <div class="section-title text-amber-500">
              <el-icon><i class="el-icon-warning" /></el-icon>
              <span>温馨提示</span>
            </div>
            <div class="tips-content">
              <p class="tip-item">
                1. 扫码登录应用创建入口：
                <el-link type="primary" href="https://open-dev.dingtalk.com/fe/app" target="_blank">
                  https://open-dev.dingtalk.com/fe/app
                </el-link>
              </p>
              <p class="tip-item">
                2. AppId和AppSecret是扫码登录应用的唯一标识，创建完成后可见
              </p>
              <p class="tip-item">
                查看路径: 钉钉开放平台>应用开发>移动接入应用>扫码登录应用授权应用的列表。
              </p>
            </div>
          </div>
        </div>
      </el-tabs>
    </el-form>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { getSystemDingTalk, setSystemDingTalk } from "@/api/gaia/system";

defineOptions({
  name: 'IntegratedDingTalk',
})

const host = ref("")
const openEdit = ref(false)
const config = ref({
  id: 0,
  classify: 1,
  status: false,
  corp_id: "",
  agent_id: "",
  app_id: "",
  app_key: "",
  app_secret: "",
  test: false,
})

// 验证配置是否有效
const isConfigValid = computed(() => {
  return !!(config.value.corp_id && config.value.agent_id && config.value.app_key && config.value.app_secret);
})

// 处理状态变更
const handleStatusChange = (val) => {
  if (val && !isConfigValid.value) {
    config.value.status = false;
    ElMessage({
      type: 'warning',
      message: '请先填写应用信息配置'
    });
    return;
  }
  update();
}

// 掩码显示文本
const openConfig = () => {
  openEdit.value = true
}

// 复制回调URL
const copyHost = () => {
  navigator.clipboard.writeText(host.value);
  ElMessage({
    type: 'success',
    message: '复制成功'
  });
}

// 测试连接
const testConnection = async () => {
  if (!isConfigValid.value) {
    ElMessage({
      type: 'warning',
      message: '请先填写完整的应用信息配置'
    });
    return;
  }

  config.value.test = true;
  const res = await setSystemDingTalk(config.value)
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '链接成功',
    })
  }
}

const initForm = async() => {
  const res = await getSystemDingTalk()
  if (res.code === 0) {
    host.value = res.data.host
    config.value = res.data.config
  }
}
initForm()

const update = async() => {
  config.value.test = false;

  if (config.value.status && !isConfigValid.value) {
    config.value.status = false;
    ElMessage({
      type: 'warning',
      message: '请先填写应用信息配置'
    });
    return;
  }

  const res = await setSystemDingTalk(config.value)
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '设置成功',
    })
    await initForm()
    openEdit.value = false
  }
}

// 跳转到钉钉安全设置页面
const goToSecuritySettings = () => {
  if (!config.value.app_id || config.value.app_id === '') {
    ElMessage({
      type: 'warning',
      message: '配置AppID后可以正常跳转',
    });
    return;
  }
  
  const url = `https://open-dev.dingtalk.com/fe/ai?hash=#/app/${config.value.app_id}/security#/app/${config.value.app_id}/security`;
  window.open(url, '_blank');
}
</script>

<style lang="scss" scoped>
.system {
  @apply bg-white p-9 rounded dark:bg-slate-900;
}

.card {
  @apply bg-white dark:bg-slate-900 rounded-lg shadow-sm border border-gray-100 dark:border-slate-700 p-5;
}

.card-header {
  @apply mb-4;
}

.card-section {
  @apply py-3;
}

.section-title {
  @apply text-lg font-medium mb-4 flex items-center;
}

.info-label {
  @apply text-gray-600 dark:text-gray-400 w-28;
}

.info-value {
  @apply font-mono text-gray-800 dark:text-gray-200;
}

.tips-content {
  @apply bg-amber-50 dark:bg-amber-900/20 p-4 rounded-lg;
}

.tip-item {
  @apply mb-2 text-gray-700 dark:text-gray-300;
}

.config-btn {
  @apply w-full justify-center;
}

.action-footer {
  @apply flex items-center gap-3;
}

:deep(.el-tabs__nav) {
  @apply mb-5;
}

:deep(.el-divider) {
  @apply my-5;
}
</style>
