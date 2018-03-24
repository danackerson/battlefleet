<template>
  <div>
    <h2>Your IP is {{ ip }}</h2>
  </div>
</template>

<style>
h1, h2 {
  font-weight: normal;
}
</style>

<script>
import axios from 'axios'

export default {
  name: 'RightPanel',
  data () {
    return {
      ip: ''
    }
  },
  mounted () {
    axios({ method: 'GET', 'url': 'https://httpbin.org/ip' }).then(result => {
      this.ip = result.data.origin
    }, error => {
      console.error(error)
    })
  },
  methods: { }
}
</script>
