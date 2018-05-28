import axios from 'axios'

export default {
  version(parent) {
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
  },
  startGame(parent, clicked) {
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
  },
  logout(parent) {
    return axios({
      method: 'POST',
      'url': parent.$store.state.serverURL + '/logout',
      'data': { user: parent.$auth.user.sub },
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
        } else if (parent.response.Message !== undefined) {
          parent.$q.notify({
            color: 'positive',
            position: 'top',
            message: parent.response.Message,
            icon: 'done'
          })
        }
      }).catch(e => parent.$q.notify({
        color: 'negative',
        position: 'top',
        message: 'Loading failed: ' + e,
        icon: 'report_problem'
    }))
  },
  login(parent) {
    return axios({
      method: 'POST',
      'url': parent.$store.state.serverURL + '/login',
      'data': { "Auth0Token": parent.$auth.user.sub },
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
        } else if (result.data.ID !== undefined) {
           parent.$store.commit('account/setCurrentGameID', result.data.ID)
           parent.$store.commit('account/setCmdrName', result.data.Account.Commander)
           parent.$store.commit('account/setAccountID', result.data.Account.ID)
           parent.$store.commit('account/setAuth0Login', result.data.Account.Auth0)
           parent.$q.notify({
             color: 'positive',
             position: 'top',
             message: parent.response.Message + result.data.Account.Commander,
             icon: 'done'
           })
        }
      }).catch(e => parent.$q.notify({
        color: 'negative',
        position: 'top',
        message: 'Loading failed: ' + e,
        icon: 'report_problem'
    }))
  },
  updateAccount(parent) {
    return axios({
      method: 'POST',
      'url': parent.$store.state.serverURL + '/updateAccount',
      'data': parent.$store.state.account,
      'headers': {
        'content-type': 'application/json'
        }
      }).then(result => {
        parent.response = JSON.parse(JSON.stringify(result.data))
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
        } else if (parent.response.Message !== undefined) {
          parent.$q.notify({
            color: 'positive',
            position: 'top',
            message: parent.response.Message,
            icon: 'done'
          })
        }
      }).catch(e => parent.$q.notify({
        color: 'negative',
        position: 'top',
        message: 'Loading failed: ' + e,
        icon: 'report_problem'
    }))
  }

}
