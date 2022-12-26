import type { PageLoad } from './$types';
import { env } from '$env/dynamic/public';

export const load: PageLoad = (async ({ params, fetch }) => {
  const res = await fetch(`${env.PUBLIC_API_URL}/story/${params.slug}`);
  const story = await res.json();
  return { story };
})
