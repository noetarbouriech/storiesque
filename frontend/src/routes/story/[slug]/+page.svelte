<script lang="ts">
    import { Heading, P, Hr, Img, Button, Input, Textarea, Dropzone, Badge, A } from 'flowbite-svelte'
    import { env } from '$env/dynamic/public';
    import { userStore } from '../../../store';
    import slugify from 'slugify';
    import type { PageData } from './$types';
    import { Play } from 'svelte-heros-v2'
	import EditButton from '$lib/EditButton.svelte';
	import AddToShelf from '$lib/AddToShelf.svelte';
	import ImageUpload from '$lib/ImageUpload.svelte';

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

    let image_url: string = data.story.has_img ? `${env.PUBLIC_IMG_URL}/story/${data.story.id}.png?${new Date().getTime()}` : "/default_story.png"

</script>

<svelte:head>
  <title>{data.story.title} - Storiesque</title>
</svelte:head>

{#if editMode}
    <Input class="text-center text-xl mb-8" type="text" name="title" id="title" bind:value={data.story.title} required />
    <ImageUpload bind:has_img={data.story.has_img} id={data.story.id} type="story" alt={data.story.title} default_img="/default_story.png" bind:image_url={image_url}/>
    <Hr class="my-8" width="w-64"><P color="text-gray-500 dark:text-gray-400">DESCRIPTION</P></Hr>
    <Textarea id="description" name="description" class="text-center" bind:value={data.story.description} required />
{:else}
    <Heading class="text-center pb-8" tag="h1">{data.story.title}</Heading>
    <Img bind:src={image_url} alt="{data.story.title} cover image" class="max-h-[360px] mx-auto rounded-lg mb-8" />
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