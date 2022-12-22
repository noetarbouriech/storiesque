<script lang="ts">
    import { Card, Avatar, Button } from "flowbite-svelte";
    export let username : String;
</script>
<Card padding='sm'>
    <div class="flex flex-col items-center">
        <Avatar size="lg"></Avatar>
        <h5 class="text-xl font-medium text-gray-900 dark:text-white">@{username}</h5>
        <div class="flex mt-2 space-x-3 lg:mt-4">
            <Button href="/user/{username}">View profile</Button>
        </div>
    </div>
</Card>