<template>
  <div class="sidePanel">
    <!-- http://quasar-framework.org/quasar-play/android/index.html#/showcase/grouping/collapsible
      Account (=> view, manage, incl delete)
      Games (=> view, manage, incl delete & highlight current)
      https://material.io/icons/
    -->
    <q-list v-if="$store.state.account.ID">
      <q-collapsible label="Account" opened group="accountMgmt" icon="face">
        <q-field dark label="Commander">
          <q-input
            type="text"
            color="yellow"
            v-model="name"
            :placeholder="getAccountCmdrName"
            :after="[
              {
                icon: 'done',
                content: !name || name.length >= 2 && name != getAccountCmdrName,
                handler () {
                  updateName()
                },
                size:'sm'
              }
            ]"
          />
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
import { mapGetters, mapMutations } from 'vuex'

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
  data () {
    return {
      name: null
    }
  },
  methods: {
    updateName () {
      if (this.name.length >= 2) {
        this.$store.commit('account/setCmdrName', this.name)
        updateAccount(this)
      } else {
        this.$q.notify('Commander name must be at least 2 characters.')
      }
    }
  },
  computed: {
    ...mapGetters({
        getAccountCmdrName: 'account/getCmdrName',
    })
  },
  mounted() {
    /*this.$store.dispatch('account/getAccount').then(response => {
      alert(JSON.stringify(response, null, '\t'))
      this.name = this.getAccountCmdrName
    })*/
  }
}
</script>
