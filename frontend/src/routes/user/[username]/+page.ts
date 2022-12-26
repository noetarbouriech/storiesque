import type { PageLoad } from './$types';

export const load: PageLoad = (async ({ params, fetch }) => {
  const res = await fetch(`http://localhost:3000/user/${params.username}`);
  const user = await res.json();
  return { user };
})

