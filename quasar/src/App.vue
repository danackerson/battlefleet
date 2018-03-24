<template>
  <div id="q-app">  <!-- matches q-app in index.template.html file //-->

    <!-- following transition enables panels to auto slide in/out on resize //-->
    <transition
      enter-active-class="animated fadeIn"
      leave-active-class="animated fadeOut"
      mode="out-in"
      :duration="300"
      @leave="resetScroll"
    >
      <router-view />
    </transition>
  </div>
</template>

<script>
export default {
  name: 'App',
  methods: {
    resetScroll (el, done) {
      document.documentElement.scrollTop = 0
      document.body.scrollTop = 0
      done()
    }
  }
}
</script>

<style>
</style>
