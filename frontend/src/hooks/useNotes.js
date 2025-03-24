import { ref, onMounted } from 'vue';
import api from '@/services/api';

export function useNotes() {
    const notes = ref([]);
    const isNotesLoading = ref(true);

    const fetching = async () => {
        try {
            const response = await api.get('/qck/qck-links');
            if (Array.isArray(response.data.data)) {
                notes.value = response.data.data; 
            } else {
                console.error("Ошибка: Ожидался массив, но получен другой формат", response.data);
            }
        } catch (error) {
            console.error('Ошибка при получении данных:', error);
            // Обработка ошибок, например, показ сообщения пользователю
        } finally {
            isNotesLoading.value = false;
        }
    };

    onMounted(fetching);

    return {
        notes,
        isNotesLoading
    };
}
