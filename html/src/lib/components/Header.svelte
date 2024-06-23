<script lang="ts">

import Banner from "./Banner.svelte";
import {Navbar, NavHamburger, NavLi, NavUl} from "flowbite-svelte";
import {loggedInStore, userStore} from "../stores/stores";
import DiscordIcon from "./DiscordIcon.svelte";

</script>


<div class="flex flex-row w-fit pl-8">
    <Banner size="normal" />
    <div class="absolute right-0 top-0">
        <Navbar class="rounded-bl-3xl">
            <NavHamburger/>
            <NavUl>
                <!-- It is important to prefix links with /#/ to prevent reloading of the entire app -->
                <NavLi href="/#/">Home</NavLi>
                {#if $loggedInStore}
                    {#if $userStore?.Role === "ADMIN" }
                         <NavLi href="/#/admin/events">Manage Events</NavLi>
                    {/if}
                    <NavLi href="/">My Team</NavLi>
                    <NavLi href="/">Profile</NavLi>
                    <NavLi href="/user/logout">Logout</NavLi>
                {:else}
                    <NavLi href="/oauth/redirect">Login with Discord <DiscordIcon/></NavLi>
                {/if}
            </NavUl>
        </Navbar>
    </div>
</div>