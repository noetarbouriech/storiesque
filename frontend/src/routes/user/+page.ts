import type { PageLoad } from './$types';

export const load: PageLoad = (async ({ fetch }) => {
  const res = await fetch('http://localhost:3000/user');
  const users = await res.json();
  return { users };
})