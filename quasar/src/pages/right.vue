<template>
  <div class="sidePanel">
    <!-- http://quasar-framework.org/quasar-play/android/index.html#/showcase/grouping/collapsible
      Account (=> view, manage, incl delete)
      Games (=> view, manage, incl delete & highlight current)
      https://material.io/icons/
    -->
    <q-input
      value="getAccountCmdrName"
    />
    <q-list v-if="$store.state.account.ID">
      <q-collapsible label="Account" opened group="accountMgmt" icon="face">
        <q-field dark label="Commander">
          <q-input
            type="text"
            color="yellow"
            v-model="commanderName"
            :placeholder="getAccountCmdrName"
            :after="[
              {
                icon: 'done',
                content: !commanderName || commanderName.length >= 2 && commanderName != getAccountCmdrName,
                handler: updateName,
                size:'sm'
              }
            ]"
          />
        </q-field>
        <q-field dark label="Auth0 Name" v-if="getAccountAuth0">
          <q-input
            type="text"
            color="white"
            :value="this.$auth.user.name"
            readonly
          />
        </q-field>
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
import { mapGetters, mapMutations } from 'vuex'

export default {
  name: 'RightPanel',
  data () {
    return {
      commanderName: null
    }
  },
  methods: {
    updateName(e) {
      alert(JSON.stringify(e))
      if (e.keyCode == 13 || e.button == 0) {
        if (this.commanderName.length >= 2) {
          alert(this.commanderName + " : " + e.keyCode + " | " + e.button)
          this.$store.commit('account/setCmdrName', this.commanderName)
          this.$helpers.updateAccount(this)
        } else {
          this.$q.notify('Commander name must be at least 2 characters.')
        }
      } else {
        this.$q.notify('Dafuq?!')
      }
    }
  },
  computed: {
    ...mapGetters({
        getAccountCmdrName: 'account/getCmdrName',
        getAccountAuth0: 'account/getAuth0',
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
