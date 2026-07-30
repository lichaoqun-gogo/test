import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

import App from './App.vue'
import router from './router'

// 引入全局样式
import '@/styles/global.scss'

const app = createApp(App)

// 注册所有图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

const pinia = createPinia()
app.use(pinia)
app.use(router)
app.use(ElementPlus, { locale: zhCn })

// 注册全局业务组件
import GradeTag from '@/components/business/GradeTag.vue'
import RiskTag from '@/components/business/RiskTag.vue'
import TicketStatusTag from '@/components/business/TicketStatusTag.vue'
import TaskStatusTag from '@/components/business/TaskStatusTag.vue'
import PreEvalProgress from '@/components/business/PreEvalProgress.vue'
import ActivityTrendChart from '@/components/business/ActivityTrendChart.vue'

app.component('GradeTag', GradeTag)
app.component('RiskTag', RiskTag)
app.component('TicketStatusTag', TicketStatusTag)
app.component('TaskStatusTag', TaskStatusTag)
app.component('PreEvalProgress', PreEvalProgress)
app.component('ActivityTrendChart', ActivityTrendChart)

app.mount('#app')
