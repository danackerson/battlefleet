<template>
  <q-layout view="hHh LpR lFr"> <!-- this controls behavior of panels, shadows: see http://quasar-framework.org/components/layout-drawer.html -->
    <q-layout-header>
      <q-toolbar
        color="primary"
        :glossy="$q.theme === 'mat'"
        :inverted="$q.theme === 'ios'"
      >
        <q-btn
          flat round dense
          @click="leftDrawerOpen = !leftDrawerOpen"
          aria-label="Menu"
        >
          <q-icon name="menu" />
        </q-btn>

        <q-toolbar-title>
          <div id="cmdrName" >
            <span v-if="$store.state.account.CmdrName">Welcome, <a :href="`/account/${$store.state.account.ID}`">{{ $store.getters['account/getCmdrName'] }}</a>!</span>
            <span v-else>Welcome, stranger!</span>
            <q-btn style="vertical-align:middle" size="sm" color="white" :label="loggedIn" @click="toggleAuth" height="10px" text-color="deep-orange" align="between" dense>
              &nbsp;<img align="center" height="16px" src="//cdn2.auth0.com/styleguide/latest/lib/logos/img/badge.png">
            </q-btn>
          </div>
        </q-toolbar-title>

        <q-btn
          flat round dense
          @click="rightDrawerOpen = !rightDrawerOpen"
          icon="menu"
        />

      </q-toolbar>
    </q-layout-header>

    <q-layout-drawer side="left" v-model="leftDrawerOpen">
      <router-view name="left"/> <!-- checkout src/router/routes.js -->
    </q-layout-drawer>

    <q-layout-drawer side="right" v-model="rightDrawerOpen">
      <router-view name="right"/> <!-- checkout src/router/routes.js -->
    </q-layout-drawer>

    <q-page-container>
      <router-view name="main"/> <!-- checkout src/router/routes.js -->
    </q-page-container>

    <q-layout-footer>
      <q-toolbar
        color="primary"
        :glossy="$q.theme === 'mat'"
        :inverted="$q.theme === 'ios'"
      >
        <q-btn
          flat
          dense
          round
          @click="leftDrawerOpen = !leftDrawerOpen"
          aria-label="Menu"
        >
          <q-icon name="menu" />
        </q-btn>

        <q-toolbar-title>
          <div id="prodName">Battlefleet</div>
          <!--<div style="float:left;">Quasar v{{ $q.version }}</div>-->
          <div id="version" slot="subtitle"><a :href="$store.state.versionURL">v{{ $store.state.versionID }}</a></div>
        </q-toolbar-title>

        <q-btn
          flat round dense
          @click="rightDrawerOpen = !rightDrawerOpen"
          icon="menu"
        />
      </q-toolbar>
    </q-layout-footer>
  </q-layout>
</template>

<script>
import axios from 'axios'

var logout = function(parent) {
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
}

var login = function(parent) {
  return axios({
    method: 'POST',
    'url': parent.$store.state.serverURL + '/login',
    'data': { input: parent.input, user: parent.$auth.user.sub },
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
      }
    }).catch(e => parent.$q.notify({
      color: 'negative',
      position: 'top',
      message: 'Loading failed: ' + e,
      icon: 'report_problem'
    }))
}

export default {
  name: 'LayoutDefault',
  data () {
    return {
      leftDrawerOpen: this.$q.platform.is.desktop,
      rightDrawerOpen: this.$q.platform.is.desktop
    }
  },
  computed: {
    loggedIn() {
      return this.$store.getters.getLoggedIn
    }
  },
  methods: {
    toggleAuth() {
      if (this.$auth.isAuthenticated()) {
        this.$auth.logout()
        logout(this)
        this.$store.commit('reinitState')
        this.$store.state.count++
      } else {
        this.$auth.login()
      }
    }
  }
}
</script>
