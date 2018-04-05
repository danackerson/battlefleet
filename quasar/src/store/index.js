import Vue from 'vue'
import Vuex from 'vuex'

import example from './module-example'
import VueNativeSock from 'vue-native-websocket'

Vue.use(Vuex)
Vue.use(VueNativeSock, 'ws://localhost:8083/wsInit', {
  store: store,
  connectManually: true,
  format: 'json',
  reconnection: true,
  reconnectionAttempts: 5,
  reconnectionDelay: 3000,
})

const store = new Vuex.Store({
  modules: {
    example
  },
  state: {
    socket: {
      isConnected: false,
      message: '',
      reconnectError: false,
    }
  },
  mutations:{
    SOCKET_ONOPEN (state, event)  {
      state.socket.isConnected = true
      this.console.log(state.socket)
    },
    SOCKET_ONCLOSE (state, event)  {
      state.socket.isConnected = false
      this.console.log(state.socket)
    },
    SOCKET_ONERROR (state, event)  {
      this.console.error(state, event)
    },
    // default handler called for all methods
    SOCKET_ONMESSAGE (state, message)  {
      state.message = message
      this.console.log(state.message)
    },
    // mutations for reconnect methods
    SOCKET_RECONNECT(state, count) {
      this.console.info(state, count)
    },
    SOCKET_RECONNECT_ERROR(state) {
      state.socket.reconnectError = true;
    },
  }
})

export default store
