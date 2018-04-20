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
          Battlefleet!
          <div slot="subtitle">Running on Quasar v{{ $q.version }}</div>
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
export default {
  name: 'LayoutDefault',
  data () {
    return {
      leftDrawerOpen: this.$q.platform.is.desktop,
      rightDrawerOpen: this.$q.platform.is.desktop
    }
  }
}
</script>

<style>
</style>
