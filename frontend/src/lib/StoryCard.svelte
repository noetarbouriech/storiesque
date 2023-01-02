<script lang="ts">
    import { Card, Button } from "flowbite-svelte";
    import slugify from 'slugify';
    import { env } from '$env/dynamic/public';
    import { onMount } from "svelte";
    export let title : string;
    export let description : string;
    export let id : number;
    export let author : string;

    if (description.length > 103) {
      description = description.substring(0,100) + "...";
    }

    let img: string;

    onMount(async () => {
        img = `${env.PUBLIC_IMG_URL}/story/${id}.png`;
        let response = await fetch(img);
        if (!response.ok) img = "/default_story.png"
    });

</script>

<Card img={img}>
  <h3 class="text-2xl font-bold tracking-tight text-gray-900 dark:text-white">{title}</h3>
  <h4 class="mb-2 text-l tracking-tight text-gray-900 dark:text-white">by <a class="font-bold" href="/user/{author}">@{author}</a></h4>
  <p class="mb-3 font-normal text-gray-700 dark:text-gray-400 leading-tight">
    {description}
  </p>
  <Button href="/story/{slugify(title,{lower: true})}-{id}" class="w-fit">
    Read <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5 ml-2"><path stroke-linecap="round" stroke-linejoin="round" d="M13.5 4.5L21 12m0 0l-7.5 7.5M21 12H3" /></svg>
  </Button>
</Card>