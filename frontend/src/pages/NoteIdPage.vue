<template>
  <div class="note-container">
    <textarea 
      v-model="editedName" 
      class="note-input name-input" 
      rows="1" 
      :maxlength="60"
      placeholder="Имя ссылки"
    ></textarea>
    <div class="divider"></div>
    <textarea 
      v-model="editedNote" 
      class="note-input" 
      rows="28" 
      placeholder="Описание ссылки"
      @input="checkTextLength"
    ></textarea>
    
    <!-- Показываем количество символов -->
    <div class="char-count">
      Количество символов: {{ editedNote.length }} / 1000
    </div>
    
    <div v-if="textTooLong" class="warning-message">
      Сохранить можно только 1000 символов. Сделайте текст короче.
    </div>
    
    <div>
      <my-button @click="updateNote" class="save-button">Сохранить изменения</my-button>
    </div>

    <!-- Сообщения об успехе или ошибке -->
    <div v-if="showSuccessMessage" class="success-message">
      Заметка успешно обновлена!
    </div>
    <div v-if="showErrorMessage" class="error-message">
      {{ errorMessage }}
    </div>
  </div>
</template>


<script>
import api from '@/services/api';

export default {
  data() {
    return {
      noteId: this.$route.params.id,
      note: {},
      editedName: '',
      editedNote: '',
      showSuccessMessage: false,
      showErrorMessage: false,
      errorMessage: '',
      textTooLong: false, // Флаг превышения длины текста
    };
  },
  created() {
    this.fetchNote();
  },
  methods: {
    async fetchNote() {
      try {
        const noteId = this.$route.params.id; 
        const response = await api.get(`/qck/qck-link/?link=${noteId}`);
        this.note = response.data;
        this.editedNote = this.note.description;
        this.editedName = this.note.name;
      } catch (error) {
        console.error('Ошибка получения заметки:', error);
        this.showErrorMessage = true;
        this.errorMessage = 'Не удалось загрузить заметку';
      }
    },
    async updateNote() {
      if (this.editedNote.length > 1000) {
        this.textTooLong = true;
        return;
      }

      try {
        const response = await api.post('/qck/update-data-link', {
          "qck_link": this.noteId, "name": this.editedName, "description": this.editedNote
        });

        // Если статус ответа успешный
        if (response.status === 200) {
          this.showSuccessMessage = true;
          setTimeout(() => {
            this.showSuccessMessage = false;
            this.$router.push('/notes'); // Переход после успешного сохранения
          }, 1000);
        } else {
          this.showErrorMessage = true;
          this.errorMessage = response.data.message || 'Неизвестная ошибка';
        }
      } catch (error) {
        console.error('Ошибка сохранения заметки:', error);
        this.showErrorMessage = true;
        this.errorMessage = 'Произошла ошибка при сохранении заметки';
      }
    },
    checkTextLength() {
      if (this.editedNote.length > 1000) {
        this.textTooLong = true;
      } else {
        this.textTooLong = false;
      }
    },
  },
};
</script>

<style>
.note-container {
  display: flex;
  flex-direction: column;
  align-items: center; /* Центрирование элементов */
  width: 100%;
  max-width: 600px; /* Максимальная ширина блока */
  margin: 0 auto; /* Центрируем сам контейнер */
}

.note-input {
  width: 100%; /* Занимает всю ширину контейнера */
  max-width: 600px; /* Ограничение максимальной ширины */
  min-width: 300px; /* Минимальная ширина для маленьких экранов */
  padding: 10px;
  font-size: 16px;
  border: 1px solid #ccc;
  border-radius: 5px;
  resize: none; /* Отключаем изменение размера вручную */
}

.name-input {
  height: 42px; /* Фиксированная высота для названия */
}

.divider {
  width: 100%;
  max-width: 600px;
  height: 2px; /* Толщина линии */
  background-color: #ccc; /* Цвет линии */
  margin: 10px 0; /* Отступы сверху и снизу */
}

.save-button {
  margin-top: 10px;
  padding: 10px 20px;
}

.success-message {
  color: green;
  font-size: 16px;
  margin-top: 10px;
}

.error-message {
  color: red;
  font-size: 16px;
  margin-top: 10px;
}

.warning-message {
  color: orange;
  font-size: 14px;
  margin-top: 5px;
}

/* Стиль для отображения количества символов */
.char-count {
  font-size: 14px;
  color: #666;
  margin-top: 5px;
}
</style>
