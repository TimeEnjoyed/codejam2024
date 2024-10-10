<script lang="ts">

import Banner from "./Banner.svelte";
import {
    Dropdown,
    DropdownDivider,
    DropdownHeader,
    DropdownItem,
    Navbar,
    NavHamburger,
    NavLi,
    NavUl
} from "flowbite-svelte";
import {activeUserStore, activeEventStore} from "../stores/stores";
import DiscordIcon from "./DiscordIcon.svelte";
import {location} from "svelte-spa-router";
import UserAvatar from "./UserAvatar.svelte";

// reactive statement - "location" will change whenever url changes
$: activeUrl = '/#' + $location;

</script>

<div class="flex flex-row w-fit pl-8">
 <!-- <div> -->
    <Banner size="normal" />
    <div class="absolute right-0 top-0 justify-between">

        <!-- We don't want to set any initial component state until the activeUserStore has been initially
             set to avoid flickering state changes while data is initially loaded -->
        {#await $activeUserStore}
            <!-- do nothing-->
        {:then activeUser}
            {#if activeUser !== null}
                <Navbar class="rounded-bl-2xl">
                    {#if activeUser != null && activeUser.user != null}
                        <div class="flex items-center md:order-2 justify-between flex-grow">
                            <div id="menu-avatar">
                                <UserAvatar user={activeUser.user} class="mx-4 cursor-pointer" size="md"></UserAvatar>
                            </div>
                            <NavHamburger class1="w-full md:flex md:w-auto md:order-1" style="background-color: #0000001a"/>
                        </div>
                        <!-- the binding to the triggeredBy element happens in the components onMount, and the target
                             element needs to already be mounted so the Dropdown can find it.  So we need to wait for
                             the activeUserStore to be initially set before these components are mounted to prevent
                             a race condition when the Dropdown doesn't bind and therefore won't work -->
                        <Dropdown class="w-44" placement="bottom" triggeredBy="#menu-avatar" trigger="hover">
                            <DropdownHeader>{activeUser.user?.DisplayName}</DropdownHeader>
                            <DropdownItem href="/#/profile">Profile</DropdownItem>
                            <DropdownDivider />
                            <DropdownItem href="/user/logout">Logout</DropdownItem>
                        </Dropdown>
                    {/if}
                    <NavUl {activeUrl} classUl="items-left px-4 py-1">
                        <!-- It is important to prefix links with /#/ to prevent reloading of the entire app -->
                        <NavLi href="/#/">Home</NavLi>
                        {#if activeUser.user !== null}
                            {#if activeUser.user?.Role === "ADMIN" }
                                <NavLi href="/#/admin/events">Manage Events</NavLi>
                            {/if}
                            {#if $activeEventStore && $activeEventStore.AllowSignups}
                                <NavLi href="/#/team">Sign-Up</NavLi>
                                <NavLi href="/#/teams">My Teams</NavLi>
                            {/if}

                        {:else}
                            <NavLi href="/oauth/redirect?redirect={$location}">Login with Discord <DiscordIcon/></NavLi>
                        {/if}
                    </NavUl>
                </Navbar>
            {/if}
        {/await}
    </div>
</div>
