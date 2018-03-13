import Vue from 'vue'
import Router from 'vue-router'
import Health from '@/components/Health'
import Logs from '@/components/Logs'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/health',
      name: 'health',
      component: Health,
      props: {registryUrl: 'http://localhost:3001/'}
    },
    {
      path: '/logs',
      name: 'logs',
      component: Logs,
      props: {registryUrl: 'http://localhost:3001/'}
    }
  ]
})
