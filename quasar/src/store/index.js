import Vue from 'vue'
import Vuex from 'vuex'

import example from './module-example'
import VueNativeSock from 'vue-native-websocket'

Vue.use(Vuex)
Vue.use(VueNativeSock, 'ws://localhost:8083/wsInit', {
  store: store,
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
      console.log(state.socket)
    },
    SOCKET_ONCLOSE (state, event)  {
      state.socket.isConnected = false
      console.log(state.socket)
    },
    SOCKET_ONERROR (state, event)  {
      console.error(state, event)
    },
    // default handler called for all methods
    SOCKET_ONMESSAGE (state, message)  {
      state.message = message
      console.log(state.message)
    },
    // mutations for reconnect methods
    SOCKET_RECONNECT(state, count) {
      console.info(state, count)
    },
    SOCKET_RECONNECT_ERROR(state) {
      state.socket.reconnectError = true;
    },
  }
})

export default store
