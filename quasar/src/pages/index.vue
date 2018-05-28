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
      <q-input
        :value="getAccountCmdrName"
      />
      <p v-if="$store.state.socket.message.time" style="text-align:center;">{{ $store.state.socket.message.time | moment("DD.MM.YY HH:mm:ss") }}</p>
      <br><br>
      <textarea v-model="JSON.stringify(response)"/>
    </div>
  </q-page>
</template>

<script>
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
      this.$helpers.startGame(this, false)
    }
    this.$helpers.version(this)
  },
  methods: {
    startGame (e) {
      if (e.keyCode == 13 || e.button == 0) {
        this.$helpers.startGame(this, true)
      }
    }
  }
}
</script>
