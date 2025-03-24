import { createWebHistory, createRouter } from "vue-router";

import MainApp from "@/pages/MainApp";
import NotePage from "@/pages/NotePage";
import NoteIdPage from "@/pages/NoteIdPage";
import RegistrationForm from '@/components/RegistrationForm.vue';
import LoginForm from '@/components/LoginForm.vue';
import ResetPassword from "@/components/ResetPassword.vue";
import NewPassword from "@/components/NewPassword.vue";
import ActivateAccount from "@/pages/ActivateAccount.vue";

const routes = [
  {
    path: '/',
    component: MainApp
  },
  {
    path: '/notes',
    component: NotePage
  },
  {
     path: '/reset-password',
     component: ResetPassword
   },
   {
    path: '/new-password',
    component: NewPassword
  },
  {
    path: '/notes/:id',
    component: NoteIdPage
  },
  {
    path: '/login',
    component: LoginForm
  },
  {
    path: '/register',
    component: RegistrationForm
  },
  {
    path: '/:code',  // Динамический маршрут для коротких ссылок
    component: MainApp
  },
  { path: '/activate', 
    component: ActivateAccount 
  },
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
});

export default router;
