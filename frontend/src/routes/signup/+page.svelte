<script lang="ts">
    import { Card, Heading, Label, Input, Button, Helper } from 'flowbite-svelte'
    import { env } from '$env/dynamic/public';
	import { goto } from '$app/navigation';
	import { ExclamationCircle } from 'svelte-heros-v2';
	import { onMount } from 'svelte';
    import { userStore } from '../../store';

    let errorMsg: string;
    let checkPass: string;

    async function signUp(e: Event): Promise<void> {
        const formData = new FormData(e.target as HTMLFormElement);
        const jsonData = JSON.stringify(Object.fromEntries(formData));
        if (formData.get("password") !== checkPass) {
            errorMsg = "Passwords do not match";
            return;
        }
        try {
            const response = await fetch(`${env.PUBLIC_API_URL}/signup`, {
                method: 'POST',
                credentials: 'include',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: jsonData
            });
            const data = await response.json();
            if (response.ok) { 
                goto("/");
            } else {
                errorMsg = data.message;
            }
        } catch (error) {
            console.error('Error:', error);
        }
    }

    onMount(async() => {
        if ($userStore.username !== "") goto("/");
    });

</script>

<svelte:head>
  <title>Sign up - Storiesque</title>
</svelte:head>

<Heading class="text-center" tag="h1">Sign up</Heading>

<Card size="lg" class="mx-auto mt-8">
    <form on:submit|preventDefault={signUp}>
        <div class="mb-6">
            <Label for="username" class="mb-2">Username</Label>
            <Input type="text" name="username" id="username" placeholder="john_doe" autocomplete="username" required />
        </div>
        <div class="mb-6">
            <Label for="email" class="mb-2">Email address</Label>
            <Input type="email" name="email" id="email" placeholder="john.doe@company.com" autocomplete="email" required />
        </div>
        <div class="mb-6">
            <Label for="password" class="mb-2">Password</Label>
            <Input type="password" name="password" id="password" placeholder="•••••••••" autocomplete="new-password" required />
        </div>
        <div class="mb-6">
            <Label for="confirm_password" class="mb-2">Confirm password</Label>
            <Input bind:value={checkPass} type="password" id="confirm_password" placeholder="•••••••••" autocomplete="new-password" required />
        </div>
        <Button type="submit">Submit</Button>
        {#if errorMsg}
            <Helper class="inline ml-2" color="red"><ExclamationCircle class="inline-block mr-1"/>{errorMsg}</Helper>
        {/if}
    </form>
</Card>