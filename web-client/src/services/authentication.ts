import { computed, ref } from "vue"
import type { Router } from "vue-router"
import router from "../router"

const jwtLocalStorageKey = "JWT"
const jwt = ref<string | null>(localStorage.getItem(jwtLocalStorageKey) || null)
const username = computed(() => {
    if (jwt.value === null) {
        return null
    }
    return JSON.parse(atob(jwt.value.split('.')[1])).name
})

export const useAuthentication = () => {
    return {
        // TODO: 401 => reset jwt
        isAuthenticated: () => jwt.value !== null,
        authenticate: () => {
            // TODO: Replace with real redirection to auth provider
            router.push({
                name: "login",
                query: {
                    token: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c',
                }
            })
        },
        setToken: (token: string) => {
            jwt.value = token
            localStorage.setItem(jwtLocalStorageKey, token)
        },
        username
    }
}

export const redirectToAuthenticationPage = (router: Router) => {
    router.beforeEach((to) => {
        if (to.name === "login") {
            const { isAuthenticated, setToken } = useAuthentication()
            if (isAuthenticated()) {
                return { name: "home" }
            }
            const token = to.query.token
            if (token && !Array.isArray(token)) {
                setToken(token)
                return { name: "home" }
            }
            return
        }
        const { isAuthenticated } = useAuthentication()
        if (!isAuthenticated()) {
            return { name: "login" }
        }
    })
}