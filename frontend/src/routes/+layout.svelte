<script>
	import '../app.postcss';
	import Header from '$lib/Header.svelte';
	import Footer from '$lib/Footer.svelte';
</script>

<div class="flex flex-col h-screen">
	<Header />
	<main class="container flex-grow mx-auto py-14 px-8 dark:text-white ">
		<slot />
	</main>
	<Footer />
</div>