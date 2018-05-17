<template>
  <q-page class="flex flex-center">
    <div style="position:absolute;top:50px;line-height:18px;" v-if="!$store.state.account.CmdrName">
      <q-field dark helper="Enter your commander's name">
        <q-input
          type="text"
          color="yellow"
          v-model="input.cmdrName"
          @keyup="startGame"
        >
          <q-btn
            icon-right="send"
            style="vertical-align:bottom"
            size="sm"
            color="white"
            label="Start"
            @click="startGame"
            text-color="deep-orange"
            align="between"
            dense
          />
        </q-input>
      </q-field>
    </div>
    <div v-else>
      <p style="text-align:center;">{{ $store.state.socket.message.time | moment("DD.MM.YY HH:mm:ss") }}</p>
      <br><br>
      <textarea v-model="JSON.stringify(response)"/>
    </div>
  </q-page>
</template>

<script>
import axios from 'axios'

var start = function(parent, clicked) {
  if (clicked && (!parent.input.cmdrName || parent.input.cmdrName.length < 2)) {
    parent.$q.notify({
      color: 'negative',
      position: 'top',
      message: 'Name must be at least 2 characters long',
      icon: 'report_problem'
    })
    return
  }

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

var version = function(parent) {
  return axios({
    method: 'POST',
    'url': parent.$store.state.serverURL + '/version',
    //'data': parent.input,
    'headers': {
      'content-type': 'application/json'
      }
    }).then(result => {
      parent.$store.state.versionURL = result.data.version
      parent.$store.state.versionID = result.data.build
    }).catch(e => parent.$q.notify({
      color: 'negative',
      position: 'top',
      message: 'Loading failed: ' + e,
      icon: 'report_problem'
    }))
}

export default {
  name: 'PageIndex',
  data () {
    return {
      input: {
        cmdrName: '',
        gameID: ''
      },
      response: ''
    }
  },
  mounted () {
    if (!this.$auth.isAuthenticated() || this.$store.state.account.ID != "") {
      start(this, false)
    }
    version(this)
  },
  methods: {
    startGame (e) {
      if (e.keyCode == 13 || e.button == 0) {
        start(this, true)
      }
    }
  }
}
</script>
