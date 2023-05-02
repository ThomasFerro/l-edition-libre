import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import "./styles.css"
import 'vuetify/styles'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import { createAuth0 } from '@auth0/auth0-vue';


const vuetify = createVuetify({
  components,
  directives,
})

const app = createApp(App)

app.use(router(app))
app.use(vuetify)

console.log(import.meta.env.VITE_AUTH0_CLIENTID)
app.use(
  createAuth0({
    domain: import.meta.env.VITE_AUTH0_DOMAIN,
    clientId: import.meta.env.VITE_AUTH0_CLIENTID,
    authorizationParams: {
      redirect_uri: `${window.location.origin}/login`
    }
  })
);

app.mount('#app')
