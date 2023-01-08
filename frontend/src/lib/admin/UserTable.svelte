<script lang="ts">
    import { Table, TableHead, TableHeadCell, TableBody, TableBodyRow, TableBodyCell, Button, Toggle, Modal, PaginationItem } from 'flowbite-svelte';
	import { onMount } from 'svelte';
    import { env } from '$env/dynamic/public';
	import { ArrowTopRightOnSquare, PencilSquare, Trash } from 'svelte-heros-v2';
	import EditUserModal from '../EditUserModal.svelte';

    let users: Array<any> = [];
    let deleteModal: boolean = false;
    let editModal: boolean = false;
    type user = {
        id: number,
        username: string,
        email: string,
        is_admin: boolean,
    }
    let page: number = 1;

    // default values
    let selectedUser: user = {
        id: 0,
        username: "",
        email: "",
        is_admin: false,
    };

    onMount(async () => {
        loadUsers(1);
    });

    async function loadUsers(page: number): Promise<void> {
        users = await (await fetch(`${env.PUBLIC_API_URL}/user?page=${page}`)).json();
    }

    // get executed every time editModal changes
    $: if (editModal === false) loadUsers(page);

    async function updateAdmin(userId: number): Promise<void> {
        try {
            const response = await fetch(`${env.PUBLIC_API_URL}/user/${userId}/admin`, {
                method: 'PATCH',
                credentials: 'include'
            });
            const data = await response.json();
            if (!response.ok) { 
                alert(data.message);
            }
        } catch (error) {
            console.error('Error:', error);
        }

    }

    async function deleteUser(userId: number): Promise<void> {
        try {
            const response = await fetch(`${env.PUBLIC_API_URL}/user/${userId}`, {
                method: 'DELETE',
                credentials: 'include'
            });
            const data = await response.json();
            if (!response.ok) { 
                alert(data.message);
            } else {
                loadUsers(page);
            }
        } catch (error) {
            console.error('Error:', error);
        }

    }

    const previous = async () => {
        if (page === 1) return;
        page--;
        loadUsers(page);
    };
    const next = async () => {
        if (users.length < 30) return;
        page++;
        loadUsers(page);
    };

</script>

<Table hoverable noborder>
    <TableHead>
        <TableHeadCell>ID</TableHeadCell>
        <TableHeadCell>Username</TableHeadCell>
        <TableHeadCell>Email</TableHeadCell>
        <TableHeadCell>Is admin ?</TableHeadCell>
        <TableHeadCell>Actions</TableHeadCell>
    </TableHead>
    <TableBody>
        {#each users as user}
            <TableBodyRow noborder>
                <TableBodyCell>{user.id}</TableBodyCell>
                <TableBodyCell>{user.username}</TableBodyCell>
                <TableBodyCell>{user.email}</TableBodyCell>
                <TableBodyCell>
                    <Toggle size="large" on:change={() => updateAdmin(user.id)} checked={user.is_admin} />
                </TableBodyCell>
                <TableBodyCell>
                    <Button href="/user/{user.username}" target="_blank" class="!p-2"><ArrowTopRightOnSquare /></Button>
                    <Button on:click={() => {
                        editModal = true; 
                        selectedUser = user
                        }} class="!p-2">
                        <PencilSquare />
                    </Button>
                    <Button on:click={() => {deleteModal = true; selectedUser.id = user.id}} class="!bg-red-600 !p-2"><Trash /></Button>
                </TableBodyCell>
            </TableBodyRow>
        {/each}
    </TableBody>
</Table>

<div class="mt-4 place-content-center flex space-x-3">
    <PaginationItem class="flex items-center" on:click={previous}>
      <svg class="mr-2 w-5 h-5" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M7.707 14.707a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 1.414L5.414 9H17a1 1 0 110 2H5.414l2.293 2.293a1 1 0 010 1.414z" clip-rule="evenodd"/></svg>
      Prev
    </PaginationItem>
    <PaginationItem class="flex items-center" on:click={next}>
      Next
      <svg class="ml-2 w-5 h-5" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M12.293 5.293a1 1 0 011.414 0l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414-1.414L14.586 11H3a1 1 0 110-2h11.586l-2.293-2.293a1 1 0 010-1.414z" clip-rule="evenodd"/></svg>
    </PaginationItem>
</div>

<Modal bind:open={deleteModal} size="xs" autoclose>
  <div class="text-center">
      <svg aria-hidden="true" class="mx-auto mb-4 w-14 h-14 text-gray-400 dark:text-gray-200" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
      <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">Are you sure you want to delete this user?</h3>
      <Button on:click={() => deleteUser(selectedUser.id)} color="red" class="mr-2">Yes, I'm sure</Button>
      <Button color='alternative'>No, cancel</Button>
  </div>
</Modal>
<EditUserModal 
    bind:open={editModal} 
    id={selectedUser.id}
    username={selectedUser.username}
    email={selectedUser.email}
/>