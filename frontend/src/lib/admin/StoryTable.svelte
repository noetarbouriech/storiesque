<script lang="ts">
    import { Table, TableHead, TableHeadCell, TableBody, TableBodyRow, TableBodyCell, Button, Toggle, Modal } from 'flowbite-svelte';
	import { onMount } from 'svelte';
    import { env } from '$env/dynamic/public';
	import { ArrowTopRightOnSquare, Trash } from 'svelte-heros-v2';
	import slugify from 'slugify';

    let stories: Array<any> = [];
    let popupModal: boolean = false;
    let selectedId: number = 0;

    onMount(async () => {
        loadStories();
    });

    async function loadStories(): Promise<void> {
        stories = await (await fetch(`${env.PUBLIC_API_URL}/story`)).json();
    }

    async function deleteStory(storyId: number): Promise<void> {
        try {
            const response = await fetch(`${env.PUBLIC_API_URL}/story/${storyId}`, {
                method: 'DELETE',
                credentials: 'include'
            });
            const data = await response.json();
            if (!response.ok) { 
                alert(data.message);
            } else {
                loadStories();
            }
        } catch (error) {
            console.error('Error:', error);
        }
    }

</script>

<Table hoverable noborder>
    <TableHead>
        <TableHeadCell>ID</TableHeadCell>
        <TableHeadCell>Title</TableHeadCell>
        <TableHeadCell>Author</TableHeadCell>
        <TableHeadCell>Actions</TableHeadCell>
    </TableHead>
    <TableBody>
        {#each stories as story}
            <TableBodyRow noborder>
                <TableBodyCell>{story.id}</TableBodyCell>
                <TableBodyCell>{story.title}</TableBodyCell>
                <TableBodyCell>{story.author_name}</TableBodyCell>
                <TableBodyCell>
                    <Button href="/story/{slugify(story.title,{lower: true})}-{story.id}" target="_blank" class="!p-2"><ArrowTopRightOnSquare /></Button>
                    <Button on:click={() => {popupModal = true; selectedId = story.id}} class="bg-red-600 !p-2"><Trash /></Button>
                </TableBodyCell>
            </TableBodyRow>
        {/each}
    </TableBody>
</Table>

<Modal bind:open={popupModal} size="xs" autoclose>
  <div class="text-center">
      <svg aria-hidden="true" class="mx-auto mb-4 w-14 h-14 text-gray-400 dark:text-gray-200" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
      <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">Are you sure you want to delete this story?</h3>
      <Button on:click={() => deleteStory(selectedId)} color="red" class="mr-2">Yes, I'm sure</Button>
      <Button color='alternative'>No, cancel</Button>
  </div>
</Modal>
