<script lang="ts">
    import { goto } from '$app/navigation';
    import UserTable from '$lib/admin/UserTable.svelte';
    import { Tabs, TabItem, Heading } from 'flowbite-svelte';
    import { userStore } from '../../store';
    import { BookOpen, UserCircle } from 'svelte-heros-v2';
    import { onMount } from 'svelte';

    onMount(async() => {
        if (!$userStore.is_admin) goto("/");
    });
</script>

<Heading class="text-center mb-10" tag="h1">Admin Section</Heading>

<Tabs style="full" defaultClass="flex rounded-lg divide-x divide-gray-200 shadow dark:divide-gray-700">
  <TabItem class="w-full" open>
    <div slot="title" class="flex justify-center items-center gap-2">
      <UserCircle />
      Users
    </div>
    <UserTable />
  </TabItem>
  <TabItem class="w-full">
    <div slot="title" class="flex justify-center items-center gap-2">
      <BookOpen />
      Stories
    </div>
    <p class="text-sm text-gray-500 dark:text-gray-400">Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>
  </TabItem>
</Tabs>