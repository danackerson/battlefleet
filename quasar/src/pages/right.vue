<template>
  <div class="sidePanel">
    <!-- http://quasar-framework.org/quasar-play/android/index.html#/showcase/grouping/collapsible
      Account (=> view, manage, incl delete)
      Games (=> view, manage, incl delete & highlight current)
      https://material.io/icons/
    -->
    <q-list>
      <q-collapsible label="Account" opened group="accountMgmt" icon="face">
        <q-field dark label="Commander">
          <q-input
            type="text"
            color="white"
            v-model="$store.state.account.CmdrName"
          >
            <q-btn
              icon-right="done"
              size="sm"
              @click="updateName"
              text-color="deep-orange"
              dense
            />
          </q-input>
        </q-field>
        <div>
          Auth0 Account?
        </div>
      </q-collapsible>
      <q-collapsible label="Games" group="accountMgmt" icon="games">
        <div>
          Game 1
        </div>
        <div>
          Game 2
        </div>
      </q-collapsible>
    </q-list>
  </div>
</template>

<script>
import axios from 'axios'

var updateAccount = function(parent) {
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

export default {
  name: 'RightPanel',
  methods: {
    updateName() {
      updateAccount(this)
    }
  }
}
</script>
