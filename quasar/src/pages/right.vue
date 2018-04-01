<template>
  <div>
    <div>
      <h3 v-if="this.$auth.isAuthenticated()">{{$auth.user.name}}'s IP is {{ ip }}</h3>
      <h3 v-else>Your IP is {{ ip }}</h3>
      <input type="button" @click="toggleAuth" v-model="authState">
    </div>
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
      ip: '',
      authState: this.$auth.isAuthenticated ? 'Logout' : 'Login'
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
    toggleAuth() {
      if (this.$auth == 'undefined' || this.$auth.isAuthenticated()) {
        this.$auth.logout()
        this.authState = 'Login'
      } else {
        this.$auth.login()
        this.authState = 'Logout'
      }
    }
  }
}
</script>
