
// 对外暴露配置路由
export const constantRoute = [
    {
        // 登录
        path: '/login',
        // component: () => import('@/views/index.vue'),
        component: () => import('@/views/login/index.vue'),
        name: 'login',
        meta: {
            title: '登录', //菜单标题
            hidden: true, //代表路由是否隐藏（true：隐藏，false：不隐藏）
            icon: 'Promotion'  // 菜单左侧的图标
        }
    },
    {
        path: '/',
        component: () => import('@/layout/index.vue'),
        name: 'layout',
        meta: {
            title: '',
            hidden: false,
            icon: ''
        },
        redirect: '/home',
        children: [
            {
                path: '/home',
                component: () => import('@/views/home/index.vue'),
                meta: {
                    title: '首页',
                    hidden: false,
                    icon: 'HomeFilled'
                }
            }
        ]
    },
    {
        path: '/periodical',
        component: () => import('@/layout/index.vue'),
        name: 'Periodical',
        meta: {
            title: '',
            hidden: false,
            icon: 'Orange'
        },
        redirect: '/periodical',
        children: [
            {
                path: '/periodical',
                component: () => import('@/views/periodical/index.vue'),
                meta: {
                    title: '期刊管理',
                    hidden: false,
                    icon: 'Orange'
                }
            }
        ]
    },
]
