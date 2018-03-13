import axios from 'axios'

function registryHealth (registryUrl) {
  return axios.get(registryUrl + '/registry/health')
    .then(response => response.body.up)
}

function getServiceTypes (registryUrl) {
  return axios.get(registryUrl + '/registry/service/types')
    .then(response => response.body)
}

function getServiceList (registryUrl, serviceType) {
  return axios.get(registryUrl + '/registry/service/types/' + serviceType)
    .then(response => response.body)
}

function getBestService (registryUrl, serviceType) {
  return axios.get(registryUrl + '/registry/service/types/' + serviceType + '/best')
    .then(response => response.body)
}

function getInstanceInfoByName (registryUrl, instanceName) {
  return axios.get(registryUrl + '/registry/service/' + instanceName)
    .then(response => response.body)
}

function getInstanceInfo (registryUrl, instanceType, instanceId) {
  return axios.get(registryUrl + '/registry/service/types/' + instanceType + '/' + instanceId)
    .then(response => response.body)
}

export default {
  registryHealth,
  getServiceTypes,
  getServiceList,
  getBestService,
  getInstanceInfoByName,
  getInstanceInfo
}
