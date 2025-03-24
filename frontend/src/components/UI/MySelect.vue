<template>
  <select :value="modelValue" @change="changeOption">
    <option disabled selected value="">Выберите из списка</option>
    <option 
      v-for="option in options" 
      :key="option.value" 
      :value="option.value"
    >
      {{ option.name }}
    </option>
  </select>
</template>

<script>
export default {
  name: 'MySelect',
  props: {
    modelValue: {
      type: String
    },
    options: {
      type: Array,
      default: () => []
    }
  },
  methods: {
    changeOption(event) {
      this.$emit('update:modelValue', event.target.value)
    }
  }
}
</script>

<style lang="scss" scoped>
select {
  width: 180px; /* Устанавливаем ширину */
  padding: 8px; /* Внутренние отступы */
  font-size: 14px; /* Размер шрифта */
  color: teal; /* Цвет текста */
  border: 2px solid teal; /* Цвет рамки */
  border-radius: 5px; /* Скругление углов */
  background-color: white; /* Фон */
  cursor: pointer; /* Курсор в виде указателя */
  outline: none; /* Убираем стандартное выделение */

  &:focus {
    border-color: darkcyan; /* Изменяем цвет при фокусе */
  }
}
</style>
