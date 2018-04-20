<template>
  <q-page class="flex flex-center">
    <p v-if="this.$store.state.socket.message != 'Space'">{{ this.$store.state.socket.message.time | moment("DD.MM.YY HH:mm:ss") }}</p>
  </q-page>
</template>

<style>
</style>

<script>
import axios from 'axios'

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
  mounted() {
    version(this)
  }
}
</script>
