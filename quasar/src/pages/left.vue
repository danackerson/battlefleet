<template>
  <div>
    <input type="text" v-model="input.firstname" placeholder="First Name" />
    <input type="text" v-model="input.lastname" placeholder="Last Name" />
    <button v-on:click="sendData()">Send</button>
    <br />
    <br />
    <textarea v-model="JSON.stringify(response)"/>
    <button v-on:click="clickButton()">Submit</button>
    {{ this.$store.state.count }}
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'LeftPanel',
  data () {
    return {
      input: {
        firstname: '',
        lastname: ''
      },
      response: ''
    }
  },
  methods: {
    sendData () {
      axios({ method: 'POST', 'url': 'https://httpbin.org/post', 'data': this.input, 'headers': { 'content-type': 'application/json' } }).then(result => {
        this.response = result.data
      }).catch(e => this.$q.notify({
        color: 'negative',
        position: 'top',
        message: 'Loading failed: ' + e,
        icon: 'report_problem'
      }))
      console.log(this.response)
    },
    clickButton () {
        this.$socket.sendObj(this.response)
    }
  }
}
</script>

<style>
textarea {
    width: 300px;
    height: 500px;
}
</style>
