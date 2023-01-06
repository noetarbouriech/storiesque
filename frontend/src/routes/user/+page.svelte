<script lang="ts">
    import { Heading, PaginationItem } from 'flowbite-svelte'
    import UserCard from '$lib/UserCard.svelte';
    import type { PageData } from './$types';
    import { env } from '$env/dynamic/public';

    export let data: PageData;
    let page: number = 1;

    const previous = async () => {
        if (page === 1) return;
        page--;
        const res = await fetch(`${env.PUBLIC_API_URL}/user?page=${page}`);
        const users = await res.json();
        data = { users };
    };
    const next = async () => {
        if (data.users.length < 30) return;
        page++;
        const res = await fetch(`${env.PUBLIC_API_URL}/user?page=${page}`);
        const users = await res.json();
        data = { users };
    };
</script>

<svelte:head>
  <title>Users - Storiesque</title>
</svelte:head>

<Heading class="text-center" tag="h1">👤 Users</Heading>

<div class="max-w-2xl mx-auto px-4 py-8 lg:max-w-7xl grid grid-cols-1 gap-y-10 sm:grid-cols-2 gap-x-8 lg:grid-cols-3 xl:grid-cols-4">
    {#each data.users as user}
        <UserCard username={user.username} id={user.id} has_img={user.has_img}></UserCard>
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
