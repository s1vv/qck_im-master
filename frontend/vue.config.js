const { defineConfig } = require('@vue/cli-service');

module.exports = defineConfig({
  transpileDependencies: true,
  devServer: {
    proxy: {
      '/api': {
        target: 'http://localhost:8081', // Замените на адрес вашего бэкенда
        changeOrigin: true,
      },
    },
    client: {
      webSocketURL: 'ws://localhost:8080/ws', // Явно указываем WebSocket URL
    },
  },
});
