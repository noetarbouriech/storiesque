<script lang="ts">
	import { A, Badge, Dropzone, Img } from "flowbite-svelte";
    import { env } from '$env/dynamic/public';

    export let has_img: boolean;
    export let id: string;
    export let type: string;
    export let alt: string;
    export let default_img: string;
    export let image_url: string = `${env.PUBLIC_IMG_URL}/${type}/${id}.png?${new Date().getTime()}`;

    let fileUploaded: FileList;

    async function uploadImage() {
        let formData: FormData = new FormData();
        formData.append("id", id);
        formData.append("type", type);
        formData.append("file", fileUploaded[0]);
        await fetch(`${env.PUBLIC_API_URL}/image/upload`, {
            method: 'POST',
            credentials: 'include',
            body: formData
        });
        image_url = `${env.PUBLIC_IMG_URL}/${type}/${id}.png?${new Date().getTime()}`; // prevent caching
        has_img = true;
    }

    async function deleteImage() {
        await fetch(`${env.PUBLIC_API_URL}/image/${type}/${id}`, {
            method: 'DELETE',
            credentials: 'include'
        });
        image_url = default_img;
        has_img = false;
    }
</script>

{#if has_img}
    <div class="relative w-fit mx-auto">
        <Img bind:src={image_url} alt="{alt} cover image" class="max-h-[360px] mx-auto rounded-lg mb-8"/>
        <Badge large rounded index ><A aClass="" on:click={deleteImage}>X</A></Badge>
    </div>
{:else}
    <Dropzone on:change={uploadImage} bind:files={fileUploaded} id='dropzone' class="max-w-[640px] mx-auto">
        <svg aria-hidden="true" class="mb-3 w-10 h-10 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"></path></svg>
        <p class="mb-2 text-sm text-gray-500 dark:text-gray-400"><span class="font-semibold">Click to upload</span> or drag and drop</p>
        <p class="text-xs text-gray-500 dark:text-gray-400">PNG or JPEG (Max. 2MB)</p>
    </Dropzone>
{/if}