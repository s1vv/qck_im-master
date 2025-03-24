<template>
  <div class="login-container">
    <form @submit.prevent="login" class="login-form">
      <my-input v-model="name" type="name" placeholder="Login" required v-focus class="input-field"/>
      <my-input v-model="password" type="password" placeholder="Password" required class="input-field"/>
      <my-button type="submit" class="submit-button">Войти</my-button>
    </form>
    <my-button @click="redirectToRegister" class="register-button">Зарегистрироваться или изменить пароль</my-button>
    <p v-if="error" class="error-message">{{ error }}</p>
    <div class="agreement">
      <h4 class="agreement-header">Пользовательское соглашение</h4>
      <p class="agreement-text">
        Сервис не предназначен для хранения персональных данных. Пользователь самостоятельно определяет содержание вводимых данных. Администрация сервиса не обрабатывает, не собирает и не анализирует введенную информацию с целью идентификации личности.
        Пользователь самостоятельно несет ответственность за вводимые в сервис данные. Администрация сервиса не контролирует и не несет ответственности за содержание информации, вводимой пользователем.
      </p>
      <h5 class="contact-info">Остались вопросы - пишите email: support@qck.im</h5>
    </div>
  </div>
</template>

<script>
import api from '@/services/api';

export default {
  name: 'LoginForm',
  data() {
    return {
      name: '',
      password: '',
      error: null
    };
  },
  
  methods: {
    async login() {
      try {
        const response = await api.post('/users/login', {
          login: this.name,
          password: this.password
        });

        localStorage.setItem('token', response.data.access_token);
        localStorage.setItem('refresh_token', response.data.refresh_token);
        this.$router.push('/');

      } catch (error) {
        this.error = error.response?.data?.message || 'Ошибка входа';
      }
    },
    redirectToRegister() {
      this.$router.push('/register');
    },
    
  }
};
</script>


<style lang="scss" scoped>
.login-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh; /* Заполняет всю высоту экрана */
  padding: 20px;
}

.header {
  color: teal;
  text-align: center;
  margin-bottom: 20px;
}

.login-form {
  display: flex;
  flex-direction: column;
  align-items: center;
  max-width: 400px;
  width: 100%;
  margin-bottom: 20px;
}

.input-field {
  width: 100%;
  max-width: 400px;
  padding: 10px;
  font-size: 16px;
  border: 1px solid #ccc;
  border-radius: 5px;
  margin-bottom: 10px;
}

.submit-button {
  width: 100%;
  max-width: 400px;
  padding: 10px;
  font-size: 16px;
  margin-top: 10px;
  border-radius: 10px;  /* Скругляем углы */
}

.error-message {
  color: red;
  font-size: 14px;
  margin-top: 10px;
  text-align: center;
}

.register-button {
  width: 100%;
  max-width: 400px;
  padding: 10px;
  font-size: 16px;
  margin-top: 20px;
  border-radius: 10px;  /* Скругляем углы */
}

.agreement {
  margin-top: 30px;
  text-align: center;
}

.agreement-header {
  color: grey;
  font-size: 18px;
}

.agreement-text {
  font-size: 14px;
  color: #333;
  margin-top: 10px;
  line-height: 1.5;
}

.contact-info {
  font-size: 14px;
  margin-top: 10px;
}
</style>
