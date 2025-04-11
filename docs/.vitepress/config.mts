import { defineConfig } from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "AIFlowy",
  description: "AIFlowy Site",
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: '开发文档', link: '/zh/development/info/what-is-aiflowy' },
      { text: '产品文档', link: '/zh/product/info/what-is-aiflowy' },
      { text: '更新记录', link: 'https://gitee.com/aiflowy/aiflowy/commits/main' },
    ],

    sidebar:{
      '/zh/development/': { base: '/zh/development/', items: sidebarDevelopment() },
      '/zh/product/': { base: '/zh/product/', items: sidebarProduct() }
    },

    socialLinks: [
      { icon: 'github', link: 'https://github.com/vuejs/vitepress' }
    ]
  }
})



function sidebarDevelopment(): DefaultTheme.SidebarItem[] {
  return [
    {
      text: '快速开始',
      collapsed: false,
      items: [
        { text: '快速开始', link: 'getting-started/getting-started' },
        { text: '目录结构', link: 'getting-started/directory-structure' },
        { text: '配置文件', link: 'getting-started/config-file' },
        { text: '项目部署', link: 'getting-started/deploy' },
        { text: '常见问题', link: 'getting-started/questions' }
      ]
    },
    {
      text: '前端相关',
      collapsed: false,
      items: [
        { text: '前端目录结构', link: 'front/front-structure' },
        { text: '路由管理', link: 'front/routes' },
        { text: '状态管理', link: 'front/state' },
        { text: 'Hooks', link: 'front/hooks' },
        { text: '组件使用', link: 'front/components' },
        { text: '国际化', link: 'front/locales' }
      ]
    },
    {
      text: '后端相关',
      collapsed: false,
      items: [
        { text: 'Controller', link: '/backend/controller' },
        { text: '权限管理', link: '/backend/permission_management' },
        { text: '验证码', link: '/backend/captcha' },
        { text: '数据字典', link: '/backend/data_dictionary' },
        { text: '文件管理', link: '/backend/file_management' },
        { text: '其他', link: '/backend/other' }
      ]
    },
    {
      text: 'AI 相关',
      collapsed: false,
      items: [
        { text: '大语言模型', link: '/ai/language-model' },
        { text: 'Bot 应用', link: '/ai/bot-application' },
        { text: '插件', link: '/ai/plugin' },
        { text: '知识库', link: '/ai/knowledge' },
        { text: 'Ollama', link: '/ai/ollama' },
      ]
    },
  ]
}

function sidebarProduct(): DefaultTheme.SidebarItem[] {
  return [
    {
      text: '简介',
      collapsed: false,
      items: [
        { text: '什么是 AIFlowy？', link: '/zh/product/info/what-is-aiflowy' },
        { text: '快速开始', link: 'getting-started' },
      ]
    },
    {
      text: 'Bot 应用',
      collapsed: false,
      items: [
        { text: '什么是 Bot', link: '/zh/product/info/what-is-aiflowy' },
        { text: '快速开始', link: 'getting-started' },
      ]
    },
    {
      text: '插件',
      collapsed: false,
      items: [
        { text: '什么是 插件？', link: '/zh/product/info/what-is-aiflowy' },
        { text: '快速开始', link: 'getting-started' },
        { text: '自定义插件', link: 'routing' },
      ]
    },
    {
      text: '知识库',
      collapsed: false,
      items: [
        { text: '什么是知识库', link: '/zh/product/info/what-is-aiflowy' },
        { text: '快速开始', link: 'getting-started' },
      ]
    },
    {
      text: '工作流',
      collapsed: false,
      items: [
        { text: '什么是工作流', link: '/zh/product/info/what-is-aiflowy' },
        { text: '快速开始', link: 'getting-started' },
      ]
    },
    {
      text: '模型管理',
      collapsed: false,
      items: [
        { text: '什么是 AIFlowy？', link: '/zh/product/info/what-is-aiflowy' },
        { text: '快速开始', link: 'getting-started' },
        { text: '路由', link: 'routing' },
        { text: '部署', link: 'deploy' }
      ]
    },
  ]
}