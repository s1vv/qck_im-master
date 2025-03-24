<template>
  <textarea
    :value="modelValue"
    @input="updateInput"
    class="textarea"
    placeholder="Пароль ссылки"
  ></textarea>
</template>

<script>
export default {
  name: 'note-input',
  props: {
    modelValue: String
  },
  methods: {
    updateInput(event) {
      this.$emit('update:modelValue', event.target.value)
    }
  }
}
</script>

<style scoped>
.textarea {
  width: 100%;
  border: 1px solid teal;
  padding: 10px 15px;
  margin-top: 10px;
  resize: vertical;
  border-radius: 5px;
}
</style>
