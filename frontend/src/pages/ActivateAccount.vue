<template>
  <div class="activation-container">
    <h2 v-if="message">{{ message }}</h2>
    <div v-if="isLoading" class="loading">Активация...</div>
  </div>
</template>

<script>

export default {
  data() {
    return {
      message: "",
      isLoading: true,
    };
  },
  async mounted() {
    try {
      const status = this.$route.query.status;
      this.showMessage(status || "Активация успешна!", true);
    } catch (error) {
      this.handleActivationError(error);
    }
  },
  methods: {
    showMessage(msg, success) {
      this.message = msg;
      this.isLoading = false;
      if (success) {
        setTimeout(() => this.$router.push("/login"), 3000);
      }
    },
    handleActivationError(error) {
      if (error.response) {
        // Сервер ответил, но статус 4xx/5xx
        if (error.response.status === 400) {
          this.showMessage("Неверный или просроченный токен.", false);
        } else if (error.response.status === 404) {
          this.showMessage("Пользователь не найден.", false);
        } else {
          this.showMessage("Ошибка сервера. Попробуйте позже.", false);
        }
      } else {
        // Ошибка сети или сервера нет
        this.showMessage("Сервер недоступен. Проверьте соединение.", false);
      }
    },
  },
};
</script>

<style scoped>
.activation-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
}

.loading {
  font-size: 18px;
  color: teal;
}
</style>
