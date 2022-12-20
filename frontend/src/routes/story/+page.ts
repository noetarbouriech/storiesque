import type { PageLoad } from './$types';

export const load: PageLoad = (async ({ fetch }) => {
  const res = await fetch('http://localhost:3000/story');
  const stories = await res.json();
  return { stories };
})