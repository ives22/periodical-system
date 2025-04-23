
import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
//  引入svg需要用到的插件
import { createSvgIconsPlugin } from "vite-plugin-svg-icons";
import * as path from "node:path";

import { viteMockServe } from 'vite-plugin-mock'
import { createHtmlPlugin } from 'vite-plugin-html'

// https://vite.dev/config/
export default defineConfig(({ command, mode }) => {
  // 获取各个环境下的对应的环境变量
  let env = loadEnv(mode, process.cwd());
  return {
    plugins: [
      vue(),
      vueDevTools(),
      createSvgIconsPlugin({
        iconDirs: [path.resolve(process.cwd(), 'src/assets/icons')],
        symbolId: 'icon-[dir]-[name]'
      }),
      viteMockServe({
        localEnabled: command === 'serve',
      }),
      // 修改HTML插件配置
      createHtmlPlugin({
        inject: {
          data: {
            title: '期刊管理系统'
          }
        }
      })
    ],
    resolve: {
      // alias: {
      //   '@': fileURLToPath(new URL('./src', import.meta.url))
      // },
      alias: {
        '@': path.resolve(__dirname, 'src')
      },
    },
    // scss 全局变量配置
    css: {
      preprocessorOptions: {
        scss: {
          // javascriptEnabled: true,
          additionalData: `@use "@/styles/variables.scss" as *;`,
        }
      }
    },
    server: {
      //代理跨域
      proxy: {
        [env.VITE_APP_BASE_API]: {
          //获取数据服务器地址的设置
          target: env.VITE_SERVE,
          //需要代理跨域
          changeOrigin: true,
          //路径重写
          rewrite: (path) => path.replace(/^\/api/, ''),
        },
      }
    }
  }
})


