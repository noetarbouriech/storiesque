import type { PageLoad } from './$types';
import { env } from '$env/dynamic/public';

export const load: PageLoad = (async ({ params, fetch }) => {
  const res_story = await fetch(`${env.PUBLIC_API_URL}/story/${params.slug}`);
  const story = await res_story.json();
  return { story };
})

