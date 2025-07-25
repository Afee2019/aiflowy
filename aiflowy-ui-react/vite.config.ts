import {defineConfig} from 'vite'
import react from '@vitejs/plugin-react'
import commonjs from "@rollup/plugin-commonjs";

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [
        react(),
        commonjs(),
    ],


    build: {
        commonjsOptions: {include: []},
    },
    server: {
        proxy: {
            '^/api': {
                target: 'http://172.16.100.1:8080/',
            },
            '^/attachment/': {
                target: 'http://127.0.0.1:8080/static/',
            },
        },
    },
})
