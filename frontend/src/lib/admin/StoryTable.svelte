<script lang="ts">
    import { Table, TableHead, TableHeadCell, TableBody, TableBodyRow, TableBodyCell, Button, Toggle, Modal, PaginationItem } from 'flowbite-svelte';
	import { onMount } from 'svelte';
    import { env } from '$env/dynamic/public';
	import { ArrowTopRightOnSquare, Trash } from 'svelte-heros-v2';
	import slugify from 'slugify';

    let stories: Array<any> = [];
    let popupModal: boolean = false;
    let selectedId: number = 0;
    let page: number = 1;

    onMount(async () => {
        loadStories(1);
    });

    async function loadStories(page: number): Promise<void> {
        stories = await (await fetch(`${env.PUBLIC_API_URL}/story?page=${page}`)).json();
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
                loadStories(page);
            }
        } catch (error) {
            console.error('Error:', error);
        }
    }

    async function updateFeatured(storyId: number): Promise<void> {
        try {
            const response = await fetch(`${env.PUBLIC_API_URL}/story/${storyId}/featured`, {
                method: 'PUT',
                credentials: 'include'
            });
            const data = await response.json();
            if (!response.ok) { 
                alert(data.message);
            }
        } catch (error) {
            console.error('Error:', error);
        }
    }

    const previous = async () => {
        if (page === 1) return;
        page--;
        loadStories(page);
    };
    const next = async () => {
        if (stories.length < 30) return;
        page++;
        loadStories(page);
    };

</script>

<Table hoverable noborder>
    <TableHead>
        <TableHeadCell>ID</TableHeadCell>
        <TableHeadCell>Title</TableHeadCell>
        <TableHeadCell>Author</TableHeadCell>
        <TableHeadCell>Is Featured ?</TableHeadCell>
        <TableHeadCell>Actions</TableHeadCell>
    </TableHead>
    <TableBody>
        {#each stories as story}
            <TableBodyRow noborder>
                <TableBodyCell>{story.id}</TableBodyCell>
                <TableBodyCell>{story.title}</TableBodyCell>
                <TableBodyCell>{story.author_name}</TableBodyCell>
                <TableBodyCell>
                    <Toggle size="large" on:change={() => updateFeatured(story.id)} checked={story.is_featured} />
                </TableBodyCell>
                <TableBodyCell>
                    <Button href="/story/{slugify(story.title,{lower: true})}-{story.id}" target="_blank" class="!p-2"><ArrowTopRightOnSquare /></Button>
                    <Button on:click={() => {popupModal = true; selectedId = story.id}} class="!bg-red-600 !p-2"><Trash /></Button>
                </TableBodyCell>
            </TableBodyRow>
        {/each}
    </TableBody>
</Table>

<div class="mt-4 place-content-center flex space-x-3">
    <PaginationItem class="flex items-center" on:click={previous}>
      <svg class="mr-2 w-5 h-5" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M7.707 14.707a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 1.414L5.414 9H17a1 1 0 110 2H5.414l2.293 2.293a1 1 0 010 1.414z" clip-rule="evenodd"/></svg>
      Prev
    </PaginationItem>
    <PaginationItem class="flex items-center" on:click={next}>
      Next
      <svg class="ml-2 w-5 h-5" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M12.293 5.293a1 1 0 011.414 0l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414-1.414L14.586 11H3a1 1 0 110-2h11.586l-2.293-2.293a1 1 0 010-1.414z" clip-rule="evenodd"/></svg>
    </PaginationItem>
</div>

<Modal bind:open={popupModal} size="xs" autoclose>
  <div class="text-center">
      <svg aria-hidden="true" class="mx-auto mb-4 w-14 h-14 text-gray-400 dark:text-gray-200" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
      <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">Are you sure you want to delete this story?</h3>
      <Button on:click={() => deleteStory(selectedId)} color="red" class="mr-2">Yes, I'm sure</Button>
      <Button color='alternative'>No, cancel</Button>
  </div>
</Modal>
