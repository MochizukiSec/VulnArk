<template>
  <!-- No changes to template section -->
</template>

<script>
export default {
  data() {
    return {
      menuItems: [
        {
          title: '仪表盘',
          icon: 'el-icon-monitor',
          route: '/dashboard'
        },
        {
          title: '漏洞管理',
          icon: 'el-icon-warning',
          route: '/vulnerabilities'
        },
        {
          title: 'AI智能分析',
          icon: 'el-icon-cpu',
          route: '/ai-analysis'
        },
        {
          title: '报告中心',
          icon: 'el-icon-document',
          children: [
            {
              title: '生成报告',
              route: '/reports/generate'
            },
            {
              title: '报告历史',
              route: '/reports/history'
            }
          ]
        },
        {
          title: '系统管理',
          icon: 'el-icon-setting',
          children: [
            {
              title: '用户管理',
              route: '/admin/users',
              requiresAdmin: true
            },
            {
              title: '系统设置',
              route: '/admin/settings',
              requiresAdmin: true
            }
          ]
        }
      ],
      // ... rest of the existing code
    };
  },
  methods: {
    // ... existing methods
  }
};
</script>

<style>
  /* No changes to style section */
</style> 