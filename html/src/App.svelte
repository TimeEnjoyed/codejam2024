<script lang="ts">
import routes from "./routes";
import Router from "svelte-spa-router";
import { Toaster } from 'svelte-french-toast'
</script>

<main>
    <Toaster />
    <Router {routes}/>
</main>