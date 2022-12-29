<script lang="ts">
    import { Button, Navbar, NavBrand, NavLi, NavUl, NavHamburger, Avatar, Dropdown, DropdownItem, DropdownHeader, DropdownDivider, DarkMode } from 'flowbite-svelte'
    import { page } from '$app/stores';
    import { userStore } from '../store';
    import { env } from '$env/dynamic/public';
	import LoginModal from './LoginModal.svelte';
	import { goto } from '$app/navigation';

	let formModal: boolean = false;

    async function logout() {
        try {
            const response = await fetch(`${env.PUBLIC_API_URL}/logout`, {
                credentials: 'include',
            });
            if (response.ok) { 
                $userStore = {
                    username: "",
                    email: ""
                };
                goto("/");
            } else {
                alert(response.json());
            }
        } catch (error) {
            console.error('Error:', error);
        }
    }
</script>

<Navbar let:hidden let:toggle>
<NavBrand href="/">
    <img src="https://hotemoji.com/images/dl/f/open-book-emoji-by-twitter.png" class="mr-3 h-6 sm:h-9" alt="Storiesque Logo"/>
    <span class="self-center whitespace-nowrap text-xl font-semibold dark:text-white max-[300px]:hidden">Storiesque</span>
</NavBrand>
{#if $userStore.username == ""}
<div class="flex items-center md:order-2">
    <DarkMode />
    <Button on:click={() => formModal = true} size="sm">Log in</Button>
    <NavHamburger on:click={toggle} class1="w-full md:flex md:w-auto md:order-1"/>
</div>
{:else}
<div class="flex items-center md:order-2">
    <DarkMode />
    <Avatar class="cursor-pointer" id="avatar-menu" />
    <NavHamburger on:click={toggle} class1="w-full md:flex md:w-auto md:order-1"/>
</div>
<Dropdown placement="bottom" triggeredBy="#avatar-menu">
    <DropdownHeader>
    <span class="block text-sm">Hello there !</span>
    <span class="block text-md font-semibold">@{$userStore.username}</span>
    </DropdownHeader>
    <DropdownItem href="/user/{$userStore.username}">My Profile</DropdownItem>
    <DropdownItem href="/library">My Library</DropdownItem>
    <DropdownDivider />
    <DropdownItem on:click={logout}>Sign out</DropdownItem>
</Dropdown>
{/if}
<NavUl {hidden}>
    <NavLi href="/" active={$page.url.pathname == "/"}>Home</NavLi>
    <NavLi href="/story" active={$page.url.pathname.startsWith("/story")}>Stories</NavLi>
    <NavLi href="/user" active={$page.url.pathname.startsWith("/user")}>Users</NavLi>
</NavUl>
</Navbar>
<LoginModal bind:open={formModal} />