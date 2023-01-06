<script lang="ts">
	import { onMount } from "svelte";
    import { env } from '$env/dynamic/public';
    import { userStore } from '../../../../store';
    import { Toast, P, A, Hr, Button, Popover, TextPlaceholder, Spinner, Span, Input, Textarea, Img, Badge } from 'flowbite-svelte'
    import type { PageData } from './$types';
    import { ArrowUturnLeft, BookOpen, InformationCircle, Plus } from 'svelte-heros-v2';
    import { page } from '$app/stores';
	import StoryCard from '$lib/StoryCard.svelte';
	import EditButton from "$lib/EditButton.svelte";
	import ImageUpload from "$lib/ImageUpload.svelte";

    export let data: PageData;
    type page = {
        id: number,
        action: string | number | undefined,
        body: string,
        has_img: boolean,
        choices: {page_id: number, action: string}[]
    }

    // default values
    let currPage: page = {
        id: 0,
        action: "",
        body: "",
        has_img: false,
        choices: []
    };
    let loading: boolean = true;
    let editMode: boolean;
    let history: Array<number> = [];
    let image_url: string; 

    async function changePage(pageId: number): Promise<void> {
        // save current page before changing page
        if (editMode) save();

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
            image_url = currPage.has_img ? `${env.PUBLIC_IMG_URL}/page/${currPage.id}.png?${new Date().getTime()}` : ""
        }, 300)
    }

    async function back(): Promise<void> {
        // save current page before going back
        if (editMode) save();

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

    async function removeChoice(choice: { page_id: number; action: string; }): Promise<void> {
        await fetch(`${env.PUBLIC_API_URL}/page/${choice.page_id}`, { 
            method: 'DELETE',
            credentials: 'include',
        })
        .then(() => {
            currPage.choices.splice(currPage.choices.indexOf(choice), 1);
            currPage.choices = currPage.choices;
        });
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

<svelte:head>
  <title>{data.story.title} - Storiesque</title>
</svelte:head>

<Toast transition={undefined} divClass="w-full max-w-fit p-4" class="mx-auto mb-8" simple>
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
        has_img={data.story.has_img}
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
            <div class="mb-8 mx-auto w-fit items-center text-center">
                <A on:click={back} class="mb-4 font-medium hover:underline"><ArrowUturnLeft /> Go back to the previous page</A>
                <P justify class="flex flex-wrap items-center w-full md:text-xl" weight="light" size="lg" color="text-gray-500 dark:text-gray-400">
                    <div class="mr-2 inline-flex items-center whitespace-nowrap">
                        <InformationCircle />You chose to:
                    </div>
                    {#if editMode}
                        <Input class="flex-1 min-w-fit" type="text" name="action" id="action" bind:value={currPage.action} required />
                    {:else}
                        <Span italic class="font-thin break-all">{currPage.action}</Span>
                    {/if}
                </P>
            </div>
        {/if}
        {#if editMode}
            <ImageUpload bind:has_img={currPage.has_img} id={String(currPage.id)} type="page" alt={String(currPage.action)} default_img="" bind:image_url={image_url}/>
            <Textarea class="mt-8 flex justify-center mx-auto max-w-[1000px]" id="body" name="body" rows=15 bind:value={currPage.body} required />
        {:else}
            {#if currPage.has_img && image_url != ""}
                <Img bind:src={image_url} alt="{currPage.action} cover image" class="max-h-[360px] mx-auto rounded-lg mb-8" />
            {/if}
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
            <div class="relative">
                <Button on:click={() => changePage(choice.page_id)} class="mb-2 w-full break-all">{choice.action}</Button>
                {#if editMode}
                    <Badge large rounded index ><A aClass="" on:click={() => removeChoice(choice)}>X</A></Badge>
                {/if}
            </div>
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