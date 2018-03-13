<template>
  <div id="health">
    <h2 v-once>{{ name }}</h2>

    <transition name="fade" mode="out-in">
      <div id="registryBlock" class="block">
        <p><b>Registry</b></p>
        <p>URL: {{ registry.url }}</p>
        <transition name="fade" mode="out-in">
          <span v-if="registry.up" class="up-badge badge">UP</span>
          <span v-else class="down-badge badge">DOWN</span>
        </transition>
      </div>
    </transition>

    <hr :width="this.getNotEmptyTypesNumber() * 300 + 'px'" color="#f0f0f0" />

    <div style="margin: auto">
      <div id="serviceList" v-for="type in types" :key="type" v-if="services[type]" class="serviceTypeBlock">
        <span style="font-weight: bold">{{ type }}</span>
        <transition-group name="fade" mode="out-in">
          <div v-for="item in services[type]" v-bind:key="item.signature">
            <service-block v-bind:service="item"></service-block>
          </div>
        </transition-group>
      </div>
    </div>
  </div>
</template>

<script>
import request from '../helpers/registry_requests'
import ServiceBlock from '@/components/ServiceBlock'
import axios from 'axios'

export default {
  name: 'health',
  props: {
    registryUrl: {
      type: String,
      required: true
    }
  },
  components: {
    'ServiceBlock': ServiceBlock
  },

  data () {
    return {
      name: 'Health',
      registry: {},
      types: [],
      services: {}
    }
  },

  methods: {
    pingAndRemoveDead: function () {
      if (!this.types || !this.services) {
        return
      }
      this.types.forEach(type => {
        let list = this.services[type]
        list.forEach((service, index) => {
          console.log('Pinging service ' + type + '#' + index)
          if (!service) {
            this.markDead(type, index)
          }
          let url = ''
          if (service.proto === 'HTTPS') {
            url = 'https://' + service.host + ':' + service.httpsPort + service.health
          } else {
            url = 'http://' + service.host + ':' + service.httpPort + service.health
          }
          axios.get(url).then(response => {
            if (!response || !response.body || !response.body.up) {
              this.markDead(type, index)
            }
          }).catch(e => {
            this.markDead(type, index)
          })
        })
      })

      this.clearLists()
    },

    markDead: function (type, index) {
      if (!this.services || !this.services[type] || typeof this.services[type] === 'undefined') {
        return
      }
      if (index >= this.services[type].length) {
        return
      }
      this.services[type][index].isDead = true
    },

    clearLists: function () {
      if (!this.types || typeof this.types === 'undefined' || !this.services) {
        return
      }
      let newServices = {}
      let newTypes = []
      this.types.forEach(type => {
        if (!this.services.hasOwnProperty(type) || typeof this.services[type] === 'undefined') {
          return
        }
        this.services[type].forEach(service => {
          if (!service.isDead) {
            if (typeof newServices[type] === 'undefined') {
              newServices[type] = []
            }
            newServices[type].push(service)
          }
        })
        if (newServices[type] && newServices[type].length > 1) {
          newTypes.push(type)
        }
      })

      this.services = newServices
      this.types = newTypes
    },

    getNotEmptyTypesNumber: function () {
      let i = 0
      this.types.forEach(item => {
        if (this.services[item] && this.services[item].length > 0) {
          i++
        }
      })
      return i
    },

    update: async function () {
      this.services = {}
      this.types = []
      try {
        const response = await request.registryHealth(this.registryUrl)
        this.registry = {
          url: this.registryUrl,
          up: response.body.up
        }
      } catch (e) {
        this.registry = {
          url: this.registryUrl,
          up: false
        }
      }

      if (!this.registry.up) {
        return
      }

      try {
        const response = await request.getServiceTypes(this.registryUrl)
        console.log(response)
        if (!response.hasOwnProperty('types') || response.types.length === 0) {
          return
        }

        let promises = []
        response.types.forEach(serviceType => {
          promises.push(new Promise((resolve, reject) => {
            request.getServiceList(this.registryUrl, serviceType.type)
              .then(response => resolve(response.body))
              .catch(e => reject(e))
          }))
        })

        Promise.all(promises).then(values => {
          if (values && values.hasOwnProperty('type')) {
            this.registry[values.type].push(values)
            if (this.types.indexOf(values.type) < 0) {
              this.types.push(values.type)
            }
          }
        })
      } catch (e) {}
    }
  },

  async created () {
    // this.update()
    // setInterval(() => this.update(), 6000)
    setInterval(() => this.pingAndRemoveDead(), 2000)
  }
}
</script>

<style scoped>
.serviceTypeBlock {
  display: inline-block;
  margin: auto 10px;
  padding: 10px 5px 5px;
  background: #f6f6f6;
  vertical-align: top;
}
</style>
