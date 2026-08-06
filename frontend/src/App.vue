<template>
  <router-view :key="layoutKey" />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
const route = useRoute()
const layoutKey = computed(() => route.matched[0]?.path || route.path)
</script>
