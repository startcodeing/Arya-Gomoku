<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script setup lang="ts">
// App.vue now serves as the main router outlet
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  background: #f5f5f5;
}

#app {
  width: 100%;
  min-height: 100vh;
}

/* Global styles */
button {
  font-family: inherit;
}

input, textarea {
  font-family: inherit;
}

/* Utility classes */
.text-center {
  text-align: center;
}

.flex {
  display: flex;
}

.flex-col {
  flex-direction: column;
}

.items-center {
  align-items: center;
}

.justify-center {
  justify-content: center;
}

.gap-4 {
  gap: 1rem;
}

.p-4 {
  padding: 1rem;
}

.m-4 {
  margin: 1rem;
}

.rounded {
  border-radius: 0.5rem;
}

.shadow {
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}
</style>