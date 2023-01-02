<script lang="ts">
    import { Heading, Avatar, Tooltip, Card } from 'flowbite-svelte';
    import { env } from '$env/dynamic/public';
    import type { PageData } from './$types';
    import { CheckBadge } from 'svelte-heros-v2';
    import StoryCard from '$lib/StoryCard.svelte';

    export let data: PageData;
</script>

<Card padding='sm' class="mx-auto">
{#if data.user.has_img}
    <Avatar src="{env.PUBLIC_IMG_URL}/user/{data.user.id}.png" class="mx-auto mb-8" size="xl" data-name=data.user.username></Avatar>
{:else}
    <Avatar class="mx-auto mb-8" size="xl" data-name=data.user.username></Avatar>
{/if}
<Heading customSize="text-3xl font-extrabold md:text-5xl" class="text-center pb-8" tag="h1">
    @{data.user.username} {#if data.user.is_admin} <CheckBadge size="50" class="inline"/><Tooltip>Admin</Tooltip>{/if}
</Heading>
<span class="text-center text-m text-gray-500 dark:text-gray-400">Joined in 2022</span>
<span class="text-center text-m text-gray-500 dark:text-gray-400">Has written {data.user.stories.length} stories</span>
</Card>

<Heading class="text-center mt-16 mb-4" tag="h1">Author of</Heading>
<div class="max-w-2xl mx-auto px-4 py-8 lg:max-w-7xl grid grid-cols-1 gap-y-10 sm:grid-cols-2 gap-x-8 lg:grid-cols-3">
    {#each data.user.stories as story}
        <StoryCard
            title={story.title}
            description={story.description}
            author={story.author_name}
            id={story.id}
            has_img={story.has_img}
        ></StoryCard>
    {/each}
</div>