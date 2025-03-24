<template>
  <form @submit.prevent="saveNote">
    <h4>Активация ссылки</h4>

    <my-input 
      v-model.trim="note.qck_link"
      type="text" 
      placeholder="Qck ссылка"
      :maxlength="8"
    />

    <my-input 
      v-focus
      v-model.trim="note.password"
      type="text"
      placeholder="Пароль ссылки"
      :maxlength="8"
    />

    <!-- Вывод ошибки -->
    <p v-if="errorMessage" class="error">{{ errorMessage }}</p>

    <my-button 
      class="btn" 
      style="align-self: flex-end; margin-top: 15px;"
      :disabled="isLoading"
    >
      {{ isLoading ? "Активация..." : "Активировать" }}
    </my-button>
  </form>
</template>

<script>
import axios from "axios";

export default {
  data() {
    return {
      note: {
        qck_link: "",
        password: "",
      },
      errorMessage: "", // Ошибка от сервера
      isLoading: false, // Флаг загрузки
    };
  },
  methods: {
    async saveNote() {
      this.errorMessage = ""; // Очищаем ошибку перед запросом

      if (!this.note.qck_link.trim() || this.note.qck_link.length!=8){
        this.errorMessage = "Должно быть 8 символов";
        return;
      }

      if (!this.note.password.trim() || this.note.password.length!=8){
        this.errorMessage = "Должно быть 8 символов";
        return;
      }

      this.isLoading = true;
      try {
        const token = localStorage.getItem('token');
        const response = await axios.post("api/qck/activate-link", this.note, {
          headers: {
            'Authorization': `Bearer ${token}` 
          }
        });

        if (response.data.success) {
          this.$emit("save", this.note);
          this.note = { qck_link: "", password: "" }; // Очистка полей
        } else {
          this.errorMessage = response.data.message || "Ошибка активации";
        }
      } catch (error) {
        if (error.response) {
          this.errorMessage = error.response.data.error || "Ошибка сервера";
        } else {
          this.errorMessage = "Сеть недоступна";
        }
      } finally {
        this.isLoading = false;
      }
    },
  },
};
</script>

<style lang="scss" scoped>
form {
  display: flex;
  flex-direction: column;
}

.error {
  color: red;
  font-size: 14px;
  margin-top: 10px;
}
</style>
