<template>
  <q-page class="flex flex-center">
    <p>{{ this.$store.state.socket.message.time | moment("DD.MM.YY HH:mm:ss") }}</p>
  </q-page>
</template>

<style>
</style>

<script>
export default {
  name: 'PageIndex'
}
</script>
