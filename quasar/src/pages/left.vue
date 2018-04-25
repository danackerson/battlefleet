<template>
  <div class="sidePanel">
    <p>
      <span v-if="this.$auth.isAuthenticated()">{{this.$auth.user.name}}'s IP is {{ ip }}</span>
      <span v-else>Your IP is {{ ip }}</span>
    </p>
    <br>
    <div v-show="$store.state.account.CurrentGameID">
    <p>Game ID: {{ $store.state.account.CurrentGameID }}</p>
    </div>
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
        this.$socket.sendObj({})
    }
  }
}
</script>
