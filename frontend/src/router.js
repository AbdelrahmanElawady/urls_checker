import { createWebHistory, createRouter } from "vue-router";
import Checker from "@/components/Checker.vue";
import Home from "@/components/Home.vue";

const routes = [
  {
    path: "/check",
    name: "Checker",
    component: Checker,
  },
  {
    path: "/",
    name: "Home",
    component: Home,
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;