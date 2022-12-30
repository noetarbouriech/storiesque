<script lang="ts">
	import { Button, Modal, Input, Label, Textarea } from 'flowbite-svelte';
    import { env } from '$env/dynamic/public';
	import { goto } from '$app/navigation';
	import slugify from 'slugify';
	import { PlusCircle } from 'svelte-heros-v2';

    export let open: boolean;

    async function createStory(e: Event): Promise<void> {
        const formData = new FormData(e.target as HTMLFormElement);
        const jsonData = JSON.stringify(Object.fromEntries(formData));
        try {
            const response = await fetch(`${env.PUBLIC_API_URL}/story`, {
                method: 'POST',
                credentials: 'include',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: jsonData
            });
            const data = await response.json();
            if (response.ok) { 
                goto(`/story/${slugify(data.title,{lower: true})}-${data.id}`);
                open = false;
            } else {
                alert(data.message);
            }
        } catch (error) {
            console.error('Error:', error);
        }
    }
</script>
<Modal bind:open={open} size="md" autoclose={false}>
    <form class="flex flex-col space-y-6" on:submit|preventDefault={createStory}>
        <h3 class="inline-flex items-center text-xl font-medium text-gray-900 dark:text-white p-0"><PlusCircle class="mr-1"/>New Story</h3>
        <Label class="space-y-2">
        <span>Title</span>
        <Input type="text" name="title" required />
        </Label>
        <Label class="space-y-2">
        <span>Description</span>
        <Textarea type="text" name="description" />
        </Label>
        <Button type="submit" class="w-full1">Create story</Button>
    </form>
</Modal>

