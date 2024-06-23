<script lang="ts">

import Page from "../components/Page.svelte";
import {onMount} from "svelte";
import {Breadcrumb, BreadcrumbItem, Button, Card, Input, Select, Spinner, Textarea} from "flowbite-svelte";
import {getActiveEvent, getEvent, putEvent} from "../services/services";
import {eventStatusStore} from "../stores/stores";
import CodeJamEvent from "../models/event";
import FormField from "../components/FormField.svelte";
import Form from "../components/Form.svelte";

export let params : any; // set by svelte-spa-router

let statusOptions = [];
let formData : CodeJamEvent | null = null;
let form : Form;
let isSaving : boolean = false;

/* Convert the API statuses to Select options */
function mapStatusOptions(statuses) {
    if (statuses) {
        statusOptions = statuses.map((status) => {
            return {
                name: status.Title,
                value: status.Id
            }
        });
    }
}

eventStatusStore.subscribe((statuses) => {
    mapStatusOptions(statuses);
});

let clearErrors : () => {};

let parseResponse : (response: object) => {};

function saveForm() {
    if (formData !== null) {
        isSaving = true;
        clearErrors();
        putEvent(formData)
            .then((response) => {
                getActiveEvent();
                parseResponse(response);
                response.json()
                    .then(() => {
                        window.location.href = '/#/admin/events';
                        isSaving = false;
                    })
                    .catch(() => {
                        isSaving = false;
                    })
            })
            .catch((err) => {
                console.error("Error saving event", err);
                isSaving = false;
            });
    }
}

onMount(() => {
    if (params) {
        getEvent(params.id)
            .then((response) => {
                response.json()
                    .then((data) => {
                        formData = data as CodeJamEvent;
                    })
            });
    }
})

</script>


<Page>
    <Breadcrumb solid class="mb-4 w-full max-w-screen-xl">
        <BreadcrumbItem href="/#/" home>Home</BreadcrumbItem>
        <BreadcrumbItem href="/#/admin/events">Manage Events</BreadcrumbItem>
        <BreadcrumbItem>Edit Event</BreadcrumbItem>
    </Breadcrumb>

    <Card size="xl" class="w-full">
        <h2>Edit Event</h2>
        {#if formData !== null}
            <div class="flex flex-col gap-8 my-8">
                <Form bind:clearErrors={clearErrors} bind:parseResponse={parseResponse}>
                    <FormField name="Status">
                        <Select id="status" items={statusOptions} bind:value={formData.StatusId}></Select>
                    </FormField>

                    <FormField name="Title">
                        <Input bind:value={formData.Title}></Input>
                    </FormField>

                    <FormField name="Description">
                        <Textarea rows=10 bind:value={formData.Description}></Textarea>
                    </FormField>

                    <FormField name="Rules">
                        <Textarea rows=10 bind:value={formData.Rules}></Textarea>
                    </FormField>
                </Form>
            </div>

            <Button on:click={saveForm} disabled={isSaving}>
                {#if isSaving}
                    <Spinner/>
                {:else}
                    Save
                {/if}
            </Button>
        {:else}
            <Spinner/>
        {/if}
    </Card>
</Page>

