<script lang="ts">
    import { Heading, P, Hr, Img, Button } from 'flowbite-svelte'
    import type { PageData } from './$types';
    import { Play } from 'svelte-heros-v2'

    export let data: PageData;
</script>

<Heading class="text-center pb-8" tag="h1">{data.story.title}</Heading>

<Img src="https://loremflickr.com/840/480" alt="{data.story.title} cover image" alignment="mx-auto" class="rounded-lg mb-8" />

<Hr class="my-8" width="w-64"><P color="text-gray-500 dark:text-gray-400">DESCRIPTION</P></Hr>

<P class="mx-auto max-w-xl" align="center" weight="light" color="text-gray-500 dark:text-gray-400">{data.story.description}</P>

<div class="flex justify-center py-8 gap-x-4">
    <Button gradient color="greenToBlue" class="w-fit font-extrabold bg-gradient-to-r to-emerald-600 from-sky-400" href="/story/{data.story.id}/read">
        <Play class="mr-1" variation="solid"/>Begin story
    </Button>
</div>