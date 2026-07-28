import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/login',
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/Login.vue'),
    },
    {
      path: "/admin",
      component: () => import("@/layouts/AdminLayout.vue"),
      meta: {
        requiresAdmin: true,
      },
      children: [
        {
          path: "",
          name: "admin-dashboard",
          component: () => import("@/views/admin/dashboard.vue"),
        },
        {
          path: "users",
          name: "users",
          component: () => import("@/views/admin/users.vue"),
        },
        {
          path: "endpoints",
          name: "endpoints",
          component: () => import("@/views/admin/endpoints.vue"),
        },
        {
          path: "groups",
          name: "groups",
          component: () => import("@/views/admin/groups.vue"),
        },
        {
          path: "config-types",
          name: "config-types",
          component: () => import("@/views/admin/configtypes.vue"),
        },
        {
          path: "config-types/:id",
          name: "config-type-edit",
          component: () => import("@/views/admin/configtype-edit.vue"),
        },
        {
          path: "config-types/create",
          name: "config-type-create",
          component: () => import("@/views/admin/configtype-create.vue"),
        },
        {
          path: "configs",
          name: "configs",
          component: () => import("@/views/admin/configs.vue"),
        },
        {
          path: "configs/:id",
          name: "admin-config-edit",
          component: () => import("@/views/config-edit.vue"),
        },

      ],
    },
    {
      path: '/user',
      component: () => import('@/layouts/UserLayout.vue'),
      meta: {
        requiresUser: true,
      },
      children: [
        {
          path: "",
          name: "user-dashboard",
          component: () => import("@/views/user/dashboard.vue"),
        },
        {
          path: 'configs',
          name: 'user-configs',
          component: () => import('@/views/user/configs.vue'),
        },
        {
          path: "configs/:id",
          name: "user-config-edit",
          component: () => import("@/views/config-edit.vue"),
        },
        {
          path: "profile",
          name: "user-profile",
          component: () => import("@/views/profile.vue"),
        },
      ],
    },

  ],
})

router.beforeEach((to) => {
  const auth = useAuthStore()

  if (to.path === '/login') {
    return true
  }

  if (!auth.isAuthenticated) {
    return '/login'
  }

  if (to.meta.requiresAdmin && !auth.user?.is_admin) {
    return '/user'
  }

  if (to.meta.requiresUser && auth.user?.is_admin) {
    return '/admin'
  }

  return true
})

export default router
