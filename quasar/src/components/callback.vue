<template>
  <div class="callback">Callback</div>
</template>

<script>
import axios from 'axios'

var loginAccount = function(parent) {
  return axios({
    method: 'POST',
    'url': parent.$store.state.serverURL + '/login',
    'data': { "Auth0Token": parent.$auth.user.sub},
    'headers': {
      'content-type': 'application/json'
      }
    }).then(result => {
      parent.response = JSON.parse(JSON.stringify(result.data))
      alert(parent.response)
      if (parent.response.Error !== undefined) {
        parent.$q.notify({
          color: 'warning',
          position: 'top',
          message: parent.response.Error + ' (' + parent.response.HTTPCode + ')',
          icon: 'report_problem'
        })
      } else if (parent.response.Message !== undefined) {
        parent.$q.notify({
          color: 'positive',
          position: 'top',
          message: parent.response.Message,
          icon: 'done'
        })
      }
  })
}

export default {
  name: 'callback',
  mounted() {
    this.$auth.handleAuthentication().then((data) => {
      if (this.$auth.isAuthenticated()) {
        //var auth0User = JSON.parse(JSON.stringify(localStorage.getItem('user')))
        this.$store.state.account.Auth0 = this.$auth.user.sub
        // TODO load Account and currentGameID if exists
        loginAccount(this)
      } else {
        this.$store.state.loginState = 'Login'
      }
      this.$router.push({ name: 'home' })
    })
  }
}
</script>
