<script lang="ts">
import routes from "./routes";
import Router from "svelte-spa-router";
import {Toaster} from "svelte-sonner";
</script>
<svelte:head>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous">
    <link href="https://fonts.googleapis.com/css2?family=Fira+Code:wght@300..700&family=Liu+Jian+Mao+Cao&family=Space+Grotesk:wght@300..700&family=Ubuntu+Mono:wght@400;700&display=swap" rel="stylesheet">
</svelte:head>

<Toaster richColors/>

<main>
    <Router {routes}/>
</main>