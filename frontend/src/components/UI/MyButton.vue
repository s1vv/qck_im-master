<template>
    <button 
      class="btn" 
      :class="{ 'btn-disabled': disabled }" 
      :disabled="disabled"
    >
      <slot></slot>
    </button>
  </template>
  
  <script>
  export default {
    name: "my-button",
    props: {
      disabled: {
        type: Boolean,
        default: false
      }
    }
  };
  </script>
  
  <style lang="scss" scoped>
  .btn {
    padding: 10px 15px;
    background: teal;
    color: white;
    border: none;
    cursor: pointer;
    font-size: 16px;
    border-radius: 5px;
    transition: background 0.3s ease-in-out;
  
    &:hover {
      background: darkcyan;
    }
  
  }
  
  .btn-disabled {
    background: gray;
    cursor: not-allowed;
  }
  </style>
  