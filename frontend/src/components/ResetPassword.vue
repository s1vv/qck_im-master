<template>
  <div class="registration-container">
    <h3 class="header">Изменить пароль</h3>
    <form @submit.prevent="resetPassword" class="form-container">
      <my-input v-model="email" type="email" placeholder="Email" required v-focus class="input-field"/>
      <my-button type="submit" class="submit-button">Изменить пароль</my-button>

      <div class="agreement">
        <h4>Что будет дальше:</h4>
        <p>
          На указанную почту будет выслано письмо со ссылкой, перейдя по которой можно сменить пароль.
        </p>
        <h5>Если остались вопросы - пишите email: support@qck.im</h5>
      </div>
      
      <div v-if="showSuccessMessage" class="success-message">
        Заметка успешно сохранена
      </div>
      <div v-if="showErrorMessage" class="error-message">
        {{ errorMessage }}
      </div>
    </form>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'ResetPasswordForm',
  data() {
    return {
      email: '',
      showSuccessMessage: false,
      showErrorMessage: false,
      errorMessage: '',
    };
  },
  methods: {
    async resetPassword() {
      try {
        const response = await axios.post('/api/users/request-password-reset', {
          email: this.email
        });
        if (!response.data.success) {
          console.log(response.data.message);
          this.showSuccessMessage = true;
          setTimeout(() => {
            this.showSuccessMessage = false;
          }, 3000);
        } else {
          this.showErrorMessage = true;
          this.errorMessage = response.data.message;
        }
      } catch (error) {
        this.showErrorMessage = true;
        this.errorMessage = error.message;
      }
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
  height: 100vh;
  padding: 20px;
  background-color: #f9f9f9;
}

.header {
  color: teal;
  text-align: center;
  margin-bottom: 20px;
}

.form-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
  max-width: 400px;
  background: white;
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
}

.input-field {
  width: 100%;
  padding: 10px;
  font-size: 16px;
  border: 2px solid teal;
  border-radius: 5px;
  outline: none;
  transition: border-color 0.3s;

  &:focus {
    border-color: darkcyan;
  }
}

.submit-button {
  width: 100%;
  margin-top: 10px;
  padding: 10px;
  font-size: 16px;
  background-color: teal;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.3s;

  &:hover {
    background-color: darkcyan;
  }
}

.agreement {
  text-align: center;
  font-size: 14px;
  color: grey;
  margin-top: 15px;
}

.success-message, .error-message {
  margin-top: 10px;
  padding: 10px;
  border-radius: 5px;
  text-align: center;
  width: 100%;
}

.success-message {
  background-color: #d4edda;
  color: #155724;
}

.error-message {
  background-color: #f8d7da;
  color: #721c24;
}

/* Адаптация для мобильных устройств */
@media (max-width: 480px) {
  .registration-container {
    height: auto;
    padding: 10px;
  }

  .form-container {
    width: 100%;
    padding: 15px;
  }

  .input-field,
  .submit-button {
    font-size: 14px;
  }
}
</style>
