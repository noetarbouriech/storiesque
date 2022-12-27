<script lang="ts">
	import { onMount } from "svelte";
    import { env } from '$env/dynamic/public';
    import { Toast, P, A, Hr, Button, Popover, TextPlaceholder, Spinner } from 'flowbite-svelte'
    import type { PageData } from './$types';
    import { BookOpen } from 'svelte-heros-v2';
    import { page } from '$app/stores';
	import StoryCard from '$lib/StoryCard.svelte';

    export let data: PageData;
    type page = {
        id: Number,
        title: String,
        body: String,
        choices: Array<any>
    }

    // default values
    let currPage: page = {
        id: 0,
        title: "none",
        body: "",
        choices: []
    };
    let loading: boolean = true;

    async function changePage(pageId: Number): Promise<void> {
        loading = true;
        setTimeout(async () => {
        currPage = await fetch(`${env.PUBLIC_API_URL}/page/${pageId}`)
        .then(r => r.json());
        loading = false;
        console.log(currPage)
        }, 300)
    }

    async function addChoice(): Promise<void> {
        currPage.choices.push(await fetch(`${env.PUBLIC_API_URL}/page/${currPage.id}`, { method: 'POST' }));
    }

    onMount(async (): Promise<void> => {
        changePage(data.story.first_page_id)
    });
</script>

<Toast divClass="w-full max-w-fit p-4" class="mx-auto mb-8" simple>
    <svelte:fragment slot="icon">
        <BookOpen />
    </svelte:fragment>
    You are reading <A id="story-title" href="/story/{$page.params.slug}" color="text-gray-900 dark:text-white" class="font-bold">{data.story.title}</A>
</Toast>

<Popover triggeredBy="#story-title" placement="bottom-start">
    <StoryCard
        title={data.story.title}
        description={data.story.description}
        author={data.story.author_name}
        id={data.story.id}
    />
</Popover>

{#if loading}
    <TextPlaceholder class="md:mx-16 lg:mx-32 xl:mx-64"/>
    <Hr class="my-4 mx-auto md:my-8" width="w-48" height="h-1" />
    <div class="mx-auto flex flex-col justify-center max-w-fit">
        <Button class="mb-2 break-all">
            <Spinner class="mr-3" size="4" color="white" />
            Loading ...
        </Button>
    </div>
{:else} 
    <P firstupper justify class="md:mx-16 lg:mx-32 xl:mx-64">{currPage.body}</P>
    <Hr class="my-4 mx-auto md:my-8" width="w-48" height="h-1" />
    <div class="mx-auto flex flex-col justify-center max-w-fit">
        {#if currPage.choices}
            {#each currPage.choices as choice}
            <Button on:click={() => changePage(choice)} class="mb-2 break-all">{choice}</Button>
            {/each}
        {/if}
        <Button on:click={addChoice} class="break-all" outline>Add a choice</Button>
    </div>
{/if}