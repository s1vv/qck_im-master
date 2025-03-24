<template>
  <nav-bar></nav-bar> 
  <div class="app">
    <router-view></router-view>
  </div>
</template>

<script>
import NavBar from "@/components/NavBar";

export default {
  components: {
    NavBar,
  }
}
</script>

<style lang="scss">
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
  font-family: "Inter", sans-serif; /* Укажите нужный шрифт */
  // color: darkslateblue;
}
.app {
  padding: 20px;
}
</style>
