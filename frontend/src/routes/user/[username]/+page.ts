import type { PageLoad } from './$types';
import { env } from '$env/dynamic/public';

export const load: PageLoad = (async ({ params, fetch }) => {
  const res = await fetch(`${env.PUBLIC_API_URL}/user/${params.username}`);
  const user = await res.json();
  return { user };
})

