import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/about',
    name: 'About',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "about" */ '../views/About.vue')
  },
  {
    path: '/employee',
    name: 'Employees',
    component: ()=> import('@/views/Employees.vue')
  },
  {
    path: '/employee/:uid',
    name: 'Employee',
    component: ()=> import('@/views/Employee.vue')
  },
  {
    path: '/team',
    name: 'Teams',
    component: ()=> import('@/views/Teams.vue')
  },
  {
    path: '/team/:uid',
    name: 'Team',
    component: ()=> import('@/views/Team.vue')
  }


]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

export default router
