<script setup lang="ts">
import { useAuthentication } from '@/services/authentication';

const { authenticate } = useAuthentication()
</script>

<template>
  <main>
    <v-btn data-test="Authenticate button" @click="authenticate()">Authenticate</v-btn>
  </main>
</template>
 