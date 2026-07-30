<template>
  <div class="sca-import">
    <el-alert title="SCA自动化解析入库" description="上传 SBOM 文件或源码包进行安全解析与合规校验" type="info" :closable="false" show-icon style="margin-bottom: 16px" />

    <el-card>
      <template #header>手动应急通道</template>
      <el-upload drag action="#" :auto-upload="false" style="max-width: 500px">
        <el-icon :size="48"><UploadFilled /></el-icon>
        <div class="el-upload__text">将 SBOM 文件或源码包拖到此处，或 <em>点击上传</em></div>
        <template #tip>
          <div class="el-upload__tip">支持 SPDX/CycloneDX 格式的 SBOM 文件</div>
        </template>
      </el-upload>
    </el-card>

    <el-divider content-position="left">解析记录</el-divider>
    <el-table :data="[]" stripe>
      <el-table-column prop="fileName" label="文件名" />
      <el-table-column prop="status" label="解析状态" />
      <el-table-column prop="componentCount" label="组件数" />
      <el-table-column prop="createdAt" label="上传时间" />
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { UploadFilled } from '@element-plus/icons-vue'
</script>
