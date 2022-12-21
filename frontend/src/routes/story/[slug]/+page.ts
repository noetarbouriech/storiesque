import type { PageLoad } from './$types';

export const load: PageLoad = (async ({ params, fetch }) => {
  const res = await fetch(`http://localhost:3000/story/${params.slug}`);
  const story = await res.json();
  return { story };
})
