<script lang="ts">
	import Page from '../components/Page.svelte';
	import { Button, Card, Input, Spinner } from 'flowbite-svelte';
	import FormField from '../components/FormField.svelte';
	import Form from '../components/Form.svelte';
	import type { ActiveUser, User } from '../models/user';
	import { putProfile } from '../services/services';
	import { activeUserStore } from '../stores/stores';

	let isSaving: boolean = false;
	let formData: User | null = null;

	$: {
		if ($activeUserStore !== null && $activeUserStore.user !== null && formData == null) {
			formData = { ...$activeUserStore.user };
		}
	}
	let clearErrors: () => {};
	let parseResponse: (response: object) => {};

	async function saveForm() {
		if (formData !== null) {
			isSaving = true;
			clearErrors();
			try {
				const response = await putProfile(formData);
				parseResponse(response);
				const responseData = await response.json();
				activeUserStore.set(<ActiveUser>{ user: <User>responseData.Data, loggedIn: true });
				isSaving = false;
			} catch (err) {
				console.error('Error saving profile: ', err);
				isSaving = false;
			}
		}
	}
</script>

<Page>
	<Card size="xl" class="w-full">
		<h2>Edit Profile</h2>
		{#if formData !== null}
			<div class="flex flex-col gap-8 my-8">
				<Form bind:clearErrors bind:parseResponse>
					<FormField label="Display Name" name="DisplayName">
						<Input bind:value={formData.DisplayName}></Input>
					</FormField>
				</Form>
			</div>

			<Button on:click={saveForm} disabled={isSaving}>
				{#if isSaving}
					<Spinner />
				{:else}
					Save
				{/if}
			</Button>
		{:else}
			<Spinner />
		{/if}
	</Card>
</Page>
