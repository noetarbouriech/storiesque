<script lang="ts">
    import { Card, Avatar, Button } from "flowbite-svelte";
    import { env } from '$env/dynamic/public';
	import { onMount } from "svelte";
    export let username: String;
    export let id: number;
    export let has_img: boolean;
</script>
<Card padding='sm'>
    <div class="flex flex-col items-center">
        {#if has_img}
            <Avatar src="{env.PUBLIC_IMG_URL}/user/{id}.png" size="lg"></Avatar>
        {:else}
            <Avatar size="lg"></Avatar>
        {/if}
        <h5 class="text-xl font-medium text-gray-900 dark:text-white">@{username}</h5>
        <div class="flex mt-2 space-x-3 lg:mt-4">
            <Button href="/user/{username}">View profile</Button>
        </div>
    </div>
</Card>