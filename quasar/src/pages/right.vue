<template>
  <div v-if="$store.state.account.CmdrName">
    <p >Welcome, {{ $store.state.account.CmdrName }} <img v-if="$auth.isAuthenticated()" height="16px" :src="$auth.user.picture">!</p>
    <br>
    <textarea v-model="JSON.stringify(response)"/>
    <br>
    <div v-show="$store.state.account.CurrentGameID">
    <p>Game ID: {{ $store.state.account.CurrentGameID }}</p>
    <p><a :href="`/account/${$store.state.account.ID}`">Commander: {{ $store.state.account.CmdrName }}</a></p>
    <!-- <p><a :href="versionURL()">{{ versionTag() }}</a></p> -->
    </div>
    Auth0: <button type="text" @click="toggleAuth" v-model="authState">{{ authState }}</button>
    <br><br>
  </div>
  <div v-else>
    <p>Welcome, stranger!</p>
    <input type="text" v-model="$store.state.account.CmdrName" placeholder="Anonymous" required/>
    <button v-if="$store.state.account.CmdrName.length > 2" v-on:click="startGame()">Join the fleet!</button>
    <label v-else>Enter Name</label>
    <br><br>
    Auth0: <button type="text" @click="toggleAuth" v-model="authState">{{ authState }}</button>
    <br><br>
    <textarea v-model="JSON.stringify(response)"/>
  </div>
</template>

<style>
h1, h2 {
  font-weight: normal;
}
div.rightPanel {
  text-align: center;
}
</style>

<script>
import axios from 'axios'

var login = function(vueX) {
  return axios({
    method: 'POST',
    'url': vueX.$store.state.serverURL + '/login',
    'data': { input: vueX.input, user: vueX.$auth.user.sub },
    'headers': {
      'content-type': 'application/json'
      }
    }).then(result => {
      vueX.response = JSON.parse(JSON.stringify(result.data))
      vueX.$store.state.count++
      if (vueX.response.Error !== undefined) {
        if (vueX.response.HTTPCode == '412') {
          vueX.response.Error = "Your session is no longer on the server. Please login or create a new game."
        }
        vueX.$q.notify({
          color: 'warning-l',
          position: 'top',
          message: vueX.response.Error + ' (' + vueX.response.HTTPCode + ')',
          icon: 'report_problem'
        })
      } else if (result.data.ID !== undefined) {
        vueX.gameInfo = vueX.response
      }
    }).catch(e => vueX.$q.notify({
      color: 'negative',
      position: 'top',
      message: 'Loading failed: ' + e,
      icon: 'report_problem'
    }))
}

var start = function(vueX) {
  return axios({
    method: 'POST',
    'url': vueX.$store.state.serverURL + '/newGame',
    'data': vueX.input,
    'headers': {
      'content-type': 'application/json'
      }
    }).then(result => {
      vueX.response = JSON.parse(JSON.stringify(result.data))
      vueX.$store.state.count++
      if (vueX.response.Error !== undefined) {
        if (vueX.response.HTTPCode == '412') {
          vueX.response.Error = "Your session is no longer on the server. Please login or create a new game."
        }
        vueX.$q.notify({
          color: 'negative',
          position: 'top',
          message: vueX.response.Error + ' (' + vueX.response.HTTPCode + ')',
          icon: 'report_problem'
        })
      } else if (result.data.ID !== undefined) {
        //vueX.gameInfo = vueX.response
        vueX.$store.commit('account/setCurrentGameID', result.data.ID)
        vueX.$store.commit('account/setCmdrName', result.data.Account.Commander)
        vueX.$store.commit('account/setAccountID', result.data.Account.ID)
      }
    }).catch(e => vueX.$q.notify({
      color: 'negative',
      position: 'top',
      message: 'Loading failed: ' + e,
      icon: 'report_problem'
    }))
}

export default {
  name: 'RightPanel',
  data () {
    return {
      authState: this.$auth.isAuthenticated() ? 'Logout' : (this.$store.getters.getAccountID != '' ? 'Save' : 'Login'),
      input: {
        cmdrName: '',
        gameID: ''
      },
      response: ''
    }
  },
  mounted () {
    // try and auto login/refresh with current game
    start(this)
  },
  methods: {
    toggleAuth() {
      if (this.$auth.isAuthenticated()) {
        this.$auth.logout()
        this.authState = 'Login'
        this.$store.state.count++
      } else {
        this.$auth.login()
        start(this)
      }
    },
    startGame () {
      start(this)
    }
  }
}
</script>


<style>
textarea {
    width: 300px;
    height: 200px;
}
</style>
