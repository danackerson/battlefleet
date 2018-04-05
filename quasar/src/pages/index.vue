<template>
  <q-page class="flex flex-center">
    <p>Space</p>
  </q-page>
</template>

<style>
</style>

<script>
export default {
  name: 'PageIndex',
  methods: {
    clickButton: function(val) {
        this.$socket.sendObj({awesome: 'data'})
    }
  }
}
</script>
