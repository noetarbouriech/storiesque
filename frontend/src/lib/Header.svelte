<script>
    import { Navbar, NavBrand, NavLi, NavUl, NavHamburger, Avatar, Dropdown, DropdownItem, DropdownHeader, DropdownDivider, DarkMode } from 'flowbite-svelte'
    import { page } from '$app/stores';
</script>

<Navbar let:hidden let:toggle>
<NavBrand href="/">
    <img src="https://flowbite.com/docs/images/logo.svg" class="mr-3 h-6 sm:h-9" alt="Flowbite Logo"/>
    <span class="self-center whitespace-nowrap text-xl font-semibold dark:text-white">Storiesque</span>
</NavBrand>
<div class="flex items-center md:order-2">
    <DarkMode />
    <Avatar id="avatar-menu">US</Avatar>
    <NavHamburger on:click={toggle} class1="w-full md:flex md:w-auto md:order-1"/>
</div>
<Dropdown placement="bottom" triggeredBy="#avatar-menu">
    <DropdownHeader>
    <span class="block text-sm">Hello there !</span>
    <span class="block text-md font-semibold">@Username</span>
    </DropdownHeader>
    <DropdownItem>Profile</DropdownItem>
    <DropdownItem>My Stories</DropdownItem>
    <DropdownItem>My Library</DropdownItem>
    <DropdownDivider />
    <DropdownItem>Sign out</DropdownItem>
</Dropdown>
<NavUl {hidden}>
    <NavLi href="/" active={$page.url.pathname == "/"}>Home</NavLi>
    <NavLi href="/story" active={$page.url.pathname.startsWith("/story")}>Stories</NavLi>
    <NavLi href="/user" active={$page.url.pathname.startsWith("/user")}>Users</NavLi>
</NavUl>
</Navbar>