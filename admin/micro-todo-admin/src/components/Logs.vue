<template>
  <div id="logs">
    <h2 v-once>{{ name }}</h2>
    <p v-once>In Development</p>
  </div>
</template>

<script>

export default {
  name: 'logs',
  props: {
    registryUrl: {
      type: String,
      required: true
    }
  },

  data () {
    return {
      name: 'Logs',
      errors: [],
      services: []
    }
  },

  created () {

  }
}
</script>

<style scoped>
</style>
