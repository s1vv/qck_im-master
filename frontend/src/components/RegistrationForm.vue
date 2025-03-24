<template>
  <div class="registration-container">
    <h3 class="header">Регистрация</h3>
    <form @submit.prevent="register" class="registration-form">
      <my-input v-model="login" type="login" placeholder="Login" required class="input-field"/>
      <my-input v-model="email" type="email" placeholder="Email" required class="input-field"/>
      <my-input v-model="password" type="password" placeholder="Password" required class="input-field"/>
      <my-input v-model="qckLink" type="text" placeholder="Qck link" required class="input-field"/>
      <my-input v-model="qckPassword" type="password" placeholder="Qck password" required class="input-field"/>
      <my-button type="submit" class="submit-button">Регистрация</my-button>
      <p v-if="error" class="error-message">{{ error }}</p>
    </form>

    <div class="reset-password-container">
      <my-button @click="redirectToResetPassword" class="reset-button">Изменить пароль</my-button>
    </div>

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
import axios from 'axios';

export default {
  name: 'RegistrationForm',
  data() {
    return {
      login: '',
      email: '',
      password: '',
      qckLink: '',
      qckPassword:'',
      error: null
    };
  },
  methods: {
    async register() {
      if (!this.login.trim() || this.login.length<3){
        this.error = "Login не менее 3 символов";
        return;
      }

      if (!this.password.trim() || this.password.length<8){
        this.error = "Пароль не менее 8 символов";
        return;
      }

      if (!this.qckLink.trim() || this.qckLink.length!=8){
        this.error = "Ссылка должна быть из 8 символов";
        return;
      }

      if (!this.qckPassword.trim() || this.qckPassword.length!=8){
        this.error = "Пароль qck link должен быть из 8 символов";
        return;
      }

      try {
        const response = await axios.post('/api/users/register', {
          login: this.login,
          email: this.email,
          password: this.password,
          qck_link: this.qckLink,
          qck_link_password: this.qckPassword
        });
        console.log(response);
        this.$router.push('/login');
      } catch (error) {
        this.error = error.response.data.message;
      }
    },
    redirectToResetPassword() {
      this.$router.push('/reset-password');
    }
  }
}
</script>

<style lang="scss" scoped>
.registration-container {
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

.registration-form {
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
}

.reset-password-container {
  margin-top: 20px;
}

.reset-button {
  width: 100%;
  max-width: 400px;
  padding: 10px;
  font-size: 16px;
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

.error-message {
  color: red;
  font-size: 14px;
  margin-top: 20px;
  text-align: center;
}
</style>
