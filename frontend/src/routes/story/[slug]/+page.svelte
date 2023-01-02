<script lang="ts">
    import { Heading, P, Hr, Img, Button, Input, Textarea, Dropzone } from 'flowbite-svelte'
    import { env } from '$env/dynamic/public';
    import { userStore } from '../../../store';
    import slugify from 'slugify';
    import type { PageData } from './$types';
    import { Play } from 'svelte-heros-v2'
	import EditButton from '$lib/EditButton.svelte';
	import AddToShelf from '$lib/AddToShelf.svelte';
	import { onMount } from 'svelte';

    export let data: PageData;
    let editMode: boolean;

    async function save() {
        await fetch(`${env.PUBLIC_API_URL}/story/${data.story.id}`, {
            method: 'PUT',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                title: data.story.title,
                description: data.story.description
            })
        })
    }

    let img: HTMLImageElement;

    onMount(() => {
        img.src = `${env.PUBLIC_IMG_URL}/story/${data.story.id}.png`;
        img.onerror = (): string => img.src = "/default_story.png"
    });

</script>

{#if editMode}
    <Input class="text-center text-xl mb-8" type="text" name="title" id="title" bind:value={data.story.title} required />
    <Dropzone id='dropzone'>
        <svg aria-hidden="true" class="mb-3 w-10 h-10 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"></path></svg>
        <p class="mb-2 text-sm text-gray-500 dark:text-gray-400"><span class="font-semibold">Click to upload</span> or drag and drop</p>
        <p class="text-xs text-gray-500 dark:text-gray-400">SVG, PNG, JPG or GIF (MAX. 800x400px)</p>
    </Dropzone>
    <Hr class="my-8" width="w-64"><P color="text-gray-500 dark:text-gray-400">DESCRIPTION</P></Hr>
    <Textarea id="description" name="description" class="text-center" bind:value={data.story.description} required />
{:else}
    <Heading class="text-center pb-8" tag="h1">{data.story.title}</Heading>
    <img bind:this={img} alt="{data.story.title} cover image" class="h-[360px] mx-auto rounded-lg mb-8" />
    <Hr class="my-8" width="w-64"><P color="text-gray-500 dark:text-gray-400">DESCRIPTION</P></Hr>
    <P class="mx-auto max-w-xl" align="center" weight="light" color="text-gray-500 dark:text-gray-400">{data.story.description}</P>
{/if}

<div class="flex justify-center items-center py-8 gap-x-4">
    <Button gradient color="greenToBlue" class="w-fit font-extrabold bg-gradient-to-r to-emerald-600 from-sky-400" href="/story/{slugify(data.story.title,{lower: true})}-{data.story.id}/read">
        <Play class="mr-1" variation="solid"/>Begin story
    </Button>
    {#if $userStore.username !== ""}
        <AddToShelf storyId={data.story.id} />
    {/if}
</div>
{#if data.story.author_name == $userStore.username || $userStore.is_admin}
    <EditButton bind:editMode={editMode} on:save={save}/>
{/if}