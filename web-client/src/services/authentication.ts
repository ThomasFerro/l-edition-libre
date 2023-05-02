import { computed, ref, type App } from "vue"
import type { NavigationGuard, Router } from "vue-router"
import router from "../router"
import { useAuth0, authGuard, createAuthGuard } from "@auth0/auth0-vue"

export const useAuthentication = () => {
    const auth0 = useAuth0()
    return {
        isAuthenticated: () => auth0.isAuthenticated,
        authenticate: () => {
            auth0.loginWithRedirect();
        },
        logout: auth0.logout,
        username: computed(() => auth0.user.value.nickname),
    }
}

export const shouldbeAuthenticated = (app: App): NavigationGuard => createAuthGuard(app)
