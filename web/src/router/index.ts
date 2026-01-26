import { createRouter, createWebHistory } from 'vue-router'
import StreamList from '../views/StreamList.vue'
import AlgorithmList from '../views/AlgorithmList.vue'
import InferenceResultList from '../views/InferenceResultList.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: '/streams'
    },
    {
      path: '/streams',
      name: 'StreamList',
      component: StreamList
    },
    {
      path: '/algorithms',
      name: 'AlgorithmList',
      component: AlgorithmList
    },
    {
      path: '/inference-results',
      name: 'InferenceResultList',
      component: InferenceResultList
    }
  ]
})

export default router
