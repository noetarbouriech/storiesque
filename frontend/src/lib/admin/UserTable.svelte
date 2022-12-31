<script lang="ts">
    import { Table, TableHead, TableHeadCell, TableBody, TableBodyRow, TableBodyCell, Button, Toggle, Modal } from 'flowbite-svelte';
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

    // default values
    let selectedUser: user = {
        id: 0,
        username: "",
        email: "",
        is_admin: false,
    };

    onMount(async () => {
        loadUsers();
    });

    async function loadUsers(): Promise<void> {
        users = await (await fetch(`${env.PUBLIC_API_URL}/user`)).json();
    }

    // get executed every time editModal changes
    $: if (editModal === false) loadUsers();

    async function updateAdmin(userId: number): Promise<void> {
        try {
            const response = await fetch(`${env.PUBLIC_API_URL}/user/${userId}/admin`, {
                method: 'PUT',
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
                loadUsers();
            }
        } catch (error) {
            console.error('Error:', error);
        }

    }

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
                    <Button on:click={() => {deleteModal = true; selectedUser.id = user.id}} class="bg-red-600 !p-2"><Trash /></Button>
                </TableBodyCell>
            </TableBodyRow>
        {/each}
    </TableBody>
</Table>

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