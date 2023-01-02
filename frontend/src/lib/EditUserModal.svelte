<script lang="ts">
	import { Button, Modal, Input, Label, Textarea } from 'flowbite-svelte';
    import { env } from '$env/dynamic/public';
    import { userStore } from '../store';
	import { User } from 'svelte-heros-v2';

    export let open: boolean;

    export let id: number;
    export let username: string;
    export let email: string;
    let password: string = "";

    async function editUser(e: Event): Promise<void> {
        const formData = new FormData(e.target as HTMLFormElement);
        const jsonData = JSON.stringify(Object.fromEntries(formData));
        try {
            const response = await fetch(`${env.PUBLIC_API_URL}/user/${id}`, {
                method: 'PUT',
                credentials: 'include',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: jsonData
            });
            const data = await response.json();
            if (response.ok) { 
                open = false;
                if (id === $userStore.id) {
                    $userStore = {
                        id: $userStore.id,
                        username: username,
                        email: email,
                        is_admin: $userStore.is_admin,
                        has_img: $userStore.has_img
                    };
                }
            } else {
                alert(data.message);
            }
        } catch (error) {
            console.error('Error:', error);
        }
    }
</script>
<Modal bind:open={open} size="md" autoclose={false}>
    <form class="flex flex-col space-y-6" on:submit|preventDefault={editUser}>
        <h3 class="inline-flex items-center text-xl font-medium text-gray-900 dark:text-white p-0"><User class="mr-1"/>Edit Account</h3>
        <Label class="space-y-2">
        <span>Username</span>
        <Input type="text" name="username" bind:value={username} required />
        </Label>
        <Label class="space-y-2">
        <span>Email</span>
        <Input type="email" name="email" bind:value={email} required />
        </Label>
        <Label class="space-y-2">
        <span>Password</span>
        <Input type="password" name="password" bind:value={password} />
        </Label>
        <Button type="submit" class="w-full1">Apply changes</Button>
    </form>
</Modal>


