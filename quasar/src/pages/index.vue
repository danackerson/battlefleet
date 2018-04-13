<template>
  <q-page class="flex flex-center">
    <p v-if="this.$store.state.socket.message != 'Space'">{{ this.$store.state.socket.message.time | moment("DD.MM.YY HH:mm:ss") }}</p>
  </q-page>
</template>

<style>
</style>

<script>
export default {
  name: 'PageIndex'
}
</script>
