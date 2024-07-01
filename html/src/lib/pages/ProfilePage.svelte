<script lang="ts">

    import Page from "../components/Page.svelte";
    import {Button, Card, Input, Spinner, Tooltip} from "flowbite-svelte";
    import FormField from "../components/FormField.svelte";
    import Form from "../components/Form.svelte";
    import type {ActiveUser, User} from "../models/user";
    import {putProfile} from "../services/services";
    import {activeUserStore} from "../stores/stores";
    import {FontAwesomeIcon} from "@fortawesome/svelte-fontawesome";
    import {faLock} from "@fortawesome/free-solid-svg-icons";

    let isSaving : boolean = false;
    let formData: User | null = null;

    $: {
        if ($activeUserStore !== null && $activeUserStore.user !== null && formData == null) {
            formData = {...$activeUserStore.user};
        }
    }
    let clearErrors: () => {};
    let parseResponse: (response: object) => {};

    async function saveForm() {
        if (formData !== null) {
            isSaving = true;
            clearErrors();
            try {
                const response = await putProfile(formData.DisplayName);
                parseResponse(response);
                const responseData = await response.json();
                activeUserStore.set(<ActiveUser>{user: <User>responseData.Data, loggedIn: true});
                isSaving = false;
            } catch(err) {
                console.error("Error saving profile: ", err);
                isSaving = false;
            }
        }
    }

</script>


<Page>

    <Card size="lg" class="w-full">
        <h2>Edit Profile</h2>
        {#await $activeUserStore}
        {:then activeUser}
            {#if formData !== null}
                {#if activeUser !== null}
                    <div class="flex flex-col gap-8 my-8">
                        <Form bind:clearErrors bind:parseResponse>
                            {#if activeUser.user?.LockDisplayName}
                                <div class="font-bold">
                                    Display Name <FontAwesomeIcon icon={faLock}/>
                                    <Tooltip>Your Display Name has been locked by an admin and may not be changed</Tooltip>
                                </div>
                                <div>{formData.DisplayName}</div>
                            {:else}
                                <FormField label="Display Name" name="DisplayName">
                                    <Input bind:value={formData.DisplayName}></Input>
                                </FormField>
                            {/if}
                        </Form>
                    </div>

                    <Button on:click={saveForm} disabled={isSaving || activeUser.user?.LockDisplayName}>
                        {#if isSaving}
                            <Spinner />
                        {:else}
                            Save
                        {/if}
                    </Button>
                {/if}
            {/if}
        {/await}
    </Card>


</Page>