import type { PageLoad } from './$types';
import { env } from '$env/dynamic/public';

export const load: PageLoad = (async ({ fetch }) => {
  const res = await fetch(`${env.PUBLIC_API_URL}/user`);
  const users = await res.json();
  return { users };
})