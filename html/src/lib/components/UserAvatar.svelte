<script lang="ts">

    import {Avatar} from "flowbite-svelte";
    import type {User} from "../models/user";

    export let user : User | null;

</script>


{#if user && user.AvatarId}
    {#if user.AvatarId && user.AvatarId.startsWith("a_")}
    <Avatar border src="https://cdn.discordapp.com/avatars/{user.ServiceUserId}/{user.AvatarId}.gif?size=48" {...$$restProps} />
    {:else}
    <Avatar border src="https://cdn.discordapp.com/avatars/{user.ServiceUserId}/{user.AvatarId}.png?size=48" {...$$restProps} />
    {/if}
{:else}
    <Avatar border {...$$restProps} />
{/if}