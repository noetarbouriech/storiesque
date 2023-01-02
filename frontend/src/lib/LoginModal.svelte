<script lang="ts">
	import { Button, Modal, Input, Label } from 'flowbite-svelte';
    import { env } from '$env/dynamic/public';
    import { userStore } from '../store';

    export let open: boolean;

    async function login(e: Event): Promise<void> {
        const formData = new FormData(e.target as HTMLFormElement);
        const jsonData = JSON.stringify(Object.fromEntries(formData));
        try {
            const response = await fetch(`${env.PUBLIC_API_URL}/login`, {
                method: 'POST',
                credentials: 'include',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: jsonData
            });
            const data = await response.json();
            if (response.ok) { 
                $userStore = {
                    id: data.id,
                    username: data.username,
                    email: data.email,
                    is_admin: data.is_admin,
                    has_img: data.has_img
                };
                open = false;
            } else {
                alert(data.message);
            }
        } catch (error) {
            console.error('Error:', error);
        }
    }
</script>
<Modal bind:open={open} size="xs" autoclose={false}>
    <form class="flex flex-col space-y-6" on:submit|preventDefault={login}>
        <h3 class="text-xl font-medium text-gray-900 dark:text-white p-0">Welcome back!</h3>
        <Label class="space-y-2">
        <span>Email</span>
        <Input type="email" name="email" placeholder="name@example.com" required />
        </Label>
        <Label class="space-y-2">
        <span>Your password</span>
        <Input type="password" name="password" placeholder="•••••" required />
        </Label>
        <Button type="submit" class="w-full1">Log In</Button>
        <div class="text-sm font-medium text-gray-500 dark:text-gray-300">
            Don't have an account? <a on:click={() => open=false} href="/signup" class="text-blue-700 hover:underline dark:text-blue-500">Sign up</a>
        </div>
    </form>
</Modal>
