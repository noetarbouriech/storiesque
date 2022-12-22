import type { PageLoad } from './$types';

export const load: PageLoad = (async ({ params, fetch }) => {
  const res_user = await fetch(`http://localhost:3000/user/${params.username}`);
  const res_stories = await fetch(`http://localhost:3000/story`);
  const user = await res_user.json();
  const stories = await res_stories.json();
  return { user, stories };
})

