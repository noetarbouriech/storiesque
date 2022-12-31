<script lang="ts">
	import { onMount } from "svelte";
    import { env } from '$env/dynamic/public';
    import { userStore } from '../../../../store';
    import { Toast, P, A, Hr, Button, Popover, TextPlaceholder, Spinner, Span, Input, Textarea } from 'flowbite-svelte'
    import type { PageData } from './$types';
    import { ArrowUturnLeft, BookOpen, InformationCircle, Plus } from 'svelte-heros-v2';
    import { page } from '$app/stores';
	import StoryCard from '$lib/StoryCard.svelte';
	import EditButton from "$lib/EditButton.svelte";

    export let data: PageData;
    type page = {
        id: number,
        action: string | number | undefined,
        body: string,
        choices: {page_id: number, action: string}[]
    }

    // default values
    let currPage: page = {
        id: 0,
        action: "",
        body: "",
        choices: []
    };
    let loading: boolean = true;
    let editMode: boolean;
    let history: Array<number> = [];

    async function changePage(pageId: number): Promise<void> {
        loading = true;
        setTimeout(async () => {
        currPage = await fetch(`${env.PUBLIC_API_URL}/page/${pageId}`, {
            credentials: 'include',
        })
        .then(r => r.json());
        loading = false;
        // add current page to history if its not already in history (like when going backward)
        // https://svelte.dev/tutorial/updating-arrays-and-objects
        if (history[history.length-1] != currPage.id) history = [...history, currPage.id];
        }, 300)
    }

    async function back(): Promise<void> {
        // remove last page from history and change to the page before this one
        history = history.slice(0, -1);
        changePage(history[history.length-1]);
    }

    async function addChoice(): Promise<void> {
        currPage.choices = [...currPage.choices,(await (await fetch(`${env.PUBLIC_API_URL}/page/${currPage.id}`, { 
            method: 'POST',
            credentials: 'include',
        })).json())];
    }

    async function save() {
        await fetch(`${env.PUBLIC_API_URL}/page/${currPage.id}`, {
            method: 'PUT',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                action: currPage.action,
                body: currPage.body
            })
        })

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
    <div class="mx-auto max-w-[1000px]">
        {#if data.story.first_page_id != currPage.id}
            <A on:click={back} class="mb-4 font-medium hover:underline"><ArrowUturnLeft /> Go back to the previous page</A>
            <P justify class="flex flex-wrap items-center w-full mb-8 md:text-xl" weight="light" size="lg" color="text-gray-500 dark:text-gray-400">
                <div class="mr-2 inline-flex items-center whitespace-nowrap">
                    <InformationCircle />You chose to:
                </div>
                {#if editMode}
                    <Input class="flex-1 min-w-fit" type="text" name="action" id="action" bind:value={currPage.action} required />
                {:else}
                    <Span italic class="font-thin break-all">{currPage.action}</Span>
                {/if}
            </P>
        {/if}
        {#if editMode}
            <Textarea class="flex justify-center mx-auto max-w-[1000px]" id="body" name="body" rows=15 bind:value={currPage.body} required />
        {:else}
            <P whitespace="preline" firstupper justify>{currPage.body}</P>
        {/if}
    </div>
    <Hr class="my-4 mx-auto md:my-8" width="w-48" height="h-1" />
    <div class="mx-auto flex flex-col justify-center max-w-fit">
        {#if currPage.choices.length === 0 && !editMode}
            <P class="md:text-xl" weight="light" size="lg" color="text-gray-500 dark:text-gray-400">
                You have reached an ending of this story. Thank you for playing !
            </P>
        {:else}
            {#each currPage.choices as choice}
            <Button on:click={() => changePage(choice.page_id)} class="mb-2 break-all">{choice.action}</Button>
            {/each}
        {/if}
        {#if editMode}
            <Button on:click={addChoice} class="break-all" outline>
                <Plus/>
                Add a choice
            </Button>
        {/if}
    </div>
    {#if data.story.author_name == $userStore.username || $userStore.is_admin}
        <EditButton bind:editMode={editMode} on:save={save}/>
    {/if}
{/if}