<template>
  <div>
    <h4 v-if="this.$auth.isAuthenticated()">{{this.$auth.user.name}}'s IP is {{ ip }}</h4>
    <h4 v-else>Your IP is {{ ip }}</h4>
    <br>
    <span>App Inits: {{ this.$store.state.count }}</span>
    <button v-on:click="updateTime()">Update Time</button>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'LeftPanel',
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
  methods: {
    updateTime () {
        this.$socket.sendObj()
    }
  }
}
</script>
