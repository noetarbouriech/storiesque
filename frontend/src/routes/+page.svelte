<script>
import StoryCard from '$lib/StoryCard.svelte';
import { Hr, A, Heading, P, Span, Button } from 'flowbite-svelte'
import { PencilSquare } from 'svelte-heros-v2'
</script>

<Heading align="center" tag="h1" class="mt-8 mb-4" customSize="text-3xl font-extrabold  md:text-5xl lg:text-6xl">Why just read when you can <Span gradient>play</Span>?</Heading>
<P whitespace="preline" align="center">Storiesque allows you to write and read interactive stories in the manner of gamebooks.
    You're not forced to follow the story anymore, you're the one who's controlling <strong>your own path</strong>!
</P>

<div class="flex justify-center py-8 gap-x-4">
    <Button gradient color="greenToBlue" class="font-extrabold bg-gradient-to-r to-emerald-600 from-sky-400">
        <PencilSquare/>Write your story
    </Button>
    <Button color="alternative" href="/story">Read a story</Button>
</div>

<Heading align="center" tag="h2" class="mt-16 mb-2">
    Some of our personal <Span underline>favorites</Span>...
</Heading>
<div class="max-w-2xl mx-auto px-4 py-6 lg:max-w-7xl grid grid-cols-1 gap-y-10 sm:grid-cols-2 gap-x-8 lg:grid-cols-3">
    <StoryCard
        title="Great story"
        description="This one is really good. I can recommend it !"
        id={1}
    ></StoryCard>
    <StoryCard
        title="Great story"
        description="This one is really good. I can recommend it !"
        id={2}
    ></StoryCard>
    <StoryCard
        title="Great story"
        description="This one is really good. I can recommend it !"
        id={3}
    ></StoryCard>
</div>

<Hr class="my-4 mx-auto md:my-10" width="w-48" height="h-1"/>

<div class="flex justify-center py-4 gap-x-4">
    <A href="https://github.com/noetarbouriech/storiesque" textColor="text-gray-500 dark:text-gray-400 dark:hover:text-white" aClass="inline-flex justify-center items-center p-5 text-base font-medium bg-gray-50 rounded-lg hover:text-gray-900 hover:bg-gray-100 dark:bg-gray-800 dark:hover:bg-gray-700">
        <svg aria-hidden="true" class="mr-1 w-5 h-5" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 19c-5 1.5-5-2.5-7-3m14 6v-3.87a3.37 3.37 0 0 0-.94-2.61c3.14-.35 6.44-1.54 6.44-7A5.44 5.44 0 0 0 20 4.77 5.07 5.07 0 0 0 19.91 1S18.73.65 16 2.48a13.38 13.38 0 0 0-7 0C6.27.65 5.09 1 5.09 1A5.07 5.07 0 0 0 5 4.77a5.44 5.44 0 0 0-1.5 3.78c0 5.42 3.3 6.61 6.44 7A3.37 3.37 0 0 0 9 18.13V22"></path></svg>
        <span class="w-full">Check out this website sourcecode on Github</span>
        <svg aria-hidden="true" class="ml-2 w-6 h-6" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M12.293 5.293a1 1 0 011.414 0l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414-1.414L14.586 11H3a1 1 0 110-2h11.586l-2.293-2.293a1 1 0 010-1.414z" clip-rule="evenodd"></path></svg>
    </A>
</div>