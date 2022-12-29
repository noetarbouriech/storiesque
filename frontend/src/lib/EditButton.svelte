<script lang="ts">
	import { Button } from "flowbite-svelte";
	import { Check, PencilSquare } from "svelte-heros-v2";
	import { createEventDispatcher } from "svelte";

    export let editMode: boolean = false;

    const dispatch = createEventDispatcher();

    // send event to parent component
    function save(): void {
        editMode = false;
        dispatch('save');
    }
</script>

{#if editMode}
    <Button on:click={save} class="!p-2 fixed right-6 bottom-6">
        <Check size="32" />
    </Button>
{:else}
    <Button on:click={() => {editMode = true}} class="!p-2 fixed right-6 bottom-6">
        <PencilSquare size="32" />
    </Button>
{/if}