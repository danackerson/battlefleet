import Vue from 'vue'

import Vuex from 'vuex'
Vue.use(Vuex)

import account from './module-account'
import VueNativeSock from 'vue-native-websocket'

import helpers from 'src/battlefleet/helpers'
const plugin = {
    install () {
        Vue.helpers = helpers
        Vue.prototype.$helpers = helpers
    }
}
Vue.use(plugin)

Vue.use(require('vue-moment')); // date formatting lib

const store = new Vuex.Store({
  modules: {
    account
  },
  state: {
    count: 0,
    versionID: '',
    versionURL: '',
    socket: {
      isConnected: false,
      message: 'Space',
      reconnectError: false,
    },
    serverURL: location.protocol + "//" + location.hostname + ':' + process.env.PORT
  },
  getters: {
    getLoggedIn: state => {
      return state.account.Auth0 ? 'Logout' : (state.account.ID ? 'Save' : 'Login')
    }
  },
  actions: {
    currentServerTime (context, content) {
      console.log('RCVD currentServerTime: ' + JSON.stringify(content))
      this.state.socket.message = content
      /* interesting code todo something on initial WS connect?
      this.$store.subscribe((mutation, state) => {
        if (mutation.type === 'SOCKET_ONOPEN' ) {
          // your code here
        }
      }*/
    }
  },
  mutations:{
    reinitState (state, event) {
      resetState()
    },
    SOCKET_ONOPEN (state, event)  {
      state.socket.isConnected = true
      console.log('WS Connected')
    },
    SOCKET_ONCLOSE (state, event)  {
      state.socket.isConnected = false
      console.log('WS Disconnected')
    },
    SOCKET_ONERROR (state, event)  {
      console.error('WS ERROR: ', state, event)
    },
    // default handler called for all methods
    SOCKET_ONMESSAGE (state, message)  {
      state.socket.message = message
      console.log('WS MSG: ', message)
    },
    // mutations for reconnect methods
    SOCKET_RECONNECT (state, count) {
      console.info(state, count)
    },
    SOCKET_RECONNECT_ERROR (state) {
      state.socket.reconnectError = true;
    }
  }
})

var websocketURL = "wss://"
if (location.protocol == "http:") {
  websocketURL = "ws://"
}
websocketURL += location.hostname + ':' + process.env.PORT
Vue.use(VueNativeSock, websocketURL + '/wsInit', {
  store: store,
  connectManually: false,
  format: 'json',
  reconnection: true,
  reconnectionAttempts: 5,
  reconnectionDelay: 3000,
})

export default store

const initialStateCopy = JSON.parse(JSON.stringify(store.state))

export function resetState() {
  store.replaceState(JSON.parse(JSON.stringify(initialStateCopy)))
}
