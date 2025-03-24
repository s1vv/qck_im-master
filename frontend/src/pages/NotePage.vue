<template>
  <div>
    <h3 style="color: teal;">Поиск по содержимому</h3>
    <my-input
      v-model="searchQuery"
      placeholder="Поиск..."
      v-focus
      :maxlength="20"
    />
    <div class="app__btns">
      <my-button 
        @click="showDialog"
      >
        Добавить ссылку
      </my-button>
      <my-select
        v-model="selectedSort"
        :options="sortOptions"
      />
    </div>
    <my-dialog v-model:show="dialogVisible">
      <note-form
      @save="saveLink"
      />
    </my-dialog>
    <note-list 
      :notes="sortedAndSearchedNotes"
      @remove="removeNote"
      v-if="!isNotesLoading"
    />
    <div v-else>Идет загрузка...</div>
    <div v-intersection="loadMoreNotes" class="observer"></div>
    

  </div>
</template>

<script>
import api from '@/services/api'; 

import NoteForm from "@/components/NoteForm";
import NoteList from "@/components/NoteList";
import { useNotes } from "@/hooks/useNotes";
import useSortedNotes from "@/hooks/useSortedNotes"
import useSortedAndSearchedNotes from "@/hooks/useSortedAndSearchedNotes"

export default {
  components: {
    NoteForm, 
    NoteList
  },
  data() {
    return {
      dialogVisible: false,
      sortOptions: [
        { value: 'name', name: 'По названию' },
        { value: 'description', name: 'По содержимому' },
      ],
      isLoadingMore: false // Инициализируем переменную isLoadingMore
    }
  },
  methods: {
    async saveLink(note) {
      try {
        const response = await api.post('qck/activate-link', note); 
        this.notes.unshift(response.data.note); // Добавляем полный объект заметки в начало списка
        this.dialogVisible = false;
      } catch (error) {
        console.error('Ошибка активации, проверьте правильность данных:', error);
      }
    },

    async removeNote(note) {
      try {
        if (note.description == ''){
          return
        }
        await api.post(`/qck/remove-link-description/`, note); 
        const updatedNotes = this.notes.map(n => 
      n.id === note.id ? { ...n, description: '' } : n
    );

    this.notes = updatedNotes;
      } catch (error) {
        console.error('Ошибка удаления описания:', error);
      }
    },
    showDialog() {
      this.dialogVisible = true;
    },
    async loadMoreNotes() {
      try {
        if (!this.isLoadingMore && this.page < this.totalPages) {
          this.isLoadingMore = true;
          this.page += 1;
          const response = await api.get('/qck/qck-links', {
            params: {
              _page: this.page,
              _limit: this.limit
            }
          });
          const responseData = response.data.data;
          this.notes = [...this.notes, ...responseData];
        }
      } catch (error) {
        console.error('Ошибка загрузки данных:', error);
      } finally {
        this.isLoadingMore = false;
      }
    },

  },
  setup() {
    const { notes, totalPages, isNotesLoading } = useNotes(10);
    const { sortedNotes, selectedSort } = useSortedNotes(notes);
    const { searchQuery, sortedAndSearchedNotes } = useSortedAndSearchedNotes(sortedNotes);

    return {
      notes,
      totalPages,
      isNotesLoading,
      sortedNotes,
      selectedSort,
      searchQuery,
      sortedAndSearchedNotes,
    }
  }
}
</script>

<style>
.app__btns {
  margin: 15px 0;
  display: flex;
  justify-content: space-between;
}
.page__wrapper {
  display: flex;
  margin-top: 15px;
}
.page {
  border: 1px solid grey;
  padding: 10px;
}
.current_page {
  border: 2px solid teal;
}
.observe {
  height: 30px;
  background: gray;
}
</style>
