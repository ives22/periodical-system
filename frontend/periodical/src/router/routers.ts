// 对外暴露配置路由（采用单根 Layout 架构，确保跨全站菜单切换时 Layout 与 Keep-Alive 缓存永不销毁重建）
export const constantRoute = [
    {
        // 登录
        path: '/login',
        component: () => import('@/views/login/index.vue'),
        name: 'login',
        meta: {
            title: '登录',
            hidden: true,
            icon: 'Promotion',
            keepAlive: false,
        },
    },
    {
        path: '/ai',
        redirect: '/ai/workspace',
        meta: {
            title: '',
            hidden: true,
        },
    },
    {
        path: '/',
        component: () => import('@/layout/index.vue'),
        name: 'layout',
        meta: {
            title: '',
            hidden: false,
            icon: '',
        },
        redirect: '/home',
        children: [
            {
                path: '/home',
                name: 'Home',
                component: () => import('@/views/home/index.vue'),
                meta: {
                    title: '首页',
                    hidden: false,
                    icon: 'HomeFilled',
                    perm: 'dashboard:view',
                    affix: true,
                    keepAlive: true,
                },
            },
            {
                path: '/periodical',
                name: 'PeriodicalManage',
                component: () => import('@/views/periodical/index.vue'),
                meta: {
                    title: '期刊管理',
                    hidden: false,
                    icon: 'Orange',
                    perm: 'periodical:view',
                    keepAlive: true,
                },
            },
            {
                path: '/ai/workspace',
                name: 'AiWorkspace',
                component: () => import('@/views/ai/workspace/index.vue'),
                meta: {
                    title: 'AI 工作台',
                    hidden: false,
                    icon: 'MagicStick',
                    perm: 'ai:workspace',
                    keepAlive: true,
                },
            },
            {
                path: '/system',
                name: 'System',
                meta: {
                    title: '系统管理',
                    hidden: false,
                    icon: 'Setting',
                    perm: 'system:view',
                },
                redirect: '/system/user',
                children: [
                    {
                        path: '/system/user',
                        name: 'SystemUser',
                        component: () => import('@/views/system/user/index.vue'),
                        meta: {
                            title: '用户管理',
                            hidden: false,
                            icon: 'User',
                            perm: 'user:view',
                            keepAlive: true,
                        },
                    },
                    {
                        path: '/system/role',
                        name: 'SystemRole',
                        component: () => import('@/views/system/role/index.vue'),
                        meta: {
                            title: '角色与权限',
                            hidden: false,
                            icon: 'Lock',
                            perm: 'role:view',
                            keepAlive: true,
                        },
                    },
                    {
                        path: '/system/ai-config',
                        name: 'SystemAiConfig',
                        component: () => import('@/views/system/ai-config/index.vue'),
                        meta: {
                            title: 'AI 模型与审计',
                            hidden: false,
                            icon: 'Cpu',
                            perm: 'ai:config',
                            keepAlive: true,
                        },
                    },
                ],
            },
        ],
    },
]
