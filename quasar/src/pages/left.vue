<template>
  <div>
    <input type="text" v-model="input.cmdrName" placeholder="Commander Name" />
    <button v-on:click="startGame()">Join the fleet!</button>
    <br />
    <br />
    <textarea v-model="JSON.stringify(response)"/>
    <button v-on:click="updateTime()">Update Time</button>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'LeftPanel',
  data () {
    return {
      input: {
        cmdrName: ''
      },
      response: ''
    }
  },
  methods: {
    startGame () {
      /*axios.get(
          'http://localhost:8083/post',
          data: {
            cmdrName: this.cmdrName
          }
        )
        .then(function (response) {*/
      axios({
        method: 'POST',
        'url': '/post',
        'data': this.input,
        'headers': {
          'content-type': 'application/json'
          }
        }).then(result => {
          this.response = result.data
        }).catch(e => this.$q.notify({
          color: 'negative',
          position: 'top',
          message: 'Loading failed: ' + e,
          icon: 'report_problem'
        }))
      console.log(this.response)
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
