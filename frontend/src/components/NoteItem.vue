<template>
  <div class="note">
    <div>
      <div><strong>Название: </strong>{{ note.name }}</div>
      <div><strong>Ссылка: </strong>{{ note.qck_link }}</div>
      <div class="note-body">{{ truncatedDescription }}</div>
    </div>
    <div class="note__btns">
      <my-button @click="$emit('remove', note)">
        Очистить
      </my-button>
      <my-button @click="openNote">
        Изменить
      </my-button>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    note: {
      type: Object,
      required: true,
    }
  },
  data() {
    return {
      charLimit: window.innerWidth > 1000 ? 200 : 60
    };
  },
  computed: {
    truncatedDescription() {
      return this.note.description.length > this.charLimit
        ? this.note.description.substring(0, this.charLimit) + '...'
        : this.note.description;
    }
  },
  methods: {
    updateCharLimit() {
      this.charLimit = window.innerWidth > 1000 ? 120 : 60;
    },
    openNote() {
      console.log(this.note.qck_link);
      this.$router.push({
        path: '/notes/' + this.note.qck_link
      });
    }
  },
  mounted() {
    window.addEventListener("resize", this.updateCharLimit);
  },
  beforeUnmount() {
    window.removeEventListener("resize", this.updateCharLimit);
  }
};
</script>

<style lang="scss" scoped>
.note {
  padding: 15px;
  border: 2px solid teal;
  margin-top: 15px;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  border-radius: 5px;
}

.note__btns {
  display: flex;
  justify-content: space-between;
  width: 100%;
}

.note__btns button:first-child {
  margin-right: auto;
  border-radius: 5px;
}

.note__btns button:last-child {
  margin-left: auto;
  border-radius: 5px;
}

.note-body {
  max-width: 100%;
  word-wrap: break-word;
  overflow-wrap: break-word;
  white-space: normal;
}
</style>
