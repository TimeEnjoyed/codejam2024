<script lang="ts">
import {
	Button,
	Card,
	Helper,
	Input,
	Radio,
	Spinner,
	Textarea
} from 'flowbite-svelte';
import toast from 'svelte-french-toast';
import Form from "../components/Form.svelte";
import FormField from "../components/FormField.svelte";
import CodeJamTeam from "../models/team";
import { postTeam } from "../services/services";
import { activeEventStore } from "../stores/stores";

interface Selections {
    value: string,
    name: string;
}

let selectedLangs: Selections[] = [];
let languages: Selections[] = [
    { value: 'py', name: 'Python'},
    { value: 'ph', name: 'PHP'},
    { value: 'js', name: 'JavaScript'},
    { value: 'ts', name: 'TypeScript'}, 
    { value: 'go', name: 'Go'}
];
let languageValues: string[] = [];

let textareaprops = {
    id: 'aboutTextArea',
    name: 'message',
    label: 'What do you want out of this team:',
    rows: 4,
    placeholder: "I want to use threejs so I can learn it!"
};

let formData: CodeJamTeam | null = new CodeJamTeam();
let isSaving: boolean = false;
let teamCreated: boolean = false;


let clearErrors: () => {};
let parseResponse: (response: object) => {};

function saveForm() {
    if (formData !== null) {
        isSaving = true;
        clearErrors();
        
        formData.EventId = $activeEventStore?.Id || ""
        // Step 1: Post Team Data API
        postTeam(formData)
            .then((response) => {
                // parseResponse(response);
                // const url = new URL(response.url);
                // const pathSegments = url.pathname.split('/');
                // const teamId = pathSegments[pathSegments.length - 1];
                response.json()
                    .then((data) => {
                        // Team creation successful, letting svelte page know:
                        teamCreated = true;
                        // Stepp 1: GET team info
                        // this uses routes.ts --> MyTeam.svelte page
                        window.location.href = `/#/team/${data.id}`;
                        toast.success("You've successfully created a team")
                        isSaving = false;
                        
                    })
                    .catch(() => {
                        isSaving = false;
                    });
            })
            .catch((err) => {
                console.error("Error saving event", err);
                isSaving = false;
            });
    }
}

</script>

<Card size="xl" class="w-full">

    {#if formData !== null}
        <div class="flex flex-col gap-8 my-8">
            <Form bind:clearErrors={clearErrors} bind:parseResponse={parseResponse}>

                <FormField label="Team Name:" name="TeamName">
                    <Input bind:value={formData.Name}></Input>
                </FormField>
                <div>
                    <Radio name="team-type" bind:group={formData.Visibility} value="public">Public Team</Radio>
                    <Helper class="ml-6 ">(If you want your team to be searchable.)</Helper>
                </div>
                <div>
                    <Radio name="team-type" bind:group={formData.Visibility} value="private">Private Team</Radio>
                    <Helper class="ml-6">(Your team will be invite only)</Helper>
                </div>

                <FormField label="Your general availability:" name="TeamAvailability">
                    <Input bind:value={formData.Availability} placeholder="example: weekends, mon (1-2pm)"></Input>
                </FormField>

                <!-- <MultiSelect id="multi-close" items={languages} bind:value={formData.Technologies} /> -->
                <FormField label="Your technologies:" name="TeamTechnologies">
                    <Input bind:value={formData.Technologies} placeholder="example: python, nextjs, django, sql"></Input>
                </FormField>

                <FormField label="What do you want out of this team?" name="Description">
                <!-- <Label for="aboutTextArea">What do you want out of this team?</Label> -->
                    <Textarea {...textareaprops} bind:value={formData.Description} />
                </FormField>
            </Form>
            

            <Button on:click={saveForm} disabled={isSaving}>
                {#if isSaving}
                    <Spinner/>
                {:else}
                    Save
                {/if}
            </Button>
        </div>
    {:else}
        <Spinner/>
    {/if}
</Card>
