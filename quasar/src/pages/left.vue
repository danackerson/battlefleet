<template>
  <div v-if="accountCmdr()">
    <p >Welcome, {{ accountCmdr() }}!</p>
    <br>
    <textarea v-model="JSON.stringify(response)"/>
    <button v-on:click="updateTime()">Update Time</button>
    <br>
    <div v-show="gameInfo.ID != null && gameInfo.ID != ''">
    <p>Game ID: {{ gameInfo.ID }}</p>
    <p><a :href="`/account/${accountID()}`">Commander: {{ accountCmdr() }}</a></p>
    <p><a :href="versionURL()">{{ versionTag() }}</a></p>
    </div>
  </div>
  <div v-else>
    <p>Welcome, stranger!</p>
    <br>
    <input type="text" v-model="input.cmdrName" placeholder="Anonymous" required/>
    <button v-if="input.cmdrName != ''" v-on:click="startGame()">Join the fleet!</button>
    <label v-else>Enter Name</label>
    <br>
  </div>
</template>

<script>
import axios from 'axios'

var start = function(vueX) {
  return axios({
    method: 'POST',
    'url': vueX.$store.state.serverURL + '/post',
    'data': vueX.input,
    'headers': {
      'content-type': 'application/json'
      }
    }).then(result => {
      vueX.response = result.data
      if (result.data.ID != 'undefined') {
        vueX.input.gameID = result.data.ID
        vueX.gameInfo = JSON.parse(JSON.stringify(result.data))
      }
      vueX.input.cmdrName = vueX.gameInfo.Account.Commander
    }).catch(e => vueX.$q.notify({
      color: 'negative',
      position: 'top',
      message: 'Loading failed: ' + e,
      icon: 'report_problem'
    }))
}

export default {
  name: 'LeftPanel',
  data () {
    return {
      input: {
        cmdrName: '',
        gameID: ''
      },
      response: '',
      gameInfo: {
        ID: '',
        Version: {
          URL: '',
          Tag: ''
        },
        GridSize: 0,
        Account: {
          ID: '',
          Commander: ''
        }
      }
    }
  },
  mounted () {
    // only interesting if we already have a BattlefleetID cookie
    if (this.$cookies.get('battlefleetID') != null) {
      start(this)
    }
  },
  methods: {
    accountID () {
      var result = ''
      if (this.gameInfo && this.gameInfo.Account) {
        result = this.gameInfo.Account.ID
      }
      return result
    },
    accountCmdr () {
      var result = ''
      if (this.gameInfo && this.gameInfo.Account) {
        result = this.gameInfo.Account.Commander
      }
      return result
    },
    versionURL () {
      var result = ''
      if (this.gameInfo && this.gameInfo.Version) {
        result = this.gameInfo.Version.URL
      }
      return result
    },
    versionTag () {
      var result = ''
      if (this.gameInfo && this.gameInfo.Version) {
        result = this.gameInfo.Version.Tag
      }
      return result
    },
    startGame () {
      start(this)
    },
    updateTime () {
        this.$socket.sendObj(this.response)
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
