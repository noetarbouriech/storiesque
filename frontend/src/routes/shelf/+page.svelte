<script lang="ts">
    import { Heading, PaginationItem } from 'flowbite-svelte'
    import StoryCard from '$lib/StoryCard.svelte';
    import { env } from '$env/dynamic/public';
	import { onMount } from 'svelte';
    import { userStore } from '../../store';
	import { goto } from '$app/navigation';

    let stories: Array<any> = [];
    let page: number = 1;

    const previous = async () => {
        if (page === 1) return;
        page--;
        const res = await fetch(`${env.PUBLIC_API_URL}/shelf?page=${page}`, {
            credentials: 'include',
        });
        stories = await res.json();
    };
    const next = async () => {
        if (stories.length < 30) return;
        page++;
        const res = await fetch(`${env.PUBLIC_API_URL}/shelf?page=${page}`, {
            credentials: 'include',
        });
        stories = await res.json();
    };

    onMount(async(): Promise<void> => {
        if ($userStore.username === "") {
            goto("/")
        } else {
            const res = await fetch(`${env.PUBLIC_API_URL}/shelf`, {
                credentials: 'include',
            });
            stories = await res.json();
        }
    });
</script>

<svelte:head>
  <title>My Shelf - Storiesque</title>
</svelte:head>

<Heading class="text-center" tag="h1">My shelf</Heading>

<div class="max-w-2xl mx-auto px-4 py-8 lg:max-w-7xl grid grid-cols-1 gap-y-10 sm:grid-cols-2 gap-x-8 lg:grid-cols-3">
    {#each stories as story}
        <StoryCard
            title={story.title}
            description={story.description}
            author={story.author_name}
            id={story.id}
            has_img={story.has_img}
        ></StoryCard>
    {/each}
</div>

<div class="place-content-center flex space-x-3">
    <PaginationItem class="flex items-center" on:click={previous}>
      <svg class="mr-2 w-5 h-5" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M7.707 14.707a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 1.414L5.414 9H17a1 1 0 110 2H5.414l2.293 2.293a1 1 0 010 1.414z" clip-rule="evenodd"/></svg>
      Prev
    </PaginationItem>
    <PaginationItem class="flex items-center" on:click={next}>
      Next
      <svg class="ml-2 w-5 h-5" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M12.293 5.293a1 1 0 011.414 0l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414-1.414L14.586 11H3a1 1 0 110-2h11.586l-2.293-2.293a1 1 0 010-1.414z" clip-rule="evenodd"/></svg>
    </PaginationItem>
</div>
