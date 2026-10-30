// ========== v-auth 按钮权限指令 ==========
// 用法：<el-button v-auth="'defect:review'">评审</el-button>
//
// 底层原理：
//   1. 后端登录后经 /auth/me 下发当前用户的权限点集合 permissions[]（由外部权限系统
//      按「角色 → 按钮权限值」配置计算），存入 stores/user.ts 的 userInfo.permissions。
//   2. 本指令在元素挂载（mounted）与更新（updated）时读取指令值（权限编码），
//      调用 store.hasPerm(code) 判断当前用户是否具备该权限。
//   3. 不具备 → 直接从 DOM 移除该元素（等价于"看不见按钮"，不占位）；
//      admin 角色默认全权限（store.hasPerm 内判断）。
//
// 对接说明：
//   - 模板中的权限编码必须与后端权限系统配置的按钮权限值一致（如 defect:review）。
//   - 权限点下发来源：后端 permission.Service（对接外部系统时替换实现，契约不变）。
import type { Directive, DirectiveBinding } from 'vue'
import { useUserStore } from '@/stores/user'

export const auth: Directive = {
  mounted(el: HTMLElement, binding: DirectiveBinding<string>) {
    if (!hasAuth(binding.value)) {
      el.parentNode?.removeChild(el)
    }
  },
  updated(el: HTMLElement, binding: DirectiveBinding<string>) {
    if (!hasAuth(binding.value)) {
      el.parentNode?.removeChild(el)
    }
  },
}

function hasAuth(code: string): boolean {
  const userStore = useUserStore()
  return userStore.hasPerm(code)
}
