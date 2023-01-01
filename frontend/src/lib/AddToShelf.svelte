<script lang="ts">
	import { Tooltip } from "flowbite-svelte";
	import { Bookmark, BookmarkSlash } from "svelte-heros-v2";
    import { env } from '$env/dynamic/public';
	import { onMount } from "svelte";

    export let storyId: number;

    let onShelf: boolean = false;

    async function isOnShelf(): Promise<void> {
        try {
            const response = await fetch(`${env.PUBLIC_API_URL}/shelf/${storyId}`, {
                method: 'GET',
                credentials: 'include',
            });
            const data = await response.json();
            if (response.ok) { 
                onShelf = true;
            } else {
                onShelf = false;
            }
        } catch (error) {
            console.error('Error:', error);
        }
    }

    async function addToShelf(): Promise<void> {
        try {
            const response = await fetch(`${env.PUBLIC_API_URL}/shelf/${storyId}`, {
                method: 'POST',
                credentials: 'include',
            });
            const data = await response.json();
            if (response.ok) { 
                onShelf = true;
            }
        } catch (error) {
            console.error('Error:', error);
        }
    }

    async function removeFromShelf(): Promise<void> {
        try {
            const response = await fetch(`${env.PUBLIC_API_URL}/shelf/${storyId}`, {
                method: 'DELETE',
                credentials: 'include',
            });
            const data = await response.json();
            if (response.ok) { 
                onShelf = false;
            }
        } catch (error) {
            console.error('Error:', error);
        }
    }

    onMount(async(): Promise<void> => isOnShelf());
</script>

{#if onShelf}
    <BookmarkSlash on:click={removeFromShelf} class="cursor-pointer" size=35/>
    <Tooltip>Remove from shelf</Tooltip>
{:else}
    <Bookmark on:click={addToShelf} class="cursor-pointer" size=35/>
    <Tooltip>Add to shelf</Tooltip>
{/if}