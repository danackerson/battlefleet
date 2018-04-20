<template>
  <div class="sidePanel">
    <div v-if="$store.state.account.CmdrName">
      <p>
        <span>Welcome, <a :href="`/account/${$store.state.account.ID}`">{{ $store.state.account.CmdrName }}</a>!</span>
        <q-btn style="vertical-align:bottom" size="xs" color="white" :label="authState" @click="toggleAuth" height="10px" text-color="deep-orange" align="between" dense>
          <img align="center" height="16px" src="//cdn2.auth0.com/styleguide/latest/lib/logos/img/badge.png">
        </q-btn>
      </p>
      <br><br>
      <textarea v-model="JSON.stringify(response)"/>
    </div>
    <div v-else>
      <p>
        <span>Welcome, stranger!</span>
        <q-btn style="vertical-align:bottom" size="xs" color="white" :label="authState" @click="toggleAuth" height="10px" text-color="deep-orange" align="between" dense>
          <img align="center" height="16px" src="//cdn2.auth0.com/styleguide/latest/lib/logos/img/badge.png">
        </q-btn>

      <q-field
        dark
        icon="airplane"
        label="Name"
        label-color="deep-orange"
        color="white"
        helper="Enter your commander's name"
       >
         <q-input
           float-label="Anonymous"
           type="text"
           v-model="input.cmdrName"
         />
       </q-field>
     </p>
      <!-- <input type="text" v-model="input.cmdrName" placeholder="Anonymous" required/>
      <button v-if="input.cmdrName.length > 2" v-on:click="startGame()">Join the fleet!</button>
      <label v-else>Enter Name</label> -->
      <br><br>
      <textarea v-model="JSON.stringify(response)"/>
    </div>
  </div>
</template>

<style>
h1, h2 {
  font-weight: normal;
}
p, span {
  padding: 10px;
}
textarea {
    width: 300px;
    height: 200px;
}
</style>

<script>
import axios from 'axios'

var login = function(parent) {
  return axios({
    method: 'POST',
    'url': parent.$store.state.serverURL + '/login',
    'data': { input: parent.input, user: parent.$auth.user.sub },
    'headers': {
      'content-type': 'application/json'
      }
    }).then(result => {
      parent.response = JSON.parse(JSON.stringify(result.data))
      parent.$store.state.count++
      if (parent.response.Error !== undefined) {
        if (parent.response.HTTPCode == '412') {
          parent.response.Error = "Your session is no longer on the server. Please login or create a new game."
        }
        parent.$q.notify({
          color: 'warning-l',
          position: 'top',
          message: parent.response.Error + ' (' + parent.response.HTTPCode + ')',
          icon: 'report_problem'
        })
      }
    }).catch(e => parent.$q.notify({
      color: 'negative',
      position: 'top',
      message: 'Loading failed: ' + e,
      icon: 'report_problem'
    }))
}

var start = function(parent) {
  return axios({
    method: 'POST',
    'url': parent.$store.state.serverURL + '/newGame',
    'data': parent.input,
    'headers': {
      'content-type': 'application/json'
      }
    }).then(result => {
      parent.response = JSON.parse(JSON.stringify(result.data))
      parent.$store.state.count++
      if (parent.response.Error !== undefined) {
        if (parent.response.HTTPCode == '412') {
          parent.response.Error = "Your session is no longer on the server. Please login or create a new game."
        }
        parent.$q.notify({
          color: 'warning',
          position: 'top',
          message: parent.response.Error + ' (' + parent.response.HTTPCode + ')',
          icon: 'report_problem'
        })
      } else if (result.data.ID !== undefined) {
        parent.$store.commit('account/setCurrentGameID', result.data.ID)
        parent.$store.commit('account/setCmdrName', result.data.Account.Commander)
        parent.$store.commit('account/setAccountID', result.data.Account.ID)
      }
    }).catch(e => parent.$q.notify({
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
      authState: this.$auth.isAuthenticated() ? 'Logout' : (this.$store.getters.getAccountID ? 'Save' : 'Login'),
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
