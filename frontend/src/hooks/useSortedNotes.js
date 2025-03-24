import {ref, computed} from 'vue'

export default function useSortedNotes(notes) {
    const selectedSort = ref('')
    const sortedNotes = computed(() => {
        return [...notes.value].sort((note1, note2) => note1[selectedSort.value]?.localeCompare(note2[selectedSort.value]))
    })

    return {
        selectedSort, sortedNotes
    }
}
