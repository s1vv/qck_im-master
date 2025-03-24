<template>
  <div v-if="notes.length > 0">
    <h3>Ваши ссылки</h3>
    <transition-group name="sort-list">
      <note-item
        v-for="(note, idx) in notes"
        :key="idx"
        :note="note"
        @remove="$emit('remove', note)"
      />
    </transition-group>
  </div>
  <h2 v-else style="color: red">
    Заметок нет
  </h2>
</template>

<script>
import NoteItem from "@/components/NoteItem";
  export default {
    components: {NoteItem},
    props: {
      notes: {
        type: Array,
        required: true
      }
    }
  }
</script>

<style lang="scss" scoped>
  .sort-list-item {
    display: inline-block;
    margin-right: 10px;
  }
  .sort-list-enter-active,
  .sort-list-leave-active {
    transition: all 0.4s ease;
  }
  .sort-list-enter-from,
  .sort-list-leave-to {
    opacity: 0;
    transform: translateX(-130px);
  }
  .sort-list-move {
    transition: transform 0.4s ease;
  }
</style>