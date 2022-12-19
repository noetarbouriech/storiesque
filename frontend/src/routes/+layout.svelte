<script>
	import '../app.postcss';
	import Header from '$lib/Header.svelte';
	import Footer from '$lib/Footer.svelte';
</script>

<Header />
<main class="container mx-auto py-16 px-8 dark:text-white ">
	<slot />
</main>
<Footer />