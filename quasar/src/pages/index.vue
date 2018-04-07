<template>
  <q-page class="flex flex-center">
    <p>{{ this.$store.state.socket.message.time }}</p>
  </q-page>
</template>

<style>
</style>

<script>
export default {
  name: 'PageIndex'
}
</script>
