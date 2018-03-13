<template>
    <div id="serviceBlock" class="block" style="width: 250px">
      <span>Proto: {{ proto }}</span><br />
      <span>Host: {{ host }}</span><br />
      <span>Port: {{ port }}</span><br />
      <span>HTTP Port: {{ httpPort }}</span><br />
      <span>HTTPS Port: {{ httpsPort }}</span><br />
      <span>Routes: </span>
      <select>
        <option v-for="route in routes" :key="route">{{ route }}</option>
      </select>
      <br />
      <span>Health route: </span>
      <span v-if="isUp" class="badge up-badge">{{ health }}</span>
      <span v-if="!isUp" class="badge down-badge">{{ health }}</span>
      <br />
      <span>Weight: {{ weight }}</span><br />
      <span>Signature: {{ signature }}</span>
    </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'service-block',
  props: {
    service: {
      type: Object,
      required: true
    }
  },

  data () {
    return {
      proto: this.service.proto,
      type: this.service.type,
      host: this.service.host,
      port: this.service.port,
      httpPort: this.service.httpPort,
      httpsPort: this.service.httpsPort,
      health: this.service.health,
      weight: this.service.weight,
      signature: this.service.signature,
      routes: this.service.routes,
      isUp: false
    }
  },

  created () {
    let url = ''
    if (this.service.proto === 'HTTPS') {
      url = 'https://' + this.service.host + ':' + this.service.httpsPort
    } else {
      url = 'http://' + this.service.host + ':' + this.service.httpPort
    }

    axios.get(url + this.service.health).then(response => {
      if (response.body && response.body.up) {
        this.isUp = true
      } else {
        this.isUp = false
      }
    }).catch(e => console.log('Service is down: ' + e))
  }
}
</script>

<style scoped>

</style>
