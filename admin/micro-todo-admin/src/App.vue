<template>
  <div id="app">
    <h1>{{ name }}</h1>
    <app-menu />
    <transition name="fade" mode="out-in">
      <router-view/>
    </transition>
  </div>
</template>

<script>
import AppMenu from '@/components/Menu'

export default {
  name: 'app',
  components: {
    AppMenu
  },

  data () {
    return {
      name: 'Micro-TODO Admin'
    }
  }
}
</script>

<style>
body {
  background: #fdfdfd;
}

#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 30px;
}

.fade-enter-active, .fade-leave-active {
  transition: opacity .5s;
}

.fade-enter, .fade-leave-to /* .fade-leave-active до версии 2.1.8 */ {
  opacity: 0;
}

#errorsList {
  display: block;
  background: darkred;
  color: white;
  margin: auto;
  width: 400px;
  text-align: center;
  padding: 5px 10px;
}

.badge {
  display: inline-block;
  padding: 4px;
  margin: 1px;
  color: white;
  border-radius: 5px;
}

.up-badge.badge {
  background: green;
}

.down-badge.badge {
  background: darkred;
}

.block {
  display: block;
  text-align: center;
  background: #efefef;
  color: black;
  border-radius: 5px;
  border-color: #dfdfdf;
  width: 500px;
  padding: 5px;
  margin: 2px auto 5px;
}
</style>
