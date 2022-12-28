<script lang="ts">
	import { onMount } from "svelte";
    import { env } from '$env/dynamic/public';
    import { userStore } from '../../../../store';
    import { Toast, P, A, Hr, Button, Popover, TextPlaceholder, Spinner, Span } from 'flowbite-svelte'
    import type { PageData } from './$types';
    import { BookOpen, InformationCircle, Plus } from 'svelte-heros-v2';
    import { page } from '$app/stores';
	import StoryCard from '$lib/StoryCard.svelte';

    export let data: PageData;
    type page = {
        id: Number,
        action: String,
        body: String,
        choices: {page_id: Number, action: String}[]
    }

    // default values
    let currPage: page = {
        id: 0,
        action: "",
        body: "",
        choices: []
    };
    let loading: boolean = true;

    async function changePage(pageId: Number): Promise<void> {
        loading = true;
        setTimeout(async () => {
        currPage = await fetch(`${env.PUBLIC_API_URL}/page/${pageId}`, {
            credentials: 'include',
        })
        .then(r => r.json());
        loading = false;
        }, 300)
    }

    async function addChoice(): Promise<void> {
        currPage.choices = [...currPage.choices,(await (await fetch(`${env.PUBLIC_API_URL}/page/${currPage.id}`, { 
            method: 'POST',
            credentials: 'include',
        })).json())];
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
    {#if data.story.first_page_id != currPage.id}
    <P justify class="md:mx-16 lg:mx-32 xl:mx-64 mb-8 md:text-xl" weight="light" size="lg" color="text-gray-500 dark:text-gray-400">
        <InformationCircle class="inline-flex"/>
        You chose to:
        <Span italic class="font-thin">{currPage.action}</Span>
    </P>
    {/if}
    <P firstupper justify class="md:mx-16 lg:mx-32 xl:mx-64">{currPage.body}</P>
    <Hr class="my-4 mx-auto md:my-8" width="w-48" height="h-1" />
    <div class="mx-auto flex flex-col justify-center max-w-fit">
        {#if currPage.choices}
            {#each currPage.choices as choice}
            <Button on:click={() => changePage(choice.page_id)} class="mb-2 break-all">{choice.action}</Button>
            {/each}
        {/if}
        {#if data.story.author_name == $userStore.username}
            <Button on:click={addChoice} class="break-all" outline>
                <Plus/>
                Add a choice
            </Button>
        {/if}
    </div>
{/if}