import { createApp } from 'vue'
import App from './App'
import commponents from '@/components/UI'
import router from '@/router/router'
import directives from './directives'

const app = createApp(App)

commponents.forEach(commponent => {
    app.component(commponent.name, commponent)
})

directives.forEach(directive => {
    app.directive(directive.name, directive)
})

app
    .use(router)
    .mount('#app');
