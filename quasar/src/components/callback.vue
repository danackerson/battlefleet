<template>
  <div class="callback">Callback</div>
</template>

<script>
export default {
  name: 'callback',
  mounted() {
    this.$auth.handleAuthentication().then((data) => {
      if (this.$auth.isAuthenticated()) {
        //var auth0User = JSON.parse(JSON.stringify(localStorage.getItem('user')))
        this.$store.state.account.Auth0 = this.$auth.user.sub
      } else {
        this.$store.state.loginState = 'Login'
      }
      this.$router.push({ name: 'home' })
    })
  }
}
</script>
