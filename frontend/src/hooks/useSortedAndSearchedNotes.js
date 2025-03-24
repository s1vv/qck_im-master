import { ref, computed } from 'vue'

export default function useSortedAndSearchedNotes(sortedNotes) {
    const searchQuery = ref('')

    const sortedAndSearchedNotes = computed(() => {
        return sortedNotes.value.filter(note => {
            if (!note) return false; // Добавляем проверку на undefined
            return note.name.toLowerCase().includes(searchQuery.value.toLowerCase())
        });
    })

    return {
        searchQuery,
        sortedAndSearchedNotes
    }
}
